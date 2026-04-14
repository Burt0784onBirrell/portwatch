// Package tagger provides rule-based labelling of port events.
//
// A Tagger is constructed from a list of Rules, each associating a
// human-readable label with a port number and an optional protocol
// ("tcp" or "udp").  When Tag is called with a slice of alert.Events,
// every event whose port and protocol match a rule has the label
// appended to its process name in the form "<process> [<label>]".
//
// Rules are evaluated in order and only the first matching rule is
// applied to each event.  Events that match no rule are returned
// unchanged.  The original event slice is never mutated.
//
// Example usage:
//
//	rules := []tagger.Rule{
//		{Label: "http",  Port: 80,  Protocol: "tcp"},
//		{Label: "https", Port: 443, Protocol: "tcp"},
//		{Label: "dns",   Port: 53,  Protocol: ""},
//	}
//	tg, err := tagger.New(rules)
//	if err != nil { ... }
//	tagged := tg.Tag(events)
package tagger
