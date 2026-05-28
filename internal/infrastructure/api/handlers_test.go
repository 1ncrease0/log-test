package api

import (
	"bytes"
	"encoding/json"
	"errors"
	"io"
	"log/slog"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/require"

	"log-parser/internal/application"
	"log-parser/internal/domain"
	"log-parser/internal/mocks"
)

func discardLogger() *slog.Logger {
	return slog.New(slog.NewTextHandler(io.Discard, nil))
}

func TestHandlers_parse(t *testing.T) {
	t.Parallel()

	cases := []struct {
		name       string
		body       string
		setup      func(*mocks.MockService)
		wantStatus int
		wantErr    string
		wantLogID  int64
	}{
		{
			name:       "empty_body",
			body:       "",
			setup:      func(_ *mocks.MockService) {},
			wantStatus: http.StatusBadRequest,
			wantErr:    "empty body",
		},
		{
			name:       "invalid_json",
			body:       "{",
			setup:      func(_ *mocks.MockService) {},
			wantStatus: http.StatusBadRequest,
			wantErr:    "invalid json",
		},
		{
			name:       "missing_path",
			body:       `{}`,
			setup:      func(_ *mocks.MockService) {},
			wantStatus: http.StatusBadRequest,
			wantErr:    "path is required",
		},
		{
			name: "created",
			body: `{"path":"logs/a.zip"}`,
			setup: func(s *mocks.MockService) {
				s.EXPECT().ProcessArchive(mock.Anything, "logs/a.zip").Return(int64(42), nil)
			},
			wantStatus: http.StatusCreated,
			wantLogID:  42,
		},
		{
			name: "duplicate_path",
			body: `{"path":"logs/dup.zip"}`,
			setup: func(s *mocks.MockService) {
				s.EXPECT().ProcessArchive(mock.Anything, "logs/dup.zip").Return(int64(0), application.ErrDuplicateLogPath)
			},
			wantStatus: http.StatusConflict,
			wantErr:    application.ErrDuplicateLogPath.Error(),
		},
		{
			name: "parse_failed",
			body: `{"path":"logs/bad.zip"}`,
			setup: func(s *mocks.MockService) {
				s.EXPECT().ProcessArchive(mock.Anything, "logs/bad.zip").Return(int64(0), errors.Join(application.ErrParseFailed, errors.New("inner")))
			},
			wantStatus: http.StatusUnprocessableEntity,
			wantErr:    "invalid or incomplete log archive",
		},
	}

	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			t.Parallel()
			svc := mocks.NewMockService(t)
			tc.setup(svc)
			h := Routes(discardLogger(), svc)

			req := httptest.NewRequest(http.MethodPost, "/api/v1/parse/", bytes.NewBufferString(tc.body))
			req.Header.Set("Content-Type", "application/json")
			rec := httptest.NewRecorder()
			h.ServeHTTP(rec, req)

			require.Equal(t, tc.wantStatus, rec.Code)
			var got map[string]any
			require.NoError(t, json.NewDecoder(rec.Body).Decode(&got))
			if tc.wantErr != "" {
				require.Equal(t, tc.wantErr, got["error"])
			}
			if tc.wantLogID != 0 {
				require.EqualValues(t, tc.wantLogID, got["log_id"])
			}
		})
	}
}

func TestHandlers_topology(t *testing.T) {
	t.Parallel()

	svc := mocks.NewMockService(t)
	svc.EXPECT().Topology(mock.Anything, int64(1)).Return(domain.Topology{
		Nodes: []domain.Node{{ID: 1, LogID: 1, NodeType: 1}},
		Groups: []domain.TopologyGroup{
			{NodeType: 1, NodeIDs: []int64{1}},
		},
	}, nil)

	h := Routes(discardLogger(), svc)
	req := httptest.NewRequest(http.MethodGet, "/api/v1/topology/1", nil)
	rec := httptest.NewRecorder()
	h.ServeHTTP(rec, req)
	require.Equal(t, http.StatusOK, rec.Code)
}

