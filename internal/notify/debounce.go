// Package notify provides debouncing middleware for alert events,
// suppressing repeated identical notifications within a configurable window.
package notify

import (
	"sync"
	"time"

	"github.com/user/portwatch/internal/alert"
)

// Debouncer suppresses duplicate events within a cooldown window.
type Debouncer struct {
	mu       sync.Mutex
	seen     map[string]time.Time
	window   time.Duration
	clock    func() time.Time
}

// NewDebouncer returns a Debouncer with the given suppression window.
func NewDebouncer(window time.Duration) *Debouncer {
	return newDebouncerWithClock(window, time.Now)
}

func newDebouncerWithClock(window time.Duration, clock func() time.Time) *Debouncer {
	return &Debouncer{
		seen:   make(map[string]time.Time),
		window: window,
		clock:  clock,
	}
}

// Filter returns only events that have not been seen within the debounce window.
func (d *Debouncer) Filter(events []alert.Event) []alert.Event {
	d.mu.Lock()
	defer d.mu.Unlock()

	now := d.clock()
	var out []alert.Event

	for _, e := range events {
		key := eventKey(e)
		if last, ok := d.seen[key]; ok && now.Sub(last) < d.window {
			continue
		}
		d.seen[key] = now
		out = append(out, e)
	}
	return out
}

// Reset clears the debounce state, allowing all events through again.
func (d *Debouncer) Reset() {
	d.mu.Lock()
	defer d.mu.Unlock()
	d.seen = make(map[string]time.Time)
}

func eventKey(e alert.Event) string {
	return e.Action + ":" + e.Port.Protocol + ":" + itoa(e.Port.Number)
}

func itoa(n int) string {
	if n == 0 {
		return "0"
	}
	buf := make([]byte, 0, 10)
	for n > 0 {
		buf = append([]byte{byte('0' + n%10)}, buf...)
		n /= 10
	}
	return string(buf)
}
