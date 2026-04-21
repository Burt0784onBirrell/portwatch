package decay

import (
	"fmt"
	"time"

	"github.com/jwhittle933/portwatch/internal/alert"
)

const defaultHalfLife = 5 * time.Minute

// ScoreFilter wraps a Decayer and drops events whose port score exceeds a
// threshold, preventing alert floods for persistently noisy ports.
type ScoreFilter struct {
	d         *Decayer
	threshold float64
	delta     float64
}

// NewScoreFilter returns a ScoreFilter that accumulates delta per event and
// suppresses events once the decayed score exceeds threshold.
func NewScoreFilter(halfLife time.Duration, threshold, delta float64) *ScoreFilter {
	return &ScoreFilter{
		d:         New(halfLife),
		threshold: threshold,
		delta:     delta,
	}
}

// FilterEvents returns only those events whose port has not exceeded the score
// threshold. Each passing event increments that port's score.
func (f *ScoreFilter) FilterEvents(events []alert.Event) []alert.Event {
	if len(events) == 0 {
		return events
	}
	out := make([]alert.Event, 0, len(events))
	for _, e := range events {
		key := portKey(e)
		current := f.d.Score(key)
		if current >= f.threshold {
			continue
		}
		f.d.Add(key, f.delta)
		out = append(out, e)
	}
	return out
}

func portKey(e alert.Event) string {
	return fmt.Sprintf("%s/%d", e.Port.Protocol, e.Port.Port)
}
