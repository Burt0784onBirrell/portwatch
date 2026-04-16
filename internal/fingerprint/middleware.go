package fingerprint

import (
	"github.com/user/portwatch/internal/scanner"
)

// Guard wraps a scanner.ScanFunc and skips invoking next if the fingerprint of
// the newly scanned PortSet matches the previous one. This avoids unnecessary
// diff and dispatch work when nothing has changed.
type Guard struct {
	last string
	scan scanner.ScanFunc
}

// NewGuard returns a Guard wrapping the provided ScanFunc.
func NewGuard(scan scanner.ScanFunc) *Guard {
	return &Guard{scan: scan}
}

// Scan executes the underlying ScanFunc. It returns the PortSet and a boolean
// indicating whether the result differs from the previous scan. On error the
// changed flag is always false.
func (g *Guard) Scan() (scanner.PortSet, bool, error) {
	ps, err := g.scan()
	if err != nil {
		return nil, false, err
	}
	current := Of(ps)
	if Equal(g.last, current) {
		return ps, false, nil
	}
	g.last = current
	return ps, true, nil
}
