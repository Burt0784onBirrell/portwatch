package tagger_test

import (
	"testing"

	"github.com/user/portwatch/internal/alert"
	"github.com/user/portwatch/internal/scanner"
	"github.com/user/portwatch/internal/tagger"
)

func makeEvent(number uint16, protocol, process string) alert.Event {
	return alert.Event{
		Action: alert.ActionOpened,
		Port: scanner.Port{
			Number:   number,
			Protocol: protocol,
			Process:  process,
		},
	}
}

func TestNew_EmptyLabelReturnsError(t *testing.T) {
	_, err := tagger.New([]tagger.Rule{{Label: "", Port: 80, Protocol: "tcp"}})
	if err == nil {
		t.Fatal("expected error for empty label")
	}
}

func TestNew_InvalidProtocolReturnsError(t *testing.T) {
	_, err := tagger.New([]tagger.Rule{{Label: "web", Port: 80, Protocol: "sctp"}})
	if err == nil {
		t.Fatal("expected error for invalid protocol")
	}
}

func TestTag_MatchingPortAndProtocol(t *testing.T) {
	tg, err := tagger.New([]tagger.Rule{{Label: "http", Port: 80, Protocol: "tcp"}})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	events := []alert.Event{makeEvent(80, "tcp", "nginx")}
	out := tg.Tag(events)
	if out[0].Port.Process != "nginx [http]" {
		t.Errorf("expected 'nginx [http]', got %q", out[0].Port.Process)
	}
}

func TestTag_NoMatchLeavesEventUnchanged(t *testing.T) {
	tg, _ := tagger.New([]tagger.Rule{{Label: "http", Port: 80, Protocol: "tcp"}})
	events := []alert.Event{makeEvent(443, "tcp", "nginx")}
	out := tg.Tag(events)
	if out[0].Port.Process != "nginx" {
		t.Errorf("expected 'nginx', got %q", out[0].Port.Process)
	}
}

func TestTag_EmptyProcessGetsLabelOnly(t *testing.T) {
	tg, _ := tagger.New([]tagger.Rule{{Label: "dns", Port: 53, Protocol: "udp"}})
	events := []alert.Event{makeEvent(53, "udp", "")}
	out := tg.Tag(events)
	if out[0].Port.Process != "[dns]" {
		t.Errorf("expected '[dns]', got %q", out[0].Port.Process)
	}
}

func TestTag_ProtocolWildcardMatchesBoth(t *testing.T) {
	tg, _ := tagger.New([]tagger.Rule{{Label: "any53", Port: 53, Protocol: ""}})
	events := []alert.Event{
		makeEvent(53, "tcp", "bind"),
		makeEvent(53, "udp", "bind"),
	}
	out := tg.Tag(events)
	for _, ev := range out {
		if ev.Port.Process != "bind [any53]" {
			t.Errorf("expected 'bind [any53]', got %q", ev.Port.Process)
		}
	}
}

func TestTag_LabelNotDuplicated(t *testing.T) {
	tg, _ := tagger.New([]tagger.Rule{{Label: "http", Port: 80, Protocol: "tcp"}})
	events := []alert.Event{makeEvent(80, "tcp", "nginx [http]")}
	out := tg.Tag(events)
	if out[0].Port.Process != "nginx [http]" {
		t.Errorf("label duplicated: %q", out[0].Port.Process)
	}
}

func TestTag_OriginalEventsUnmodified(t *testing.T) {
	tg, _ := tagger.New([]tagger.Rule{{Label: "http", Port: 80, Protocol: "tcp"}})
	orig := []alert.Event{makeEvent(80, "tcp", "nginx")}
	tg.Tag(orig)
	if orig[0].Port.Process != "nginx" {
		t.Error("original events should not be modified")
	}
}
