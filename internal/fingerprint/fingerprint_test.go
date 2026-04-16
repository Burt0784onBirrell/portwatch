package fingerprint_test

import (
	"testing"

	"github.com/user/portwatch/internal/fingerprint"
	"github.com/user/portwatch/internal/scanner"
)

func makePortSet(entries ...scanner.Port) scanner.PortSet {
	ps := make(scanner.PortSet)
	for _, p := range entries {
		ps[p] = struct{}{}
	}
	return ps
}

func TestOf_EmptyPortSet_ReturnsSentinel(t *testing.T) {
	got := fingerprint.Of(scanner.PortSet{})
	if got != "0000000000000000" {
		t.Fatalf("expected sentinel, got %q", got)
	}
}

func TestOf_SamePortSet_ProducesSameFingerprint(t *testing.T) {
	ps := makePortSet(
		scanner.Port{Port: 80, Protocol: "tcp"},
		scanner.Port{Port: 443, Protocol: "tcp"},
	)
	a := fingerprint.Of(ps)
	b := fingerprint.Of(ps)
	if a != b {
		t.Fatalf("expected identical fingerprints, got %q and %q", a, b)
	}
}

func TestOf_DifferentPortSets_ProduceDifferentFingerprints(t *testing.T) {
	ps1 := makePortSet(scanner.Port{Port: 80, Protocol: "tcp"})
	ps2 := makePortSet(scanner.Port{Port: 8080, Protocol: "tcp"})
	if fingerprint.Of(ps1) == fingerprint.Of(ps2) {
		t.Fatal("expected different fingerprints for different port sets")
	}
}

func TestOf_OrderIndependent(t *testing.T) {
	a := makePortSet(
		scanner.Port{Port: 22, Protocol: "tcp"},
		scanner.Port{Port: 80, Protocol: "tcp"},
	)
	b := makePortSet(
		scanner.Port{Port: 80, Protocol: "tcp"},
		scanner.Port{Port: 22, Protocol: "tcp"},
	)
	if fingerprint.Of(a) != fingerprint.Of(b) {
		t.Fatal("fingerprint should be order-independent")
	}
}

func TestOf_LengthIs16(t *testing.T) {
	ps := makePortSet(scanner.Port{Port: 443, Protocol: "tcp"})
	got := fingerprint.Of(ps)
	if len(got) != 16 {
		t.Fatalf("expected length 16, got %d", len(got))
	}
}

func TestEqual_IdenticalStrings(t *testing.T) {
	if !fingerprint.Equal("abc123", "abc123") {
		t.Fatal("expected Equal to return true for identical strings")
	}
}

func TestEqual_DifferentStrings(t *testing.T) {
	if fingerprint.Equal("abc123", "xyz789") {
		t.Fatal("expected Equal to return false for different strings")
	}
}
