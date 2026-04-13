package state

import (
	"os"
	"path/filepath"
	"testing"
	"time"
)

func writeTempState(t *testing.T, content string) string {
	t.Helper()
	dir := t.TempDir()
	p := filepath.Join(dir, "state.json")
	if err := os.WriteFile(p, []byte(content), 0o600); err != nil {
		t.Fatalf("writeTempState: %v", err)
	}
	return p
}

func TestIsStale_FreshFile(t *testing.T) {
	p := writeTempState(t, `{}`)
	opts := CleanupOptions{MaxAge: time.Hour}
	stale, err := IsStale(p, opts)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if stale {
		t.Error("expected fresh file to not be stale")
	}
}

func TestIsStale_OldFile(t *testing.T) {
	p := writeTempState(t, `{}`)
	past := time.Now().Add(-2 * time.Hour)
	if err := os.Chtimes(p, past, past); err != nil {
		t.Fatalf("chtimes: %v", err)
	}
	opts := CleanupOptions{MaxAge: time.Hour}
	stale, err := IsStale(p, opts)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if !stale {
		t.Error("expected old file to be stale")
	}
}

func TestIsStale_MissingFile(t *testing.T) {
	opts := DefaultCleanupOptions()
	stale, err := IsStale("/nonexistent/path/state.json", opts)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if stale {
		t.Error("missing file should not be reported as stale")
	}
}

func TestRemoveIfStale_RemovesOldFile(t *testing.T) {
	p := writeTempState(t, `{}`)
	past := time.Now().Add(-48 * time.Hour)
	_ = os.Chtimes(p, past, past)
	opts := CleanupOptions{MaxAge: time.Hour}
	removed, err := RemoveIfStale(p, opts)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if !removed {
		t.Error("expected stale file to be removed")
	}
	if _, err := os.Stat(p); !os.IsNotExist(err) {
		t.Error("file should no longer exist on disk")
	}
}

func TestRemoveIfStale_KeepsFreshFile(t *testing.T) {
	p := writeTempState(t, `{}`)
	opts := CleanupOptions{MaxAge: time.Hour}
	removed, err := RemoveIfStale(p, opts)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if removed {
		t.Error("expected fresh file to be kept")
	}
}
