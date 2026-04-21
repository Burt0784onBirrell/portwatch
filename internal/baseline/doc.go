// Package baseline captures a reference "known-good" set of open ports and
// exposes helpers for detecting deviations from that reference state.
//
// Typical usage:
//
//	b := baseline.New()
//
//	// Capture the initial port state on startup.
//	b.Capture(initialPorts)
//
//	// On each scan cycle, check for deviations.
//	events := b.Deviation(currentPorts)
//
// The DeviationFilter pipeline stage can be inserted into an event pipeline
// to suppress events that match the baseline automatically.
package baseline
