package trend

import "github.com/jwhittle933/portwatch/internal/alert"

// WithTracking returns a slice of events unchanged but records the count
// of events into the provided Tracker so callers can observe the trend
// without altering the pipeline.
func WithTracking(t *Tracker, events []alert.Event) []alert.Event {
	if len(events) == 0 {
		return events
	}
	t.Record(len(events))
	return events
}
