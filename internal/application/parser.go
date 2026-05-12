package application

import "log-parser/internal/domain"

type ParseResult struct {
	Nodes       []domain.Node
	Ports       []domain.Port
	SwitchInfos []domain.NodeSwitchInfo
	SystemInfos []domain.NodeSystemInfo
	SharpInfos  []domain.NodeSharpInfo
}

type Parser interface {
	Parse(archivePath string) (ParseResult, error)
}
