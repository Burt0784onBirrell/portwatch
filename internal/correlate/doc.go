// Package correlate provides burst-window correlation for port-change events.
//
// When portwatch detects simultaneous (or near-simultaneous) changes across
// multiple ports — for example during a service restart — the individual
// events are tagged with a shared correlation_id so that downstream
// notifiers, dashboards, or log aggregators can group them into a single
// incident rather than treating each port change as an isolated alert.
//
// Usage:
//
//	corr := correlate.New(2 * time.Second)
//	annotated := corr.Annotate(events)
//
// All events returned within the same 2-second burst window will carry the
// same correlation_id tag value. A new ID is generated automatically once
// the window expires and the next batch of events arrives.
package correlate
