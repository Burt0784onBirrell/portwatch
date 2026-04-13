// Package pluginapi defines the extension interface that third-party
// notifier plugins must satisfy, along with a registry for managing them
// at runtime.
package pluginapi

import (
	"errors"
	"fmt"
	"sync"
)

// Notifier is the interface that every plugin must implement.
type Notifier interface {
	// Name returns a unique human-readable identifier for the plugin.
	Name() string
	// Notify delivers a batch of events to the plugin.
	Notify(events []Event) error
}

// Event is a minimal, plugin-facing representation of a port change.
type Event struct {
	Action  string // "opened" or "closed"
	Port    uint16
	Protocol string
	Process  string
}

// ErrDuplicatePlugin is returned when a plugin with the same name is
// registered more than once.
var ErrDuplicatePlugin = errors.New("pluginapi: duplicate plugin name")

// ErrPluginNotFound is returned when a requested plugin is not registered.
var ErrPluginNotFound = errors.New("pluginapi: plugin not found")

// Registry holds a set of registered notifier plugins.
type Registry struct {
	mu      sync.RWMutex
	plugins map[string]Notifier
}

// NewRegistry returns an empty, ready-to-use Registry.
func NewRegistry() *Registry {
	return &Registry{plugins: make(map[string]Notifier)}
}

// Register adds a plugin to the registry. It returns ErrDuplicatePlugin if a
// plugin with the same name has already been registered.
func (r *Registry) Register(p Notifier) error {
	r.mu.Lock()
	defer r.mu.Unlock()
	if _, exists := r.plugins[p.Name()]; exists {
		return fmt.Errorf("%w: %s", ErrDuplicatePlugin, p.Name())
	}
	r.plugins[p.Name()] = p
	return nil
}

// Get retrieves a plugin by name. Returns ErrPluginNotFound if absent.
func (r *Registry) Get(name string) (Notifier, error) {
	r.mu.RLock()
	defer r.mu.RUnlock()
	p, ok := r.plugins[name]
	if !ok {
		return nil, fmt.Errorf("%w: %s", ErrPluginNotFound, name)
	}
	return p, nil
}

// All returns a snapshot of all currently registered plugins.
func (r *Registry) All() []Notifier {
	r.mu.RLock()
	defer r.mu.RUnlock()
	out := make([]Notifier, 0, len(r.plugins))
	for _, p := range r.plugins {
		out = append(out, p)
	}
	return out
}
