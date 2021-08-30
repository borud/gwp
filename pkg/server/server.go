package server

import (
	"io"
	"log"

	"github.com/borud/gwp/pkg/gwpb"
)

// Server is the server end of the Gateway Protocol
type Server struct{}

// Connect sets up a bidirectional packet stream
func (s Server) Connect(srv gwpb.Gateways_ConnectServer) error {
	ctx := srv.Context()

	log.Printf("connection %+v", srv)
	defer log.Printf("disconnect %+v", srv)

	for {
		// Deal with possible cancelled context first
		select {
		case <-ctx.Done():
			return ctx.Err()
		default:
		}

		packet, err := srv.Recv()
		if err == io.EOF {
			// return will close stream from server side
			log.Println("exit")
			return nil
		}

		if err != nil {
			log.Printf("receive error %v", err)
			continue
		}

		log.Printf("Packet> %s", packet.String())
	}
}
