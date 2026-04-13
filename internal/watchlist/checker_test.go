package watchlist_test

import (
	"testing"

	"github.com/user/portwatch/internal/scanner"
	"github.com/user/portwatch/internal/watchlist"
)

func TestChecker_NoMissingPorts(t *testing.T) {
	wl, _ := watchlist.New([]string{"22/tcp"})
	ch := watchlist.NewChecker(wl)
	ps := makePortSet(scanner.Port{Port: 22, Protocol: "tcp"})

	events := ch.Check(ps)
	if len(events) != 0 {
		t.Errorf("expected no events, got %d", len(events))
	}
}

func TestChecker_MissingPortProducesClosedEvent(t *testing.T) {
	wl, _ := watchlist.New([]string{"22/tcp"})
	ch := watchlist.NewChecker(wl)

	events := ch.Check(make(scanner.PortSet))
	if len(events) != 1 {
		t.Fatalf("expected 1 event, got %d", len(events))
	}
	if events[0].Action != "closed" {
		t.Errorf("expected action 'closed', got %q", events[0].Action)
	}
	if events[0].Port.Port != 22 {
		t.Errorf("expected port 22, got %d", events[0].Port.Port)
	}
}

func TestChecker_MultipleMissingPorts(t *testing.T) {
	wl, _ := watchlist.New([]string{"22/tcp", "443/tcp", "53/udp"})
	ch := watchlist.NewChecker(wl)

	events := ch.Check(make(scanner.PortSet))
	if len(events) != 3 {
		t.Fatalf("expected 3 events, got %d", len(events))
	}
	for _, e := range events {
		if e.Action != "closed" {
			t.Errorf("expected action 'closed', got %q", e.Action)
		}
		if e.Message == "" {
			t.Error("expected non-empty message")
		}
	}
}

func TestChecker_EmptyWatchlist(t *testing.T) {
	wl, _ := watchlist.New(nil)
	ch := watchlist.NewChecker(wl)

	events := ch.Check(make(scanner.PortSet))
	if len(events) != 0 {
		t.Errorf("expected no events for empty watchlist, got %d", len(events))
	}
}
