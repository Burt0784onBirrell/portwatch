package filter

import (
	"fmt"
	"strconv"
	"strings"
)

// RuleAction defines whether a rule allows or denies a port.
type RuleAction int

const (
	ActionAllow RuleAction = iota
	ActionDeny
)

// Rule represents a parsed filter rule with an action and port range.
type Rule struct {
	Action  RuleAction
	Low     uint16
	High    uint16
	Protocol string // "tcp", "udp", or "" for any
}

// ParseRule parses a rule string of the form:
//   [allow:|deny:]<port|port-range>[/protocol]
// Examples: "deny:8080", "allow:1000-2000/tcp", "deny:22/tcp"
func ParseRule(s string) (Rule, error) {
	s = strings.TrimSpace(s)
	if s == "" {
		return Rule{}, fmt.Errorf("empty rule")
	}

	action := ActionDeny
	body := s

	if strings.HasPrefix(s, "allow:") {
		action = ActionAllow
		body = s[len("allow:"):]
	} else if strings.HasPrefix(s, "deny:") {
		action = ActionDeny
		body = s[len("deny:"):]
	}

	proto := ""
	if idx := strings.LastIndex(body, "/"); idx != -1 {
		proto = strings.ToLower(body[idx+1:])
		if proto != "tcp" && proto != "udp" {
			return Rule{}, fmt.Errorf("unknown protocol %q", proto)
		}
		body = body[:idx]
	}

	var low, high uint16
	if idx := strings.Index(body, "-"); idx != -1 {
		l, err := parsePort(body[:idx])
		if err != nil {
			return Rule{}, err
		}
		h, err := parsePort(body[idx+1:])
		if err != nil {
			return Rule{}, err
		}
		if l > h {
			return Rule{}, fmt.Errorf("invalid range %d-%d", l, h)
		}
		low, high = l, h
	} else {
		p, err := parsePort(body)
		if err != nil {
			return Rule{}, err
		}
		low, high = p, p
	}

	return Rule{Action: action, Low: low, High: high, Protocol: proto}, nil
}

func parsePort(s string) (uint16, error) {
	n, err := strconv.ParseUint(strings.TrimSpace(s), 10, 16)
	if err != nil || n == 0 {
		return 0, fmt.Errorf("invalid port %q", s)
	}
	return uint16(n), nil
}
