package parser

import (
	"archive/zip"
	"io"
	"log/slog"
	"os"
	"path/filepath"
	"strconv"
	"strings"
	"testing"

	"github.com/stretchr/testify/require"

	"log-parser/internal/application"
)

func diagnosticPortRow(nodeGUID, portGUID string, portNum int, capMsk int64) []string {
	r := minimalPortRow()
	r[portColIndex("NodeGuid")] = nodeGUID
	r[portColIndex("PortGuid")] = portGUID
	r[portColIndex("PortNum")] = strconv.Itoa(portNum)
	r[portColIndex("CapMsk")] = strconv.FormatInt(capMsk, 10)
	return r
}

func diagnosticSwitchHeader() []string {
	return []string{
		"NodeGUID", "LinearFDBCap", "RandomFDBCap", "MCastFDBCap", "LinearFDBTop", "DefPort", "DefMCastPriPort",
		"DefMCastNotPriPort", "LifeTimeValue", "PortStateChange", "OptimizedSLVLMapping", "LidsPerPort", "PartEnfCap",
		"InbEnfCap", "OutbEnfCap", "FilterRawInbCap", "FilterRawOutbCap", "ENP0", "MCastFDBTop",
	}
}

func diagnosticSwitchRow(nodeGUID string) []string {
	return []string{
		nodeGUID, "49152", "0", "8192", "78", "0", "255", "255", "18", "0", "3", "0", "32", "1", "1", "1", "1", "0", "49183",
	}
}

func diagnosticSystemHeader() []string {
	return []string{"NodeGuid", "SerialNumber", "PartNumber", "Revision", "ProductName"}
}

func diagnosticSystemRow(nodeGUID, serial, product string) []string {
	return []string{nodeGUID, serial, "PN-1", "A", product}
}

func buildValidDiagnosticDBCSV(t *testing.T) []byte {
	t.Helper()
	nodes := [][]string{
		{"NodeGUID", "NodeDesc", "NodeType", "NumPorts", "ClassVersion", "BaseVersion", "SystemImageGUID", "PortGUID"},
		{"0xhost1", "HOST_1", "1", "1", "1", "1", "0xhost1", "0xhost1"},
		{"0xswa", "SW_A", "2", "2", "1", "1", "0xswa", "0xswa"},
		{"0xswb", "SW_B", "2", "2", "1", "1", "0xswb", "0xswb"},
	}
	ports := [][]string{portCSVHeader()}
	ports = append(ports, diagnosticPortRow("0xhost1", "0xhost1", 1, 2807162954))
	ports = append(ports, diagnosticPortRow("0xswa", "0xswa", 0, 3847280712))
	ports = append(ports, diagnosticPortRow("0xswa", "0xswa", 1, 100))
	ports = append(ports, diagnosticPortRow("0xswb", "0xswb", 0, 200))

	swH := diagnosticSwitchHeader()
	switches := [][]string{swH, diagnosticSwitchRow("0xswa"), diagnosticSwitchRow("0xswb")}

	sysH := diagnosticSystemHeader()
	systems := [][]string{sysH, diagnosticSystemRow("0xswa", "SN-A", "ProdA"), diagnosticSystemRow("0xswb", "SN-B", "ProdB")}

	var b strings.Builder
	b.WriteString(writeCSVSection(t, "NODES", nodes))
	b.WriteString(writeCSVSection(t, "PORTS", ports))
	b.WriteString(writeCSVSection(t, "SWITCHES", switches))
	b.WriteString(writeCSVSection(t, "SYSTEM_GENERAL_INFORMATION", systems))
	return []byte(b.String())
}

func buildValidDiagnosticSharpInfo() []byte {
	return []byte(strings.Join([]string{
		"SW_GUID=swa",
		"endianness = 10",
		"enable_endianness_per_job = 0",
		"reproducibility_disable = 0",
		"",
		"SW_GUID=swb",
		"endianness = 0",
		"enable_endianness_per_job = 1",
		"reproducibility_disable = 2",
		"",
	}, "\n"))
}

