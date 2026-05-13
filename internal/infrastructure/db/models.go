package db

import (
	"time"

	"log-parser/internal/domain"
)

type logRow struct {
	ID         int64     `db:"id"`
	Path       string    `db:"path"`
	Status     string    `db:"status"`
	NodeCount  int       `db:"node_count"`
	PortCount  int       `db:"port_count"`
	UploadedAt time.Time `db:"uploaded_at"`
}

func (r logRow) toDomain() domain.Log {
	return domain.Log{
		ID:         r.ID,
		Path:       r.Path,
		Status:     domain.LogStatus(r.Status),
		NodeCount:  r.NodeCount,
		PortCount:  r.PortCount,
		UploadedAt: r.UploadedAt,
	}
}

type nodeRow struct {
	ID              int64  `db:"id"`
	LogID           int64  `db:"log_id"`
	NodeGUID        string `db:"node_guid"`
	NodeDesc        string `db:"node_desc"`
	NodeType        int    `db:"node_type"`
	NumPorts        int    `db:"num_ports"`
	ClassVersion    int    `db:"class_version"`
	BaseVersion     int    `db:"base_version"`
	SystemImageGUID string `db:"system_image_guid"`
	PortGUID        string `db:"port_guid"`
}

func (r nodeRow) toDomain() domain.Node {
	return domain.Node{
		ID:              r.ID,
		LogID:           r.LogID,
		NodeGUID:        r.NodeGUID,
		NodeDesc:        r.NodeDesc,
		NodeType:        r.NodeType,
		NumPorts:        r.NumPorts,
		ClassVersion:    r.ClassVersion,
		BaseVersion:     r.BaseVersion,
		SystemImageGUID: r.SystemImageGUID,
		PortGUID:        r.PortGUID,
	}
}

func toNodes(rows []nodeRow) []domain.Node {
	out := make([]domain.Node, len(rows))
	for i, r := range rows {
		out[i] = r.toDomain()
	}
	return out
}

type portRow struct {
	ID     int64 `db:"id"`
	LogID  int64 `db:"log_id"`
	NodeID int64 `db:"node_id"`

	PortGUID string `db:"port_guid"`
	PortNum  int    `db:"port_num"`

	MKey                                string `db:"m_key"`
	GIDPrfx                             string `db:"gid_prfx"`
	MSMLID                              int    `db:"msm_lid"`
	LID                                 int    `db:"lid"`
	CapMsk                              int64  `db:"cap_msk"`
	MKeyLeasePeriod                     int    `db:"m_key_lease_period"`
	DiagCode                            int    `db:"diag_code"`
	LinkWidthActv                       int    `db:"link_width_actv"`
	LinkWidthSup                        int    `db:"link_width_sup"`
	LinkWidthEn                         int    `db:"link_width_en"`
	LocalPortNum                        int    `db:"local_port_num"`
	LinkSpeedEn                         int    `db:"link_speed_en"`
	LinkSpeedActv                       int    `db:"link_speed_actv"`
	LMC                                 int    `db:"lmc"`
	MKeyProtBits                        int    `db:"m_key_prot_bits"`
	LinkDownDefState                    int    `db:"link_down_def_state"`
	PortPhyState                        int    `db:"port_phy_state"`
	PortState                           int    `db:"port_state"`
	LinkSpeedSup                        int    `db:"link_speed_sup"`
	VLArbHighCap                        int    `db:"vl_arb_high_cap"`
	VLHighLimit                         int    `db:"vl_high_limit"`
	InitType                            int    `db:"init_type"`
	VLCap                               int    `db:"vl_cap"`
	MSMSL                               int    `db:"msmsl"`
	NMTU                                int    `db:"nmtu"`
	FilterRawOutb                       int    `db:"filter_raw_outb"`
	FilterRawInb                        int    `db:"filter_raw_inb"`
	PartEnfOutb                         int    `db:"part_enf_outb"`
	PartEnfInb                          int    `db:"part_enf_inb"`
	OpVLs                               int    `db:"op_vls"`
	HoQLife                             int    `db:"hoq_life"`
	VLStallCnt                          int    `db:"vl_stall_cnt"`
	MTUCap                              int    `db:"mtu_cap"`
	InitTypeReply                       int    `db:"init_type_reply"`
	VLArbLowCap                         int    `db:"vl_arb_low_cap"`
	PKeyViolations                      int    `db:"pkey_violations"`
	MKeyViolations                      int    `db:"mkey_violations"`
	SubnTmo                             int    `db:"subn_tmo"`
	MulticastPKeyTrapSuppressionEnabled int    `db:"multicast_pkey_trap_suppression_enabled"`
	ClientReregister                    int    `db:"client_reregister"`
	GUIDCap                             int    `db:"guid_cap"`
	QKeyViolations                      int    `db:"qkey_violations"`
	MaxCreditHint                       int    `db:"max_credit_hint"`
	OverrunErrs                         int    `db:"overrun_errs"`
	LocalPhyError                       int    `db:"local_phy_error"`
	RespTimeValue                       int    `db:"resp_time_value"`
	LinkRoundTripLatency                int    `db:"link_round_trip_latency"`
	OOOSLMask                           string `db:"ooosl_mask"`
	CapMsk2                             *int   `db:"cap_msk2"`
	FECActv                             *int   `db:"fec_actv"`
	RetransActv                         *int   `db:"retrans_actv"`
}

