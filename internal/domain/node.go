package domain

type Node struct {
	ID              int64
	LogID           int64
	NodeGUID        string
	NodeDesc        string
	NodeType        int
	NumPorts        int
	ClassVersion    int
	BaseVersion     int
	SystemImageGUID string
	PortGUID        string
}
