package domain

type NodeDetail struct {
	Node
	SwitchInfo *NodeSwitchInfo
	SystemInfo *NodeSystemInfo
	SharpInfo  *NodeSharpInfo
}
