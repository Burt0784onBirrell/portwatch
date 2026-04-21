// Package decay implements exponential score decay for port-level event scoring.
//
// Ports that generate a high volume of events accumulate a score that grows
// with each event. The score decays over time according to a configurable
// half-life, so a port that was briefly noisy will eventually be eligible for
// alerting again once activity subsides.
//
// # Usage
//
//	d := decay.New(5 * time.Minute)
//	score := d.Add("tcp/8080", 1.0) // returns new score after decay
//	current := d.Score("tcp/8080")  // read-only decayed value
//	d.Reset("tcp/8080")             // clear entry
//
// # Middleware
//
// ScoreFilter wraps the Decayer as an event pipeline stage:
//
//	f := decay.NewScoreFilter(5*time.Minute, 10.0, 1.0)
//	filtered := f.FilterEvents(events)
package decay
