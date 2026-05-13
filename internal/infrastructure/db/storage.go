package db

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"log/slog"

	"github.com/jackc/pgx/v5/pgconn"
	_ "github.com/jackc/pgx/v5/stdlib"
	"github.com/jmoiron/sqlx"

	"log-parser/internal/application"
	"log-parser/internal/domain"
)

type DB struct {
	log  *slog.Logger
	conn *sqlx.DB
}

func New(log *slog.Logger, address string) (*DB, error) {
	db, err := sqlx.Connect("pgx", address)
	if err != nil {
		log.Error("connection problem", "address", address, "error", err)
		return nil, err
	}
	return &DB{log: log, conn: db}, nil
}

func (db *DB) Close() error {
	return db.conn.Close()
}

func (db *DB) CreateLog(ctx context.Context, path string) (int64, error) {
	const q = `
		INSERT INTO logs (path, status, node_count, port_count)
		VALUES ($1, $2, 0, 0)
		RETURNING id`
	var id int64
	if err := db.conn.QueryRowContext(ctx, q, path, string(domain.LogStatusPending)).Scan(&id); err != nil {
		var pgErr *pgconn.PgError
		if errors.As(err, &pgErr) && pgErr.Code == "23505" && pgErr.ConstraintName == "logs_path_key" {
			return 0, fmt.Errorf("%w", application.ErrDuplicateLogPath)
		}
		return 0, fmt.Errorf("insert log: %w", err)
	}
	return id, nil
}

func (db *DB) SetStatus(ctx context.Context, logID int64, status domain.LogStatus) error {
	const q = `UPDATE logs SET status = $1 WHERE id = $2`
	if _, err := db.conn.ExecContext(ctx, q, string(status), logID); err != nil {
		return fmt.Errorf("set log status: %w", err)
	}
	return nil
}

func (db *DB) SaveResult(ctx context.Context, logID int64, result application.ParseResult) error {
	tx, err := db.conn.BeginTxx(ctx, nil)
	if err != nil {
		return fmt.Errorf("begin tx: %w", err)
	}
	defer func() {
		if rbErr := tx.Rollback(); rbErr != nil && !errors.Is(rbErr, sql.ErrTxDone) {
			db.log.Warn("transaction rollback failed", "log_id", logID, "error", rbErr)
		}
	}()

	guidToID, err := insertNodes(ctx, tx, logID, result.Nodes)
	if err != nil {
		return err
	}
	if err := insertPorts(ctx, tx, logID, guidToID, result.Ports); err != nil {
		return err
	}
	if err := insertSwitchInfos(ctx, tx, guidToID, result.SwitchInfos); err != nil {
		return err
	}
	if err := insertSystemInfos(ctx, tx, guidToID, result.SystemInfos); err != nil {
		return err
	}
	if err := insertSharpInfos(ctx, tx, guidToID, result.SharpInfos); err != nil {
		return err
	}

	const upd = `UPDATE logs SET status = $1, node_count = $2, port_count = $3 WHERE id = $4`
	if _, err := tx.ExecContext(ctx, upd,
		string(domain.LogStatusDone),
		len(result.Nodes),
		len(result.Ports),
		logID,
	); err != nil {
		return fmt.Errorf("update log: %w", err)
	}

	if err := tx.Commit(); err != nil {
		return fmt.Errorf("commit: %w", err)
	}
	return nil
}

func (db *DB) Log(ctx context.Context, id int64) (domain.Log, error) {
	var row logRow
	if err := db.conn.GetContext(ctx, &row, `SELECT * FROM logs WHERE id = $1`, id); err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return domain.Log{}, application.ErrNotFound
		}
		return domain.Log{}, fmt.Errorf("log %d: %w", id, err)
	}
	return row.toDomain(), nil
}

func (db *DB) Nodes(ctx context.Context, logID int64) ([]domain.Node, error) {
	var rows []nodeRow
	if err := db.conn.SelectContext(ctx, &rows, `SELECT * FROM nodes WHERE log_id = $1`, logID); err != nil {
		return nil, fmt.Errorf("nodes log=%d: %w", logID, err)
	}
	return toNodes(rows), nil
}

func (db *DB) Node(ctx context.Context, nodeID int64) (domain.Node, error) {
	var row nodeRow
	if err := db.conn.GetContext(ctx, &row, `SELECT * FROM nodes WHERE id = $1`, nodeID); err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return domain.Node{}, application.ErrNotFound
		}
		return domain.Node{}, fmt.Errorf("node %d: %w", nodeID, err)
	}
	return row.toDomain(), nil
}

