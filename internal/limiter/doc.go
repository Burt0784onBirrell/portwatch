// Package limiter implements a sliding-window token-bucket rate limiter for
// portwatch alert events.
//
// # Overview
//
// A Limiter is created with a maximum burst count and a rolling time window.
// Each call to Allow consumes one token; once the bucket is full for the
// current window, Allow returns false until older timestamps fall outside the
// window.
//
// # Middleware
//
// FilterEvents wraps the limiter so it can be used as a pipeline stage that
// accepts a []alert.Event and returns the subset that passed the rate check.
//
// # Thread Safety
//
// All exported methods are safe for concurrent use.
package limiter
