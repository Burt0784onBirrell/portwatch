package output

import (
	"fmt"
	"os"
	"sync"

	"github.com/user/portwatch/internal/alert"
)

// FileNotifier writes alert events to a file on disk.
type FileNotifier struct {
	mu        sync.Mutex
	path      string
	formatter *Formatter
}

// NewFileNotifier creates a FileNotifier that appends events to the given path.
// The file is created if it does not exist.
func NewFileNotifier(path string, format string) (*FileNotifier, error) {
	f, err := os.OpenFile(path, os.O_CREATE|os.O_APPEND|os.O_WRONLY, 0o644)
	if err != nil {
		return nil, fmt.Errorf("file notifier: open %q: %w", path, err)
	}
	f.Close()

	return &FileNotifier{
		path:      path,
		formatter: NewFormatter(format),
	}, nil
}

// Notify appends each event to the log file, one line per event.
func (fn *FileNotifier) Notify(events []alert.Event) error {
	if len(events) == 0 {
		return nil
	}

	fn.mu.Lock()
	defer fn.mu.Unlock()

	f, err := os.OpenFile(fn.path, os.O_APPEND|os.O_WRONLY, 0o644)
	if err != nil {
		return fmt.Errorf("file notifier: open %q: %w", fn.path, err)
	}
	defer f.Close()

	for _, ev := range events {
		line, err := fn.formatter.Format(ev)
		if err != nil {
			return fmt.Errorf("file notifier: format: %w", err)
		}
		if _, err := fmt.Fprintln(f, line); err != nil {
			return fmt.Errorf("file notifier: write: %w", err)
		}
	}
	return nil
}
