// Package profile provides named scanning profiles that bundle together
// filter rules, alert thresholds, and output preferences so operators can
// switch between configurations (e.g. "dev", "prod") without editing the
// main config file.
package profile

import (
	"errors"
	"fmt"
)

// Profile holds a named collection of settings that influence how portwatch
// behaves during a scan cycle.
type Profile struct {
	// Name is the unique identifier for this profile.
	Name string

	// FilterRules is a slice of rule strings in the same format accepted by
	// internal/filter (e.g. "deny:tcp:22", "allow:udp:53").
	FilterRules []string

	// AlertCooldownSecs is the minimum number of seconds between repeated
	// alerts for the same port event.
	AlertCooldownSecs int

	// Tags is an optional set of key/value labels attached to every event
	// emitted while this profile is active.
	Tags map[string]string
}

// Registry stores named profiles and allows lookup by name.
type Registry struct {
	profiles map[string]Profile
}

// NewRegistry returns an empty Registry.
func NewRegistry() *Registry {
	return &Registry{profiles: make(map[string]Profile)}
}

// Register adds p to the registry. It returns an error if a profile with the
// same name already exists or if the profile name is empty.
func (r *Registry) Register(p Profile) error {
	if p.Name == "" {
		return errors.New("profile name must not be empty")
	}
	if _, exists := r.profiles[p.Name]; exists {
		return fmt.Errorf("profile %q already registered", p.Name)
	}
	r.profiles[p.Name] = p
	return nil
}

// Get returns the profile with the given name. It returns an error if no such
// profile has been registered.
func (r *Registry) Get(name string) (Profile, error) {
	p, ok := r.profiles[name]
	if !ok {
		return Profile{}, fmt.Errorf("profile %q not found", name)
	}
	return p, nil
}

// Names returns the sorted list of registered profile names.
func (r *Registry) Names() []string {
	names := make([]string, 0, len(r.profiles))
	for n := range r.profiles {
		names = append(names, n)
	}
	return names
}
