package enrichment_test

import (
	"testing"

	"github.com/user/portwatch/internal/alert"
	"github.com/user/portwatch/internal/enrichment"
	"github.com/user/portwatch/internal/scanner"
)

func makeEvent(port uint16, proto, process string) alert.Event {
	return alert.Event{
		Action: alert.Opened,
		Port: scanner.Port{
			Number:   int(port),
			Protocol: proto,
			Process:  process,
		},
	}
}

func TestEnrich_WellKnownPortNoProcess(t *testing.T) {
	e := enrichment.New()
	events := []alert.Event{makeEvent(80, "tcp", "")}
	out := e.Enrich(events)
	if out[0].Port.Process != "http" {
		t.Errorf("expected process=http, got %q", out[0].Port.Process)
	}
}

func TestEnrich_WellKnownPortWithProcess(t *testing.T) {
	e := enrichment.New()
	events := []alert.Event{makeEvent(443, "tcp", "nginx")}
	out := e.Enrich(events)
	if out[0].Port.Process != "nginx (https)" {
		t.Errorf("expected 'nginx (https)', got %q", out[0].Port.Process)
	}
}

func TestEnrich_UnknownPortUnchanged(t *testing.T) {
	e := enrichment.New()
	events := []alert.Event{makeEvent(9999, "tcp", "custom")}
	out := e.Enrich(events)
	if out[0].Port.Process != "custom" {
		t.Errorf("expected process unchanged, got %q", out[0].Port.Process)
	}
}

func TestEnrich_AlreadyTaggedProcessNotDuplicated(t *testing.T) {
	e := enrichment.New()
	// process already contains the service name
	events := []alert.Event{makeEvent(22, "tcp", "sshd (ssh)")}
	out := e.Enrich(events)
	if out[0].Port.Process != "sshd (ssh)" {
		t.Errorf("expected no duplicate tag, got %q", out[0].Port.Process)
	}
}

func TestEnrich_EmptySlice(t *testing.T) {
	e := enrichment.New()
	out := e.Enrich([]alert.Event{})
	if len(out) != 0 {
		t.Errorf("expected empty output, got %d events", len(out))
	}
}

func TestNewWithMap_CustomMapping(t *testing.T) {
	e := enrichment.NewWithMap(map[uint16]string{1234: "myservice"})
	events := []alert.Event{makeEvent(1234, "tcp", "")}
	out := e.Enrich(events)
	if out[0].Port.Process != "myservice" {
		t.Errorf("expected myservice, got %q", out[0].Port.Process)
	}
}
