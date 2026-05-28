package application

import (
	"context"
	"errors"
	"fmt"
	"log/slog"
	"sort"
	"time"

	"log-parser/internal/domain"
)

type service struct {
	log    *slog.Logger
	store  Store
	parser Parser
}

func NewService(log *slog.Logger, store Store, parser Parser) Service {
	return &service{log: log, store: store, parser: parser}
}

func (s *service) ProcessArchive(ctx context.Context, archiveRel string) (int64, error) {
	path, err := s.parser.ResolveArchive(archiveRel)
	if err != nil {
		return 0, fmt.Errorf("resolve archive: %w", err)
	}

	logID, err := s.store.CreateLog(ctx, path)
	if err != nil {
		if !errors.Is(err, ErrDuplicateLogPath) {
			return 0, fmt.Errorf("%w: %w", ErrPersistFailed, err)
		}
		existing, errByPath := s.store.LogByPath(ctx, path)
		if errByPath != nil {
			return 0, fmt.Errorf("%w: %w", ErrPersistFailed, errByPath)
		}
		if existing.Status != domain.LogStatusFailed {
			return 0, fmt.Errorf("%w", ErrDuplicateLogPath)
		}
		logID = existing.ID
		if setErr := s.store.SetStatus(ctx, logID, domain.LogStatusPending); setErr != nil {
			return 0, fmt.Errorf("%w: %w", ErrPersistFailed, setErr)
		}
	}

	parseStart := time.Now()
	res, err := s.parser.Parse(path)
	parseDuration := time.Since(parseStart)
	if err != nil {
		s.log.Warn("parse failed", "log_id", logID, "error", err)

		if setErr := s.store.SetStatus(ctx, logID, domain.LogStatusFailed); setErr != nil {
			s.log.Error("set status failed", "log_id", logID, "error", setErr)
		}
		return 0, errors.Join(ErrParseFailed, err)
	}

	s.log.Info("parse complete", "log_id", logID, "parse_duration", parseDuration)

	if saveErr := s.store.SaveResult(ctx, logID, res); saveErr != nil {
		s.log.Error("archive persist failed", "log_id", logID, "error", saveErr)
		if setErr := s.store.SetStatus(ctx, logID, domain.LogStatusFailed); setErr != nil {
			s.log.Error("set status failed", "log_id", logID, "error", setErr)
		}
		return 0, fmt.Errorf("%w: %w", ErrPersistFailed, saveErr)
	}

	s.log.Info("archive processed",
		"log_id", logID,
		"path", path,
		"nodes", len(res.Nodes),
		"ports", len(res.Ports),
	)
	return logID, nil
}

func (s *service) Topology(ctx context.Context, logID int64) (domain.Topology, error) {
	dlog, err := s.store.Log(ctx, logID)
	if err != nil {
		return domain.Topology{}, fmt.Errorf("log: %w", err)
	}
	if dlog.Status != domain.LogStatusDone {
		return domain.Topology{}, ErrTopologyNotReady
	}

	nodes, err := s.store.Nodes(ctx, logID)
	if err != nil {
		return domain.Topology{}, fmt.Errorf("nodes: %w", err)
	}

	buckets := make(map[int][]int64)
	for _, n := range nodes {
		buckets[n.NodeType] = append(buckets[n.NodeType], n.ID)
	}

	types := make([]int, 0, len(buckets))
	for t := range buckets {
		types = append(types, t)
	}
	sort.Ints(types)

	groups := make([]domain.TopologyGroup, 0, len(types))
	for _, t := range types {
		ids := buckets[t]
		sort.Slice(ids, func(i, j int) bool { return ids[i] < ids[j] })
		groups = append(groups, domain.TopologyGroup{
			NodeType: t,
			NodeIDs:  ids,
		})
	}

	return domain.Topology{Nodes: nodes, Groups: groups}, nil
}

func (s *service) NodeDetail(ctx context.Context, nodeID int64) (domain.NodeDetail, error) {
	d, err := s.store.NodeDetail(ctx, nodeID)
	if err != nil {
		return domain.NodeDetail{}, fmt.Errorf("node detail: %w", err)
	}
	return d, nil
}

func (s *service) Ports(ctx context.Context, nodeID int64) ([]domain.Port, error) {
	if _, err := s.store.Node(ctx, nodeID); err != nil {
		return nil, fmt.Errorf("node: %w", err)
	}
	ports, err := s.store.Ports(ctx, nodeID)
	if err != nil {
		return nil, fmt.Errorf("ports: %w", err)
	}
	return ports, nil
}

func (s *service) LogMeta(ctx context.Context, logID int64) (domain.Log, error) {
	dlog, err := s.store.Log(ctx, logID)
	if err != nil {
		return domain.Log{}, fmt.Errorf("log: %w", err)
	}
	return dlog, nil
}
