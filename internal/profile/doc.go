// Package profile provides named scanning profiles for portwatch.
//
// A Profile bundles together filter rules, alert cooldown settings, and
// arbitrary key/value tags. Operators can define multiple profiles in their
// configuration and activate one by name, making it straightforward to run
// portwatch with different behaviours in development vs. production
// environments without maintaining separate config files.
//
// Basic usage:
//
//	reg := profile.NewRegistry()
//	_ = reg.Register(profile.Profile{
//		Name:              "prod",
//		FilterRules:       []string{"deny:tcp:22"},
//		AlertCooldownSecs: 60,
//		Tags:              map[string]string{"env": "production"},
//	})
//	p, _ := reg.Get("prod")
package profile
