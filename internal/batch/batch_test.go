package batch

import (
	"testing"
	"time"

	"github.com/patrickward/portwatch/internal/alert"
	"github.com/patrickward/portwatch/internal/scanner"
)

func makeEvent(port uint16) alert.Event {
	return alert.Event{
		Action: alert.ActionOpened,
		Port:   scanner.Port{Port: port, Protocol: "tcp"},
	}
}

func TestAdd_BuffersUntilMaxSize(t *testing.T) {
	b := New(3, time.Minute)
	if _, ok := b.Add(makeEvent(80)); ok {
		t.Fatal("expected no flush after first event")
	}
	if _, ok := b.Add(makeEvent(443)); ok {
		t.Fatal("expected no flush after second event")
	}
	events, ok := b.Add(makeEvent(8080))
	if !ok {
		t.Fatal("expected flush at maxSize")
	}
	if len(events) != 3 {
		t.Fatalf("want 3 events, got %d", len(events))
	}
}

func TestAdd_FlushesWhenWindowExpires(t *testing.T) {
	now := time.Now()
	b := newWithClock(10, 50*time.Millisecond, func() time.Time { return now })
	b.Add(makeEvent(80))

	// advance past the window
	now = now.Add(100 * time.Millisecond)
	events, ok := b.Add(makeEvent(443))
	if !ok {
		t.Fatal("expected flush after window elapsed")
	}
	if len(events) != 2 {
		t.Fatalf("want 2 events, got %d", len(events))
	}
}

func TestFlush_EmptyReturnsNoop(t *testing.T) {
	b := New(5, time.Minute)
	if _, ok := b.Flush(); ok {
		t.Fatal("flush on empty batcher should return false")
	}
}

func TestFlush_DrainsBuffer(t *testing.T) {
	b := New(5, time.Minute)
	b.Add(makeEvent(22))
	b.Add(makeEvent(80))
	events, ok := b.Flush()
	if !ok || len(events) != 2 {
		t.Fatalf("want 2 events flushed, got %d ok=%v", len(events), ok)
	}
	if _, ok2 := b.Flush(); ok2 {
		t.Fatal("second flush should be noop")
	}
}

func TestReady_TrueAfterWindowExpires(t *testing.T) {
	now := time.Now()
	b := newWithClock(10, 30*time.Millisecond, func() time.Time { return now })
	if b.Ready() {
		t.Fatal("empty batcher should not be ready")
	}
	b.Add(makeEvent(80))
	if b.Ready() {
		t.Fatal("should not be ready before window")
	}
	now = now.Add(60 * time.Millisecond)
	if !b.Ready() {
		t.Fatal("should be ready after window")
	}
}

func TestNew_ClampsMaxSizeToOne(t *testing.T) {
	b := New(0, time.Second)
	events, ok := b.Add(makeEvent(80))
	if !ok || len(events) != 1 {
		t.Fatal("maxSize 0 should be clamped to 1 and flush immediately")
	}
}
