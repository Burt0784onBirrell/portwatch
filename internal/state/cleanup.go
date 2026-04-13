package state

import (
	"os"
	"time"
)

// CleanupOptions configures state file retention behaviour.
type CleanupOptions struct {
	// MaxAge is the maximum age of a state file before it is considered stale.
	MaxAge time.Duration
}

// DefaultCleanupOptions returns sensible defaults for state cleanup.
func DefaultCleanupOptions() CleanupOptions {
	return CleanupOptions{
		MaxAge: 24 * time.Hour,
	}
}

// IsStale reports whether the state file at path is older than opts.MaxAge.
// If the file does not exist, IsStale returns false without error.
func IsStale(path string, opts CleanupOptions) (bool, error) {
	info, err := os.Stat(path)
	if err != nil {
		if os.IsNotExist(err) {
			return false, nil
		}
		return false, err
	}
	return time.Since(info.ModTime()) > opts.MaxAge, nil
}

// RemoveIfStale deletes the state file at path when it is stale according to
// opts. It returns true when the file was removed, false when it was kept or
// absent.
func RemoveIfStale(path string, opts CleanupOptions) (bool, error) {
	stale, err := IsStale(path, opts)
	if err != nil || !stale {
		return false, err
	}
	if err := os.Remove(path); err != nil && !os.IsNotExist(err) {
		return false, err
	}
	return true, nil
}
