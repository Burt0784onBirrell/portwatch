package masking_test

import (
	"testing"

	"github.com/user/portwatch/internal/alert"
	"github.com/user/portwatch/internal/masking"
	"github.com/user/portwatch/internal/scanner"
)

func makeEvent(process string) alert.Event {
	return alert.Event{
		Action: "opened",
		Port: scanner.Port{
			Number:   8080,
			Protocol: "tcp",
			Process:  process,
		},
	}
}

func TestMaskIP_ValidIPv4_Prefix24(t *testing.T) {
	m := masking.New(24)
	got := m.MaskIP("192.168.1.42")
	want := "192.168.1.0/24"
	if got != want {
		t.Errorf("MaskIP = %q; want %q", got, want)
	}
}

func TestMaskIP_ZeroPrefixZerosAll(t *testing.T) {
	m := masking.New(0)
	got := m.MaskIP("10.20.30.40")
	want := "0.0.0.0/0"
	if got != want {
		t.Errorf("MaskIP = %q; want %q", got, want)
	}
}

func TestMaskIP_FullPrefixPreservesAddress(t *testing.T) {
	m := masking.New(32)
	got := m.MaskIP("10.0.0.1")
	want := "10.0.0.1/32"
	if got != want {
		t.Errorf("MaskIP = %q; want %q", got, want)
	}
}

func TestMaskIP_InvalidAddressReturnsPlaceholder(t *testing.T) {
	m := masking.New(24)
	got := m.MaskIP("not-an-ip")
	if got != "<masked>" {
		t.Errorf("MaskIP = %q; want <masked>", got)
	}
}

func TestMaskIP_ClampsNegativePrefix(t *testing.T) {
	m := masking.New(-5)
	got := m.MaskIP("172.16.0.1")
	want := "0.0.0.0/0"
	if got != want {
		t.Errorf("MaskIP = %q; want %q", got, want)
	}
}

func TestApplyToEvents_MasksIPProcess(t *testing.T) {
	m := masking.New(16)
	events := []alert.Event{makeEvent("192.168.5.99")}
	out := m.ApplyToEvents(events)
	if len(out) != 1 {
		t.Fatalf("expected 1 event, got %d", len(out))
	}
	want := "192.168.0.0/16"
	if out[0].Port.Process != want {
		t.Errorf("Process = %q; want %q", out[0].Port.Process, want)
	}
}

func TestApplyToEvents_LeavesNonIPProcessUnchanged(t *testing.T) {
	m := masking.New(24)
	events := []alert.Event{makeEvent("nginx")}
	out := m.ApplyToEvents(events)
	if out[0].Port.Process != "nginx" {
		t.Errorf("Process = %q; want nginx", out[0].Port.Process)
	}
}

func TestApplyToEvents_DoesNotMutateOriginal(t *testing.T) {
	m := masking.New(24)
	original := makeEvent("10.0.0.1")
	events := []alert.Event{original}
	m.ApplyToEvents(events)
	if events[0].Port.Process != "10.0.0.1" {
		t.Error("original event was mutated")
	}
}

func TestApplyToEvents_EmptySlice(t *testing.T) {
	m := masking.New(24)
	out := m.ApplyToEvents(nil)
	if len(out) != 0 {
		t.Errorf("expected empty slice, got %d events", len(out))
	}
}
