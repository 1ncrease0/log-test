package api

import (
	"encoding/json"
	"net/http"
	"time"

	"log-parser/internal/application"
	"log-parser/internal/domain"
)

const maxJSONBody = 1 << 20

type errJSON struct {
	Error string `json:"error"`
}

func writeJSON(w http.ResponseWriter, status int, v any) {
	buf, err := json.Marshal(v)
	if err != nil {
		w.Header().Set("Content-Type", "application/json; charset=utf-8")
		w.WriteHeader(http.StatusInternalServerError)
		_, _ = w.Write([]byte(`{"error":"internal"}`))
		return
	}
	w.Header().Set("Content-Type", "application/json; charset=utf-8")
	w.WriteHeader(status)
	_, _ = w.Write(buf)
}

func writeError(w http.ResponseWriter, status int, msg string) {
	writeJSON(w, status, errJSON{Error: msg})
}

type parseResponse struct {
	LogID int64 `json:"log_id"`
}

type logMetaJSON struct {
	ID         int64  `json:"id"`
	Path       string `json:"path"`
	Status     string `json:"status"`
	NodeCount  int    `json:"node_count"`
	PortCount  int    `json:"port_count"`
	UploadedAt string `json:"uploaded_at"`
}

func logToJSON(l domain.Log) logMetaJSON {
	return logMetaJSON{
		ID:         l.ID,
		Path:       l.Path,
		Status:     string(l.Status),
		NodeCount:  l.NodeCount,
		PortCount:  l.PortCount,
		UploadedAt: l.UploadedAt.UTC().Format(time.RFC3339Nano),
	}
}

type nodeJSON struct {
	ID              int64  `json:"id"`
	LogID           int64  `json:"log_id"`
	NodeGUID        string `json:"node_guid"`
	NodeDesc        string `json:"node_desc"`
	NodeType        int    `json:"node_type"`
	NumPorts        int    `json:"num_ports"`
	ClassVersion    int    `json:"class_version"`
	BaseVersion     int    `json:"base_version"`
	SystemImageGUID string `json:"system_image_guid"`
	PortGUID        string `json:"port_guid"`
}

func nodeToJSON(n domain.Node) nodeJSON {
	return nodeJSON{
		ID:              n.ID,
		LogID:           n.LogID,
		NodeGUID:        n.NodeGUID,
		NodeDesc:        n.NodeDesc,
		NodeType:        n.NodeType,
		NumPorts:        n.NumPorts,
		ClassVersion:    n.ClassVersion,
		BaseVersion:     n.BaseVersion,
		SystemImageGUID: n.SystemImageGUID,
		PortGUID:        n.PortGUID,
	}
}

type topologyGroupJSON struct {
	NodeType int     `json:"node_type"`
	NodeIDs  []int64 `json:"node_ids"`
}

type topologyJSON struct {
	Nodes  []nodeJSON          `json:"nodes"`
	Groups []topologyGroupJSON `json:"groups"`
}

func topologyToJSON(t application.Topology) topologyJSON {
	out := topologyJSON{
		Nodes:  make([]nodeJSON, len(t.Nodes)),
		Groups: make([]topologyGroupJSON, len(t.Groups)),
	}
	for i := range t.Nodes {
		out.Nodes[i] = nodeToJSON(t.Nodes[i])
	}
	for i := range t.Groups {
		out.Groups[i] = topologyGroupJSON{
			NodeType: t.Groups[i].NodeType,
			NodeIDs:  t.Groups[i].NodeIDs,
		}
	}
	return out
}

type switchInfoJSON struct {
	LinearFDBCap         int `json:"linear_fdb_cap"`
	RandomFDBCap         int `json:"random_fdb_cap"`
	MCastFDBCap          int `json:"mcast_fdb_cap"`
	LinearFDBTop         int `json:"linear_fdb_top"`
	DefPort              int `json:"def_port"`
	DefMCastPriPort      int `json:"def_mcast_pri_port"`
	DefMCastNotPriPort   int `json:"def_mcast_not_pri_port"`
	LifeTimeValue        int `json:"life_time_value"`
	PortStateChange      int `json:"port_state_change"`
	OptimizedSLVLMapping int `json:"optimized_s_lvl_mapping"`
	LidsPerPort          int `json:"lids_per_port"`
	PartEnfCap           int `json:"part_enf_cap"`
	InbEnfCap            int `json:"inb_enf_cap"`
	OutbEnfCap           int `json:"outb_enf_cap"`
	FilterRawInbCap      int `json:"filter_raw_inb_cap"`
	FilterRawOutbCap     int `json:"filter_raw_outb_cap"`
	ENP0                 int `json:"enp0"`
	MCastFDBTop          int `json:"mcast_fdb_top"`
}

