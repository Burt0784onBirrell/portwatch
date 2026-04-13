package config

import (
	"os"
	"time"

	"gopkg.in/yaml.v3"
)

// Config holds the top-level portwatch configuration.
type Config struct {
	ScanInterval time.Duration `yaml:"scan_interval"`
	Ports        PortsConfig   `yaml:"ports"`
	Log          LogConfig     `yaml:"log"`
}

// PortsConfig controls which ports are monitored.
type PortsConfig struct {
	Allowlist []uint16 `yaml:"allowlist"`
	Ignore    []uint16 `yaml:"ignore"`
}

// LogConfig controls log output.
type LogConfig struct {
	Level  string `yaml:"level"`
	Format string `yaml:"format"`
}

// DefaultConfig returns a Config populated with sensible defaults.
func DefaultConfig() *Config {
	return &Config{
		ScanInterval: 15 * time.Second,
		Log: LogConfig{
			Level:  "info",
			Format: "text",
		},
	}
}

// Load reads a YAML config file from path and merges it with defaults.
func Load(path string) (*Config, error) {
	cfg := DefaultConfig()

	f, err := os.Open(path)
	if err != nil {
		return nil, err
	}
	defer f.Close()

	dec := yaml.NewDecoder(f)
	dec.KnownFields(true)
	if err := dec.Decode(cfg); err != nil {
		return nil, err
	}

	if err := cfg.validate(); err != nil {
		return nil, err
	}

	return cfg, nil
}

func (c *Config) validate() error {
	if c.ScanInterval < time.Second {
		return ErrInvalidInterval
	}
	return nil
}
