// Package quota provides a per-key event quota enforcer backed by a rolling
// time window.
//
// Use New to create a Quota with a maximum event count and window duration,
// then call Allow with an arbitrary string key to check whether the next event
// is within budget.
//
// The FilterEvents middleware helper integrates the quota with the standard
// []alert.Event pipeline used throughout portwatch:
//
//	events = quota.FilterEvents(q, events)
//
// Quota is safe for concurrent use.
package quota
