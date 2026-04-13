package pipeline_test

import (
	"testing"

	"github.com/user/portwatch/internal/alert"
	"github.com/user/portwatch/internal/pipeline"
	"github.com/user/portwatch/internal/scanner"
)

func makeEvent(port uint16, action alert.Action) alert.Event {
	return alert.Event{
		Port:   scanner.Port{Port: port, Proto: "tcp"},
		Action: action,
	}
}

func TestPipeline_EmptyStages_ReturnsInput(t *testing.T) {
	p := pipeline.New()
	events := []alert.Event{makeEvent(80, alert.Opened)}
	out := p.Run(events)
	if len(out) != 1 {
		t.Fatalf("expected 1 event, got %d", len(out))
	}
}

func TestPipeline_StagesAppliedInOrder(t *testing.T) {
	var order []int
	s1 := func(evs []alert.Event) []alert.Event { order = append(order, 1); return evs }
	s2 := func(evs []alert.Event) []alert.Event { order = append(order, 2); return evs }
	s3 := func(evs []alert.Event) []alert.Event { order = append(order, 3); return evs }

	p := pipeline.New(s1, s2, s3)
	p.Run([]alert.Event{makeEvent(443, alert.Opened)})

	if len(order) != 3 || order[0] != 1 || order[1] != 2 || order[2] != 3 {
		t.Fatalf("unexpected stage order: %v", order)
	}
}

func TestPipeline_ShortCircuitsOnEmptySlice(t *testing.T) {
	called := false
	drop := func(_ []alert.Event) []alert.Event { return nil }
	guard := func(evs []alert.Event) []alert.Event { called = true; return evs }

	p := pipeline.New(drop, guard)
	p.Run([]alert.Event{makeEvent(22, alert.Closed)})

	if called {
		t.Fatal("expected second stage to be skipped after empty result")
	}
}

func TestPipeline_Add_AppendStage(t *testing.T) {
	p := pipeline.New()
	if p.Len() != 0 {
		t.Fatalf("expected 0 stages, got %d", p.Len())
	}
	p.Add(func(evs []alert.Event) []alert.Event { return evs })
	if p.Len() != 1 {
		t.Fatalf("expected 1 stage, got %d", p.Len())
	}
}

func TestPipeline_FilterStage_ReducesEvents(t *testing.T) {
	filterClosed := func(evs []alert.Event) []alert.Event {
		var out []alert.Event
		for _, e := range evs {
			if e.Action == alert.Opened {
				out = append(out, e)
			}
		}
		return out
	}

	p := pipeline.New(filterClosed)
	input := []alert.Event{
		makeEvent(80, alert.Opened),
		makeEvent(443, alert.Closed),
		makeEvent(8080, alert.Opened),
	}
	out := p.Run(input)
	if len(out) != 2 {
		t.Fatalf("expected 2 events after filter, got %d", len(out))
	}
}
