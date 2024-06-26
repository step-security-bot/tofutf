package xslog

import (
	"fmt"
	"log/slog"
	"os"

	"github.com/spf13/pflag"
)

type (
	Format string
	Config struct {
		Verbosity int
		Format    string
	}
)

const (
	DefaultFormat Format = "default"
	TextFormat    Format = "text"
	JSONFormat    Format = "json"
)

// NewConfigFromFlags adds flags to the given flagset, and, after the
// flagset is parsed by the caller, the flags populate the returned logger
// config.
func NewConfigFromFlags(flags *pflag.FlagSet) *Config {
	cfg := Config{}
	flags.IntVarP(&cfg.Verbosity, "v", "v", 0, "Logging level")
	flags.StringVar(&cfg.Format, "log-format", string(DefaultFormat), "Logging format: text or json")
	return &cfg
}

func New(cfg *Config) (*slog.Logger, error) {
	var h slog.Handler

	level := toSlogLevel(cfg.Verbosity)

	switch Format(cfg.Format) {
	case DefaultFormat:
		h = slog.NewJSONHandler(os.Stdout, &slog.HandlerOptions{Level: level, AddSource: true})
	case TextFormat:
		h = slog.NewTextHandler(os.Stdout, &slog.HandlerOptions{Level: level, AddSource: true})
	case JSONFormat:
		h = slog.NewJSONHandler(os.Stdout, &slog.HandlerOptions{Level: level, AddSource: true})
	default:
		return &slog.Logger{}, fmt.Errorf("unrecognised logging format: %s", cfg.Format)
	}

	return slog.New(h), nil
}

// toSlogLevel converts a logr v-level to a slog level.
func toSlogLevel(verbosity int) slog.Level {
	if verbosity <= 0 {
		return slog.LevelInfo
	}
	return slog.Level(-4 - (verbosity - 1))
}
