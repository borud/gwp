package gwp

import (
	"sync"
	"testing"

	"github.com/borud/gwp/pkg/gwpb"
	"github.com/borud/gwp/pkg/testutil"
	"github.com/stretchr/testify/assert"
)

func TestServerSmoke(t *testing.T) {
	udpListener, err := NewUDPListener(testutil.LocalHostFreeUDPPort(t), 1)
	assert.Nil(t, err)
	assert.NotNil(t, udpListener)

	// test with no listeners
	assert.ErrorIs(t, (&Server{}).Start(), ErrNoListenersDefined)

	// test with listener, but no handler
	assert.ErrorIs(t, (&Server{Listeners: []Listener{udpListener}}).Start(), ErrNoHandlerDefined)

	// test with Listeners and Handler
	server := &Server{
		Listeners: []Listener{udpListener},
		Handler: func(r Request) error {
			// do nothing
			return nil
		},
	}

	err = server.Start()
	assert.Nil(t, err)

	assert.ErrorIs(t, server.Start(), ErrServerAlreadyStarted)

	assert.Nil(t, server.Shutdown())
	assert.ErrorIs(t, server.Shutdown(), ErrAlreadyShutDown)
}

func TestServerAndClient(t *testing.T) {
	addr := testutil.LocalHostFreeUDPPort(t)

	udpListener, err := NewUDPListener(addr, 1)
	assert.Nil(t, err)
	assert.NotNil(t, udpListener)

	numMessages := 5

	var messageCountWG sync.WaitGroup
	var shutdownCallbacksWG sync.WaitGroup

	messageCountWG.Add(numMessages)

	server := &Server{
		Listeners: []Listener{udpListener},
		Handler: func(r Request) error {
			messageCountWG.Done()
			return nil
		},
		ShutdownCallbacks: []func(){
			func() {
				shutdownCallbacksWG.Done()
			},
			func() {
				shutdownCallbacksWG.Done()
			},
		},
	}

	shutdownCallbacksWG.Add(len(server.ShutdownCallbacks))

	go func() {
		assert.Nil(t, server.Start())
	}()

	// connect a client
	conn, err := NewUDPConnection(addr, 1)
	assert.Nil(t, err)
	assert.NotNil(t, conn)

	for i := 0; i < numMessages; i++ {
		conn.Send(&gwpb.Packet{
			From: &gwpb.Address{
				Addr: &gwpb.Address_B32{
					B32: 1234,
				},
			},
			Payload: &gwpb.Packet_Data{
				Data: &gwpb.Data{
					Type: 1,
					Data: []byte{3, 4, 5},
				},
			},
		})
	}

	// will hang here if we don't get all the messages
	messageCountWG.Wait()

	assert.Nil(t, conn.Close())
	assert.Nil(t, server.Shutdown())

	// will hang here if all the shutdown callbacks aren't called
	shutdownCallbacksWG.Wait()
}
