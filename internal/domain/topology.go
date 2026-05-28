package domain

type TopologyGroup struct {
	NodeType int
	NodeIDs  []int64
}

type Topology struct {
	Nodes  []Node
	Groups []TopologyGroup
}
