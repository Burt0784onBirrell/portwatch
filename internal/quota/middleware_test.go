package quota

import (
	"testing"
	"time"

	"github.com/user/portwatch/internal/alert"
	"github.com/user/portwatch/internal/scanner"
)

func makeEvent(number int, protocol string) alert.Event {
	return alert.Event{
		Port: scanner.Port{Number: number, Protocol: protocol},
		Action: alert.ActionOpened,
	}
}

func TestFilterEvents_AllowsWhenUnderQuota(t *testing.T) {
	q := New(5, time.Minute)
	events := []alert.Event{makeEvent(80, "tcp"), makeEvent(443, "tcp")}
	out := FilterEvents(q, events)
	if len(out) != 2 {
		t.Fatalf("expected 2 events, got %d", len(out))
	}
}

func TestFilterEvents_DropsWhenOverQuota(t *testing.T) {
	q := New(1, time.Minute)
	events := []alert.Event{makeEvent(80, "tcp"), makeEvent(80, "tcp")}
	out := FilterEvents(q, events)
	if len(out) != 1 {
		t.Fatalf("expected 1 event, got %d", len(out))
	}
}

func TestFilterEvents_EmptyInputReturnsEmpty(t *testing.T) {
	q := New(10, time.Minute)
	out := FilterEvents(q, nil)
	if len(out) != 0 {
		t.Fatalf("expected empty, got %d", len(out))
	}
}

func TestFilterEvents_DifferentPortsAreIndependent(t *testing.T) {
	q := New(1, time.Minute)
	events := []alert.Event{
		makeEvent(80, "tcp"),
		makeEvent(443, "tcp"),
		makeEvent(80, "tcp"), // over quota for 80
	}
	out := FilterEvents(q, events)
	if len(out) != 2 {
		t.Fatalf("expected 2 events, got %d", len(out))
	}
}

func TestKeyForEvent_IncludesProtocol(t *testing.T) {
	e1 := makeEvent(80, "tcp")
	e2 := makeEvent(80, "udp")
	if KeyForEvent(e1) == KeyForEvent(e2) {
		t.Fatal("expected different keys for same port with different protocols")
	}
}
