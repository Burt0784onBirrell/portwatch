package output

import (
	"fmt"
	"io"
	"os"

	"github.com/user/portwatch/internal/alert"
)

// ConsoleNotifier writes formatted alert events to an io.Writer (default: stdout).
type ConsoleNotifier struct {
	writer    io.Writer
	formatter *Formatter
}

// NewConsoleNotifier creates a ConsoleNotifier that writes to stdout using the
// given format and optional timestamp prefix.
func NewConsoleNotifier(format Format, timestamp bool) *ConsoleNotifier {
	return &ConsoleNotifier{
		writer:    os.Stdout,
		formatter: NewFormatter(format, timestamp),
	}
}

// NewConsoleNotifierWithWriter creates a ConsoleNotifier with a custom writer,
// useful for testing.
func NewConsoleNotifierWithWriter(w io.Writer, format Format, timestamp bool) *ConsoleNotifier {
	return &ConsoleNotifier{
		writer:    w,
		formatter: NewFormatter(format, timestamp),
	}
}

// Notify implements alert.Notifier and prints each event to the configured writer.
func (c *ConsoleNotifier) Notify(events []alert.Event) error {
	for _, e := range events {
		line := c.formatter.FormatEvent(e)
		if _, err := fmt.Fprintln(c.writer, line); err != nil {
			return err
		}
	}
	return nil
}
