package suppress

import (
	"testing"
	"time"

	"github.com/joshbeard/portwatch/internal/alert"
	"github.com/joshbeard/portwatch/internal/scanner"
)

func advanceable(t *testing.T) (*Store, *time.Time) {
	t.Helper()
	now := time.Now()
	s := newWithClock(func() time.Time { return now })
	return s, &now
}

func TestSuppress_IsSuppressedAfterSet(t *testing.T) {
	s, _ := advanceable(t)
	s.Suppress("tcp:8080:opened", 5*time.Minute)
	if !s.IsSuppressed("tcp:8080:opened") {
		t.Fatal("expected key to be suppressed")
	}
}

func TestSuppress_NotSuppressedByDefault(t *testing.T) {
	s := New()
	if s.IsSuppressed("tcp:9090:opened") {
		t.Fatal("expected key to not be suppressed")
	}
}

func TestSuppress_ExpiresAfterDuration(t *testing.T) {
	s, now := advanceable(t)
	s.Suppress("tcp:8080:opened", 1*time.Minute)
	*now = now.Add(2 * time.Minute)
	if s.IsSuppressed("tcp:8080:opened") {
		t.Fatal("expected suppression to have expired")
	}
}

func TestSuppress_LiftRemovesEntry(t *testing.T) {
	s := New()
	s.Suppress("tcp:443:closed", 10*time.Minute)
	s.Lift("tcp:443:closed")
	if s.IsSuppressed("tcp:443:closed") {
		t.Fatal("expected suppression to be lifted")
	}
}

func TestSuppress_LenTracksEntries(t *testing.T) {
	s := New()
	if s.Len() != 0 {
		t.Fatalf("expected 0, got %d", s.Len())
	}
	s.Suppress("tcp:80:opened", time.Minute)
	s.Suppress("udp:53:opened", time.Minute)
	if s.Len() != 2 {
		t.Fatalf("expected 2, got %d", s.Len())
	}
}

func makeEvent(proto string, port uint16, action string) alert.Event {
	return alert.Event{
		Port:   scanner.Port{Protocol: proto, Number: port},
		Action: action,
	}
}

func TestFilterEvents_AllowsUnsuppressed(t *testing.T) {
	s := New()
	events := []alert.Event{makeEvent("tcp", 80, "opened")}
	out := FilterEvents(s, events)
	if len(out) != 1 {
		t.Fatalf("expected 1 event, got %d", len(out))
	}
}

func TestFilterEvents_BlocksSuppressed(t *testing.T) {
	s := New()
	e := makeEvent("tcp", 80, "opened")
	s.Suppress(KeyForEvent(e), time.Minute)
	out := FilterEvents(s, []alert.Event{e})
	if len(out) != 0 {
		t.Fatalf("expected 0 events, got %d", len(out))
	}
}

func TestFilterEvents_EmptyInput(t *testing.T) {
	s := New()
	out := FilterEvents(s, nil)
	if len(out) != 0 {
		t.Fatalf("expected 0 events, got %d", len(out))
	}
}