func writeDiagnosticZip(t *testing.T, zipPath string, members map[string][]byte) {
	t.Helper()
	f, err := os.Create(zipPath)
	require.NoError(t, err)
	zw := zip.NewWriter(f)
	for name, payload := range members {
		w, werr := zw.Create(name)
		require.NoError(t, werr)
		_, werr = w.Write(payload)
		require.NoError(t, werr)
	}
	require.NoError(t, zw.Close())
	require.NoError(t, f.Close())
}

func parseDiagnosticZipInTemp(t *testing.T, archiveFileName string, members map[string][]byte) (application.ParseResult, error) {
	t.Helper()
	tmp := t.TempDir()
	require.NoError(t, os.MkdirAll(filepath.Join(tmp, "data"), 0o755))
	zipPath := filepath.Join(tmp, "data", archiveFileName)
	writeDiagnosticZip(t, zipPath, members)
	t.Chdir(tmp)
	p := New(slog.New(slog.NewTextHandler(io.Discard, nil)))
	abs, err := p.ResolveArchive(archiveFileName)
	if err != nil {
		var zero application.ParseResult
		return zero, err
	}
	return p.Parse(abs)
}

func TestParser_ParseDiagnosticArchiveSucceeds(t *testing.T) {
	db := buildValidDiagnosticDBCSV(t)
	sh := buildValidDiagnosticSharpInfo()
	res, err := parseDiagnosticZipInTemp(t, "valid.zip", map[string][]byte{
		"ibdiagnet2.db_csv":        db,
		"ibdiagnet2.sharp_an_info": sh,
	})
	require.NoError(t, err)
	require.Len(t, res.Nodes, 3)
	require.Len(t, res.Ports, 4)
	require.Len(t, res.SwitchInfos, 2)
	require.Len(t, res.SystemInfos, 2)
	require.Len(t, res.SharpInfos, 2)

	var hostFound bool
	for _, n := range res.Nodes {
		if n.NodeGUID == "0xhost1" && n.NodeType == 1 && n.NodeDesc == "HOST_1" {
			hostFound = true
			break
		}
	}
	require.True(t, hostFound)

	var capOK bool
	for _, p := range res.Ports {
		if p.NodeGUID == "0xhost1" && p.PortNum == 1 && p.CapMsk == 2807162954 {
			capOK = true
			break
		}
	}
	require.True(t, capOK)

	byGUID := make(map[string]struct{ e, en, r int })
	for _, s := range res.SharpInfos {
		byGUID[s.NodeGUID] = struct{ e, en, r int }{s.Endianness, s.EnableEndiannessPerJob, s.ReproducibilityDisable}
	}
	require.Equal(t, struct{ e, en, r int }{10, 0, 0}, byGUID["0xswa"])
	require.Equal(t, struct{ e, en, r int }{0, 1, 2}, byGUID["0xswb"])

	var swAOK bool
	for _, s := range res.SwitchInfos {
		if s.NodeGUID == "0xswa" && s.LinearFDBCap == 49152 {
			swAOK = true
			break
		}
	}
	require.True(t, swAOK)

	var sysBOK bool
	for _, s := range res.SystemInfos {
		if s.NodeGUID == "0xswb" && s.SerialNumber == "SN-B" && s.ProductName == "ProdB" {
			sysBOK = true
			break
		}
	}
	require.True(t, sysBOK)
}

func TestParser_ParseDiagnosticArchiveSucceedsWithNestedZipPaths(t *testing.T) {
	db := buildValidDiagnosticDBCSV(t)
	sh := buildValidDiagnosticSharpInfo()
	res, err := parseDiagnosticZipInTemp(t, "nested.zip", map[string][]byte{
		"log/ibdiagnet2.db_csv":        db,
		"log/ibdiagnet2.sharp_an_info": sh,
	})
	require.NoError(t, err)
	require.Len(t, res.Nodes, 3)
	require.Len(t, res.SharpInfos, 2)
}

