package domain

type ParseResult struct {
	Nodes       []Node
	Ports       []Port
	SwitchInfos []NodeSwitchInfo
	SystemInfos []NodeSystemInfo
	SharpInfos  []NodeSharpInfo
}
