package scanner

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func makePort(proto, addr string, port int) Port {
	return Port{Protocol: proto, Address: addr, Port: port}
}

func TestCompare_NoChanges(t *testing.T) {
	ports := []Port{
		makePort("tcp", "127.0.0.1", 80),
		makePort("tcp", "127.0.0.1", 443),
	}
	prev := PortSetFromSlice(ports)
	curr := PortSetFromSlice(ports)

	diff := Compare(prev, curr)
	assert.False(t, diff.HasChanges())
	assert.Empty(t, diff.Opened)
	assert.Empty(t, diff.Closed)
}

func TestCompare_PortOpened(t *testing.T) {
	prev := PortSetFromSlice([]Port{makePort("tcp", "127.0.0.1", 80)})
	curr := PortSetFromSlice([]Port{
		makePort("tcp", "127.0.0.1", 80),
		makePort("tcp", "127.0.0.1", 8080),
	})

	diff := Compare(prev, curr)
	assert.True(t, diff.HasChanges())
	assert.Len(t, diff.Opened, 1)
	assert.Equal(t, 8080, diff.Opened[0].Port)
	assert.Empty(t, diff.Closed)
}

func TestCompare_PortClosed(t *testing.T) {
	prev := PortSetFromSlice([]Port{
		makePort("tcp", "127.0.0.1", 80),
		makePort("tcp", "127.0.0.1", 443),
	})
	curr := PortSetFromSlice([]Port{makePort("tcp", "127.0.0.1", 80)})

	diff := Compare(prev, curr)
	assert.True(t, diff.HasChanges())
	assert.Empty(t, diff.Opened)
	assert.Len(t, diff.Closed, 1)
	assert.Equal(t, 443, diff.Closed[0].Port)
}

func TestCompare_BothOpenedAndClosed(t *testing.T) {
	prev := PortSetFromSlice([]Port{makePort("tcp", "127.0.0.1", 22)})
	curr := PortSetFromSlice([]Port{makePort("tcp", "127.0.0.1", 2222)})

	diff := Compare(prev, curr)
	assert.True(t, diff.HasChanges())
	assert.Len(t, diff.Opened, 1)
	assert.Len(t, diff.Closed, 1)
}

func TestPortSetFromSlice_Empty(t *testing.T) {
	set := PortSetFromSlice(nil)
	assert.NotNil(t, set)
	assert.Empty(t, set)
}

func TestPort_String(t *testing.T) {
	p := makePort("tcp", "127.0.0.1", 443)
	assert.Equal(t, "tcp://127.0.0.1:443", p.String())
}
