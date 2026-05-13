package api

import (
	"encoding/json"
	"errors"
	"io"
	"log/slog"
	"net/http"
	"strconv"

	"log-parser/internal/application"
)

type Handlers struct {
	svc application.Service
	log *slog.Logger
}

func NewHandlers(log *slog.Logger, svc application.Service) *Handlers {
	return &Handlers{svc: svc, log: log}
}

type parseRequestBody struct {
	Path string `json:"path"`
}

func (h *Handlers) parse(w http.ResponseWriter, r *http.Request) {
	r.Body = http.MaxBytesReader(w, r.Body, maxJSONBody)
	var body parseRequestBody
	if err := json.NewDecoder(r.Body).Decode(&body); err != nil {
		if errors.Is(err, io.EOF) {
			writeError(w, http.StatusBadRequest, "empty body")
			return
		}
		writeError(w, http.StatusBadRequest, "invalid json")
		return
	}
	if body.Path == "" {
		writeError(w, http.StatusBadRequest, "path is required")
		return
	}

	logID, err := h.svc.ProcessArchive(r.Context(), body.Path)
	if err != nil {
		h.handleProcessArchiveError(w, err)
		return
	}
	writeJSON(w, http.StatusCreated, parseResponse{LogID: logID})
}

func (h *Handlers) handleProcessArchiveError(w http.ResponseWriter, err error) {
	status, msg, logErr := mapErrors(err)
	if logErr {
		h.log.Error("process archive failed", "error", err)
	}
	writeError(w, status, msg)
}

func (h *Handlers) topology(w http.ResponseWriter, r *http.Request) {
	logID, ok := pathInt64(r, "log_id")
	if !ok {
		writeError(w, http.StatusBadRequest, "invalid log_id")
		return
	}
	t, err := h.svc.Topology(r.Context(), logID)
	if err != nil {
		if errors.Is(err, application.ErrNotFound) {
			writeError(w, http.StatusNotFound, "log not found")
			return
		}
		if errors.Is(err, application.ErrTopologyNotReady) {
			writeError(w, http.StatusConflict, err.Error())
			return
		}
		h.log.Error("topology", "error", err)
		writeError(w, http.StatusInternalServerError, http.StatusText(http.StatusInternalServerError))
		return
	}
	writeJSON(w, http.StatusOK, topologyToJSON(t))
}

func (h *Handlers) node(w http.ResponseWriter, r *http.Request) {
	nodeID, ok := pathInt64(r, "node_id")
	if !ok {
		writeError(w, http.StatusBadRequest, "invalid node_id")
		return
	}
	d, err := h.svc.NodeDetail(r.Context(), nodeID)
	if err != nil {
		if errors.Is(err, application.ErrNotFound) {
			writeError(w, http.StatusNotFound, "node not found")
			return
		}
		h.log.Error("node detail", "error", err)
		writeError(w, http.StatusInternalServerError, http.StatusText(http.StatusInternalServerError))
		return
	}
	writeJSON(w, http.StatusOK, nodeDetailToJSON(d))
}

func (h *Handlers) ports(w http.ResponseWriter, r *http.Request) {
	nodeID, ok := pathInt64(r, "node_id")
	if !ok {
		writeError(w, http.StatusBadRequest, "invalid node_id")
		return
	}
	ports, err := h.svc.Ports(r.Context(), nodeID)
	if err != nil {
		if errors.Is(err, application.ErrNotFound) {
			writeError(w, http.StatusNotFound, "node not found")
			return
		}
		h.log.Error("ports", "error", err)
		writeError(w, http.StatusInternalServerError, http.StatusText(http.StatusInternalServerError))
		return
	}
	out := make([]portJSON, len(ports))
	for i := range ports {
		out[i] = portToJSON(ports[i])
	}
	writeJSON(w, http.StatusOK, out)
}

func (h *Handlers) logMeta(w http.ResponseWriter, r *http.Request) {
	logID, ok := pathInt64(r, "log_id")
	if !ok {
		writeError(w, http.StatusBadRequest, "invalid log_id")
		return
	}
	l, err := h.svc.LogMeta(r.Context(), logID)
	if err != nil {
		if errors.Is(err, application.ErrNotFound) {
			writeError(w, http.StatusNotFound, "log not found")
			return
		}
		h.log.Error("log meta", "error", err)
		writeError(w, http.StatusInternalServerError, http.StatusText(http.StatusInternalServerError))
		return
	}
	writeJSON(w, http.StatusOK, logToJSON(l))
}

func pathInt64(r *http.Request, name string) (int64, bool) {
	s := r.PathValue(name)
	if s == "" {
		return 0, false
	}
	v, err := strconv.ParseInt(s, 10, 64)
	if err != nil || v < 1 {
		return 0, false
	}
	return v, true
}
