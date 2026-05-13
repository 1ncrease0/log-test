package application

import (
	"context"
	"log-parser/internal/domain"
)

type ParseResult struct {
	Nodes       []domain.Node
	Ports       []domain.Port
	SwitchInfos []domain.NodeSwitchInfo
	SystemInfos []domain.NodeSystemInfo
	SharpInfos  []domain.NodeSharpInfo
}

type NodeDetail struct {
	domain.Node
	SwitchInfo *domain.NodeSwitchInfo
	SystemInfo *domain.NodeSystemInfo
	SharpInfo  *domain.NodeSharpInfo
}

type TopologyGroup struct {
	NodeType int
	NodeIDs  []int64
}

type Topology struct {
	Nodes  []domain.Node
	Groups []TopologyGroup
}

type Service interface {
	ProcessArchive(ctx context.Context, rel string) (int64, error)
	Topology(ctx context.Context, logID int64) (Topology, error)
	NodeDetail(ctx context.Context, nodeID int64) (NodeDetail, error)
	Ports(ctx context.Context, nodeID int64) ([]domain.Port, error)
	LogMeta(ctx context.Context, logID int64) (domain.Log, error)
}

type Parser interface {
	ResolveArchive(rel string) (string, error)
	Parse(archivePath string) (ParseResult, error)
}

type Store interface {
	CreateLog(ctx context.Context, path string) (int64, error)
	SaveResult(ctx context.Context, logID int64, result ParseResult) error
	SetStatus(ctx context.Context, logID int64, status domain.LogStatus) error

	Log(ctx context.Context, id int64) (domain.Log, error)
	Nodes(ctx context.Context, logID int64) ([]domain.Node, error)
	Node(ctx context.Context, nodeID int64) (domain.Node, error)
	NodeDetail(ctx context.Context, nodeID int64) (NodeDetail, error)
	Ports(ctx context.Context, nodeID int64) ([]domain.Port, error)
}
