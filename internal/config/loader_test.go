package config_test

import (
	"testing"

	"github.com/user/portwatch/internal/config"
)

func TestLoadOrDefault_MissingFile(t *testing.T) {
	cfg, err := config.LoadOrDefault("/does/not/exist.yaml")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if cfg == nil {
		t.Fatal("expected non-nil default config")
	}
}

func TestLoadOrDefault_ValidFile(t *testing.T) {
	path := writeTempConfig(t, "scan_interval: 20s\n")
	cfg, err := config.LoadOrDefault(path)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if cfg.ScanInterval.Seconds() != 20 {
		t.Errorf("expected 20s, got %v", cfg.ScanInterval)
	}
}

func TestLoadOrDefault_InvalidFile(t *testing.T) {
	path := writeTempConfig(t, "scan_interval: 100ms\n")
	_, err := config.LoadOrDefault(path)
	if err == nil {
		t.Fatal("expected validation error")
	}
}
