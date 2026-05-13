package application

import (
	"context"
	"errors"

	"log-parser/internal/domain"
)

var ErrDuplicateLogPath = errors.New("log path already exists")

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

type Parser interface {
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
