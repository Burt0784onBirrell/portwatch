package output_test

import (
	"strings"
	"testing"

	"github.com/user/portwatch/internal/alert"
	"github.com/user/portwatch/internal/output"
	"github.com/user/portwatch/internal/scanner"
)

func makeEvent(kind, proto string, number int, process string) alert.Event {
	return alert.Event{
		Kind: kind,
		Port: scanner.Port{Protocol: proto, Number: number, Process: process},
	}
}

func TestFormatter_TextFormat_Opened(t *testing.T) {
	f := output.NewFormatter(output.FormatText, false)
	e := makeEvent("opened", "tcp", 8080, "nginx")
	result := f.FormatEvent(e)
	if !strings.Contains(result, "[OPENED]") {
		t.Errorf("expected [OPENED] in output, got: %s", result)
	}
	if !strings.Contains(result, "tcp/8080") {
		t.Errorf("expected tcp/8080 in output, got: %s", result)
	}
	if !strings.Contains(result, "nginx") {
		t.Errorf("expected process name in output, got: %s", result)
	}
}

func TestFormatter_TextFormat_NoProcess(t *testing.T) {
	f := output.NewFormatter(output.FormatText, false)
	e := makeEvent("closed", "udp", 53, "")
	result := f.FormatEvent(e)
	if strings.Contains(result, "(") {
		t.Errorf("expected no process parentheses, got: %s", result)
	}
}

func TestFormatter_JSONFormat(t *testing.T) {
	f := output.NewFormatter(output.FormatJSON, false)
	e := makeEvent("opened", "tcp", 443, "sshd")
	result := f.FormatEvent(e)
	if !strings.HasPrefix(result, "{") || !strings.HasSuffix(result, "}") {
		t.Errorf("expected JSON object, got: %s", result)
	}
	if !strings.Contains(result, `"kind":"opened"`) {
		t.Errorf("expected kind field, got: %s", result)
	}
	if !strings.Contains(result, `"port":443`) {
		t.Errorf("expected port field, got: %s", result)
	}
}

func TestFormatter_TextFormat_WithTimestamp(t *testing.T) {
	f := output.NewFormatter(output.FormatText, true)
	e := makeEvent("opened", "tcp", 9090, "")
	result := f.FormatEvent(e)
	// RFC3339 timestamps start with a 4-digit year
	if len(result) < 20 || result[4] != '-' {
		t.Errorf("expected timestamp prefix, got: %s", result)
	}
}

func TestFormatter_JSONFormat_WithTimestamp(t *testing.T) {
	f := output.NewFormatter(output.FormatJSON, true)
	e := makeEvent("closed", "tcp", 22, "sshd")
	result := f.FormatEvent(e)
	if !strings.Contains(result, `"timestamp"`) {
		t.Errorf("expected timestamp field in JSON, got: %s", result)
	}
}