func TestHandlers_topology_conflict(t *testing.T) {
	t.Parallel()

	svc := mocks.NewMockService(t)
	svc.EXPECT().Topology(mock.Anything, int64(2)).Return(domain.Topology{}, application.ErrTopologyNotReady)

	h := Routes(discardLogger(), svc)
	req := httptest.NewRequest(http.MethodGet, "/api/v1/topology/2", nil)
	rec := httptest.NewRecorder()
	h.ServeHTTP(rec, req)
	require.Equal(t, http.StatusConflict, rec.Code)
}

func TestHandlers_topology_invalidID(t *testing.T) {
	t.Parallel()

	svc := mocks.NewMockService(t)
	h := Routes(discardLogger(), svc)
	req := httptest.NewRequest(http.MethodGet, "/api/v1/topology/0", nil)
	rec := httptest.NewRecorder()
	h.ServeHTTP(rec, req)
	require.Equal(t, http.StatusBadRequest, rec.Code)
}

func TestHandlers_node_notFound(t *testing.T) {
	t.Parallel()

	svc := mocks.NewMockService(t)
	svc.EXPECT().NodeDetail(mock.Anything, int64(9)).Return(domain.NodeDetail{}, application.ErrNotFound)

	h := Routes(discardLogger(), svc)
	req := httptest.NewRequest(http.MethodGet, "/api/v1/node/9", nil)
	rec := httptest.NewRecorder()
	h.ServeHTTP(rec, req)
	require.Equal(t, http.StatusNotFound, rec.Code)
}

func TestHandlers_ports_success(t *testing.T) {
	t.Parallel()

	svc := mocks.NewMockService(t)
	svc.EXPECT().Ports(mock.Anything, int64(3)).Return([]domain.Port{{ID: 1, PortNum: 2}}, nil)

	h := Routes(discardLogger(), svc)
	req := httptest.NewRequest(http.MethodGet, "/api/v1/port/3", nil)
	rec := httptest.NewRecorder()
	h.ServeHTTP(rec, req)
	require.Equal(t, http.StatusOK, rec.Code)
	var out []map[string]any
	require.NoError(t, json.NewDecoder(rec.Body).Decode(&out))
	require.Len(t, out, 1)
}

func TestHandlers_logMeta(t *testing.T) {
	t.Parallel()

	ts := time.Date(2020, 1, 2, 3, 4, 5, 0, time.UTC)
	svc := mocks.NewMockService(t)
	svc.EXPECT().LogMeta(mock.Anything, int64(4)).Return(domain.Log{
		ID:         4,
		Path:       "/p",
		Status:     domain.LogStatusDone,
		NodeCount:  2,
		PortCount:  3,
		UploadedAt: ts,
	}, nil)

	h := Routes(discardLogger(), svc)
	req := httptest.NewRequest(http.MethodGet, "/api/v1/log/4", nil)
	rec := httptest.NewRecorder()
	h.ServeHTTP(rec, req)
	require.Equal(t, http.StatusOK, rec.Code)
}

func TestPathInt64(t *testing.T) {
	t.Parallel()

	cases := []struct {
		name   string
		raw    string
		wantOK bool
		wantV  int64
	}{
		{name: "ok", raw: "12", wantOK: true, wantV: 12},
		{name: "zero_invalid", raw: "0", wantOK: false},
		{name: "negative_invalid", raw: "-1", wantOK: false},
		{name: "not_a_number", raw: "abc", wantOK: false},
		{name: "empty", raw: "", wantOK: false},
	}

	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			t.Parallel()
			req := httptest.NewRequest(http.MethodGet, "/api/v1/topology/"+tc.raw, nil)
			req.SetPathValue("log_id", tc.raw)
			v, ok := pathInt64(req, "log_id")
			require.Equal(t, tc.wantOK, ok)
			if tc.wantOK {
				require.Equal(t, tc.wantV, v)
			}
		})
	}
}
