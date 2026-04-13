// Package pipeline provides a composable event processing pipeline
// that chains multiple event transformation steps together.
package pipeline

import "github.com/user/portwatch/internal/alert"

// Stage is a function that transforms a slice of events.
type Stage func([]alert.Event) []alert.Event

// Pipeline applies a sequence of stages to a set of events in order.
type Pipeline struct {
	stages []Stage
}

// New creates a new Pipeline with the given stages.
func New(stages ...Stage) *Pipeline {
	return &Pipeline{stages: stages}
}

// Add appends one or more stages to the pipeline.
func (p *Pipeline) Add(stages ...Stage) {
	p.stages = append(p.stages, stages...)
}

// Run passes events through each stage in order, returning the final result.
// If any stage returns an empty slice, processing stops early.
func (p *Pipeline) Run(events []alert.Event) []alert.Event {
	result := events
	for _, stage := range p.stages {
		if len(result) == 0 {
			return result
		}
		result = stage(result)
	}
	return result
}

// Len returns the number of stages in the pipeline.
func (p *Pipeline) Len() int {
	return len(p.stages)
}
