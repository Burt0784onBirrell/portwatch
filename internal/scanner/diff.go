package scanner

// Diff holds the result of comparing two port snapshots.
type Diff struct {
	Opened []Port
	Closed []Port
}

// HasChanges returns true when at least one port opened or closed.
func (d Diff) HasChanges() bool {
	return len(d.Opened) > 0 || len(d.Closed) > 0
}

// Compare returns a Diff between a previous and current set of open ports.
// previous and current are maps produced by PortSetFromSlice.
func Compare(previous, current map[string]Port) Diff {
	var diff Diff

	// Ports present in current but not in previous → newly opened.
	for key, port := range current {
		if _, exists := previous[key]; !exists {
			diff.Opened = append(diff.Opened, port)
		}
	}

	// Ports present in previous but not in current → newly closed.
	for key, port := range previous {
		if _, exists := current[key]; !exists {
			diff.Closed = append(diff.Closed, port)
		}
	}

	return diff
}
