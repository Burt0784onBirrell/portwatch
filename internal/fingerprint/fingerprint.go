// Package fingerprint produces a stable string identity for a port snapshot,
// allowing callers to detect whether the observed port set has changed between
// two scan cycles without performing a full diff.
package fingerprint

import (
	"crypto/sha256"
	"encoding/hex"
	"fmt"
	"sort"
	"strings"

	"github.com/user/portwatch/internal/scanner"
)

// Of returns a short hex fingerprint derived from the sorted string
// representations of every port in ps. An empty PortSet always returns the
// same sentinel value so callers can treat it as a valid fingerprint.
func Of(ps scanner.PortSet) string {
	if len(ps) == 0 {
		return "0000000000000000"
	}

	keys := make([]string, 0, len(ps))
	for p := range ps {
		keys = append(keys, fmt.Sprintf("%s/%d", p.Protocol, p.Port))
	}
	sort.Strings(keys)

	h := sha256.Sum256([]byte(strings.Join(keys, "|")))
	return hex.EncodeToString(h[:])[:16]
}

// Equal reports whether two fingerprints are identical.
func Equal(a, b string) bool {
	return a == b
}
