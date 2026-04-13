package redact_test

import (
	"testing"

	"github.com/yourusername/portwatch/internal/redact"
)

func TestString_SafeValuePassesThrough(t *testing.T) {
	r := redact.NewDefault()
	got := r.String("nginx")
	if got != "nginx" {
		t.Fatalf("expected %q, got %q", "nginx", got)
	}
}

func TestString_SensitiveValueIsRedacted(t *testing.T) {
	r := redact.NewDefault()
	for _, s := range []string{"mypassword", "db_secret", "API_TOKEN", "apikey123"} {
		got := r.String(s)
		if got != "[REDACTED]" {
			t.Errorf("String(%q) = %q, want [REDACTED]", s, got)
		}
	}
}

func TestString_CaseInsensitive(t *testing.T) {
	r := redact.NewDefault()
	got := r.String("MyPassWord")
	if got != "[REDACTED]" {
		t.Fatalf("expected [REDACTED], got %q", got)
	}
}

func TestProcessName_OnlySensitiveTokenRedacted(t *testing.T) {
	r := redact.NewDefault()
	cmd := "myapp --host localhost --password s3cr3t"
	got := r.ProcessName(cmd)
	want := "myapp --host localhost --password [REDACTED]"
	if got != want {
		t.Fatalf("ProcessName() = %q, want %q", got, want)
	}
}

func TestProcessName_NoSensitiveTokens(t *testing.T) {
	r := redact.NewDefault()
	cmd := "nginx -g 'daemon off;'"
	got := r.ProcessName(cmd)
	if got != cmd {
		t.Fatalf("expected command unchanged, got %q", got)
	}
}

func TestNew_CustomPatternsAndPlaceholder(t *testing.T) {
	r := redact.New([]string{"banana"}, "***")
	if got := r.String("banana_split"); got != "***" {
		t.Errorf("expected ***, got %q", got)
	}
	if got := r.String("apple"); got != "apple" {
		t.Errorf("expected apple, got %q", got)
	}
}

func TestDefaultPatterns_ArePresent(t *testing.T) {
	if len(redact.DefaultPatterns) == 0 {
		t.Fatal("DefaultPatterns must not be empty")
	}
}
