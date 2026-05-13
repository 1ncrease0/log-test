package main

import (
	"context"
	"errors"
	"fmt"
	"os"

	"log-parser/internal/application"
	"log-parser/internal/config"
	"log-parser/internal/domain"
	"log-parser/internal/infrastructure/db"
	"log-parser/internal/infrastructure/parser"
	"log-parser/pkg/logger"
)

func main() {
	cfg := config.MustLoad()
	log := logger.New(cfg.LogLevel)
	log.Debug("debug messages enabled")

	path := "data/log.zip"
	if len(os.Args) > 1 {
		path = os.Args[1]
	}

	ctx := context.Background()

	store, err := db.New(log, cfg.DBAddress)
	if err != nil {
		fmt.Fprintf(os.Stderr, "db connect: %v\n", err)
		os.Exit(1)
	}
	defer func() {
		if err := store.Close(); err != nil {
			log.Warn("db close", "error", err)
		}
	}()

	if err := store.Migrate(); err != nil {
		fmt.Fprintf(os.Stderr, "migrations: %v\n", err)
		os.Exit(1)
	}
	log.Info("migrations applied")

	logID, err := store.CreateLog(ctx, path)
	if err != nil {
		if errors.Is(err, application.ErrDuplicateLogPath) {
			fmt.Fprintf(os.Stderr, "duplicate log path (already in database): %s\n", path)
			os.Exit(2)
		}
		fmt.Fprintf(os.Stderr, "create log: %v\n", err)
		os.Exit(1)
	}
	log.Info("log record created", "log_id", logID, "path", path)

	p := parser.New(log)
	res, err := p.Parse(path)
	if err != nil {
		if sErr := store.SetStatus(ctx, logID, domain.LogStatusFailed); sErr != nil {
			log.Error("set status failed after parse error", "log_id", logID, "error", sErr)
		}
		fmt.Fprintf(os.Stderr, "parse failed: %v\n", err)
		os.Exit(1)
	}

	if err := store.SaveResult(ctx, logID, res); err != nil {
		if sErr := store.SetStatus(ctx, logID, domain.LogStatusFailed); sErr != nil {
			log.Error("set status failed after save error", "log_id", logID, "error", sErr)
		}
		fmt.Fprintf(os.Stderr, "save to db: %v\n", err)
		os.Exit(1)
	}

	log.Info("parse and save completed",
		"log_id", logID,
		"path", path,
		"nodes", len(res.Nodes),
		"ports", len(res.Ports),
		"switch_infos", len(res.SwitchInfos),
		"system_infos", len(res.SystemInfos),
		"sharp_infos", len(res.SharpInfos),
	)
}
