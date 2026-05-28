package application

import (
	"context"
	"log-parser/internal/domain"
)

type Service interface {
	ProcessArchive(ctx context.Context, rel string) (int64, error)
	Topology(ctx context.Context, logID int64) (domain.Topology, error)
	NodeDetail(ctx context.Context, nodeID int64) (domain.NodeDetail, error)
	Ports(ctx context.Context, nodeID int64) ([]domain.Port, error)
	LogMeta(ctx context.Context, logID int64) (domain.Log, error)
}

type Parser interface {
	ResolveArchive(rel string) (string, error)
	Parse(archivePath string) (domain.ParseResult, error)
}

type Store interface {
	CreateLog(ctx context.Context, path string) (int64, error)
	LogByPath(ctx context.Context, path string) (domain.Log, error)
	SaveResult(ctx context.Context, logID int64, result domain.ParseResult) error
	SetStatus(ctx context.Context, logID int64, status domain.LogStatus) error

	Log(ctx context.Context, id int64) (domain.Log, error)
	Nodes(ctx context.Context, logID int64) ([]domain.Node, error)
	Node(ctx context.Context, nodeID int64) (domain.Node, error)
	NodeDetail(ctx context.Context, nodeID int64) (domain.NodeDetail, error)
	Ports(ctx context.Context, nodeID int64) ([]domain.Port, error)
}
