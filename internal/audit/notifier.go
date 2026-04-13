package audit

import (
	"fmt"

	"github.com/robjkc/portwatch/internal/alert"
)

// Notifier implements alert.Notifier and writes every received event to an
// audit Logger.
type Notifier struct {
	logger *Logger
}

// NewNotifier wraps logger in an alert.Notifier.
func NewNotifier(logger *Logger) *Notifier {
	return &Notifier{logger: logger}
}

// Notify converts each alert.Event to an Entry and logs it.
func (n *Notifier) Notify(events []alert.Event) error {
	for _, ev := range events {
		e := Entry{
			Timestamp: ev.Timestamp,
			Action:    ev.Action,
			Proto:     ev.Port.Proto,
			Port:      ev.Port.Number,
			PID:       ev.Port.PID,
			Process:   ev.Port.Process,
		}
		if err := n.logger.Log(e); err != nil {
			return fmt.Errorf("audit notifier: %w", err)
		}
	}
	return nil
}
