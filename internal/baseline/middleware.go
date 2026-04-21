package baseline

import (
	"github.com/yourorg/portwatch/internal/alert"
	"github.com/yourorg/portwatch/internal/scanner"
)

// DeviationFilter is a pipeline stage that filters events so that only those
// representing a deviation from the baseline are forwarded.
// Events that match the baseline (i.e. ports that are expected to be open)
// are dropped; unexpected opens/closes are passed through unchanged.
type DeviationFilter struct {
	baseline *Baseline
}

// NewDeviationFilter returns a DeviationFilter backed by the given Baseline.
func NewDeviationFilter(b *Baseline) *DeviationFilter {
	return &DeviationFilter{baseline: b}
}

// Apply returns only those events that represent a deviation from the baseline.
// If no baseline has been set all events are forwarded unchanged.
func (f *DeviationFilter) Apply(events []alert.Event) []alert.Event {
	if !f.baseline.IsSet() {
		return events
	}
	snap := f.baseline.Snapshot()
	out := events[:0:0]
	for _, ev := range events {
		port := scanner.Port{
			Port:     ev.Port.Port,
			Protocol: ev.Port.Protocol,
			Process:  ev.Port.Process,
			PID:      ev.Port.PID,
		}
		_, inBaseline := snap[port]
		// A port that is open and was expected → not a deviation.
		if ev.Action == alert.ActionOpened && inBaseline {
			continue
		}
		out = append(out, ev)
	}
	return out
}
