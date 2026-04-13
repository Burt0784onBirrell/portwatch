// Package ratelimit provides a cooldown-based rate limiter for portwatch alert
// events. It prevents alert storms caused by rapidly flapping ports by
// suppressing repeated notifications for the same port within a configurable
// time window.
//
// Usage:
//
//	limiter := ratelimit.New(30 * time.Second)
//	allowed := ratelimit.FilterEvents(limiter, events)
//	// dispatch only `allowed` events
package ratelimit
