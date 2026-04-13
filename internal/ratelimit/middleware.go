package ratelimit

import (
	"fmt"

	"github.com/iamcathal/portwatch/internal/alert"
	"github.com/iamcathal/portwatch/internal/scanner"
)

// eventKey returns a stable string key for a port event used as the limiter
// lookup key.
func eventKey(e alert.Event) string {
	return fmt.Sprintf("%s:%d:%s", e.Port.Protocol, e.Port.Number, e.Action)
}

// FilterEvents removes events that are suppressed by the rate limiter and
// returns only those that are permitted to be dispatched.
func FilterEvents(l *Limiter, events []alert.Event) []alert.Event {
	allowed := make([]alert.Event, 0, len(events))
	for _, e := range events {
		if l.Allow(eventKey(e)) {
			allowed = append(allowed, e)
		}
	}
	return allowed
}

// KeyForPort returns the rate-limit key used for a scanner.Port, combining
// protocol and port number. Useful for pre-seeding the limiter.
func KeyForPort(p scanner.Port, action string) string {
	return fmt.Sprintf("%s:%d:%s", p.Protocol, p.Number, action)
}
