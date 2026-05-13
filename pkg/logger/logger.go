package logger

import (
	"log/slog"
	"os"
	"strings"
)

const (
	levelInfo  = "INFO"
	levelDebug = "DEBUG"
	levelError = "ERROR"
	levelWarn  = "WARN"
)

func New(level string) *slog.Logger {
	var lvl slog.Level

	switch strings.ToUpper(level) {
	case levelDebug:
		lvl = slog.LevelDebug
	case levelInfo:
		lvl = slog.LevelInfo
	case levelWarn:
		lvl = slog.LevelWarn
	case levelError:
		lvl = slog.LevelError
	default:
		lvl = slog.LevelInfo
	}

	log := slog.New(slog.NewTextHandler(
		os.Stdout,
		&slog.HandlerOptions{
			Level:     lvl,
			AddSource: lvl == slog.LevelDebug,
		},
	))

	return log
}
