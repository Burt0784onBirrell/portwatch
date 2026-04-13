package snapshot_test

import (
	"context"
	"errors"
	"log"
	"os"
	"sync/atomic"
	"testing"
	"time"

	"github.com/user/portwatch/internal/scanner"
	"github.com/user/portwatch/internal/snapshot"
)

// stubWriter records Save calls.
type stubWriter struct {
	calls atomic.Int32
	err   error
}

func (s *stubWriter) Save(_ scanner.PortSet) error {
	s.calls.Add(1)
	return s.err
}

func silentLogger() *log.Logger {
	return log.New(os.Discard, "", 0)
}

func TestManager_CapturesOnTick(t *testing.T) {
	w := &stubWriter{}
	src := func(_ context.Context) (scanner.PortSet, error) {
		return scanner.PortSet{}, nil
	}

	mgr := snapshot.New(src, nil, 20*time.Millisecond, silentLogger())
	// Replace internal writer via the exported constructor path — use stub directly.
	mgr2 := snapshot.NewWithWriter(src, w, 20*time.Millisecond, silentLogger())

	ctx, cancel := context.WithTimeout(context.Background(), 70*time.Millisecond)
	defer cancel()
	_ = mgr
	mgr2.Run(ctx)

	if w.calls.Load()  {
		t.Fatalf("expected at least 2 saves, got %d", w.calls.Load())
	}
}

funcopsOnContextCancel(t *testing.T) {
	w := &stubWriter{}
	src := func(_ context.Context) (scanner.PortSet, error) {
		return scanner.PortSet{}, nil
	}

	mgr := snapshot.NewWithWriter(src, w, 10*time.Millisecond, silentLogger())
	ctx, cancel := context.WithCancel(context.Background())
	cancel()

	done := make(chan struct{})
	go func() { mgr.Run(ctx); close(done) }()

	select {
	case <-done:
	case <-time.After(200 * time.Millisecond):
		t.Fatal("Run did not stop after context cancellation")
	}
}

func TestManager_LogsSourceError(t *testing.T) {
	w := &stubWriter{}
	src := func(_ context.Context) (scanner.PortSet, error) {
		return nil, errors.New("scan failed")
	}

	mgr := snapshot.NewWithWriter(src, w, 10*time.Millisecond, silentLogger())
	ctx, cancel := context.WithTimeout(context.Background(), 35*time.Millisecond)
	defer cancel()
	mgr.Run(ctx)

	if w.calls.Load() != 0 {
		t.Fatalf("expected 0 saves on source error, got %d", w.calls.Load())
	}
}

func TestManager_LogsWriteError(t *testing.T) {
	w := &stubWriter{err: errors.New("disk full")}
	src := func(_ context.Context) (scanner.PortSet, error) {
		return scanner.PortSet{}, nil
	}

	mgr := snapshot.NewWithWriter(src, w, 10*time.Millisecond, silentLogger())
	ctx, cancel := context.WithTimeout(context.Background(), 35*time.Millisecond)
	defer cancel()
	mgr.Run(ctx) // must not panic

	if w.calls.Load() == 0 {
		t.Fatal("expected at least one save attempt")
	}
}
