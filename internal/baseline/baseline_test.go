package baseline_test

import (
	"testing"
	"time"

	"github.com/yourorg/portwatch/internal/alert"
	"github.com/yourorg/portwatch/internal/baseline"
	"github.com/yourorg/portwatch/internal/scanner"
)

func makePortSet(ports ...scanner.Port) scanner.PortSet {
	ps := scanner.PortSet{}
	for _, p := range ports {
		ps[p] = struct{}{}
	}
	return ps
}

func makePort(port uint16, proto string) scanner.Port {
	return scanner.Port{Port: port, Protocol: proto}
}

func TestNew_IsNotSet(t *testing.T) {
	b := baseline.New()
	if b.IsSet() {
		t.Fatal("expected baseline to be unset initially")
	}
}

func TestCapture_SetsBaseline(t *testing.T) {
	b := baseline.New()
	ps := makePortSet(makePort(80, "tcp"))
	before := time.Now()
	b.Capture(ps)
	if !b.IsSet() {
		t.Fatal("expected baseline to be set after Capture")
	}
	if b.CapturedAt().Before(before) {
		t.Error("CapturedAt should be >= time before Capture")
	}
}

func TestDeviation_NoBaseline_ReturnsNil(t *testing.T) {
	b := baseline.New()
	current := makePortSet(makePort(443, "tcp"))
	if got := b.Deviation(current); got != nil {
		t.Errorf("expected nil without baseline, got %v", got)
	}
}

func TestDeviation_NoChange_ReturnsEmpty(t *testing.T) {
	b := baseline.New()
	ps := makePortSet(makePort(80, "tcp"))
	b.Capture(ps)
	events := b.Deviation(ps)
	if len(events) != 0 {
		t.Errorf("expected no events, got %d", len(events))
	}
}

func TestDeviation_NewPort_ReturnsOpenedEvent(t *testing.T) {
	b := baseline.New()
	base := makePortSet(makePort(80, "tcp"))
	b.Capture(base)
	current := makePortSet(makePort(80, "tcp"), makePort(8080, "tcp"))
	events := b.Deviation(current)
	if len(events) != 1 {
		t.Fatalf("expected 1 event, got %d", len(events))
	}
	if events[0].Action != alert.ActionOpened {
		t.Errorf("expected Opened, got %s", events[0].Action)
	}
	if events[0].Port.Port != 8080 {
		t.Errorf("expected port 8080, got %d", events[0].Port.Port)
	}
}

func TestDeviation_ClosedPort_ReturnsClosedEvent(t *testing.T) {
	b := baseline.New()
	base := makePortSet(makePort(80, "tcp"), makePort(443, "tcp"))
	b.Capture(base)
	current := makePortSet(makePort(80, "tcp"))
	events := b.Deviation(current)
	if len(events) != 1 {
		t.Fatalf("expected 1 event, got %d", len(events))
	}
	if events[0].Action != alert.ActionClosed {
		t.Errorf("expected Closed, got %s", events[0].Action)
	}
}

func TestSnapshot_ReturnsCopy(t *testing.T) {
	b := baseline.New()
	ps := makePortSet(makePort(22, "tcp"))
	b.Capture(ps)
	snap := b.Snapshot()
	// Mutate the snapshot; the baseline should be unaffected.
	delete(snap, makePort(22, "tcp"))
	if len(b.Snapshot()) != 1 {
		t.Error("mutating snapshot should not affect baseline")
	}
}
