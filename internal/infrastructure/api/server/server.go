package server

import (
	"context"
	"errors"
	"log/slog"
	"net/http"
	"time"
)

const (
	defaultAddress         = ":8080"
	defaultReadTimeout     = 10 * time.Second
	defaultWriteTimeout    = 10 * time.Second
	defaultShutdownTimeout = 3 * time.Second
	defaultIdleTimeout     = 60 * time.Second
)

type Server struct {
	server *http.Server
	notify chan error

	shutdownTimeout time.Duration
	log             *slog.Logger
}

func New(l *slog.Logger, mux http.Handler, opts ...Option) *Server {
	s := &Server{
		notify:          make(chan error, 1),
		shutdownTimeout: defaultShutdownTimeout,
		log:             l,
	}

	s.server = &http.Server{
		Addr:         defaultAddress,
		Handler:      mux,
		ReadTimeout:  defaultReadTimeout,
		WriteTimeout: defaultWriteTimeout,
		IdleTimeout:  defaultIdleTimeout,
	}

	for _, opt := range opts {
		opt(s)
	}

	return s
}

func (s *Server) Start() {
	go func() {
		s.log.Info("http server started", slog.String("addr", s.server.Addr))

		if err := s.server.ListenAndServe(); !errors.Is(err, http.ErrServerClosed) {
			s.notify <- err
		}

		close(s.notify)
	}()
}

func (s *Server) Notify() <-chan error {
	return s.notify
}

func (s *Server) Shutdown() error {
	ctx, cancel := context.WithTimeout(context.Background(), s.shutdownTimeout)
	defer cancel()

	if err := s.server.Shutdown(ctx); err != nil {
		s.log.Error("http server shutdown error", slog.String("err", err.Error()))
		return err
	}

	s.log.Info("http server stopped gracefully")

	return nil
}
