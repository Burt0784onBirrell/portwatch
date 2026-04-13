package pluginapi_test

import (
	"errors"
	"testing"

	"github.com/user/portwatch/internal/pluginapi"
)

// stub is a minimal Notifier used in tests.
type stub struct {
	name   string
	called int
	fail   bool
}

func (s *stub) Name() string { return s.name }
func (s *stub) Notify(events []pluginapi.Event) error {
	s.called++
	if s.fail {
		return errors.New("stub error")
	}
	return nil
}

func TestRegistry_RegisterAndGet(t *testing.T) {
	reg := pluginapi.NewRegistry()
	p := &stub{name: "my-plugin"}

	if err := reg.Register(p); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	got, err := reg.Get("my-plugin")
	if err != nil {
		t.Fatalf("Get failed: %v", err)
	}
	if got.Name() != "my-plugin" {
		t.Errorf("expected my-plugin, got %s", got.Name())
	}
}

func TestRegistry_DuplicateReturnsError(t *testing.T) {
	reg := pluginapi.NewRegistry()
	p := &stub{name: "dup"}

	_ = reg.Register(p)
	err := reg.Register(p)

	if !errors.Is(err, pluginapi.ErrDuplicatePlugin) {
		t.Errorf("expected ErrDuplicatePlugin, got %v", err)
	}
}

func TestRegistry_GetMissingReturnsError(t *testing.T) {
	reg := pluginapi.NewRegistry()

	_, err := reg.Get("ghost")

	if !errors.Is(err, pluginapi.ErrPluginNotFound) {
		t.Errorf("expected ErrPluginNotFound, got %v", err)
	}
}

func TestRegistry_All_ReturnsAllPlugins(t *testing.T) {
	reg := pluginapi.NewRegistry()
	_ = reg.Register(&stub{name: "alpha"})
	_ = reg.Register(&stub{name: "beta"})

	all := reg.All()
	if len(all) != 2 {
		t.Errorf("expected 2 plugins, got %d", len(all))
	}
}

func TestRegistry_All_EmptyRegistry(t *testing.T) {
	reg := pluginapi.NewRegistry()

	if got := reg.All(); len(got) != 0 {
		t.Errorf("expected empty slice, got %d elements", len(got))
	}
}
