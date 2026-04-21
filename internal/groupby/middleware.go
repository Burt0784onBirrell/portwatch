package groupby

import "github.com/user/portwatch/internal/alert"

// FlattenGroups converts a slice of Groups back into a flat event slice,
// preserving per-group order and concatenating groups in their original order.
// This is useful when a downstream stage expects a plain []alert.Event.
func FlattenGroups(groups []Group) []alert.Event {
	var out []alert.Event
	for _, g := range groups {
		out = append(out, g.Events...)
	}
	return out
}

// GroupFilter applies fn to each Group and keeps only the events for which fn
// returns true. Groups that become empty are omitted from the result.
func GroupFilter(groups []Group, fn func(Group, alert.Event) bool) []Group {
	var result []Group
	for _, grp := range groups {
		var kept []alert.Event
		for _, ev := range grp.Events {
			if fn(grp, ev) {
				kept = append(kept, ev)
			}
		}
		if len(kept) > 0 {
			result = append(result, Group{Key: grp.Key, Events: kept})
		}
	}
	return result
}
