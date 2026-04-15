// Package masking provides IP address masking for port events,
// allowing sensitive source addresses to be anonymised before
// they are forwarded to notifiers or written to audit logs.
package masking

import (
	"fmt"
	"net"

	"github.com/user/portwatch/internal/alert"
)

// Masker anonymises IP addresses contained in alert events.
type Masker struct {
	prefixLen int // bits to preserve (e.g. 24 keeps the /24 network)
}

// New returns a Masker that zeros all host bits beyond prefixLen.
// prefixLen is clamped to [0, 32].
func New(prefixLen int) *Masker {
	if prefixLen < 0 {
		prefixLen = 0
	}
	if prefixLen > 32 {
		prefixLen = 32
	}
	return &Masker{prefixLen: prefixLen}
}

// MaskIP applies the configured prefix mask to a raw IPv4 string.
// If the address cannot be parsed it is replaced with "<masked>".
func (m *Masker) MaskIP(raw string) string {
	ip := net.ParseIP(raw)
	if ip == nil {
		return "<masked>"
	}
	ip = ip.To4()
	if ip == nil {
		// IPv6 — return network prefix notation only
		return "<masked-ipv6>"
	}
	mask := net.CIDRMask(m.prefixLen, 32)
	masked := ip.Mask(mask)
	return fmt.Sprintf("%s/%d", masked.String(), m.prefixLen)
}

// ApplyToEvents returns a new slice of events with the Process field
// left intact but the source address (stored in Process when it begins
// with a digit) masked. In portwatch events the remote address is
// carried inside Port.Process for raw-socket listeners; this helper
// masks that field when it looks like an IP.
func (m *Masker) ApplyToEvents(events []alert.Event) []alert.Event {
	out := make([]alert.Event, len(events))
	for i, e := range events {
		copy := e
		if looksLikeIP(e.Port.Process) {
			copy.Port.Process = m.MaskIP(e.Port.Process)
		}
		out[i] = copy
	}
	return out
}

func looksLikeIP(s string) bool {
	if len(s) == 0 {
		return false
	}
	return s[0] >= '0' && s[0] <= '9'
}
