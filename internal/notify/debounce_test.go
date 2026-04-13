package notify

import (
	"testing"
	"time"

	"github.com/user/portwatch/internal/alert"
	"github.com/user/portwatch/internal/scanner"
)

func makeEvent(action, proto string, port int) alert.Event {
	return alert.Event{
		Action: action,
		Port:   scanner.Port{Protocol: proto, Number: port},
	}
}

func TestDebouncer_FirstEventAlwaysPasses(t *testing.T) {
	d := NewDebouncer(5 * time.Second)
	events := []alert.Event{makeEvent("opened", "tcp", 80)}
	out := d.Filter(events)
	if len(out) != 1 {
		t.Fatalf("expected 1 event, got %d", len(out))
	}
}

func TestDebouncer_DuplicateWithinWindowSuppressed(t *testing.T) {
	now := time.Now()
	d := newDebouncerWithClock(5*time.Second, func() time.Time { return now })

	events := []alert.Event{makeEvent("opened", "tcp", 80)}
	d.Filter(events)

	out := d.Filter(events)
	if len(out) != 0 {
		t.Fatalf("expected 0 events, got %d", len(out))
	}
}

func TestDebouncer_PassesAfterWindowExpires(t *testing.T) {
	now := time.Now()
	d := newDebouncerWithClock(5*time.Second, func() time.Time { return now })

	events := []alert.Event{makeEvent("opened", "tcp", 80)}
	d.Filter(events)

	now = now.Add(6 * time.Second)
	out := d.Filter(events)
	if len(out) != 1 {
		t.Fatalf("expected 1 event after window, got %d", len(out))
	}
}

func TestDebouncer_DifferentActionsAreIndependent(t *testing.T) {
	d := NewDebouncer(5 * time.Second)

	opened := []alert.Event{makeEvent("opened", "tcp", 443)}
	closed := []alert.Event{makeEvent("closed", "tcp", 443)}

	d.Filter(opened)
	out := d.Filter(closed)
	if len(out) != 1 {
		t.Fatalf("expected closed event to pass, got %d", len(out))
	}
}

func TestDebouncer_Reset_AllowsRepeat(t *testing.T) {
	now := time.Now()
	d := newDebouncerWithClock(1*time.Minute, func() time.Time { return now })

	events := []alert.Event{makeEvent("opened", "udp", 53)}
	d.Filter(events)
	d.Reset()

	out := d.Filter(events)
	if len(out) != 1 {
		t.Fatalf("expected event after reset, got %d", len(out))
	}
}

func TestDebouncer_EmptyInputReturnsEmpty(t *testing.T) {
	d := NewDebouncer(5 * time.Second)
	out := d.Filter(nil)
	if len(out) != 0 {
		t.Fatalf("expected empty output, got %d", len(out))
	}
}
