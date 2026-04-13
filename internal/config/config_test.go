package config_test

import (
	"os"
	"testing"
	"time"

	"github.com/user/portwatch/internal/config"
)

func writeTempConfig(t *testing.T, content string) string {
	t.Helper()
	f, err := os.CreateTemp(t.TempDir(), "portwatch-*.yaml")
	if err != nil {
		t.Fatalf("create temp file: %v", err)
	}
	if _, err := f.WriteString(content); err != nil {
		t.Fatalf("write temp file: %v", err)
	}
	f.Close()
	return f.Name()
}

func TestDefaultConfig(t *testing.T) {
	cfg := config.DefaultConfig()
	if cfg.ScanInterval != 15*time.Second {
		t.Errorf("expected 15s, got %v", cfg.ScanInterval)
	}
	if cfg.Log.Level != "info" {
		t.Errorf("expected log level 'info', got %q", cfg.Log.Level)
	}
}

func TestLoad_ValidFile(t *testing.T) {
	path := writeTempConfig(t, `
scan_interval: 30s
ports:
  allowlist: [80, 443]
  ignore: [8080]
log:
  level: debug
  format: json
`)
	cfg, err := config.Load(path)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if cfg.ScanInterval != 30*time.Second {
		t.Errorf("expected 30s, got %v", cfg.ScanInterval)
	}
	if len(cfg.Ports.Allowlist) != 2 {
		t.Errorf("expected 2 allowlist ports, got %d", len(cfg.Ports.Allowlist))
	}
	if cfg.Log.Format != "json" {
		t.Errorf("expected format 'json', got %q", cfg.Log.Format)
	}
}

func TestLoad_FileNotFound(t *testing.T) {
	_, err := config.Load("/nonexistent/portwatch.yaml")
	if err == nil {
		t.Fatal("expected error for missing file")
	}
}

func TestLoad_InvalidInterval(t *testing.T) {
	path := writeTempConfig(t, "scan_interval: 500ms\n")
	_, err := config.Load(path)
	if err != config.ErrInvalidInterval {
		t.Errorf("expected ErrInvalidInterval, got %v", err)
	}
}

func TestLoad_UnknownField(t *testing.T) {
	path := writeTempConfig(t, "unknown_field: true\n")
	_, err := config.Load(path)
	if err == nil {
		t.Fatal("expected error for unknown field")
	}
}
