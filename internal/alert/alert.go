package alert

import (
	"fmt"
	"io"
	"os"
	"time"

	"github.com/user/portwatch/internal/scanner"
)

// Level represents the severity of an alert.
type Level string

const (
	LevelInfo  Level = "INFO"
	LevelWarn  Level = "WARN"
	LevelAlert Level = "ALERT"
)

// Event describes a single port change event.
type Event struct {
	Timestamp time.Time
	Level     Level
	Message   string
	Port      scanner.Port
}

// Notifier sends alert events to a destination.
type Notifier interface {
	Notify(e Event) error
}

// LogNotifier writes alerts as formatted lines to an io.Writer.
type LogNotifier struct {
	Out io.Writer
}

// NewLogNotifier returns a LogNotifier that writes to stdout by default.
func NewLogNotifier(out io.Writer) *LogNotifier {
	if out == nil {
		out = os.Stdout
	}
	return &LogNotifier{Out: out}
}

// Notify formats and writes the event.
func (l *LogNotifier) Notify(e Event) error {
	_, err := fmt.Fprintf(
		l.Out,
		"[%s] %s  proto=%-3s port=%d pid=%d\n",
		e.Timestamp.Format(time.RFC3339),
		e.Level,
		e.Port.Proto,
		e.Port.Number,
		e.Port.PID,
	)
	return err
}

// BuildEvents converts diff results into a slice of alert Events.
func BuildEvents(opened, closed []scanner.Port) []Event {
	now := time.Now()
	events := make([]Event, 0, len(opened)+len(closed))

	for _, p := range opened {
		events = append(events, Event{
			Timestamp: now,
			Level:     LevelAlert,
			Message:   fmt.Sprintf("port opened: %d/%s (pid %d)", p.Number, p.Proto, p.PID),
			Port:      p,
		})
	}

	for _, p := range closed {
		events = append(events, Event{
			Timestamp: now,
			Level:     LevelInfo,
			Message:   fmt.Sprintf("port closed: %d/%s (pid %d)", p.Number, p.Proto, p.PID),
			Port:      p,
		})
	}

	return events
}
