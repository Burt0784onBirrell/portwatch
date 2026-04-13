package filter

import (
	"testing"
)

func TestParseRule_SinglePort(t *testing.T) {
	r, err := ParseRule("deny:22")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if r.Action != ActionDeny || r.Low != 22 || r.High != 22 || r.Protocol != "" {
		t.Errorf("unexpected rule: %+v", r)
	}
}

func TestParseRule_AllowWithProtocol(t *testing.T) {
	r, err := ParseRule("allow:443/tcp")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if r.Action != ActionAllow || r.Low != 443 || r.High != 443 || r.Protocol != "tcp" {
		t.Errorf("unexpected rule: %+v", r)
	}
}

func TestParseRule_Range(t *testing.T) {
	r, err := ParseRule("deny:1000-2000/udp")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if r.Action != ActionDeny || r.Low != 1000 || r.High != 2000 || r.Protocol != "udp" {
		t.Errorf("unexpected rule: %+v", r)
	}
}

func TestParseRule_DefaultActionIsDeny(t *testing.T) {
	r, err := ParseRule("8080")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if r.Action != ActionDeny {
		t.Errorf("expected deny, got %v", r.Action)
	}
}

func TestParseRule_InvalidRange(t *testing.T) {
	_, err := ParseRule("deny:2000-1000")
	if err == nil {
		t.Error("expected error for inverted range")
	}
}

func TestParseRule_UnknownProtocol(t *testing.T) {
	_, err := ParseRule("deny:80/icmp")
	if err == nil {
		t.Error("expected error for unknown protocol")
	}
}

func TestParseRule_EmptyString(t *testing.T) {
	_, err := ParseRule("")
	if err == nil {
		t.Error("expected error for empty rule")
	}
}

func TestParseRule_InvalidPort(t *testing.T) {
	_, err := ParseRule("deny:notaport")
	if err == nil {
		t.Error("expected error for non-numeric port")
	}
}

func TestParseRule_ZeroPort(t *testing.T) {
	_, err := ParseRule("deny:0")
	if err == nil {
		t.Error("expected error for port 0")
	}
}
