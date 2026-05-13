package api

import (
	"log/slog"
	"net/http"

	"log-parser/internal/application"
)

func Routes(log *slog.Logger, svc application.Service) http.Handler {
	h := NewHandlers(log, svc)
	mux := http.NewServeMux()
	mux.HandleFunc("POST /api/v1/parse/", h.parse)
	mux.HandleFunc("GET /api/v1/topology/{log_id}", h.topology)
	mux.HandleFunc("GET /api/v1/node/{node_id}", h.node)
	mux.HandleFunc("GET /api/v1/port/{node_id}", h.ports)
	mux.HandleFunc("GET /api/v1/log/{log_id}", h.logMeta)
	return Chain(Recovery(log), Logger(log))(mux)
}