func (r portRow) toDomain() domain.Port {
	return domain.Port{
		ID:                                  r.ID,
		LogID:                               r.LogID,
		PortGUID:                            r.PortGUID,
		PortNum:                             r.PortNum,
		MKey:                                r.MKey,
		GIDPrfx:                             r.GIDPrfx,
		MSMLID:                              r.MSMLID,
		LID:                                 r.LID,
		CapMsk:                              r.CapMsk,
		MKeyLeasePeriod:                     r.MKeyLeasePeriod,
		DiagCode:                            r.DiagCode,
		LinkWidthActv:                       r.LinkWidthActv,
		LinkWidthSup:                        r.LinkWidthSup,
		LinkWidthEn:                         r.LinkWidthEn,
		LocalPortNum:                        r.LocalPortNum,
		LinkSpeedEn:                         r.LinkSpeedEn,
		LinkSpeedActv:                       r.LinkSpeedActv,
		LMC:                                 r.LMC,
		MKeyProtBits:                        r.MKeyProtBits,
		LinkDownDefState:                    r.LinkDownDefState,
		PortPhyState:                        r.PortPhyState,
		PortState:                           r.PortState,
		LinkSpeedSup:                        r.LinkSpeedSup,
		VLArbHighCap:                        r.VLArbHighCap,
		VLHighLimit:                         r.VLHighLimit,
		InitType:                            r.InitType,
		VLCap:                               r.VLCap,
		MSMSL:                               r.MSMSL,
		NMTU:                                r.NMTU,
		FilterRawOutb:                       r.FilterRawOutb,
		FilterRawInb:                        r.FilterRawInb,
		PartEnfOutb:                         r.PartEnfOutb,
		PartEnfInb:                          r.PartEnfInb,
		OpVLs:                               r.OpVLs,
		HoQLife:                             r.HoQLife,
		VLStallCnt:                          r.VLStallCnt,
		MTUCap:                              r.MTUCap,
		InitTypeReply:                       r.InitTypeReply,
		VLArbLowCap:                         r.VLArbLowCap,
		PKeyViolations:                      r.PKeyViolations,
		MKeyViolations:                      r.MKeyViolations,
		SubnTmo:                             r.SubnTmo,
		MulticastPKeyTrapSuppressionEnabled: r.MulticastPKeyTrapSuppressionEnabled,
		ClientReregister:                    r.ClientReregister,
		GUIDCap:                             r.GUIDCap,
		QKeyViolations:                      r.QKeyViolations,
		MaxCreditHint:                       r.MaxCreditHint,
		OverrunErrs:                         r.OverrunErrs,
		LocalPhyError:                       r.LocalPhyError,
		RespTimeValue:                       r.RespTimeValue,
		LinkRoundTripLatency:                r.LinkRoundTripLatency,
		OOOSLMask:                           r.OOOSLMask,
		CapMsk2:                             r.CapMsk2,
		FECActv:                             r.FECActv,
		RetransActv:                         r.RetransActv,
	}
}

