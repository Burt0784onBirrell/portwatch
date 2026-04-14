package dedup

import (
	"testing"
	"time"

	"github.com/user/portwatch/internal/alert"
	"github.com/user/portwatch/internal/scanner"
)

func makeEvent(action alert.Action, number uint16, proto string) alert.Event {
	return alert.Event{
		Action: action,
		Port:   scanner.Port{Number: number, Protocol: proto},
	}
}

func TestFilter_FirstEventAlwaysPasses(t *testing.T) {
	dd := New(time.Minute)
	events := []alert.Event{makeEvent(alert.Opened, 80, "tcp")}
	out := dd.Filter(events)
	if len(out) != 1 {
		t.Fatalf("expected 1 event, got %d", len(out))
	}
}

func TestFilter_DuplicateWithinWindowSuppressed(t *testing.T) {
	now := time.Now()
	dd := newWithClock(time.Minute, func() time.Time { return now })

	events := []alert.Event{makeEvent(alert.Opened, 80, "tcp")}
	dd.Filter(events)
	out := dd.Filter(events)
	if len(out) != 0 {
		t.Fatalf("expected 0 events, got %d", len(out))
	}
}

func TestFilter_PassesAfterWindowExpires(t *testing.T) {
	now := time.Now()
	dd := newWithClock(time.Minute, func() time.Time { return now })

	events := []alert.Event{makeEvent(alert.Opened, 443, "tcp")}
	dd.Filter(events)

	now = now.Add(2 * time.Minute)
	out := dd.Filter(events)
	if len(out) != 1 {
		t.Fatalf("expected 1 event after window, got %d", len(out))
	}
}

func TestFilter_DifferentActionsAreIndependent(t *testing.T) {
	dd := New(time.Minute)
	opened := makeEvent(alert.Opened, 22, "tcp")
	closed := makeEvent(alert.Closed, 22, "tcp")

	out := dd.Filter([]alert.Event{opened, closed})
	if len(out) != 2 {
		t.Fatalf("expected 2 events, got %d", len(out))
	}
}

func TestFlush_RemovesStaleEntries(t *testing.T) {
	now := time.Now()
	dd := newWithClock(time.Minute, func() time.Time { return now })

	events := []alert.Event{makeEvent(alert.Opened, 8080, "tcp")}
	dd.Filter(events)

	now = now.Add(2 * time.Minute)
	dd.Flush()

	dd.mu.Lock()
	l := len(dd.seen)
	dd.mu.Unlock()

	if l != 0 {
		t.Fatalf("expected empty store after flush, got %d entries", l)
	}
}

func TestFilter_EmptyInputReturnsEmpty(t *testing.T) {
	dd := New(time.Minute)
	out := dd.Filter(nil)
	if len(out) != 0 {
		t.Fatalf("expected empty output, got %d", len(out))
	}
}
