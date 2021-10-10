package gwp

import (
	"github.com/borud/gwp/pkg/gwpb"
)

// Common contains the definition of the functions that are common
// to both Listener and Connection
type Common interface {
	Requests() <-chan Request
	Close() error
}

// Listener listens for incoming connections, accepts these and receives
// data from them.
type Listener interface {
	Common
}

// Connection represent a connection to a peer.
type Connection interface {
	Common
	Send(Packet *gwpb.Packet) error
}

const (
	// MaxPacketSize is the maximum packet size we wish to emit. We also use this value to
	// warn of incoming packets that are oversize.
	// We have deliberately set this very low so that we feel the pain of this being too
	// low early.  At some later point we will turn this into a configuration variable
	// so that when we have a transport that can deal with large packets, or we have a stream
	// transport, we can make use of larger permissible message size.
	MaxPacketSize = 250
)
