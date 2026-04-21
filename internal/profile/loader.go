package profile

import (
	"fmt"
	"os"

	"gopkg.in/yaml.v3"
)

// fileProfiles is the on-disk representation of a profiles YAML file.
type fileProfiles struct {
	Profiles []Profile `yaml:"profiles"`
}

// LoadFile reads a YAML file at path and registers every profile it contains
// into reg. The file must contain a top-level "profiles" key whose value is a
// sequence of profile objects.
//
// Example file:
//
//	profiles:
//	  - name: dev
//	    filter_rules: []
//	    alert_cooldown_secs: 5
//	  - name: prod
//	    filter_rules: ["deny:tcp:22"]
//	    alert_cooldown_secs: 60
//	    tags:
//	      env: production
func LoadFile(path string, reg *Registry) error {
	data, err := os.ReadFile(path)
	if err != nil {
		return fmt.Errorf("profile: read file: %w", err)
	}

	var fp fileProfiles
	if err := yaml.Unmarshal(data, &fp); err != nil {
		return fmt.Errorf("profile: parse yaml: %w", err)
	}

	for _, p := range fp.Profiles {
		if err := reg.Register(p); err != nil {
			return fmt.Errorf("profile: register %q: %w", p.Name, err)
		}
	}
	return nil
}
