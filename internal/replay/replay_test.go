package replay_test

import (
	"context"
	"errors"
	"testing"
	"time"

	"github.com/user/portwatch/internal/alert"
	"github.com/user/portwatch/internal/replay"
	"github.com/user/portwatch/internal/scanner"
)

// stubSource implements replay.Source.
type stubSource struct {
	ports scanner.PortSet
	err   error
}

func (s *stubSource) Load() (scanner.PortSet, error) { return s.ports, s.err }

// stubSink implements replay.Sink.
type stubSink struct {
	events []alert.Event
	err    error
}

func (s *stubSink) Dispatch(_ context.Context, events []alert.Event) error {
	s.events = append(s.events, events...)
	return s.err
}

func makePortSet(ports ...scanner.Port) scanner.PortSet {
	ps := make(scanner.PortSet)
	for _, p := range ports {
		ps[p] = struct{}{}
	}
	return ps
}

func TestReplay_EmptyStateIsNoop(t *testing.T) {
	sink := &stubSink{}
	r := replay.New(&stubSource{ports: scanner.PortSet{}}, sink, replay.DefaultOptions())
	if err := r.Run(context.Background()); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(sink.events) != 0 {
		t.Errorf("expected no events, got %d", len(sink.events))
	}
}

func TestReplay_DispatchesOpenedEventPerPort(t *testing.T) {
	p1 := scanner.Port{Port: 80, Protocol: "tcp"}
	p2 := scanner.Port{Port: 443, Protocol: "tcp"}
	src := &stubSource{ports: makePortSet(p1, p2)}
	sink := &stubSink{}

	ts := time.Date(2024, 1, 1, 0, 0, 0, 0, time.UTC)
	opts := replay.Options{AsOf: ts}
	r := replay.New(src, sink, opts)

	if err := r.Run(context.Background()); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(sink.events) != 2 {
		t.Fatalf("expected 2 events, got %d", len(sink.events))
	}
	for _, ev := range sink.events {
		if ev.Action != alert.ActionOpened {
			t.Errorf("expected ActionOpened, got %v", ev.Action)
		}
		if !ev.Timestamp.Equal(ts) {
			t.Errorf("expected timestamp %v, got %v", ts, ev.Timestamp)
		}
	}
}

func TestReplay_SourceErrorPropagates(t *testing.T) {
	srcErr := errors.New("disk read failure")
	r := replay.New(&stubSource{err: srcErr}, &stubSink{}, replay.DefaultOptions())
	err := r.Run(context.Background())
	if err == nil {
		t.Fatal("expected error, got nil")
	}
	if !errors.Is(err, srcErr) {
		t.Errorf("expected wrapped srcErr, got %v", err)
	}
}

func TestReplay_SinkErrorPropagates(t *testing.T) {
	p := scanner.Port{Port: 22, Protocol: "tcp"}
	sinkErr := errors.New("network timeout")
	sink := &stubSink{err: sinkErr}
	r := replay.New(&stubSource{ports: makePortSet(p)}, sink, replay.DefaultOptions())
	err := r.Run(context.Background())
	if err == nil {
		t.Fatal("expected error, got nil")
	}
	if !errors.Is(err, sinkErr) {
		t.Errorf("expected wrapped sinkErr, got %v", err)
	}
}
