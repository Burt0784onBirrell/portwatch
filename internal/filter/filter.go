package filter

import "github.com/user/portwatch/internal/scanner"

// Rule defines a single port filter rule.
type Rule struct {
	Port     uint16
	Protocol string // "tcp" or "udp", empty means both
}

// Filter decides which ports should be monitored based on a set of rules.
type Filter struct {
	allowList []Rule
	denyList  []Rule
}

// New returns a Filter with the provided allow and deny lists.
// If allowList is empty, all ports are allowed unless denied.
func New(allowList, denyList []Rule) *Filter {
	return &Filter{
		allowList: allowList,
		denyList:  denyList,
	}
}

// Apply returns only the ports that pass the filter rules.
func (f *Filter) Apply(ports []scanner.Port) []scanner.Port {
	var result []scanner.Port
	for _, p := range ports {
		if f.isDenied(p) {
			continue
		}
		if len(f.allowList) > 0 && !f.isAllowed(p) {
			continue
		}
		result = append(result, p)
	}
	return result
}

func (f *Filter) isAllowed(p scanner.Port) bool {
	for _, r := range f.allowList {
		if matches(r, p) {
			return true
		}
	}
	return false
}

func (f *Filter) isDenied(p scanner.Port) bool {
	for _, r := range f.denyList {
		if matches(r, p) {
			return true
		}
	}
	return false
}

func matches(r Rule, p scanner.Port) bool {
	if r.Port != 0 && r.Port != p.Number {
		return false
	}
	if r.Protocol != "" && r.Protocol != p.Protocol {
		return false
	}
	return true
}
