package server

import (
	"io"
	"log"
	"time"

	"github.com/borud/gwp/pkg/gwpb"
)

// Server is the server end of the Gateway Protocol
type Server struct{}

// Connect sets up a bidirectional packet stream
func (s Server) Connect(srv gwpb.Gateways_ConnectServer) error {
	ctx := srv.Context()

	log.Printf("connection %+v", srv)
	defer log.Printf("disconnect %+v", srv)

	// Jut periodically send stuff to the gateway
	go func() {
		for {
			payload := &gwpb.Packet_Config{
				Config: &gwpb.Config{
					Config: map[string]*gwpb.Value{
						"foo": {Value: &gwpb.Value_FloatVal{FloatVal: 23.5}},
						"bar": {Value: &gwpb.Value_Int32Val{Int32Val: 12345678}},
						"baz": {Value: &gwpb.Value_StringVal{StringVal: "wohoo"}},
					},
				},
			}

			err := srv.Send(&gwpb.Packet{
				Timestamp: 0,
				From:      &gwpb.Address{},
				To: &gwpb.Address{
					NodeId: 0,
					Addr: &gwpb.Address_B32{
						B32: 123,
					},
				},
				Payload: payload,
			})
			if err != nil {
				log.Printf("exiting sender: %v", err)
				return
			}

			time.Sleep(1000 * time.Millisecond)
		}
	}()

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
