package digest_test

import (
	"testing"

	"github.com/user/portwatch/internal/digest"
	"github.com/user/portwatch/internal/scanner"
)

func makePort(proto string, port uint16) scanner.Port {
	return scanner.Port{Protocol: proto, Port: port}
}

func makePortSet(ports ...scanner.Port) scanner.PortSet {
	ps := make(scanner.PortSet)
	for _, p := range ports {
		ps[p] = struct{}{}
	}
	return ps
}

func TestOf_EmptyPortSet_MatchesConstant(t *testing.T) {
	d := digest.Of(make(scanner.PortSet))
	if d != digest.Empty {
		t.Errorf("expected Empty digest, got %s", d)
	}
}

func TestOf_SamePortSet_ProducesSameDigest(t *testing.T) {
	ps := makePortSet(makePort("tcp", 80), makePort("tcp", 443))
	a := digest.Of(ps)
	b := digest.Of(ps)
	if !digest.Equal(a, b) {
		t.Errorf("expected equal digests, got %s and %s", a, b)
	}
}

func TestOf_DifferentPortSets_ProduceDifferentDigests(t *testing.T) {
	ps1 := makePortSet(makePort("tcp", 80))
	ps2 := makePortSet(makePort("tcp", 8080))
	if digest.Equal(digest.Of(ps1), digest.Of(ps2)) {
		t.Error("expected different digests for different port sets")
	}
}

func TestOf_OrderIndependent(t *testing.T) {
	// Build two port sets with the same ports added in different iteration
	// orders (map iteration is random, so we add them explicitly).
	ps1 := makePortSet(makePort("tcp", 22), makePort("tcp", 80), makePort("udp", 53))
	ps2 := makePortSet(makePort("udp", 53), makePort("tcp", 22), makePort("tcp", 80))

	if !digest.Equal(digest.Of(ps1), digest.Of(ps2)) {
		t.Error("digest should be order-independent")
	}
}

func TestOf_AddingPort_ChangesDigest(t *testing.T) {
	base := makePortSet(makePort("tcp", 80))
	before := digest.Of(base)

	base[makePort("tcp", 443)] = struct{}{}
	after := digest.Of(base)

	if digest.Equal(before, after) {
		t.Error("expected digest to change after adding a port")
	}
}

func TestEqual_Reflexive(t *testing.T) {
	ps := makePortSet(makePort("tcp", 9090))
	d := digest.Of(ps)
	if !digest.Equal(d, d) {
		t.Error("digest should equal itself")
	}
}
