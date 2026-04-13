// Package version exposes build-time version information for portwatch.
package version

import "fmt"

// These variables are set at build time via -ldflags.
var (
	// Version is the semantic version string, e.g. "1.2.3".
	Version = "dev"

	// Commit is the short Git SHA of the build.
	Commit = "unknown"

	// Date is the ISO-8601 build timestamp.
	Date = "unknown"
)

// Info holds all version metadata.
type Info struct {
	Version string `json:"version"`
	Commit  string `json:"commit"`
	Date    string `json:"date"`
}

// Get returns the current build Info.
func Get() Info {
	return Info{
		Version: Version,
		Commit:  Commit,
		Date:    Date,
	}
}

// String returns a human-readable one-line summary.
func (i Info) String() string {
	return fmt.Sprintf("portwatch %s (commit %s, built %s)", i.Version, i.Commit, i.Date)
}
