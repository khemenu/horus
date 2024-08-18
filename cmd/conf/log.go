package conf

import (
	"context"
	"log/slog"
	"os"
)

type discardHandler struct{}

func (discardHandler) Enabled(context.Context, slog.Level) bool  { return false }
func (discardHandler) Handle(context.Context, slog.Record) error { return nil }
func (d discardHandler) WithAttrs([]slog.Attr) slog.Handler      { return d }
func (d discardHandler) WithGroup(string) slog.Handler           { return d }

type LogConfig struct {
	Enabled *bool      `yaml:"enabled"`
	Format  string     `yaml:"format"` // "text" | "json"
	Level   slog.Level `yaml:"level"`  // "debug" | "info" | "warn" | "error"
}

func (c *LogConfig) NewLogger() *slog.Logger {
	if c.Enabled != nil {
		if !*c.Enabled {
			return slog.New(discardHandler{})
		}
	}

	opts := &slog.HandlerOptions{
		Level: c.Level,
	}

	var logger *slog.Logger
	switch c.Format {
	case "text":
		logger = slog.New(slog.NewTextHandler(os.Stderr, opts))

	case "json":
		logger = slog.New(slog.NewJSONHandler(os.Stderr, opts))

	default:
		logger = slog.New(slog.NewJSONHandler(os.Stderr, opts))
		logger.Warn("fallback logger format to \"json\" since configured format is unknown", slog.String("given", c.Format))
	}

	return logger
}
