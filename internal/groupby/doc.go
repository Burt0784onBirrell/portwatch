// Package groupby partitions alert.Event slices into named buckets using a
// caller-supplied key function.
//
// Usage:
//
//	g, err := groupby.New(groupby.ByAction)
//	if err != nil { ... }
//	groups := g.Apply(events)
//	for _, grp := range groups {
//	    fmt.Println(grp.Key, len(grp.Events))
//	}
//
// Built-in key functions:
//
//	ByAction   — groups by "opened" / "closed"
//	ByProtocol — groups by "tcp" / "udp"
//
// Custom key functions can be provided to group by any event attribute.
package groupby
