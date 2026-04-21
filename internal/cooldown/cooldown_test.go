package cooldown

import (
	"testing"
	"time"

	"github.com/dkrichards86/portwatch/internal/alert"
	"github.com/dkrichards86/portwatch/internal/scanner"
)

func makeEvent(port uint16, proto, action string) alert.Event {
	return alert.Event{
		Port:   scanner.Port{Port: port, Protocol: proto},
		Action: action,
	}
}

func TestAllow_FirstCallAlwaysPasses(t *testing.T) {
	tr := New(5 * time.Second)
	if !tr.Allow("tcp:80:opened") {
		t.Fatal("expected first call to be allowed")
	}
}

func TestAllow_SecondCallWithinPeriodBlocked(t *testing.T) {
	tr := New(5 * time.Second)
	tr.Allow("tcp:80:opened")
	if tr.Allow("tcp:80:opened") {
		t.Fatal("expected second call within period to be blocked")
	}
}

func TestAllow_PassesAfterPeriodExpires(t *testing.T) {
	now := time.Now()
	tr := newWithClock(5*time.Second, func() time.Time { return now })
	tr.Allow("tcp:80:opened")
	now = now.Add(6 * time.Second)
	if !tr.Allow("tcp:80:opened") {
		t.Fatal("expected call after period to be allowed")
	}
}

func TestAllow_DifferentKeysAreIndependent(t *testing.T) {
	tr := New(5 * time.Second)
	tr.Allow("tcp:80:opened")
	if !tr.Allow("tcp:443:opened") {
		t.Fatal("expected different key to be allowed")
	}
}

func TestReset_AllowsKeyImmediately(t *testing.T) {
	tr := New(5 * time.Second)
	tr.Allow("tcp:80:opened")
	tr.Reset("tcp:80:opened")
	if !tr.Allow("tcp:80:opened") {
		t.Fatal("expected key to be allowed after reset")
	}
}

func TestLen_TracksEntries(t *testing.T) {
	tr := New(5 * time.Second)
	if tr.Len() != 0 {
		t.Fatalf("expected 0 entries, got %d", tr.Len())
	}
	tr.Allow("tcp:80:opened")
	tr.Allow("tcp:443:opened")
	if tr.Len() != 2 {
		t.Fatalf("expected 2 entries, got %d", tr.Len())
	}
}

func TestFilterEvents_AllowsOnFirstSeen(t *testing.T) {
	tr := New(5 * time.Second)
	events := []alert.Event{makeEvent(80, "tcp", "opened")}
	out := FilterEvents(tr, events)
	if len(out) != 1 {
		t.Fatalf("expected 1 event, got %d", len(out))
	}
}

func TestFilterEvents_SuppressesDuplicateWithinPeriod(t *testing.T) {
	tr := New(5 * time.Second)
	events := []alert.Event{makeEvent(80, "tcp", "opened")}
	FilterEvents(tr, events)
	out := FilterEvents(tr, events)
	if len(out) != 0 {
		t.Fatalf("expected 0 events, got %d", len(out))
	}
}

func TestFilterEvents_EmptyInputReturnsEmpty(t *testing.T) {
	tr := New(5 * time.Second)
	out := FilterEvents(tr, nil)
	if len(out) != 0 {
		t.Fatalf("expected empty output, got %d", len(out))
	}
}
