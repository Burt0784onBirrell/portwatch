package watchlist_test

import (
	"testing"

	"github.com/user/portwatch/internal/scanner"
	"github.com/user/portwatch/internal/watchlist"
)

func makePortSet(ports ...scanner.Port) scanner.PortSet {
	ps := make(scanner.PortSet)
	for _, p := range ports {
		ps[p] = struct{}{}
	}
	return ps
}

func TestNew_ValidRules(t *testing.T) {
	wl, err := watchlist.New([]string{"22/tcp", "443/tcp", "53/udp"})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(wl.Entries()) != 3 {
		t.Fatalf("expected 3 entries, got %d", len(wl.Entries()))
	}
}

func TestNew_DefaultsToTCP(t *testing.T) {
	wl, err := watchlist.New([]string{"8080"})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if wl.Entries()[0].Protocol != "tcp" {
		t.Errorf("expected protocol tcp, got %s", wl.Entries()[0].Protocol)
	}
}

func TestNew_InvalidPort(t *testing.T) {
	_, err := watchlist.New([]string{"notaport"})
	if err == nil {
		t.Fatal("expected error for invalid port, got nil")
	}
}

func TestNew_InvalidProtocol(t *testing.T) {
	_, err := watchlist.New([]string{"80/sctp"})
	if err == nil {
		t.Fatal("expected error for unknown protocol, got nil")
	}
}

func TestMissingFrom_AllPresent(t *testing.T) {
	wl, _ := watchlist.New([]string{"22/tcp", "443/tcp"})
	ps := makePortSet(
		scanner.Port{Port: 22, Protocol: "tcp"},
		scanner.Port{Port: 443, Protocol: "tcp"},
	)
	if missing := wl.MissingFrom(ps); len(missing) != 0 {
		t.Errorf("expected no missing ports, got %v", missing)
	}
}

func TestMissingFrom_SomeMissing(t *testing.T) {
	wl, _ := watchlist.New([]string{"22/tcp", "443/tcp"})
	ps := makePortSet(scanner.Port{Port: 22, Protocol: "tcp"})
	missing := wl.MissingFrom(ps)
	if len(missing) != 1 {
		t.Fatalf("expected 1 missing entry, got %d", len(missing))
	}
	if missing[0].Port != 443 {
		t.Errorf("expected missing port 443, got %d", missing[0].Port)
	}
}

func TestMissingFrom_EmptyPortSet(t *testing.T) {
	wl, _ := watchlist.New([]string{"80/tcp", "53/udp"})
	missing := wl.MissingFrom(make(scanner.PortSet))
	if len(missing) != 2 {
		t.Errorf("expected 2 missing entries, got %d", len(missing))
	}
}
