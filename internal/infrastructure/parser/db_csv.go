package parser

import (
	"bufio"
	"encoding/csv"
	"fmt"
	"strconv"
	"strings"

	"log-parser/internal/domain"
)

type dbCSVResult struct {
	nodes       []domain.Node
	ports       []domain.Port
	switchInfos []domain.NodeSwitchInfo
	systemInfos []domain.NodeSystemInfo
}

func parseDBCSV(data []byte) (dbCSVResult, error) {
	sections, err := splitSections(data)
	if err != nil {
		return dbCSVResult{}, err
	}

	for _, name := range []string{"NODES", "PORTS"} {
		if _, ok := sections[name]; !ok {
			return dbCSVResult{}, fmt.Errorf("missing required section START_%s", name)
		}
	}

	var r dbCSVResult

	if r.nodes, err = parseCSVRows(sections["NODES"], parseNodeRow); err != nil {
		return dbCSVResult{}, fmt.Errorf("section NODES: %w", err)
	}
	if r.ports, err = parseCSVRows(sections["PORTS"], parsePortRow); err != nil {
		return dbCSVResult{}, fmt.Errorf("section PORTS: %w", err)
	}
	if s, ok := sections["SWITCHES"]; ok {
		if r.switchInfos, err = parseCSVRows(s, parseSwitchRow); err != nil {
			return dbCSVResult{}, fmt.Errorf("section SWITCHES: %w", err)
		}
	}
	if s, ok := sections["SYSTEM_GENERAL_INFORMATION"]; ok {
		if r.systemInfos, err = parseCSVRows(s, parseSystemInfoRow); err != nil {
			return dbCSVResult{}, fmt.Errorf("section SYSTEM_GENERAL_INFORMATION: %w", err)
		}
	}

	return r, nil
}

type rowParser struct {
	row []string
	idx map[string]int
	err error
}

func newRowParser(row []string, idx map[string]int) *rowParser {
	return &rowParser{row: row, idx: idx}
}

func (r *rowParser) readString(dst *string, col string) {
	if r.err != nil {
		return
	}
	*dst, r.err = getCol(r.row, r.idx, col)
}

func (r *rowParser) readInt(dst *int, col string) {
	if r.err != nil {
		return
	}
	*dst, r.err = intCol(r.row, r.idx, col)
}

func (r *rowParser) readInt64(dst *int64, col string) {
	if r.err != nil {
		return
	}
	*dst, r.err = int64Col(r.row, r.idx, col)
}

func (r *rowParser) readIntOpt(dst **int, col string) {
	if r.err != nil {
		return
	}
	*dst, r.err = optIntCol(r.row, r.idx, col)
}

func parseNodeRow(row []string, idx map[string]int) (domain.Node, error) {
	var n domain.Node
	r := newRowParser(row, idx)
	r.readString(&n.NodeGUID, "NodeGUID")
	r.readString(&n.NodeDesc, "NodeDesc")
	r.readInt(&n.NodeType, "NodeType")
	r.readInt(&n.NumPorts, "NumPorts")
	r.readInt(&n.ClassVersion, "ClassVersion")
	r.readInt(&n.BaseVersion, "BaseVersion")
	r.readString(&n.SystemImageGUID, "SystemImageGUID")
	r.readString(&n.PortGUID, "PortGUID")
	return n, r.err
}