func switchInfoToJSON(s domain.NodeSwitchInfo) switchInfoJSON {
	return switchInfoJSON{
		LinearFDBCap:         s.LinearFDBCap,
		RandomFDBCap:         s.RandomFDBCap,
		MCastFDBCap:          s.MCastFDBCap,
		LinearFDBTop:         s.LinearFDBTop,
		DefPort:              s.DefPort,
		DefMCastPriPort:      s.DefMCastPriPort,
		DefMCastNotPriPort:   s.DefMCastNotPriPort,
		LifeTimeValue:        s.LifeTimeValue,
		PortStateChange:      s.PortStateChange,
		OptimizedSLVLMapping: s.OptimizedSLVLMapping,
		LidsPerPort:          s.LidsPerPort,
		PartEnfCap:           s.PartEnfCap,
		InbEnfCap:            s.InbEnfCap,
		OutbEnfCap:           s.OutbEnfCap,
		FilterRawInbCap:      s.FilterRawInbCap,
		FilterRawOutbCap:     s.FilterRawOutbCap,
		ENP0:                 s.ENP0,
		MCastFDBTop:          s.MCastFDBTop,
	}
}

type systemInfoJSON struct {
	SerialNumber string `json:"serial_number"`
	PartNumber   string `json:"part_number"`
	Revision     string `json:"revision"`
	ProductName  string `json:"product_name"`
}

func systemInfoToJSON(s domain.NodeSystemInfo) systemInfoJSON {
	return systemInfoJSON{
		SerialNumber: s.SerialNumber,
		PartNumber:   s.PartNumber,
		Revision:     s.Revision,
		ProductName:  s.ProductName,
	}
}

type sharpInfoJSON struct {
	Endianness             int `json:"endianness"`
	EnableEndiannessPerJob int `json:"enable_endianness_per_job"`
	ReproducibilityDisable int `json:"reproducibility_disable"`
}

func sharpInfoToJSON(s domain.NodeSharpInfo) sharpInfoJSON {
	return sharpInfoJSON{
		Endianness:             s.Endianness,
		EnableEndiannessPerJob: s.EnableEndiannessPerJob,
		ReproducibilityDisable: s.ReproducibilityDisable,
	}
}

type nodeDetailJSON struct {
	nodeJSON
	SwitchInfo *switchInfoJSON `json:"switch_info,omitempty"`
	SystemInfo *systemInfoJSON `json:"system_info,omitempty"`
	SharpInfo  *sharpInfoJSON  `json:"sharp_info,omitempty"`
}

func nodeDetailToJSON(d application.NodeDetail) nodeDetailJSON {
	out := nodeDetailJSON{nodeJSON: nodeToJSON(d.Node)}
	if d.SwitchInfo != nil {
		v := switchInfoToJSON(*d.SwitchInfo)
		out.SwitchInfo = &v
	}
	if d.SystemInfo != nil {
		v := systemInfoToJSON(*d.SystemInfo)
		out.SystemInfo = &v
	}
	if d.SharpInfo != nil {
		v := sharpInfoToJSON(*d.SharpInfo)
		out.SharpInfo = &v
	}
	return out
}

