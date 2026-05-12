package domain

type NodeSharpInfo struct {
	ID                     int64
	NodeGUID               string
	Endianness             int
	EnableEndiannessPerJob int
	ReproducibilityDisable int
}
