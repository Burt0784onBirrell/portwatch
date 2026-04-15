package limiter

import (
	"testing"
	"time"

	"github.com/rnemeth90/portwatch/internal/alert"
	"github.com/rnemeth90/portwatch/internal/scanner"
)

func makeEvent(port uint16) alert.Event {
	return alert.Event{
		Action: alert.Opened,
		Port:   scanner.Port{Number: port, Protocol: "tcp"},
	}
}

func TestFilterEvents_AllowsWhenUnderLimit(t *testing.T) {
	l := New(10, time.Second)
	events := []alert.Event{makeEvent(80), makeEvent(443), makeEvent(8080)}
	out := FilterEvents(l, events)
	if len(out) != 3 {
		t.Fatalf("expected 3 events, got %d", len(out))
	}
}

func TestFilterEvents_DropsWhenOverLimit(t *testing.T) {
	l := New(2, time.Second)
	events := []alert.Event{makeEvent(80), makeEvent(443), makeEvent(8080)}
	out := FilterEvents(l, events)
	if len(out) != 2 {
		t.Fatalf("expected 2 events after limit, got %d", len(out))
	}
}

func TestFilterEvents_EmptyInputReturnsEmpty(t *testing.T) {
	l := New(5, time.Second)
	out := FilterEvents(l, []alert.Event{})
	if out == nil || len(out) != 0 {
		t.Fatal("expected empty non-nil slice")
	}
}

func TestFilterEvents_ResetRestoresCapacity(t *testing.T) {
	l := New(1, time.Second)
	events := []alert.Event{makeEvent(80), makeEvent(443)}
	FilterEvents(l, events) // exhausts budget
	l.Reset()
	out := FilterEvents(l, []alert.Event{makeEvent(22)})
	if len(out) != 1 {
		t.Fatalf("expected 1 event after reset, got %d", len(out))
	}
}
