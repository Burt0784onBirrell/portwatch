// Package cooldown implements a per-key quiet-period gate.
//
// A Tracker records the last time each key was allowed through and
// suppresses subsequent activations until the configured period has
// elapsed. This is useful when the same port-change event is seen
// repeatedly within a short window and downstream notifiers should
// only receive it once.
//
// Usage:
//
//	tracker := cooldown.New(30 * time.Second)
//	filtered := cooldown.FilterEvents(tracker, events)
//
// Keys are derived from the port number, protocol, and action so that
// an "opened" and a "closed" event for the same port are treated as
// independent activations.
package cooldown
