package filter_test

import (
	"testing"

	"github.com/user/portwatch/internal/filter"
	"github.com/user/portwatch/internal/scanner"
)

func makePort(number uint16, proto string) scanner.Port {
	return scanner.Port{Number: number, Protocol: proto}
}

func TestFilter_NoRules_AllowsAll(t *testing.T) {
	f := filter.New(nil, nil)
	ports := []scanner.Port{makePort(80, "tcp"), makePort(443, "tcp"), makePort(53, "udp")}
	got := f.Apply(ports)
	if len(got) != len(ports) {
		t.Fatalf("expected %d ports, got %d", len(ports), len(got))
	}
}

func TestFilter_DenyList_BlocksPort(t *testing.T) {
	deny := []filter.Rule{{Port: 80, Protocol: "tcp"}}
	f := filter.New(nil, deny)
	ports := []scanner.Port{makePort(80, "tcp"), makePort(443, "tcp")}
	got := f.Apply(ports)
	if len(got) != 1 || got[0].Number != 443 {
		t.Fatalf("expected only port 443, got %+v", got)
	}
}

func TestFilter_AllowList_RestrictsToAllowed(t *testing.T) {
	allow := []filter.Rule{{Port: 443, Protocol: "tcp"}}
	f := filter.New(allow, nil)
	ports := []scanner.Port{makePort(80, "tcp"), makePort(443, "tcp"), makePort(8080, "tcp")}
	got := f.Apply(ports)
	if len(got) != 1 || got[0].Number != 443 {
		t.Fatalf("expected only port 443, got %+v", got)
	}
}

func TestFilter_DenyOverridesAllow(t *testing.T) {
	allow := []filter.Rule{{Port: 443, Protocol: "tcp"}, {Port: 80, Protocol: "tcp"}}
	deny := []filter.Rule{{Port: 80, Protocol: "tcp"}}
	f := filter.New(allow, deny)
	ports := []scanner.Port{makePort(80, "tcp"), makePort(443, "tcp")}
	got := f.Apply(ports)
	if len(got) != 1 || got[0].Number != 443 {
		t.Fatalf("expected only port 443, got %+v", got)
	}
}

func TestFilter_ProtocolWildcard(t *testing.T) {
	deny := []filter.Rule{{Port: 53}} // no protocol = match both
	f := filter.New(nil, deny)
	ports := []scanner.Port{makePort(53, "tcp"), makePort(53, "udp"), makePort(80, "tcp")}
	got := f.Apply(ports)
	if len(got) != 1 || got[0].Number != 80 {
		t.Fatalf("expected only port 80, got %+v", got)
	}
}

func TestFilter_EmptyPorts(t *testing.T) {
	f := filter.New(nil, nil)
	got := f.Apply([]scanner.Port{})
	if len(got) != 0 {
		t.Fatalf("expected empty result, got %+v", got)
	}
}
