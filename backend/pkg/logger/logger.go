package logger

import (
	"log/slog"
	"os"

	"github.com/kingqaquuu/stackbill/internal/config"
)

// Init configures the global slog logger based on the application config.
func Init(cfg *config.LogConfig) {
	level := parseLevel(cfg.Level)
	handlerOpts := &slog.HandlerOptions{Level: level}

	var handler slog.Handler
	if cfg.Format == "json" {
		handler = slog.NewJSONHandler(os.Stdout, handlerOpts)
	} else {
		handler = slog.NewTextHandler(os.Stdout, handlerOpts)
	}

	slog.SetDefault(slog.New(handler))
	slog.Info("logger initialized", "level", cfg.Level, "format", cfg.Format)
}

func parseLevel(s string) slog.Level {
	switch s {
	case "debug":
		return slog.LevelDebug
	case "warn":
		return slog.LevelWarn
	case "error":
		return slog.LevelError
	default:
		return slog.LevelInfo
	}
}