func parsePortRow(row []string, idx map[string]int) (domain.Port, error) {
	var p domain.Port
	r := newRowParser(row, idx)
	r.readString(&p.NodeGUID, "NodeGuid")
	r.readString(&p.PortGUID, "PortGuid")
	r.readInt(&p.PortNum, "PortNum")
	r.readString(&p.MKey, "MKey")
	r.readString(&p.GIDPrfx, "GIDPrfx")
	r.readInt(&p.MSMLID, "MSMLID")
	r.readInt(&p.LID, "LID")
	r.readInt64(&p.CapMsk, "CapMsk")
	r.readInt(&p.MKeyLeasePeriod, "M_KeyLeasePeriod")
	r.readInt(&p.DiagCode, "DiagCode")
	r.readInt(&p.LinkWidthActv, "LinkWidthActv")
	r.readInt(&p.LinkWidthSup, "LinkWidthSup")
	r.readInt(&p.LinkWidthEn, "LinkWidthEn")
	r.readInt(&p.LocalPortNum, "LocalPortNum")
	r.readInt(&p.LinkSpeedEn, "LinkSpeedEn")
	r.readInt(&p.LinkSpeedActv, "LinkSpeedActv")
	r.readInt(&p.LMC, "LMC")
	r.readInt(&p.MKeyProtBits, "MKeyProtBits")
	r.readInt(&p.LinkDownDefState, "LinkDownDefState")
	r.readInt(&p.PortPhyState, "PortPhyState")
	r.readInt(&p.PortState, "PortState")
	r.readInt(&p.LinkSpeedSup, "LinkSpeedSup")
	r.readInt(&p.VLArbHighCap, "VLArbHighCap")
	r.readInt(&p.VLHighLimit, "VLHighLimit")
	r.readInt(&p.InitType, "InitType")
	r.readInt(&p.VLCap, "VLCap")
	r.readInt(&p.MSMSL, "MSMSL")
	r.readInt(&p.NMTU, "NMTU")
	r.readInt(&p.FilterRawOutb, "FilterRawOutb")
	r.readInt(&p.FilterRawInb, "FilterRawInb")
	r.readInt(&p.PartEnfOutb, "PartEnfOutb")
	r.readInt(&p.PartEnfInb, "PartEnfInb")
	r.readInt(&p.OpVLs, "OpVLs")
	r.readInt(&p.HoQLife, "HoQLife")
	r.readInt(&p.VLStallCnt, "VLStallCnt")
	r.readInt(&p.MTUCap, "MTUCap")
	r.readInt(&p.InitTypeReply, "InitTypeReply")
	r.readInt(&p.VLArbLowCap, "VLArbLowCap")
	r.readInt(&p.PKeyViolations, "PKeyViolations")
	r.readInt(&p.MKeyViolations, "MKeyViolations")
	r.readInt(&p.SubnTmo, "SubnTmo")
	r.readInt(&p.MulticastPKeyTrapSuppressionEnabled, "MulticastPKeyTrapSuppressionEnabled")
	r.readInt(&p.ClientReregister, "ClientReregister")
	r.readInt(&p.GUIDCap, "GUIDCap")
	r.readInt(&p.QKeyViolations, "QKeyViolations")
	r.readInt(&p.MaxCreditHint, "MaxCreditHint")
	r.readInt(&p.OverrunErrs, "OverrunErrs")
	r.readInt(&p.LocalPhyError, "LocalPhyError")
	r.readInt(&p.RespTimeValue, "RespTimeValue")
	r.readInt(&p.LinkRoundTripLatency, "LinkRoundTripLatency")
	r.readString(&p.OOOSLMask, "OOOSLMask")
	r.readIntOpt(&p.CapMsk2, "CapMsk2")
	r.readIntOpt(&p.FECActv, "FECActv")
	r.readIntOpt(&p.RetransActv, "RetransActv")
	return p, r.err
}

func parseSwitchRow(row []string, idx map[string]int) (domain.NodeSwitchInfo, error) {
	var s domain.NodeSwitchInfo
	r := newRowParser(row, idx)
	r.readString(&s.NodeGUID, "NodeGUID")
	r.readInt(&s.LinearFDBCap, "LinearFDBCap")
	r.readInt(&s.RandomFDBCap, "RandomFDBCap")
	r.readInt(&s.MCastFDBCap, "MCastFDBCap")
	r.readInt(&s.LinearFDBTop, "LinearFDBTop")
	r.readInt(&s.DefPort, "DefPort")
	r.readInt(&s.DefMCastPriPort, "DefMCastPriPort")
	r.readInt(&s.DefMCastNotPriPort, "DefMCastNotPriPort")
	r.readInt(&s.LifeTimeValue, "LifeTimeValue")
	r.readInt(&s.PortStateChange, "PortStateChange")
	r.readInt(&s.OptimizedSLVLMapping, "OptimizedSLVLMapping")
	r.readInt(&s.LidsPerPort, "LidsPerPort")
	r.readInt(&s.PartEnfCap, "PartEnfCap")
	r.readInt(&s.InbEnfCap, "InbEnfCap")
	r.readInt(&s.OutbEnfCap, "OutbEnfCap")
	r.readInt(&s.FilterRawInbCap, "FilterRawInbCap")
	r.readInt(&s.FilterRawOutbCap, "FilterRawOutbCap")
	r.readInt(&s.ENP0, "ENP0")
	r.readInt(&s.MCastFDBTop, "MCastFDBTop")
	return s, r.err
}

