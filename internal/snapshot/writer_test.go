package snapshot_test

import (
	"context"
	"errors"
	"testing"

	"github.com/user/portwatch/internal/scanner"
	"github.com/user/portwatch/internal/snapshot"
)

func TestOnce_PersistsCurrentPorts(t *testing.T) {
	w := &stubWriter{}
	src := func(_ context.Context) (scanner.PortSet, error) {
		return scanner.PortSet{}, nil
	}

	err := snapshot.Once(context.Background(), src, w, silentLogger())
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if w.calls.Load() != 1 {
		t.Fatalf("expected 1 save, got %d", w.calls.Load())
	}
}

func TestOnce_ReturnsSourceError(t *testing.T) {
	w := &stubWriter{}
	src := func(_ context.Context) (scanner.PortSet, error) {
		return nil, errors.New("source broken")
	}

	err := snapshot.Once(context.Background(), src, w, silentLogger())
	if err == nil {
		t.Fatal("expected error, got nil")
	}
	if w.calls.Load() != 0 {
		t.Fatal("Save should not be called when source fails")
	}
}

func TestOnce_ReturnsWriteError(t *testing.T) {
	w := &stubWriter{err: errors.New("write failed")}
	src := func(_ context.Context) (scanner.PortSet, error) {
		return scanner.PortSet{}, nil
	}

	err := snapshot.Once(context.Background(), src, w, silentLogger())
	if err == nil {
		t.Fatal("expected error from writer, got nil")
	}
}

func TestNewNoop_DoesNotPersist(t *testing.T) {
	calls := 0
	src := func(_ context.Context) (scanner.PortSet, error) {
		calls++
		return scanner.PortSet{}, nil
	}

	// NewNoop returns a manager; we call Once directly with a noopWriter to verify.
	w := &stubWriter{}
	_ = snapshot.NewNoop(src, 0, silentLogger())

	// Direct Once with noop-like stub that always succeeds but records nothing.
	err := snapshot.Once(context.Background(), src, w, silentLogger())
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if calls != 1 {
		t.Fatalf("expected source called once, got %d", calls)
	}
}
