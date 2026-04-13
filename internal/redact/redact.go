// Package redact provides utilities for scrubbing sensitive values
// (e.g. environment-variable names, process command lines) before they
// are emitted in alerts or audit logs.
package redact

import "strings"

// DefaultPatterns is the set of substrings that trigger redaction when
// found (case-insensitive) inside a process name or argument string.
var DefaultPatterns = []string{
	"password",
	"passwd",
	"secret",
	"token",
	"apikey",
	"api_key",
	"credential",
}

// Redactor scrubs sensitive substrings from arbitrary strings.
type Redactor struct {
	patterns []string
	placeholder string
}

// New returns a Redactor that replaces any value whose key matches one
// of the given patterns with placeholder.
func New(patterns []string, placeholder string) *Redactor {
	normalised := make([]string, len(patterns))
	for i, p := range patterns {
		normalised[i] = strings.ToLower(p)
	}
	return &Redactor{patterns: normalised, placeholder: placeholder}
}

// NewDefault returns a Redactor using DefaultPatterns and "[REDACTED]".
func NewDefault() *Redactor {
	return New(DefaultPatterns, "[REDACTED]")
}

// String returns s unchanged if it contains no sensitive pattern,
// otherwise it returns the placeholder.
func (r *Redactor) String(s string) string {
	lower := strings.ToLower(s)
	for _, p := range r.patterns {
		if strings.Contains(lower, p) {
			return r.placeholder
		}
	}
	return s
}

// ProcessName applies redaction to a process command-line string.
// Each whitespace-separated token is evaluated independently so that
// only the sensitive argument is replaced, not the entire command.
func (r *Redactor) ProcessName(cmd string) string {
	tokens := strings.Fields(cmd)
	for i, t := range tokens {
		tokens[i] = r.String(t)
	}
	return strings.Join(tokens, " ")
}
