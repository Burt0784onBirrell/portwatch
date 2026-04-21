package groupby_test

import (
	"testing"

	"github.com/user/portwatch/internal/alert"
	"github.com/user/portwatch/internal/groupby"
	"github.com/user/portwatch/internal/scanner"
)

func makeEvent(action, proto string, port uint16) alert.Event {
	return alert.Event{
		Action: action,
		Port:   scanner.Port{Port: port, Protocol: proto},
	}
}

func TestNew_NilKeyFuncReturnsError(t *testing.T) {
	_, err := groupby.New(nil)
	if err == nil {
		t.Fatal("expected error for nil key func")
	}
}

func TestNew_ValidKeyFunc(t *testing.T) {
	g, err := groupby.New(groupby.ByAction)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if g == nil {
		t.Fatal("expected non-nil Grouper")
	}
}

func TestApply_EmptyInputReturnsNil(t *testing.T) {
	g, _ := groupby.New(groupby.ByAction)
	if got := g.Apply(nil); got != nil {
		t.Fatalf("expected nil, got %v", got)
	}
}

func TestApply_GroupsByAction(t *testing.T) {
	g, _ := groupby.New(groupby.ByAction)
	events := []alert.Event{
		makeEvent("opened", "tcp", 80),
		makeEvent("closed", "tcp", 443),
		makeEvent("opened", "tcp", 8080),
	}

	groups := g.Apply(events)
	if len(groups) != 2 {
		t.Fatalf("expected 2 groups, got %d", len(groups))
	}
	if groups[0].Key != "opened" || len(groups[0].Events) != 2 {
		t.Errorf("unexpected first group: %+v", groups[0])
	}
	if groups[1].Key != "closed" || len(groups[1].Events) != 1 {
		t.Errorf("unexpected second group: %+v", groups[1])
	}
}

func TestApply_GroupsByProtocol(t *testing.T) {
	g, _ := groupby.New(groupby.ByProtocol)
	events := []alert.Event{
		makeEvent("opened", "tcp", 80),
		makeEvent("opened", "udp", 53),
		makeEvent("closed", "tcp", 443),
	}

	groups := g.Apply(events)
	if len(groups) != 2 {
		t.Fatalf("expected 2 groups, got %d", len(groups))
	}
	if groups[0].Key != "tcp" || len(groups[0].Events) != 2 {
		t.Errorf("unexpected tcp group: %+v", groups[0])
	}
	if groups[1].Key != "udp" || len(groups[1].Events) != 1 {
		t.Errorf("unexpected udp group: %+v", groups[1])
	}
}

func TestApply_PreservesInsertionOrder(t *testing.T) {
	g, _ := groupby.New(groupby.ByAction)
	events := []alert.Event{
		makeEvent("closed", "tcp", 22),
		makeEvent("opened", "tcp", 80),
	}

	groups := g.Apply(events)
	if groups[0].Key != "closed" {
		t.Errorf("expected first group to be 'closed', got %q", groups[0].Key)
	}
}

func TestApply_CustomKeyFunc(t *testing.T) {
	keyByPort := func(ev alert.Event) string {
		if ev.Port.Port < 1024 {
			return "privileged"
		}
		return "ephemeral"
	}
	g, _ := groupby.New(keyByPort)
	events := []alert.Event{
		makeEvent("opened", "tcp", 80),
		makeEvent("opened", "tcp", 8080),
		makeEvent("opened", "tcp", 443),
	}

	groups := g.Apply(events)
	if len(groups) != 2 {
		t.Fatalf("expected 2 groups, got %d", len(groups))
	}
	if groups[0].Key != "privileged" || len(groups[0].Events) != 2 {
		t.Errorf("unexpected privileged group: %+v", groups[0])
	}
	if groups[1].Key != "ephemeral" || len(groups[1].Events) != 1 {
		t.Errorf("unexpected ephemeral group: %+v", groups[1])
	}
}
