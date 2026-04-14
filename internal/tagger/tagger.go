// Package tagger assigns user-defined labels to port events based on
// configurable matching rules, making it easier to classify traffic in
// downstream notifiers and audit logs.
package tagger

import (
	"fmt"
	"strings"

	"github.com/user/portwatch/internal/alert"
)

// Rule pairs a label with the port/protocol pattern it matches.
type Rule struct {
	Label    string
	Port     uint16
	Protocol string // "tcp", "udp", or "" for any
}

// Tagger holds a set of labelling rules.
type Tagger struct {
	rules []Rule
}

// New creates a Tagger from a slice of rules. Rules are applied in order;
// the first match wins.
func New(rules []Rule) (*Tagger, error) {
	for i, r := range rules {
		if strings.TrimSpace(r.Label) == "" {
			return nil, fmt.Errorf("tagger: rule %d has an empty label", i)
		}
		proto := strings.ToLower(r.Protocol)
		if proto != "" && proto != "tcp" && proto != "udp" {
			return nil, fmt.Errorf("tagger: rule %d has invalid protocol %q", i, r.Protocol)
		}
		rules[i].Protocol = proto
	}
	return &Tagger{rules: rules}, nil
}

// Tag returns a copy of events with a label appended to the process name
// for any event whose port and protocol match a rule.  Events that match
// no rule are returned unchanged.
func (t *Tagger) Tag(events []alert.Event) []alert.Event {
	out := make([]alert.Event, len(events))
	for i, ev := range events {
		out[i] = ev
		for _, r := range t.rules {
			if r.Port != ev.Port.Number {
				continue
			}
			if r.Protocol != "" && r.Protocol != strings.ToLower(ev.Port.Protocol) {
				continue
			}
			suffix := "[" + r.Label + "]"
			if !strings.Contains(ev.Port.Process, suffix) {
				cp := out[i].Port
				if cp.Process == "" {
					cp.Process = suffix
				} else {
					cp.Process = cp.Process + " " + suffix
				}
				out[i].Port = cp
			}
			break
		}
	}
	return out
}