func parseSystemInfoRow(row []string, idx map[string]int) (domain.NodeSystemInfo, error) {
	var s domain.NodeSystemInfo
	r := newRowParser(row, idx)
	r.readString(&s.NodeGUID, "NodeGuid")
	r.readString(&s.SerialNumber, "SerialNumber")
	r.readString(&s.PartNumber, "PartNumber")
	r.readString(&s.Revision, "Revision")
	r.readString(&s.ProductName, "ProductName")
	return s, r.err
}

func parseCSVRows[T any](section string, parseRow func([]string, map[string]int) (T, error)) ([]T, error) {
	header, rows, err := parseCSVSection(section)
	if err != nil {
		return nil, err
	}
	idx := colIdx(header)
	out := make([]T, 0, len(rows))
	for i, row := range rows {
		item, err := parseRow(row, idx)
		if err != nil {
			return nil, fmt.Errorf("row %d: %w", i+1, err)
		}
		out = append(out, item)
	}
	return out, nil
}

func splitSections(data []byte) (map[string]string, error) {
	sections := make(map[string]string)
	var current string
	var buf strings.Builder

	scanner := bufio.NewScanner(strings.NewReader(string(data)))
	for scanner.Scan() {
		line := scanner.Text()
		switch {
		case strings.HasPrefix(line, "START_"):
			current = strings.TrimPrefix(line, "START_")
			buf.Reset()
		case strings.HasPrefix(line, "END_"):
			if current != "" {
				sections[current] = buf.String()
				current = ""
			}
		default:
			if current != "" {
				buf.WriteString(line)
				buf.WriteByte('\n')
			}
		}
	}
	return sections, scanner.Err()
}

func parseCSVSection(data string) ([]string, [][]string, error) {
	r := csv.NewReader(strings.NewReader(data))
	r.TrimLeadingSpace = true
	all, err := r.ReadAll()
	if err != nil {
		return nil, nil, fmt.Errorf("csv: %w", err)
	}
	if len(all) < 1 {
		return nil, nil, fmt.Errorf("csv: empty section")
	}
	return all[0], all[1:], nil
}

func colIdx(header []string) map[string]int {
	idx := make(map[string]int, len(header))
	for i, h := range header {
		idx[strings.TrimSpace(h)] = i
	}
	return idx
}

func getCol(row []string, idx map[string]int, col string) (string, error) {
	i, ok := idx[col]
	if !ok {
		return "", fmt.Errorf("column %q not found", col)
	}
	if i >= len(row) {
		return "", fmt.Errorf("column %q: index %d out of range (row len %d)", col, i, len(row))
	}
	return strings.TrimSpace(row[i]), nil
}

func intCol(row []string, idx map[string]int, col string) (int, error) {
	s, err := getCol(row, idx, col)
	if err != nil {
		return 0, err
	}
	n, err := strconv.Atoi(s)
	if err != nil {
		return 0, fmt.Errorf("column %q: %w", col, err)
	}
	return n, nil
}

func int64Col(row []string, idx map[string]int, col string) (int64, error) {
	s, err := getCol(row, idx, col)
	if err != nil {
		return 0, err
	}
	n, err := strconv.ParseInt(s, 10, 64)
	if err != nil {
		u, err2 := strconv.ParseUint(s, 10, 64)
		if err2 != nil {
			return 0, fmt.Errorf("column %q: %w", col, err)
		}
		return int64(u), nil
	}
	return n, nil
}

func optIntCol(row []string, idx map[string]int, col string) (*int, error) {
	s, err := getCol(row, idx, col)
	if err != nil {
		return nil, err
	}
	if s == "N/A" {
		return nil, nil
	}
	n, err := strconv.Atoi(s)
	if err != nil {
		return nil, fmt.Errorf("column %q: %w", col, err)
	}
	return &n, nil
}
