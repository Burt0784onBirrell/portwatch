package state_test

import (
	"os"
	"path/filepath"
	"testing"
	"time"

	"github.com/yourorg/portwatch/internal/scanner"
	"github.com/yourorg/portwatch/internal/state"
)

func makePort(proto string, port uint16, pid int, process string) scanner.Port {
	return scanner.Port{Protocol: proto, Port: port, PID: pid, Process: process}
}

func TestStore_SaveAndLoad_RoundTrip(t *testing.T) {
	dir := t.TempDir()
	path := filepath.Join(dir, "state.json")
	store := state.NewStore(path)

	original := scanner.PortSetFromSlice([]scanner.Port{
		makePort("tcp", 80, 100, "nginx"),
		makePort("tcp", 443, 101, "nginx"),
	})

	if err := store.Save(original); err != nil {
		t.Fatalf("Save failed: %v", err)
	}

	snap, err := store.Load()
	if err != nil {
		t.Fatalf("Load failed: %v", err)
	}

	if len(snap.Ports) != 2 {
		t.Errorf("expected 2 ports, got %d", len(snap.Ports))
	}
	if snap.UpdatedAt.IsZero() {
		t.Error("expected UpdatedAt to be set")
	}
}

func TestStore_Load_MissingFile_ReturnsEmpty(t *testing.T) {
	store := state.NewStore("/tmp/portwatch_nonexistent_state.json")
	_ = os.Remove("/tmp/portwatch_nonexistent_state.json")

	snap, err := store.Load()
	if err != nil {
		t.Fatalf("expected no error for missing file, got: %v", err)
	}
	if len(snap.Ports) != 0 {
		t.Errorf("expected empty snapshot, got %d ports", len(snap.Ports))
	}
}

func TestSnapshot_ToPortSet(t *testing.T) {
	snap := state.Snapshot{
		Ports: []scanner.Port{
			makePort("udp", 53, 200, "dnsmasq"),
		},
		UpdatedAt: time.Now(),
	}

	ps := snap.ToPortSet()
	if len(ps) != 1 {
		t.Errorf("expected 1 port in set, got %d", len(ps))
	}
}

func TestStore_Load_CorruptFile_ReturnsError(t *testing.T) {
	dir := t.TempDir()
	path := filepath.Join(dir, "bad.json")
	if err := os.WriteFile(path, []byte("not-json{"), 0o600); err != nil {
		t.Fatal(err)
	}

	store := state.NewStore(path)
	_, err := store.Load()
	if err == nil {
		t.Error("expected error for corrupt JSON, got nil")
	}
}
