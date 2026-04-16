// Package normalize provides event field normalization for consistent downstream processing.
package normalize

import (
	"strings"

	"github.com/joshbeard/portwatch/internal/alert"
)

// Option configures the Normalizer.
type Option func(*Normalizer)

// Normalizer applies normalization rules to alert events.
type Normalizer struct {
	lowercaseProcess bool
	trimProcess      bool
	defaultProtocol  string
}

// New creates a Normalizer with the given options.
func New(opts ...Option) *Normalizer {
	n := &Normalizer{
		lowercaseProcess: true,
		trimProcess:      true,
		defaultProtocol:  "tcp",
	}
	for _, o := range opts {
		o(n)
	}
	return n
}

// WithLowercaseProcess controls whether process names are lowercased.
func WithLowercaseProcess(v bool) Option {
	return func(n *Normalizer) { n.lowercaseProcess = v }
}

// WithDefaultProtocol sets the fallback protocol when a port has none.
func WithDefaultProtocol(p string) Option {
	return func(n *Normalizer) { n.defaultProtocol = p }
}

// Apply normalizes a slice of events, returning a new slice.
func (n *Normalizer) Apply(events []alert.Event) []alert.Event {
	out := make([]alert.Event, len(events))
	for i, e := range events {
		out[i] = n.normalize(e)
	}
	return out
}

func (n *Normalizer) normalize(e alert.Event) alert.Event {
	if n.trimProcess {
		e.Port.Process = strings.TrimSpace(e.Port.Process)
	}
	if n.lowercaseProcess {
		e.Port.Process = strings.ToLower(e.Port.Process)
	}
	if e.Port.Protocol == "" {
		e.Port.Protocol = n.defaultProtocol
	}
	return e
}
