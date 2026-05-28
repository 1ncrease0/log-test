package application_test

import (
	"context"
	"errors"
	"io"
	"log/slog"
	"testing"

	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/require"

	"log-parser/internal/application"
	"log-parser/internal/domain"
	"log-parser/internal/mocks"
)

func testLogger() *slog.Logger {
	return slog.New(slog.NewTextHandler(io.Discard, nil))
}

func TestService_ProcessArchive_Success(t *testing.T) {
	t.Parallel()

	st := mocks.NewMockStore(t)
	pr := mocks.NewMockParser(t)

	pr.EXPECT().ResolveArchive("logs/a.zip").Return("/abs/data/logs/a.zip", nil)
	st.EXPECT().CreateLog(mock.Anything, "/abs/data/logs/a.zip").Return(int64(10), nil)
	pr.EXPECT().Parse("/abs/data/logs/a.zip").Return(domain.ParseResult{
		Nodes: []domain.Node{{ID: 1, NodeGUID: "0x1"}},
	}, nil)
	st.EXPECT().SaveResult(mock.Anything, int64(10), mock.MatchedBy(func(r domain.ParseResult) bool {
		return len(r.Nodes) == 1 && r.Nodes[0].NodeGUID == "0x1"
	})).Return(nil)

	svc := application.NewService(testLogger(), st, pr)
	id, err := svc.ProcessArchive(context.Background(), "logs/a.zip")
	require.NoError(t, err)
	require.Equal(t, int64(10), id)
}

func TestService_ProcessArchive_ResolveFails(t *testing.T) {
	t.Parallel()

	pr := mocks.NewMockParser(t)
	pr.EXPECT().ResolveArchive("x").Return("", errors.New("boom"))

	svc := application.NewService(testLogger(), mocks.NewMockStore(t), pr)
	_, err := svc.ProcessArchive(context.Background(), "x")
	require.Error(t, err)
	require.Contains(t, err.Error(), "resolve archive")
}

func TestService_ProcessArchive_DuplicatePath(t *testing.T) {
	t.Parallel()

	st := mocks.NewMockStore(t)
	pr := mocks.NewMockParser(t)
	pr.EXPECT().ResolveArchive("logs/a.zip").Return("/abs/a.zip", nil)
	st.EXPECT().CreateLog(mock.Anything, "/abs/a.zip").Return(int64(0), application.ErrDuplicateLogPath)
	st.EXPECT().LogByPath(mock.Anything, "/abs/a.zip").Return(domain.Log{
		ID:     7,
		Path:   "/abs/a.zip",
		Status: domain.LogStatusDone,
	}, nil)

	svc := application.NewService(testLogger(), st, pr)
	_, err := svc.ProcessArchive(context.Background(), "logs/a.zip")
	require.Error(t, err)
	require.ErrorIs(t, err, application.ErrDuplicateLogPath)
}

func TestService_ProcessArchive_RetryAfterFailed(t *testing.T) {
	t.Parallel()

	st := mocks.NewMockStore(t)
	pr := mocks.NewMockParser(t)
	pr.EXPECT().ResolveArchive("logs/a.zip").Return("/abs/a.zip", nil)
	st.EXPECT().CreateLog(mock.Anything, "/abs/a.zip").Return(int64(0), application.ErrDuplicateLogPath)
	st.EXPECT().LogByPath(mock.Anything, "/abs/a.zip").Return(domain.Log{
		ID:     7,
		Path:   "/abs/a.zip",
		Status: domain.LogStatusFailed,
	}, nil)
	st.EXPECT().SetStatus(mock.Anything, int64(7), domain.LogStatusPending).Return(nil)
	pr.EXPECT().Parse("/abs/a.zip").Return(domain.ParseResult{
		Nodes: []domain.Node{{NodeGUID: "0x1"}},
	}, nil)
	st.EXPECT().SaveResult(mock.Anything, int64(7), mock.Anything).Return(nil)

	svc := application.NewService(testLogger(), st, pr)
	id, err := svc.ProcessArchive(context.Background(), "logs/a.zip")
	require.NoError(t, err)
	require.Equal(t, int64(7), id)
}

func TestService_ProcessArchive_CreateLogPersistFailed(t *testing.T) {
	t.Parallel()

	st := mocks.NewMockStore(t)
	pr := mocks.NewMockParser(t)
	pr.EXPECT().ResolveArchive("logs/a.zip").Return("/abs/a.zip", nil)
	st.EXPECT().CreateLog(mock.Anything, "/abs/a.zip").Return(int64(0), errors.New("db down"))

	svc := application.NewService(testLogger(), st, pr)
	_, err := svc.ProcessArchive(context.Background(), "logs/a.zip")
	require.Error(t, err)
	require.ErrorIs(t, err, application.ErrPersistFailed)
}

