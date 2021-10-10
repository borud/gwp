package testutil

import (
	"net"
	"testing"

	"github.com/stretchr/testify/assert"
)

// FreeUDPPort returns a random free UDP port.
func FreeUDPPort(t *testing.T) int {
	lis, err := net.ListenUDP("udp", &net.UDPAddr{
		IP:   []byte{0, 0, 0, 0},
		Port: 0,
	})
	assert.Nil(t, err)
	addr, ok := lis.LocalAddr().(*net.UDPAddr)
	assert.True(t, ok)
	assert.Nil(t, lis.Close())
	return addr.Port
}
