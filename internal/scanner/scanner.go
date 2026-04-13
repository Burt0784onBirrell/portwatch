package scanner

import (
	"fmt"
	"net"
	"strconv"
	"strings"
	"time"
)

// Port represents a single open port entry.
type Port struct {
	Protocol string
	Address  string
	Port     int
}

func (p Port) String() string {
	return fmt.Sprintf("%s://%s:%d", p.Protocol, p.Address, p.Port)
}

// Scanner scans for open ports on the local machine.
type Scanner struct {
	StartPort int
	EndPort   int
	Timeout   time.Duration
	Protocols []string
}

// NewScanner creates a Scanner with sensible defaults.
func NewScanner() *Scanner {
	return &Scanner{
		StartPort: 1,
		EndPort:   65535,
		Timeout:   500 * time.Millisecond,
		Protocols: []string{"tcp"},
	}
}

// Scan returns all open ports in the configured range.
func (s *Scanner) Scan() ([]Port, error) {
	var open []Port
	for _, proto := range s.Protocols {
		for port := s.StartPort; port <= s.EndPort; port++ {
			addr := net.JoinHostPort("127.0.0.1", strconv.Itoa(port))
			conn, err := net.DialTimeout(proto, addr, s.Timeout)
			if err != nil {
				continue
			}
			conn.Close()
			open = append(open, Port{
				Protocol: strings.ToLower(proto),
				Address:  "127.0.0.1",
				Port:     port,
			})
		}
	}
	return open, nil
}

// PortSetFromSlice converts a slice of Ports into a map keyed by Port.String()
// for O(1) lookup.
func PortSetFromSlice(ports []Port) map[string]Port {
	set := make(map[string]Port, len(ports))
	for _, p := range ports {
		set[p.String()] = p
	}
	return set
}
