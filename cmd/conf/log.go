package conf

import (
	"log/slog"
	"os"
)

type LogConfig struct {
	Format string `yaml:"format"` // "text" | "json"
}

func (c *LogConfig) NewLogger() *slog.Logger {
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
