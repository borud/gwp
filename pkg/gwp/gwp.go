package gwp

import (
	"net"
	"time"

	"github.com/borud/gwp/pkg/gwpb"
)

// Request represents an incoming packet.
type Request struct {
	Peer       Connection
	RemoteAddr net.Addr
	Packet     *gwpb.Packet
	Timestamp  time.Time
}

// Handler takes care of incoming requests
type Handler func(r Request)
