package config

import (
	"fmt"
	"os"
	"time"

	"gopkg.in/yaml.v3"
)

// Config holds all runtime configuration for portwatch.
type Config struct {
	Interval  time.Duration `yaml:"interval"`
	LogFile   string        `yaml:"log_file"`
	StateFile string        `yaml:"state_file"`
	Webhook   WebhookConfig `yaml:"webhook"`
	Filter    FilterConfig  `yaml:"filter"`
}

// WebhookConfig holds optional webhook notification settings.
type WebhookConfig struct {
	URL     string `yaml:"url"`
	Enabled bool   `yaml:"enabled"`
}

// FilterConfig holds port allow/deny rule strings.
type FilterConfig struct {
	Allow []string `yaml:"allow"`
	Deny  []string `yaml:"deny"`
}

// DefaultConfig returns a Config populated with sensible defaults.
func DefaultConfig() Config {
	return Config{
		Interval:  15 * time.Second,
		StateFile: "/tmp/portwatch.state",
	}
}

// Load reads a YAML config file from path and returns a validated Config.
func Load(path string) (Config, error) {
	data, err := os.ReadFile(path)
	if err != nil {
		return Config{}, fmt.Errorf("%w: %s", ErrFileNotFound, path)
	}

	cfg := DefaultConfig()
	if err := yaml.Unmarshal(data, &cfg); err != nil {
		return Config{}, fmt.Errorf("config: parse error: %w", err)
	}

	if cfg.Interval <= 0 {
		return Config{}, fmt.Errorf("config: %w", ErrInvalidInterval)
	}

	return cfg, nil
}
