package domain

type NodeType int

const (
	NodeTypeHost   NodeType = 1
	NodeTypeSwitch NodeType = 2
)

type Node struct {
	ID              int64
	LogID           int64
	NodeGUID        string
	NodeDesc        string
	NodeType        NodeType
	NumPorts        int
	ClassVersion    int
	BaseVersion     int
	SystemImageGUID string
	PortGUID        string
}
