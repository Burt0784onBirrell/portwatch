// Package audit provides an append-only audit trail of every port-change
// event observed by the daemon.  Entries are written to a JSON-lines file so
// they can be ingested by external tooling without any special parser.
package audit

import (
	"encoding/json"
	"fmt"
	"io"
	"os"
	"sync"
	"time"
)

// Entry is a single audit-log record.
type Entry struct {
	Timestamp time.Time `json:"timestamp"`
	Action    string    `json:"action"`
	Proto     string    `json:"proto"`
	Port      uint16    `json:"port"`
	PID       int       `json:"pid,omitempty"`
	Process   string    `json:"process,omitempty"`
}

// Logger writes audit entries to an underlying io.Writer.
type Logger struct {
	mu  sync.Mutex
	out io.Writer
}

// New returns a Logger that appends JSON-lines to path, creating the file if
// it does not already exist.
func New(path string) (*Logger, error) {
	f, err := os.OpenFile(path, os.O_CREATE|os.O_APPEND|os.O_WRONLY, 0o640)
	if err != nil {
		return nil, fmt.Errorf("audit: open %s: %w", path, err)
	}
	return NewWithWriter(f), nil
}

// NewWithWriter returns a Logger that writes to w.  Useful in tests.
func NewWithWriter(w io.Writer) *Logger {
	return &Logger{out: w}
}

// Log serialises e as a JSON line and appends it to the underlying writer.
func (l *Logger) Log(e Entry) error {
	l.mu.Lock()
	defer l.mu.Unlock()

	if e.Timestamp.IsZero() {
		e.Timestamp = time.Now().UTC()
	}

	b, err := json.Marshal(e)
	if err != nil {
		return fmt.Errorf("audit: marshal: %w", err)
	}
	b = append(b, '\n')

	_, err = l.out.Write(b)
	return err
}
