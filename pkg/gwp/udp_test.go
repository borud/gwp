package gwp

import (
	"fmt"
	"testing"
	"time"

	"github.com/borud/gwp/pkg/gwpb"
	"github.com/borud/gwp/pkg/testutil"
	"github.com/stretchr/testify/assert"
)

func TestUDP(t *testing.T) {
	addrStr := fmt.Sprintf(":%d", testutil.FreeUDPPort(t))

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
	conn, err := NewUDPConnection(addrStr, 1)
	assert.Nil(t, err)
	assert.NotNil(t, conn)

	go func() {
		for i := 0; i < numMessages; i++ {
			err := conn.Send(&gwpb.Packet{
				To: &gwpb.Address{
					Addr: &gwpb.Address_Name{
						Name: "some name",
					},
				},
				Payload: &gwpb.Packet_Data{
					Data: &gwpb.Data{
						Type: 1,
						Data: []byte{1, 2, 3, 4},
					},
				},
			})
			assert.Nil(t, err)
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
