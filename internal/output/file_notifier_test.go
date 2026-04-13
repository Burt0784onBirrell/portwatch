package output_test

import (
	"os"
	"path/filepath"
	"strings"
	"testing"
	"time"

	"github.com/user/portwatch/internal/alert"
	"github.com/user/portwatch/internal/output"
	"github.com/user/portwatch/internal/scanner"
)

func makeFileEvent(kind alert.EventKind, port uint16) alert.Event {
	return alert.Event{
		Kind:      kind,
		Timestamp: time.Date(2024, 1, 1, 12, 0, 0, 0, time.UTC),
		Port: scanner.Port{
			Number:   port,
			Protocol: "tcp",
			Process:  "nginx",
		},
	}
}

func TestFileNotifier_CreatesFile(t *testing.T) {
	dir := t.TempDir()
	path := filepath.Join(dir, "portwatch.log")

	_, err := output.NewFileNotifier(path, "text")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if _, err := os.Stat(path); os.IsNotExist(err) {
		t.Fatal("expected file to be created")
	}
}

func TestFileNotifier_WritesEvents(t *testing.T) {
	dir := t.TempDir()
	path := filepath.Join(dir, "portwatch.log")

	fn, err := output.NewFileNotifier(path, "text")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	events := []alert.Event{
		makeFileEvent(alert.EventOpened, 80),
		makeFileEvent(alert.EventClosed, 443),
	}
	if err := fn.Notify(events); err != nil {
		t.Fatalf("notify error: %v", err)
	}

	data, err := os.ReadFile(path)
	if err != nil {
		t.Fatalf("read error: %v", err)
	}
	content := string(data)
	if !strings.Contains(content, "80") {
		t.Errorf("expected port 80 in output, got: %s", content)
	}
	if !strings.Contains(content, "443") {
		t.Errorf("expected port 443 in output, got: %s", content)
	}
}

func TestFileNotifier_EmptyEventsIsNoop(t *testing.T) {
	dir := t.TempDir()
	path := filepath.Join(dir, "portwatch.log")

	fn, err := output.NewFileNotifier(path, "text")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if err := fn.Notify(nil); err != nil {
		t.Errorf("expected no error for empty events, got: %v", err)
	}
	info, _ := os.Stat(path)
	if info.Size() != 0 {
		t.Errorf("expected empty file, got size %d", info.Size())
	}
}

func TestFileNotifier_InvalidPath(t *testing.T) {
	_, err := output.NewFileNotifier("/nonexistent/dir/portwatch.log", "text")
	if err == nil {
		t.Fatal("expected error for invalid path")
	}
}
