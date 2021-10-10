package transport

import (
	"fmt"
	"net"
	"testing"
	"time"

	"github.com/borud/gwp/pkg/gwpb"
	"github.com/stretchr/testify/assert"
)

func TestUDP(t *testing.T) {
	lis, err := net.ListenUDP("udp", &net.UDPAddr{
		IP:   []byte{0, 0, 0, 0},
		Port: 0,
	})
	assert.Nil(t, err)
	addr, ok := lis.LocalAddr().(*net.UDPAddr)
	assert.True(t, ok)
	assert.Nil(t, lis.Close())

	addrStr := fmt.Sprintf(":%d", addr.Port)

	// Set up listener
	listener, err := NewUDPListener(addrStr, 1)
	assert.Nil(t, err)

	// Set up reception and counting of messages
	numMessages := 20
	waitCh := make(chan struct{})
	go func() {
		count := 0
		for range listener.Requests() {
			count++
			if count == numMessages {
				close(waitCh)
			}
		}
	}()

	// fire up client and send messages
	conn, err := NewUDPConnection(addrStr, 10)
	assert.Nil(t, err)

	go func() {
		for i := 0; i < numMessages; i++ {
			conn.Send(&gwpb.Packet{
				To: &gwpb.Address{
					Addr: &gwpb.Address_Name{
						Name: "some name",
					},
				},
				Id:         uint32(i + 1),
				RequireAck: false,
				Payload: &gwpb.Packet_Data{
					Data: &gwpb.Data{
						Type: 1,
						Id:   1,
						Data: []byte{1, 2, 3, 4},
					},
				},
			})
		}
		assert.Nil(t, conn.Close())
	}()

	// Wait for completion of reception or timeout, whichever comes first
	select {
	case <-waitCh:
	case <-time.After(100 * time.Millisecond):
		assert.Fail(t, "test timed out")
	}

	assert.Nil(t, listener.Close())
}
