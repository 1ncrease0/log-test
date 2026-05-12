package domain

type NodeSwitchInfo struct {
	ID                   int64
	NodeGUID             string
	LinearFDBCap         int
	RandomFDBCap         int
	MCastFDBCap          int
	LinearFDBTop         int
	DefPort              int
	DefMCastPriPort      int
	DefMCastNotPriPort   int
	LifeTimeValue        int
	PortStateChange      int
	OptimizedSLVLMapping int
	LidsPerPort          int
	PartEnfCap           int
	InbEnfCap            int
	OutbEnfCap           int
	FilterRawInbCap      int
	FilterRawOutbCap     int
	ENP0                 int
	MCastFDBTop          int
}
