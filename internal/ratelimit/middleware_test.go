package ratelimit_test

import (
	"testing"
	"time"

	"github.com/iamcathal/portwatch/internal/alert"
	"github.com/iamcathal/portwatch/internal/ratelimit"
	"github.com/iamcathal/portwatch/internal/scanner"
)

func makeEvent(proto string, number uint16, action string) alert.Event {
	return alert.Event{
		Port:   scanner.Port{Protocol: proto, Number: number},
		Action: action,
	}
}

func TestFilterEvents_AllowsOnFirstSeen(t *testing.T) {
	l := ratelimit.New(time.Hour)
	events := []alert.Event{
		makeEvent("tcp", 80, "opened"),
		makeEvent("tcp", 443, "opened"),
	}
	got := ratelimit.FilterEvents(l, events)
	if len(got) != 2 {
		t.Fatalf("expected 2 events, got %d", len(got))
	}
}

func TestFilterEvents_SuppressesDuplicateWithinCooldown(t *testing.T) {
	l := ratelimit.New(time.Hour)
	e := makeEvent("tcp", 8080, "opened")
	ratelimit.FilterEvents(l, []alert.Event{e}) // first pass — records key
	got := ratelimit.FilterEvents(l, []alert.Event{e})
	if len(got) != 0 {
		t.Fatalf("expected 0 events after cooldown, got %d", len(got))
	}
}

func TestFilterEvents_EmptyInput(t *testing.T) {
	l := ratelimit.New(time.Second)
	got := ratelimit.FilterEvents(l, nil)
	if len(got) != 0 {
		t.Fatalf("expected empty result for nil input, got %d", len(got))
	}
}

func TestFilterEvents_DifferentActionsAreIndependent(t *testing.T) {
	l := ratelimit.New(time.Hour)
	opened := makeEvent("tcp", 9000, "opened")
	closed := makeEvent("tcp", 9000, "closed")
	ratelimit.FilterEvents(l, []alert.Event{opened})
	got := ratelimit.FilterEvents(l, []alert.Event{closed})
	if len(got) != 1 {
		t.Fatalf("expected closed event to pass independently, got %d", len(got))
	}
}

func TestKeyForPort_Format(t *testing.T) {
	p := scanner.Port{Protocol: "udp", Number: 53}
	key := ratelimit.KeyForPort(p, "opened")
	want := "udp:53:opened"
	if key != want {
		t.Fatalf("expected %q, got %q", want, key)
	}
}
