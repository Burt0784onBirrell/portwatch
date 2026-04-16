package normalize_test

import (
	"testing"

	"github.com/joshbeard/portwatch/internal/alert"
	"github.com/joshbeard/portwatch/internal/normalize"
	"github.com/joshbeard/portwatch/internal/scanner"
)

func makeEvent(process, protocol string) alert.Event {
	return alert.Event{
		Action: alert.ActionOpened,
		Port: scanner.Port{
			Number:   8080,
			Protocol: protocol,
			Process:  process,
		},
	}
}

func TestApply_LowercasesProcess(t *testing.T) {
	n := normalize.New()
	out := n.Apply([]alert.Event{makeEvent("NGINX", "tcp")})
	if out[0].Port.Process != "nginx" {
		t.Errorf("expected nginx, got %q", out[0].Port.Process)
	}
}

func TestApply_TrimsProcess(t *testing.T) {
	n := normalize.New()
	out := n.Apply([]alert.Event{makeEvent("  sshd  ", "tcp")})
	if out[0].Port.Process != "sshd" {
		t.Errorf("expected sshd, got %q", out[0].Port.Process)
	}
}

func TestApply_DefaultProtocol(t *testing.T) {
	n := normalize.New()
	out := n.Apply([]alert.Event{makeEvent("sshd", "")})
	if out[0].Port.Protocol != "tcp" {
		t.Errorf("expected tcp, got %q", out[0].Port.Protocol)
	}
}

func TestApply_ExistingProtocolUnchanged(t *testing.T) {
	n := normalize.New()
	out := n.Apply([]alert.Event{makeEvent("sshd", "udp")})
	if out[0].Port.Protocol != "udp" {
		t.Errorf("expected udp, got %q", out[0].Port.Protocol)
	}
}

func TestApply_EmptyInput(t *testing.T) {
	n := normalize.New()
	out := n.Apply([]alert.Event{})
	if len(out) != 0 {
		t.Errorf("expected empty slice")
	}
}

func TestApply_WithLowercaseDisabled(t *testing.T) {
	n := normalize.New(normalize.WithLowercaseProcess(false))
	out := n.Apply([]alert.Event{makeEvent("NGINX", "tcp")})
	if out[0].Port.Process != "NGINX" {
		t.Errorf("expected NGINX, got %q", out[0].Port.Process)
	}
}

func TestApply_CustomDefaultProtocol(t *testing.T) {
	n := normalize.New(normalize.WithDefaultProtocol("udp"))
	out := n.Apply([]alert.Event{makeEvent("app", "")})
	if out[0].Port.Protocol != "udp" {
		t.Errorf("expected udp, got %q", out[0].Port.Protocol)
	}
}
