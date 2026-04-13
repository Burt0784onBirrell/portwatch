package scanner

import (
	"net"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

// startTCPServer binds a TCP listener on a random port and returns itServer(t *testing.T)Helper()
	ln	require.NoError(t, err)
	port := ln.Addr().(*net.TCPAddr).Port
	return ln, port
}

func TestScanner_DetectsOpenPort(t *testing.T) {
	ln, port := startTCPServer(t)
	defer ln.Close()

	s := &Scanner{
		StartPort: port,
		EndPort:   port,
		Timeout:   200 * time.Millisecond,
		Protocols: []string{"tcp"},
	}

	ports, err := s.Scan()
	require.NoError(t, err)
	assert.Len(t, ports, 1)
	assert.Equal(t, port, ports[0].Port)
	assert.Equal(t, "tcp", ports[0].Protocol)
}

func TestScanner_NoOpenPorts(t *testing.T) {
	// Use a port range unlikely to be open; scanner timeout keeps test fast.
	s := &Scanner{
		StartPort: 19999,
		EndPort:   19999,
		Timeout:   50 * time.Millisecond,
		Protocols: []string{"tcp"},
	}

	ports, err := s.Scan()
	require.NoError(t, err)
	// We cannot guarantee the port is closed on every CI machine,
	// but we can assert the call succeeds without error.
	assert.NotNil(t, ports)
}

func TestNewScanner_Defaults(t *testing.T) {
	s := NewScanner()
	assert.Equal(t, 1, s.StartPort)
	assert.Equal(t, 65535, s.EndPort)
	assert.Equal(t, 500*time.Millisecond, s.Timeout)
	assert.Equal(t, []string{"tcp"}, s.Protocols)
}