func TestParser_ParseDiagnosticArchiveFailsWithoutDbCsv(t *testing.T) {
	sh := buildValidDiagnosticSharpInfo()
	_, err := parseDiagnosticZipInTemp(t, "no_db.zip", map[string][]byte{
		"ibdiagnet2.sharp_an_info": sh,
	})
	require.Error(t, err)
	require.Contains(t, err.Error(), "ibdiagnet2.db_csv")
}

func TestParser_ParseDiagnosticArchiveFailsWithoutSharpInfo(t *testing.T) {
	db := buildValidDiagnosticDBCSV(t)
	_, err := parseDiagnosticZipInTemp(t, "no_sharp.zip", map[string][]byte{
		"ibdiagnet2.db_csv": db,
	})
	require.Error(t, err)
	require.Contains(t, err.Error(), "ibdiagnet2.sharp_an_info")
}

func TestParser_ParseDiagnosticArchiveFailsWhenNodesSectionMissing(t *testing.T) {
	ports := append([][]string{portCSVHeader()}, diagnosticPortRow("0xhost1", "0xhost1", 1, 1))
	var b strings.Builder
	b.WriteString(writeCSVSection(t, "PORTS", ports))
	db := []byte(b.String())
	sh := buildValidDiagnosticSharpInfo()
	_, err := parseDiagnosticZipInTemp(t, "no_nodes.zip", map[string][]byte{
		"ibdiagnet2.db_csv":        db,
		"ibdiagnet2.sharp_an_info": sh,
	})
	require.Error(t, err)
	require.Contains(t, err.Error(), "NODES")
}

func TestParser_ParseDiagnosticArchiveFailsOnInvalidPortInteger(t *testing.T) {
	nodes := [][]string{
		{"NodeGUID", "NodeDesc", "NodeType", "NumPorts", "ClassVersion", "BaseVersion", "SystemImageGUID", "PortGUID"},
		{"0xhost1", "HOST_1", "1", "1", "1", "1", "0xhost1", "0xhost1"},
	}
	bad := diagnosticPortRow("0xhost1", "0xhost1", 1, 1)
	bad[portColIndex("PortNum")] = "not-a-number"
	ports := append([][]string{portCSVHeader()}, bad)
	var b strings.Builder
	b.WriteString(writeCSVSection(t, "NODES", nodes))
	b.WriteString(writeCSVSection(t, "PORTS", ports))
	db := []byte(b.String())
	sh := buildValidDiagnosticSharpInfo()
	_, err := parseDiagnosticZipInTemp(t, "bad_port.zip", map[string][]byte{
		"ibdiagnet2.db_csv":        db,
		"ibdiagnet2.sharp_an_info": sh,
	})
	require.Error(t, err)
	require.Contains(t, err.Error(), "PORTS")
}

func TestParser_ParseDiagnosticArchiveFailsOnInvalidSharpKeyValueLine(t *testing.T) {
	db := buildValidDiagnosticDBCSV(t)
	sh := []byte("SW_GUID=x\nthis-is-not-a-key-value-pair\n")
	_, err := parseDiagnosticZipInTemp(t, "bad_sharp.zip", map[string][]byte{
		"ibdiagnet2.db_csv":        db,
		"ibdiagnet2.sharp_an_info": sh,
	})
	require.Error(t, err)
	require.Contains(t, err.Error(), "invalid key=value")
}

func TestParser_ParseDiagnosticArchiveFailsOnNonIntegerSharpField(t *testing.T) {
	db := buildValidDiagnosticDBCSV(t)
	sh := []byte("SW_GUID=x\nendianness=notint\n")
	_, err := parseDiagnosticZipInTemp(t, "bad_sharp_int.zip", map[string][]byte{
		"ibdiagnet2.db_csv":        db,
		"ibdiagnet2.sharp_an_info": sh,
	})
	require.Error(t, err)
	require.Contains(t, err.Error(), "invalid int")
}
