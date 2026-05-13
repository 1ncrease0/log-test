package parser

import (
	"bytes"
	"encoding/csv"
	"strings"
	"testing"

	"github.com/stretchr/testify/require"
)

func portCSVHeader() []string {
	return []string{
		"NodeGuid", "PortGuid", "PortNum", "MKey", "GIDPrfx", "MSMLID", "LID", "CapMsk", "M_KeyLeasePeriod",
		"DiagCode", "LinkWidthActv", "LinkWidthSup", "LinkWidthEn", "LocalPortNum", "LinkSpeedEn", "LinkSpeedActv",
		"LMC", "MKeyProtBits", "LinkDownDefState", "PortPhyState", "PortState", "LinkSpeedSup", "VLArbHighCap",
		"VLHighLimit", "InitType", "VLCap", "MSMSL", "NMTU", "FilterRawOutb", "FilterRawInb", "PartEnfOutb",
		"PartEnfInb", "OpVLs", "HoQLife", "VLStallCnt", "MTUCap", "InitTypeReply", "VLArbLowCap", "PKeyViolations",
		"MKeyViolations", "SubnTmo", "MulticastPKeyTrapSuppressionEnabled", "ClientReregister", "GUIDCap",
		"QKeyViolations", "MaxCreditHint", "OverrunErrs", "LocalPhyError", "RespTimeValue", "LinkRoundTripLatency",
		"OOOSLMask", "CapMsk2", "FECActv", "RetransActv",
	}
}

func portColIndex(name string) int {
	for i, col := range portCSVHeader() {
		if col == name {
			return i
		}
	}
	panic("unknown column " + name)
}

func minimalPortRow() []string {
	h := portCSVHeader()
	row := make([]string, len(h))
	for i := range row {
		row[i] = "0"
	}
	row[portColIndex("NodeGuid")] = "0xnode"
	row[portColIndex("PortGuid")] = "0xport"
	row[portColIndex("PortNum")] = "1"
	row[portColIndex("MKey")] = ""
	row[portColIndex("GIDPrfx")] = ""
	row[portColIndex("OOOSLMask")] = ""
	row[portColIndex("CapMsk2")] = csvOptionalNA
	row[portColIndex("FECActv")] = csvOptionalNA
	row[portColIndex("RetransActv")] = csvOptionalNA
	return row
}

func writeCSVSection(t *testing.T, name string, rows [][]string) string {
	t.Helper()

	var buf bytes.Buffer
	buf.WriteString("START_" + name + "\n")
	w := csv.NewWriter(&buf)
	for _, r := range rows {
		require.NoError(t, w.Write(r))
	}
	w.Flush()
	require.NoError(t, w.Error())
	buf.WriteString("END_" + name + "\n")
	return buf.String()
}

func TestParseDBCSV_MissingRequiredSection(t *testing.T) {
	t.Parallel()

	_, err := parseDBCSV([]byte("START_PORTS\na\nEND_PORTS\n"))
	require.Error(t, err)
	require.Contains(t, err.Error(), "missing required section START_NODES")
}

func TestParseDBCSV_MinimalSuccess(t *testing.T) {
	t.Parallel()

	nodes := [][]string{
		{"NodeGUID", "NodeDesc", "NodeType", "NumPorts", "ClassVersion", "BaseVersion", "SystemImageGUID", "PortGUID"},
		{"0x10", "host-a", "1", "2", "3", "4", "sys-img", "0xport"},
	}
	ports := append([][]string{portCSVHeader()}, minimalPortRow())

	var b strings.Builder
	b.WriteString(writeCSVSection(t, "NODES", nodes))
	b.WriteString(writeCSVSection(t, "PORTS", ports))

	res, err := parseDBCSV([]byte(b.String()))
	require.NoError(t, err)
	require.Len(t, res.nodes, 1)
	require.Equal(t, "0x10", res.nodes[0].NodeGUID)
	require.Len(t, res.ports, 1)
	require.Equal(t, "0xnode", res.ports[0].NodeGUID)
	require.Equal(t, "0xport", res.ports[0].PortGUID)
}

func TestParseDBCSV_InvalidPortRow(t *testing.T) {
	t.Parallel()

	nodes := [][]string{
		{"NodeGUID", "NodeDesc", "NodeType", "NumPorts", "ClassVersion", "BaseVersion", "SystemImageGUID", "PortGUID"},
		{"0x10", "host-a", "1", "2", "3", "4", "sys-img", "0xport"},
	}
	badRow := minimalPortRow()
	badRow[portColIndex("PortNum")] = "not-int"
	ports := append([][]string{portCSVHeader()}, badRow)

	var b strings.Builder
	b.WriteString(writeCSVSection(t, "NODES", nodes))
	b.WriteString(writeCSVSection(t, "PORTS", ports))

	_, err := parseDBCSV([]byte(b.String()))
	require.Error(t, err)
	require.Contains(t, err.Error(), "section PORTS")
}
