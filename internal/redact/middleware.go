package redact

import "github.com/yourusername/portwatch/internal/alert"

// ApplyToEvents returns a new slice of alert.Event values where the
// Process field of each port has been passed through the Redactor.
// The original slice is never modified.
func (r *Redactor) ApplyToEvents(events []alert.Event) []alert.Event {
	if len(events) == 0 {
		return events
	}
	out := make([]alert.Event, len(events))
	for i, ev := range events {
		cp := ev
		cp.Port.Process = r.ProcessName(ev.Port.Process)
		out[i] = cp
	}
	return out
}