func toPorts(rows []portRow) []domain.Port {
	out := make([]domain.Port, len(rows))
	for i, r := range rows {
		out[i] = r.toDomain()
	}
	return out
}

type nodeSwitchInfoRow struct {
	NodeID               int64 `db:"node_id"`
	LinearFDBCap         int   `db:"linear_fdb_cap"`
	RandomFDBCap         int   `db:"random_fdb_cap"`
	MCastFDBCap          int   `db:"mcast_fdb_cap"`
	LinearFDBTop         int   `db:"linear_fdb_top"`
	DefPort              int   `db:"def_port"`
	DefMCastPriPort      int   `db:"def_mcast_pri_port"`
	DefMCastNotPriPort   int   `db:"def_mcast_not_pri_port"`
	LifeTimeValue        int   `db:"life_time_value"`
	PortStateChange      int   `db:"port_state_change"`
	OptimizedSLVLMapping int   `db:"optimized_s_lvl_mapping"`
	LidsPerPort          int   `db:"lids_per_port"`
	PartEnfCap           int   `db:"part_enf_cap"`
	InbEnfCap            int   `db:"inb_enf_cap"`
	OutbEnfCap           int   `db:"outb_enf_cap"`
	FilterRawInbCap      int   `db:"filter_raw_inb_cap"`
	FilterRawOutbCap     int   `db:"filter_raw_outb_cap"`
	ENP0                 int   `db:"enp0"`
	MCastFDBTop          int   `db:"mcast_fdb_top"`
}

func (r nodeSwitchInfoRow) toDomain() domain.NodeSwitchInfo {
	return domain.NodeSwitchInfo{
		NodeGUID:             "",
		LinearFDBCap:         r.LinearFDBCap,
		RandomFDBCap:         r.RandomFDBCap,
		MCastFDBCap:          r.MCastFDBCap,
		LinearFDBTop:         r.LinearFDBTop,
		DefPort:              r.DefPort,
		DefMCastPriPort:      r.DefMCastPriPort,
		DefMCastNotPriPort:   r.DefMCastNotPriPort,
		LifeTimeValue:        r.LifeTimeValue,
		PortStateChange:      r.PortStateChange,
		OptimizedSLVLMapping: r.OptimizedSLVLMapping,
		LidsPerPort:          r.LidsPerPort,
		PartEnfCap:           r.PartEnfCap,
		InbEnfCap:            r.InbEnfCap,
		OutbEnfCap:           r.OutbEnfCap,
		FilterRawInbCap:      r.FilterRawInbCap,
		FilterRawOutbCap:     r.FilterRawOutbCap,
		ENP0:                 r.ENP0,
		MCastFDBTop:          r.MCastFDBTop,
	}
}

type nodeSystemInfoRow struct {
	NodeID       int64  `db:"node_id"`
	SerialNumber string `db:"serial_number"`
	PartNumber   string `db:"part_number"`
	Revision     string `db:"revision"`
	ProductName  string `db:"product_name"`
}

func (r nodeSystemInfoRow) toDomain() domain.NodeSystemInfo {
	return domain.NodeSystemInfo{
		NodeGUID:     "",
		SerialNumber: r.SerialNumber,
		PartNumber:   r.PartNumber,
		Revision:     r.Revision,
		ProductName:  r.ProductName,
	}
}

type nodeSharpInfoRow struct {
	NodeID                 int64 `db:"node_id"`
	Endianness             int   `db:"endianness"`
	EnableEndiannessPerJob int   `db:"enable_endianness_per_job"`
	ReproducibilityDisable int   `db:"reproducibility_disable"`
}

func (r nodeSharpInfoRow) toDomain() domain.NodeSharpInfo {
	return domain.NodeSharpInfo{
		NodeGUID:               "",
		Endianness:             r.Endianness,
		EnableEndiannessPerJob: r.EnableEndiannessPerJob,
		ReproducibilityDisable: r.ReproducibilityDisable,
	}
}
