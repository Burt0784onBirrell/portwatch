package config

import (
	"errors"
	"os"
)

// LoadOrDefault attempts to load config from path.
// If the file does not exist, DefaultConfig is returned without error.
// Any other error is propagated to the caller.
func LoadOrDefault(path string) (*Config, error) {
	cfg, err := Load(path)
	if err != nil {
		if errors.Is(err, os.ErrNotExist) {
			return DefaultConfig(), nil
		}
		return nil, err
	}
	return cfg, nil
}
