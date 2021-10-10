package gwp

import (
	"net"
	"time"

	"github.com/borud/gwp/pkg/gwpb"
)

// Connection represents a Gateway Protocol Connection.
type Connection interface {
	Requests() <-chan Request
	Close() error
}

type ClientConnection interface {
	Send(Packet *gwpb.Packet) error
	Connection
}

// Request represents an incoming packet.
type Request struct {
	RemoteAddr net.Addr
	Packet     *gwpb.Packet
	Timestamp  time.Time
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
