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

// Send sends a packet to the peer we got the request from.
func (r *Request) Send(p *gwpb.Packet) error {
	return r.Peer.Send(p)
}

// Handler takes care of incoming requests
type Handler func(r Request)