type portJSON struct {
	ID                                  int64  `json:"id"`
	LogID                               int64  `json:"log_id"`
	NodeGUID                            string `json:"node_guid"`
	PortGUID                            string `json:"port_guid"`
	PortNum                             int    `json:"port_num"`
	MKey                                string `json:"m_key"`
	GIDPrfx                             string `json:"gid_prfx"`
	MSMLID                              int    `json:"msm_lid"`
	LID                                 int    `json:"lid"`
	CapMsk                              int64  `json:"cap_msk"`
	MKeyLeasePeriod                     int    `json:"m_key_lease_period"`
	DiagCode                            int    `json:"diag_code"`
	LinkWidthActv                       int    `json:"link_width_actv"`
	LinkWidthSup                        int    `json:"link_width_sup"`
	LinkWidthEn                         int    `json:"link_width_en"`
	LocalPortNum                        int    `json:"local_port_num"`
	LinkSpeedEn                         int    `json:"link_speed_en"`
	LinkSpeedActv                       int    `json:"link_speed_actv"`
	LMC                                 int    `json:"lmc"`
	MKeyProtBits                        int    `json:"m_key_prot_bits"`
	LinkDownDefState                    int    `json:"link_down_def_state"`
	PortPhyState                        int    `json:"port_phy_state"`
	PortState                           int    `json:"port_state"`
	LinkSpeedSup                        int    `json:"link_speed_sup"`
	VLArbHighCap                        int    `json:"vl_arb_high_cap"`
	VLHighLimit                         int    `json:"vl_high_limit"`
	InitType                            int    `json:"init_type"`
	VLCap                               int    `json:"vl_cap"`
	MSMSL                               int    `json:"msmsl"`
	NMTU                                int    `json:"nmtu"`
	FilterRawOutb                       int    `json:"filter_raw_outb"`
	FilterRawInb                        int    `json:"filter_raw_inb"`
	PartEnfOutb                         int    `json:"part_enf_outb"`
	PartEnfInb                          int    `json:"part_enf_inb"`
	OpVLs                               int    `json:"op_vls"`
	HoQLife                             int    `json:"hoq_life"`
	VLStallCnt                          int    `json:"vl_stall_cnt"`
	MTUCap                              int    `json:"mtu_cap"`
	InitTypeReply                       int    `json:"init_type_reply"`
	VLArbLowCap                         int    `json:"vl_arb_low_cap"`
	PKeyViolations                      int    `json:"pkey_violations"`
	MKeyViolations                      int    `json:"mkey_violations"`
	SubnTmo                             int    `json:"subn_tmo"`
	MulticastPKeyTrapSuppressionEnabled int    `json:"multicast_pkey_trap_suppression_enabled"`
	ClientReregister                    int    `json:"client_reregister"`
	GUIDCap                             int    `json:"guid_cap"`
	QKeyViolations                      int    `json:"qkey_violations"`
	MaxCreditHint                       int    `json:"max_credit_hint"`
	OverrunErrs                         int    `json:"overrun_errs"`
	LocalPhyError                       int    `json:"local_phy_error"`
	RespTimeValue                       int    `json:"resp_time_value"`
	LinkRoundTripLatency                int    `json:"link_round_trip_latency"`
	OOOSLMask                           string `json:"ooosl_mask"`
	CapMsk2                             *int   `json:"cap_msk2,omitempty"`
	FECActv                             *int   `json:"fec_actv,omitempty"`
	RetransActv                         *int   `json:"retrans_actv,omitempty"`
}

func portToJSON(p domain.Port) portJSON {
	return portJSON{
		ID:                                  p.ID,
		LogID:                               p.LogID,
		NodeGUID:                            p.NodeGUID,
		PortGUID:                            p.PortGUID,
		PortNum:                             p.PortNum,
		MKey:                                p.MKey,
		GIDPrfx:                             p.GIDPrfx,
		MSMLID:                              p.MSMLID,
		LID:                                 p.LID,
		CapMsk:                              p.CapMsk,
		MKeyLeasePeriod:                     p.MKeyLeasePeriod,
		DiagCode:                            p.DiagCode,
		LinkWidthActv:                       p.LinkWidthActv,
		LinkWidthSup:                        p.LinkWidthSup,
		LinkWidthEn:                         p.LinkWidthEn,
		LocalPortNum:                        p.LocalPortNum,
		LinkSpeedEn:                         p.LinkSpeedEn,
		LinkSpeedActv:                       p.LinkSpeedActv,
		LMC:                                 p.LMC,
		MKeyProtBits:                        p.MKeyProtBits,
		LinkDownDefState:                    p.LinkDownDefState,
		PortPhyState:                        p.PortPhyState,
		PortState:                           p.PortState,
		LinkSpeedSup:                        p.LinkSpeedSup,
		VLArbHighCap:                        p.VLArbHighCap,
		VLHighLimit:                         p.VLHighLimit,
		InitType:                            p.InitType,
		VLCap:                               p.VLCap,
		MSMSL:                               p.MSMSL,
		NMTU:                                p.NMTU,
		FilterRawOutb:                       p.FilterRawOutb,
		FilterRawInb:                        p.FilterRawInb,
		PartEnfOutb:                         p.PartEnfOutb,
		PartEnfInb:                          p.PartEnfInb,
		OpVLs:                               p.OpVLs,
		HoQLife:                             p.HoQLife,
		VLStallCnt:                          p.VLStallCnt,
		MTUCap:                              p.MTUCap,
		InitTypeReply:                       p.InitTypeReply,
		VLArbLowCap:                         p.VLArbLowCap,
		PKeyViolations:                      p.PKeyViolations,
		MKeyViolations:                      p.MKeyViolations,
		SubnTmo:                             p.SubnTmo,
		MulticastPKeyTrapSuppressionEnabled: p.MulticastPKeyTrapSuppressionEnabled,
		ClientReregister:                    p.ClientReregister,
		GUIDCap:                             p.GUIDCap,
		QKeyViolations:                      p.QKeyViolations,
		MaxCreditHint:                       p.MaxCreditHint,
		OverrunErrs:                         p.OverrunErrs,
		LocalPhyError:                       p.LocalPhyError,
		RespTimeValue:                       p.RespTimeValue,
		LinkRoundTripLatency:                p.LinkRoundTripLatency,
		OOOSLMask:                           p.OOOSLMask,
		CapMsk2:                             p.CapMsk2,
		FECActv:                             p.FECActv,
		RetransActv:                         p.RetransActv,
	}
}