func (db *DB) NodeDetail(ctx context.Context, nodeID int64) (application.NodeDetail, error) {
	node, err := db.Node(ctx, nodeID)
	if err != nil {
		return application.NodeDetail{}, err
	}

	detail := application.NodeDetail{Node: node}

	var sw nodeSwitchInfoRow
	if err := db.conn.GetContext(ctx, &sw, `SELECT * FROM nodes_switch_info WHERE node_id = $1`, nodeID); err != nil {
		if !errors.Is(err, sql.ErrNoRows) {
			return application.NodeDetail{}, fmt.Errorf("nodes_switch_info node_id=%d: %w", nodeID, err)
		}
	} else {
		v := sw.toDomain()
		detail.SwitchInfo = &v
	}

	var sys nodeSystemInfoRow
	if err := db.conn.GetContext(ctx, &sys, `SELECT * FROM nodes_system_info WHERE node_id = $1`, nodeID); err != nil {
		if !errors.Is(err, sql.ErrNoRows) {
			return application.NodeDetail{}, fmt.Errorf("nodes_system_info node_id=%d: %w", nodeID, err)
		}
	} else {
		v := sys.toDomain()
		detail.SystemInfo = &v
	}

	var sharp nodeSharpInfoRow
	if err := db.conn.GetContext(ctx, &sharp, `SELECT * FROM nodes_sharp_info WHERE node_id = $1`, nodeID); err != nil {
		if !errors.Is(err, sql.ErrNoRows) {
			return application.NodeDetail{}, fmt.Errorf("nodes_sharp_info node_id=%d: %w", nodeID, err)
		}
	} else {
		v := sharp.toDomain()
		detail.SharpInfo = &v
	}

	return detail, nil
}

func (db *DB) Ports(ctx context.Context, nodeID int64) ([]domain.Port, error) {
	var rows []portRow
	if err := db.conn.SelectContext(ctx, &rows, `SELECT * FROM ports WHERE node_id = $1`, nodeID); err != nil {
		return nil, fmt.Errorf("ports node=%d: %w", nodeID, err)
	}
	return toPorts(rows), nil
}

func insertNodes(ctx context.Context, tx *sqlx.Tx, logID int64, nodes []domain.Node) (map[string]int64, error) {
	const q = `
		INSERT INTO nodes (
			log_id, node_guid, node_desc, node_type, num_ports,
			class_version, base_version, system_image_guid, port_guid
		) VALUES ($1,$2,$3,$4,$5,$6,$7,$8,$9)
		RETURNING id`

	out := make(map[string]int64, len(nodes))
	for _, n := range nodes {
		var id int64
		err := tx.QueryRowContext(ctx, q,
			logID, n.NodeGUID, n.NodeDesc, n.NodeType, n.NumPorts,
			n.ClassVersion, n.BaseVersion, n.SystemImageGUID, n.PortGUID,
		).Scan(&id)
		if err != nil {
			return nil, fmt.Errorf("insert node %q: %w", n.NodeGUID, err)
		}
		out[n.NodeGUID] = id
	}
	return out, nil
}

func insertPorts(ctx context.Context, tx *sqlx.Tx, logID int64, guidToID map[string]int64, ports []domain.Port) error {
	const q = `
		INSERT INTO ports (
			log_id, node_id, port_guid, port_num,
			m_key, gid_prfx, msm_lid, lid, cap_msk, m_key_lease_period, diag_code,
			link_width_actv, link_width_sup, link_width_en, local_port_num, link_speed_en, link_speed_actv,
			lmc, m_key_prot_bits, link_down_def_state, port_phy_state, port_state, link_speed_sup,
			vl_arb_high_cap, vl_high_limit, init_type, vl_cap, msmsl, nmtu,
			filter_raw_outb, filter_raw_inb, part_enf_outb, part_enf_inb, op_vls, hoq_life, vl_stall_cnt,
			mtu_cap, init_type_reply, vl_arb_low_cap, pkey_violations, mkey_violations, subn_tmo,
			multicast_pkey_trap_suppression_enabled, client_reregister, guid_cap, qkey_violations, max_credit_hint,
			overrun_errs, local_phy_error, resp_time_value, link_round_trip_latency, ooosl_mask,
			cap_msk2, fec_actv, retrans_actv
		) VALUES (
			$1,$2,$3,$4,
			$5,$6,$7,$8,$9,$10,$11,
			$12,$13,$14,$15,$16,$17,
			$18,$19,$20,$21,$22,$23,
			$24,$25,$26,$27,$28,$29,
			$30,$31,$32,$33,$34,$35,$36,
			$37,$38,$39,$40,$41,$42,
			$43,$44,$45,$46,$47,
			$48,$49,$50,$51,$52,
			$53,$54,$55
		)`

	for _, p := range ports {
		nodeID, ok := guidToID[p.NodeGUID]
		if !ok {
			return fmt.Errorf("port references unknown node_guid %q", p.NodeGUID)
		}
		_, err := tx.ExecContext(ctx, q,
			logID, nodeID, p.PortGUID, p.PortNum,
			p.MKey, p.GIDPrfx, p.MSMLID, p.LID, p.CapMsk, p.MKeyLeasePeriod, p.DiagCode,
			p.LinkWidthActv, p.LinkWidthSup, p.LinkWidthEn, p.LocalPortNum, p.LinkSpeedEn, p.LinkSpeedActv,
			p.LMC, p.MKeyProtBits, p.LinkDownDefState, p.PortPhyState, p.PortState, p.LinkSpeedSup,
			p.VLArbHighCap, p.VLHighLimit, p.InitType, p.VLCap, p.MSMSL, p.NMTU,
			p.FilterRawOutb, p.FilterRawInb, p.PartEnfOutb, p.PartEnfInb, p.OpVLs, p.HoQLife, p.VLStallCnt,
			p.MTUCap, p.InitTypeReply, p.VLArbLowCap, p.PKeyViolations, p.MKeyViolations, p.SubnTmo,
			p.MulticastPKeyTrapSuppressionEnabled, p.ClientReregister, p.GUIDCap, p.QKeyViolations, p.MaxCreditHint,
			p.OverrunErrs, p.LocalPhyError, p.RespTimeValue, p.LinkRoundTripLatency, p.OOOSLMask,
			p.CapMsk2, p.FECActv, p.RetransActv,
		)
		if err != nil {
			return fmt.Errorf("insert port node=%q port_num=%d: %w", p.NodeGUID, p.PortNum, err)
		}
	}
	return nil
}

