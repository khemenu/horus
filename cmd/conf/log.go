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
	Enabled *bool  `yaml:"enabled"`
	Format  string `yaml:"format"` // "text" | "json"
}

func (c *LogConfig) NewLogger() *slog.Logger {
	if c.Enabled != nil {
		if !*c.Enabled {
			return slog.New(discardHandler{})
		}
	}

	var logger *slog.Logger
	switch c.Format {
	case "text":
		logger = slog.New(slog.NewTextHandler(os.Stderr, nil))

	case "json":
		logger = slog.New(slog.NewJSONHandler(os.Stderr, nil))

	default:
		panic("unreachable")
	}

	return logger
}
