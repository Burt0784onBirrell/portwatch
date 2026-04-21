package profile_test

import (
	"os"
	"path/filepath"
	"testing"

	"github.com/user/portwatch/internal/profile"
)

func TestRegistry_RegisterAndGet(t *testing.T) {
	reg := profile.NewRegistry()
	p := profile.Profile{Name: "dev", AlertCooldownSecs: 5}
	if err := reg.Register(p); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	got, err := reg.Get("dev")
	if err != nil {
		t.Fatalf("Get returned error: %v", err)
	}
	if got.Name != "dev" {
		t.Errorf("expected name dev, got %q", got.Name)
	}
}

func TestRegistry_DuplicateReturnsError(t *testing.T) {
	reg := profile.NewRegistry()
	p := profile.Profile{Name: "dev"}
	_ = reg.Register(p)
	if err := reg.Register(p); err == nil {
		t.Fatal("expected error for duplicate profile, got nil")
	}
}

func TestRegistry_EmptyNameReturnsError(t *testing.T) {
	reg := profile.NewRegistry()
	err := reg.Register(profile.Profile{Name: ""})
	if err == nil {
		t.Fatal("expected error for empty name, got nil")
	}
}

func TestRegistry_GetMissingReturnsError(t *testing.T) {
	reg := profile.NewRegistry()
	_, err := reg.Get("missing")
	if err == nil {
		t.Fatal("expected error for missing profile, got nil")
	}
}

func TestRegistry_NamesReturnsAllRegistered(t *testing.T) {
	reg := profile.NewRegistry()
	_ = reg.Register(profile.Profile{Name: "a"})
	_ = reg.Register(profile.Profile{Name: "b"})
	names := reg.Names()
	if len(names) != 2 {
		t.Errorf("expected 2 names, got %d", len(names))
	}
}

func TestLoadFile_RegistersProfiles(t *testing.T) {
	const yaml = `profiles:
  - name: dev
    alert_cooldown_secs: 5
  - name: prod
    alert_cooldown_secs: 60
    tags:
      env: production
`
	dir := t.TempDir()
	path := filepath.Join(dir, "profiles.yaml")
	if err := os.WriteFile(path, []byte(yaml), 0o644); err != nil {
		t.Fatal(err)
	}

	reg := profile.NewRegistry()
	if err := profile.LoadFile(path, reg); err != nil {
		t.Fatalf("LoadFile error: %v", err)
	}

	prod, err := reg.Get("prod")
	if err != nil {
		t.Fatalf("Get prod: %v", err)
	}
	if prod.Tags["env"] != "production" {
		t.Errorf("expected tag env=production, got %q", prod.Tags["env"])
	}
}

func TestLoadFile_MissingFileReturnsError(t *testing.T) {
	reg := profile.NewRegistry()
	err := profile.LoadFile("/no/such/file.yaml", reg)
	if err == nil {
		t.Fatal("expected error for missing file, got nil")
	}
}