func insertSwitchInfos(ctx context.Context, tx *sqlx.Tx, guidToID map[string]int64, infos []domain.NodeSwitchInfo) error {
	const q = `
		INSERT INTO nodes_switch_info (
			node_id,
			linear_fdb_cap, random_fdb_cap, mcast_fdb_cap, linear_fdb_top,
			def_port, def_mcast_pri_port, def_mcast_not_pri_port, life_time_value, port_state_change,
			optimized_s_lvl_mapping, lids_per_port, part_enf_cap, inb_enf_cap, outb_enf_cap,
			filter_raw_inb_cap, filter_raw_outb_cap, enp0, mcast_fdb_top
		) VALUES (
			$1,$2,$3,$4,$5,$6,$7,$8,$9,$10,$11,$12,$13,$14,$15,$16,$17,$18,$19
		)`

	for _, s := range infos {
		nodeID, ok := guidToID[s.NodeGUID]
		if !ok {
			continue
		}
		_, err := tx.ExecContext(ctx, q,
			nodeID,
			s.LinearFDBCap, s.RandomFDBCap, s.MCastFDBCap, s.LinearFDBTop,
			s.DefPort, s.DefMCastPriPort, s.DefMCastNotPriPort, s.LifeTimeValue, s.PortStateChange,
			s.OptimizedSLVLMapping, s.LidsPerPort, s.PartEnfCap, s.InbEnfCap, s.OutbEnfCap,
			s.FilterRawInbCap, s.FilterRawOutbCap, s.ENP0, s.MCastFDBTop,
		)
		if err != nil {
			return fmt.Errorf("insert switch info node=%q: %w", s.NodeGUID, err)
		}
	}
	return nil
}

func insertSystemInfos(ctx context.Context, tx *sqlx.Tx, guidToID map[string]int64, infos []domain.NodeSystemInfo) error {
	const q = `
		INSERT INTO nodes_system_info (node_id, serial_number, part_number, revision, product_name)
		VALUES ($1, $2, $3, $4, $5)`

	for _, s := range infos {
		nodeID, ok := guidToID[s.NodeGUID]
		if !ok {
			continue
		}
		if _, err := tx.ExecContext(ctx, q, nodeID, s.SerialNumber, s.PartNumber, s.Revision, s.ProductName); err != nil {
			return fmt.Errorf("insert system info node=%q: %w", s.NodeGUID, err)
		}
	}
	return nil
}

func insertSharpInfos(ctx context.Context, tx *sqlx.Tx, guidToID map[string]int64, infos []domain.NodeSharpInfo) error {
	const q = `
		INSERT INTO nodes_sharp_info (node_id, endianness, enable_endianness_per_job, reproducibility_disable)
		VALUES ($1, $2, $3, $4)`

	for _, s := range infos {
		nodeID, ok := guidToID[s.NodeGUID]
		if !ok {
			continue
		}
		if _, err := tx.ExecContext(ctx, q, nodeID, s.Endianness, s.EnableEndiannessPerJob, s.ReproducibilityDisable); err != nil {
			return fmt.Errorf("insert sharp info node=%q: %w", s.NodeGUID, err)
		}
	}
	return nil
}
