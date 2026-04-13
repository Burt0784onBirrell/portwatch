// Package pipeline provides a lightweight, composable event processing
// pipeline for portwatch.
//
// A Pipeline is constructed from a sequence of Stage functions, each of
// which receives a slice of alert.Event values and returns a (possibly
// filtered or mutated) slice.
//
// Stages are applied in the order they were added. Processing short-
// circuits if any stage returns an empty slice, avoiding unnecessary
// work downstream.
//
// Example usage:
//
//	 p := pipeline.New(
//	     ratelimit.FilterEvents,
//	     redact.ApplyToEvents,
//	     sampler.Filter,
//	 )
//	 out := p.Run(events)
package pipeline
