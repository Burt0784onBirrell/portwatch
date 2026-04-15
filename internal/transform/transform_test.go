package transform_test

import (
	"strings"
	"testing"

	"github.com/user/portwatch/internal/alert"
	"github.com/user/portwatch/internal/scanner"
	"github.com/user/portwatch/internal/transform"
)

func makeEvent(action, process string, port uint16) alert.Event {
	return alert.Event{
		Action:  action,
		Process: process,
		Port:    scanner.Port{Port: port, Protocol: "tcp"},
	}
}

func TestApply_NoFuncs_ReturnsOriginal(t *testing.T) {
	tr := transform.New()
	events := []alert.Event{makeEvent("opened", "nginx", 80)}
	out := tr.Apply(events)
	if len(out) != 1 || out[0].Process != "nginx" {
		t.Fatalf("expected unchanged events, got %+v", out)
	}
}

func TestApply_EmptyEvents_ReturnsEmpty(t *testing.T) {
	tr := transform.New(func(e alert.Event) alert.Event { return e })
	out := tr.Apply(nil)
	if out != nil {
		t.Fatalf("expected nil, got %v", out)
	}
}

func TestApply_TransformsProcess(t *testing.T) {
	upperFn := func(e alert.Event) alert.Event {
		e.Process = strings.ToUpper(e.Process)
		return e
	}
	tr := transform.New(upperFn)
	events := []alert.Event{makeEvent("opened", "nginx", 80)}
	out := tr.Apply(events)
	if out[0].Process != "NGINX" {
		t.Fatalf("expected NGINX, got %s", out[0].Process)
	}
}

func TestApply_FuncsAppliedInOrder(t *testing.T) {
	var order []string
	first := func(e alert.Event) alert.Event { order = append(order, "first"); return e }
	second := func(e alert.Event) alert.Event { order = append(order, "second"); return e }
	tr := transform.New(first, second)
	tr.Apply([]alert.Event{makeEvent("opened", "sshd", 22)})
	if len(order) != 2 || order[0] != "first" || order[1] != "second" {
		t.Fatalf("unexpected order: %v", order)
	}
}

func TestAdd_AppendsFn(t *testing.T) {
	tr := transform.New()
	tr.Add(func(e alert.Event) alert.Event { e.Process = "added"; return e })
	out := tr.Apply([]alert.Event{makeEvent("opened", "original", 443)})
	if out[0].Process != "added" {
		t.Fatalf("expected 'added', got %s", out[0].Process)
	}
}

func TestCompose_ChainsTransformers(t *testing.T) {
	a := transform.New(func(e alert.Event) alert.Event { e.Process = "A"; return e })
	b := transform.New(func(e alert.Event) alert.Event { e.Process += "B"; return e })
	c := a.Compose(b)
	out := c.Apply([]alert.Event{makeEvent("opened", "", 8080)})
	if out[0].Process != "AB" {
		t.Fatalf("expected AB, got %s", out[0].Process)
	}
	// originals must be unmodified
	orig := a.Apply([]alert.Event{makeEvent("opened", "", 8080)})
	if orig[0].Process != "A" {
		t.Fatalf("compose mutated original transformer")
	}
}
