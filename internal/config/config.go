package config

import (
	"errors"
	"os"
	"time"

	"gopkg.in/yaml.v3"
)

// FilterRule mirrors filter.Rule for YAML unmarshalling.
type FilterRule struct {
	Port     uint16 `yaml:"port"`
	Protocol string `yaml:"protocol"`
}

// Config holds the full portwatch configuration.
type Config struct {
	Interval  time.Duration `yaml:"interval"`
	LogFile   string        `yaml:"log_file"`
	AllowList []FilterRule  `yaml:"allow_list"`
	DenyList  []FilterRule  `yaml:"deny_list"`
}

// DefaultConfig returns a Config with sensible defaults.
func DefaultConfig() Config {
	return Config{
		Interval: 5 * time.Second,
		LogFile:  "",
	}
}

// Load reads and parses a YAML config file from the given path.
func Load(path string) (Config, error) {
	data, err := os.ReadFile(path)
	if err != nil {
		if errors.Is(err, os.ErrNotExist) {
			return Config{}, ErrConfigNotFound
		}
		return Config{}, err
	}

	cfg := DefaultConfig()
	if err := yaml.Unmarshal(data, &cfg); err != nil {
		return Config{}, err
	}

	if cfg.Interval <= 0 {
		return Config{}, ErrInvalidInterval
	}

	return cfg, nil
}
