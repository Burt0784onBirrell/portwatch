package state_test

import (
	"os"
	"path/filepath"
	"testing"
	"time"

	"github.com/stevezaluk/portwatch/internal/state"
)

// TestStore_StaleStateIsRemovedAndReinitialised verifies that when a persisted
// state file is stale the caller can detect this, remove it, and start fresh.
func TestStore_StaleStateIsRemovedAndReinitialised(t *testing.T) {
	dir := t.TempDir()
	path := filepath.Join(dir, "state.json")

	// Bootstrap an initial store and persist some content.
	store := state.NewStore(path)
	if err := store.Save(state.Snapshot{}); err != nil {
		t.Fatalf("Save: %v", err)
	}

	// Wind back the modification time so the file appears stale.
	past := time.Now().Add(-48 * time.Hour)
	if err := os.Chtimes(path, past, past); err != nil {
		t.Fatalf("Chtimes: %v", err)
	}

	opts := state.CleanupOptions{MaxAge: time.Hour}
	removed, err := state.RemoveIfStale(path, opts)
	if err != nil {
		t.Fatalf("RemoveIfStale: %v", err)
	}
	if !removed {
		t.Fatal("expected stale state file to be removed")
	}

	// A new store pointing at the same path should load an empty snapshot.
	freshStore := state.NewStore(path)
	snap, err := freshStore.Load()
	if err != nil {
		t.Fatalf("Load after removal: %v", err)
	}
	if len(snap.ToPortSet()) != 0 {
		t.Errorf("expected empty port set after stale removal, got %d ports", len(snap.ToPortSet()))
	}
}
