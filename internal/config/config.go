package config

import (
	"errors"
	"os"
	"time"

	"github.com/BurntSushi/toml"
)

// Config holds all runtime configuration for portwatch.
type Config struct {
	Interval  time.Duration `toml:"interval"`
	StatePath string        `toml:"state_path"`
	LogLevel  string        `toml:"log_level"`
	Filter    FilterConfig  `toml:"filter"`
}

// FilterConfig holds allow/deny rule lists.
type FilterConfig struct {
	Allow []string `toml:"allow"`
	Deny  []string `toml:"deny"`
}

// DefaultConfig returns a Config populated with sensible defaults.
func DefaultConfig() Config {
	return Config{
		Interval:  15 * time.Second,
		StatePath: "/var/lib/portwatch/state.json",
		LogLevel:  "info",
	}
}

// Load reads a TOML config file from path and merges it over defaults.
func Load(path string) (Config, error) {
	cfg := DefaultConfig()

	data, err := os.ReadFile(path)
	if err != nil {
		if errors.Is(err, os.ErrNotExist) {
			return cfg, ErrConfigNotFound
		}
		return cfg, err
	}

	if _, err := toml.Decode(string(data), &cfg); err != nil {
		return cfg, err
	}

	if cfg.Interval <= 0 {
		return cfg, ErrInvalidInterval
	}

	return cfg, nil
}