func TestService_ProcessArchive_ParseFailsSetsFailed(t *testing.T) {
	t.Parallel()

	st := mocks.NewMockStore(t)
	pr := mocks.NewMockParser(t)
	pr.EXPECT().ResolveArchive("logs/a.zip").Return("/abs/a.zip", nil)
	st.EXPECT().CreateLog(mock.Anything, "/abs/a.zip").Return(int64(5), nil)
	pr.EXPECT().Parse("/abs/a.zip").Return(domain.ParseResult{}, errors.New("parse boom"))
	st.EXPECT().SetStatus(mock.Anything, int64(5), domain.LogStatusFailed).Return(nil)

	svc := application.NewService(testLogger(), st, pr)
	_, err := svc.ProcessArchive(context.Background(), "logs/a.zip")
	require.Error(t, err)
	require.ErrorIs(t, err, application.ErrParseFailed)
}

func TestService_ProcessArchive_SaveFailsSetsFailed(t *testing.T) {
	t.Parallel()

	st := mocks.NewMockStore(t)
	pr := mocks.NewMockParser(t)
	pr.EXPECT().ResolveArchive("logs/a.zip").Return("/abs/a.zip", nil)
	st.EXPECT().CreateLog(mock.Anything, "/abs/a.zip").Return(int64(5), nil)
	pr.EXPECT().Parse("/abs/a.zip").Return(domain.ParseResult{Nodes: []domain.Node{{}}}, nil)
	st.EXPECT().SaveResult(mock.Anything, int64(5), mock.Anything).Return(errors.New("save boom"))
	st.EXPECT().SetStatus(mock.Anything, int64(5), domain.LogStatusFailed).Return(nil)

	svc := application.NewService(testLogger(), st, pr)
	_, err := svc.ProcessArchive(context.Background(), "logs/a.zip")
	require.Error(t, err)
	require.ErrorIs(t, err, application.ErrPersistFailed)
}

func TestService_Topology_NotReady(t *testing.T) {
	t.Parallel()

	st := mocks.NewMockStore(t)
	st.EXPECT().Log(mock.Anything, int64(1)).Return(domain.Log{Status: domain.LogStatusPending}, nil)

	svc := application.NewService(testLogger(), st, mocks.NewMockParser(t))
	_, err := svc.Topology(context.Background(), 1)
	require.ErrorIs(t, err, application.ErrTopologyNotReady)
}

func TestService_Topology_GroupsSorted(t *testing.T) {
	t.Parallel()

	st := mocks.NewMockStore(t)
	st.EXPECT().Log(mock.Anything, int64(2)).Return(domain.Log{Status: domain.LogStatusDone}, nil)
	st.EXPECT().Nodes(mock.Anything, int64(2)).Return([]domain.Node{
		{ID: 10, NodeType: 2},
		{ID: 5, NodeType: 1},
		{ID: 7, NodeType: 2},
	}, nil)

	svc := application.NewService(testLogger(), st, mocks.NewMockParser(t))
	topo, err := svc.Topology(context.Background(), 2)
	require.NoError(t, err)
	require.Len(t, topo.Groups, 2)
	require.Equal(t, 1, topo.Groups[0].NodeType)
	require.Equal(t, []int64{5}, topo.Groups[0].NodeIDs)
	require.Equal(t, 2, topo.Groups[1].NodeType)
	require.Equal(t, []int64{7, 10}, topo.Groups[1].NodeIDs)
}

func TestService_Ports_NodeNotFound(t *testing.T) {
	t.Parallel()

	st := mocks.NewMockStore(t)
	st.EXPECT().Node(mock.Anything, int64(99)).Return(domain.Node{}, application.ErrNotFound)

	svc := application.NewService(testLogger(), st, mocks.NewMockParser(t))
	_, err := svc.Ports(context.Background(), 99)
	require.Error(t, err)
	require.ErrorIs(t, err, application.ErrNotFound)
}

func TestService_Ports_Success(t *testing.T) {
	t.Parallel()

	st := mocks.NewMockStore(t)
	st.EXPECT().Node(mock.Anything, int64(1)).Return(domain.Node{ID: 1}, nil)
	st.EXPECT().Ports(mock.Anything, int64(1)).Return([]domain.Port{{ID: 2}}, nil)

	svc := application.NewService(testLogger(), st, mocks.NewMockParser(t))
	ports, err := svc.Ports(context.Background(), 1)
	require.NoError(t, err)
	require.Len(t, ports, 1)
	require.Equal(t, int64(2), ports[0].ID)
}
