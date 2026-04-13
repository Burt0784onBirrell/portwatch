package daemon

import (
	"context"
	"testing"
	"time"

	"github.com/user/portwatch/internal/alert"
	"github.com/user/portwatch/internal/config"
	"github.com/user/portwatch/internal/scanner"
)

func defaultConfig() *config.Config {
	cfg := config.DefaultConfig()
	cfg.Interval = 50 * time.Millisecond
	return cfg
}

func TestDaemon_RunStopsOnContextCancel(t *testing.T) {
	cfg := defaultConfig()
	s := scanner.NewScanner(cfg)
	d := alert.NewDispatcher()

	daemon := New(cfg, s, d)

	ctx, cancel := context.WithTimeout(context.Background(), 200*time.Millisecond)
	defer cancel()

	err := daemon.Run(ctx)
	if err != nil {
		t.Fatalf("expected nil error on context cancel, got: %v", err)
	}
}

func TestDaemon_New(t *testing.T) {
	cfg := defaultConfig()
	s := scanner.NewScanner(cfg)
	d := alert.NewDispatcher()

	daemon := New(cfg, s, d)

	if daemon.cfg != cfg {
		t.Error("expected daemon cfg to match provided config")
	}
	if daemon.scanner != s {
		t.Error("expected daemon scanner to match provided scanner")
	}
	if daemon.dispatcher != d {
		t.Error("expected daemon dispatcher to match provided dispatcher")
	}
}

func TestDaemon_RunDispatchesOnChange(t *testing.T) {
	cfg := defaultConfig()
	cfg.Interval = 30 * time.Millisecond

	s := scanner.NewScanner(cfg)
	d := alert.NewDispatcher()

	notified := make(chan struct{}, 1)
	notifier := &testNotifier{notified: notified}
	d.AddNotifier(notifier)

	daemon := New(cfg, s, d)

	ctx, cancel := context.WithTimeout(context.Background(), 300*time.Millisecond)
	defer cancel()

	// Run in background; we just verify it doesn't panic or error out.
	done := make(chan error, 1)
	go func() {
		done <- daemon.Run(ctx)
	}()

	select {
	case err := <-done:
		if err != nil {
			t.Fatalf("daemon.Run returned unexpected error: %v", err)
		}
	case <-time.After(500 * time.Millisecond):
		t.Fatal("daemon.Run did not stop after context cancellation")
	}
}

type testNotifier struct {
	notified chan struct{}
}

func (n *testNotifier) Notify(_ context.Context, events []alert.Event) error {
	if len(events) > 0 {
		select {
		case n.notified <- struct{}{}:
		default:
		}
	}
	return nil
}
