package server

import "time"

type Option func(*Server)

func WithAddress(addr string) Option {
	return func(s *Server) {
		s.server.Addr = addr
	}
}

func WithReadTimeout(t time.Duration) Option {
	return func(s *Server) {
		s.server.ReadTimeout = t
	}
}

func WithWriteTimeout(t time.Duration) Option {
	return func(s *Server) {
		s.server.WriteTimeout = t
	}
}

func WithShutdownTimeout(t time.Duration) Option {
	return func(s *Server) {
		s.shutdownTimeout = t
	}
}

func WithIdleTimeout(t time.Duration) Option {
	return func(s *Server) {
		s.server.IdleTimeout = t
	}
}
