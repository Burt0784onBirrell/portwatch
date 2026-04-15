package export_test

import (
	"bytes"
	"encoding/json"
	"strings"
	"testing"

	"github.com/user/portwatch/internal/export"
	"github.com/user/portwatch/internal/scanner"
)

func makePortSet(entries ...scanner.Port) scanner.PortSet {
	return scanner.PortSetFromSlice(entries)
}

func TestNew_ValidFormats(t *testing.T) {
	for _, f := range []export.Format{export.FormatJSON, export.FormatCSV} {
		_, err := export.New(f, &bytes.Buffer{})
		if err != nil {
			t.Errorf("expected no error for format %q, got %v", f, err)
		}
	}
}

func TestNew_InvalidFormat(t *testing.T) {
	_, err := export.New("xml", &bytes.Buffer{})
	if err == nil {
		t.Fatal("expected error for unsupported format")
	}
}

func TestWrite_JSON_ContainsPort(t *testing.T) {
	var buf bytes.Buffer
	e, _ := export.New(export.FormatJSON, &buf)

	ps := makePortSet(scanner.Port{Protocol: "tcp", Port: 8080, Process: "nginx", PID: 42})
	if err := e.Write(ps); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	var records []export.Record
	if err := json.Unmarshal(buf.Bytes(), &records); err != nil {
		t.Fatalf("invalid JSON output: %v", err)
	}
	if len(records) != 1 {
		t.Fatalf("expected 1 record, got %d", len(records))
	}
	r := records[0]
	if r.Port != 8080 || r.Protocol != "tcp" || r.Process != "nginx" || r.PID != 42 {
		t.Errorf("record fields mismatch: %+v", r)
	}
}

func TestWrite_JSON_EmptyPortSet(t *testing.T) {
	var buf bytes.Buffer
	e, _ := export.New(export.FormatJSON, &buf)

	if err := e.Write(makePortSet()); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if !strings.Contains(buf.String(), "[]") {
		t.Errorf("expected empty JSON array, got: %s", buf.String())
	}
}

func TestWrite_CSV_ContainsHeader(t *testing.T) {
	var buf bytes.Buffer
	e, _ := export.New(export.FormatCSV, &buf)

	ps := makePortSet(scanner.Port{Protocol: "udp", Port: 53, Process: "dns", PID: 7})
	if err := e.Write(ps); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	output := buf.String()
	if !strings.Contains(output, "timestamp,protocol,port,process,pid") {
		t.Errorf("CSV header missing: %s", output)
	}
	if !strings.Contains(output, "udp") || !strings.Contains(output, "53") {
		t.Errorf("CSV row missing expected fields: %s", output)
	}
}

func TestWrite_CSV_EmptyPortSet(t *testing.T) {
	var buf bytes.Buffer
	e, _ := export.New(export.FormatCSV, &buf)

	if err := e.Write(makePortSet()); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	lines := strings.Split(strings.TrimSpace(buf.String()), "\n")
	if len(lines) != 1 {
		t.Errorf("expected only header line, got %d lines", len(lines))
	}
}
