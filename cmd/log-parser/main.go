package main

import (
	"context"
	"errors"
	"fmt"
	"log/slog"
	"os"
	"os/signal"
	"syscall"

	"log-parser/internal/application"
	"log-parser/internal/config"
	"log-parser/internal/infrastructure/api"
	"log-parser/internal/infrastructure/api/server"
	"log-parser/internal/infrastructure/db"
	"log-parser/internal/infrastructure/parser"
	"log-parser/pkg/logger"
)

func main() {
	cfg := config.MustLoad()
	log := logger.New(cfg.LogLevel)
	if err := run(cfg, log); err != nil {
		log.Error("run", "error", err)
		os.Exit(1)
	}
}

func run(cfg config.Config, log *slog.Logger) error {
	log.Debug("debug messages enabled")

	store, err := db.New(log, cfg.DBAddress)
	if err != nil {
		return fmt.Errorf("db connect: %w", err)
	}
	defer func() {
		if closeErr := store.Close(); closeErr != nil {
			log.Warn("db close", "error", closeErr)
		}
	}()

	if migrateErr := store.Migrate(); migrateErr != nil {
		return fmt.Errorf("migrations: %w", migrateErr)
	}

	svc := application.NewService(log, store, parser.New(log))
	mux := api.Routes(log, svc)
	srv := server.New(log, mux, server.WithAddress(":"+cfg.HTTPConfig.Port))

	ctx, stop := signal.NotifyContext(context.Background(), syscall.SIGINT, syscall.SIGTERM)
	defer stop()

	srv.Start()

	var runErr error
	select {
	case srvErr := <-srv.Notify():
		if srvErr != nil {
			runErr = fmt.Errorf("http server: %w", srvErr)
		}
	case <-ctx.Done():
		log.Info("shutdown signal")
	}

	if shutdownErr := srv.Shutdown(); shutdownErr != nil {
		return errors.Join(runErr, fmt.Errorf("shutdown: %w", shutdownErr))
	}
	return runErr
}
