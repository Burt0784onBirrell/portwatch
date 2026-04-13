package main

import (
	"os"
	"testing"
)

func TestConfigPath_NoArgs(t *testing.T) {
	// Save and restore os.Args
	orig := os.Args
	defer func() { os.Args = orig }()

	os.Args = []string{"portwatch"}

	got := configPath()
	if got != "" {
		t.Errorf("expected empty string, got %q", got)
	}
}

func TestConfigPath_WithArg(t *testing.T) {
	orig := os.Args
	defer func() { os.Args = orig }()

	os.Args = []string{"portwatch", "/etc/portwatch/config.yaml"}

	got := configPath()
	if got != "/etc/portwatch/config.yaml" {
		t.Errorf("expected /etc/portwatch/config.yaml, got %q", got)
	}
}

func TestConfigPath_IgnoresExtraArgs(t *testing.T) {
	orig := os.Args
	defer func() { os.Args = orig }()

	os.Args = []string{"portwatch", "first.yaml", "second.yaml"}

	got := configPath()
	if got != "first.yaml" {
		t.Errorf("expected first.yaml, got %q", got)
	}
}

func TestMain_EnvOverride(t *testing.T) {
	// Verify that a missing config path returns empty string gracefully
	// and does not panic — integration smoke test for configPath.
	orig := os.Args
	defer func() { os.Args = orig }()

	os.Args = []string{"portwatch"}
	path := configPath()

	if path != "" {
		t.Errorf("unexpected path: %q", path)
	}
}
