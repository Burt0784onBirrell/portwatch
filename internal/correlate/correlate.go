// Package correlate groups related alert events into correlated bursts.
// When multiple ports change within a short window, they are assigned
// a shared correlation ID so downstream notifiers can group them.
package correlate

import (
	"crypto/rand"
	"encoding/hex"
	"sync"
	"time"

	"github.com/user/portwatch/internal/alert"
)

// Correlator assigns a shared ID to events that arrive within Window of
// each other. A new ID is minted whenever the burst window expires.
type Correlator struct {
	mu      sync.Mutex
	window  time.Duration
	current string
	lastAt  time.Time
	clock   func() time.Time
}

// New returns a Correlator with the given burst window.
func New(window time.Duration) *Correlator {
	return newWithClock(window, time.Now)
}

func newWithClock(window time.Duration, clock func() time.Time) *Correlator {
	return &Correlator{window: window, clock: clock}
}

// Annotate sets the "correlation_id" tag on every event in the slice.
// Events that arrive within Window of the previous call share the same ID.
func (c *Correlator) Annotate(events []alert.Event) []alert.Event {
	if len(events) == 0 {
		return events
	}

	c.mu.Lock()
	now := c.clock()
	if c.current == "" || now.Sub(c.lastAt) > c.window {
		c.current = newID()
	}
	c.lastAt = now
	id := c.current
	c.mu.Unlock()

	out := make([]alert.Event, len(events))
	for i, ev := range events {
		tags := make(map[string]string, len(ev.Tags)+1)
		for k, v := range ev.Tags {
			tags[k] = v
		}
		tags["correlation_id"] = id
		ev.Tags = tags
		out[i] = ev
	}
	return out
}

func newID() string {
	b := make([]byte, 6)
	_, _ = rand.Read(b)
	return hex.EncodeToString(b)
}
