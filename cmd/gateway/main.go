package main

import (
	"context"
	"io"
	"log"
	"math/rand"
	"os"
	"time"

	"github.com/borud/gwp/pkg/gwpb"
	"github.com/jessevdk/go-flags"
	"google.golang.org/grpc"
)

var opt struct {
	GRPCAddr string `long:"grpc-addr" default:":5011" description:"gRPC listen port"`
}

func init() {
	parser := flags.NewParser(&opt, flags.Default)
	if _, err := parser.Parse(); err != nil {
		switch flagsErr := err.(type) {
		case flags.ErrorType:
			if flagsErr == flags.ErrHelp {
				os.Exit(0)
			}
			os.Exit(1)
		default:
			os.Exit(1)
		}
	}
}

func main() {
	conn, err := grpc.Dial(opt.GRPCAddr, grpc.WithInsecure())
	if err != nil {
		log.Printf("error connecting to %s", opt.GRPCAddr)
	}

	client := gwpb.NewGatewaysClient(conn)

	stream, err := client.Connect(context.Background())
	if err != nil {
		log.Fatalf("openn stream error %v", err)
	}

	// Sender
	go func() {
		for {
			sender := rand.Intn(100-1) + 1
			from := &gwpb.Address{NodeId: 1, Addr: &gwpb.Address_B32{B32: uint32(sender)}}

			packet := &gwpb.Packet{
				Timestamp: uint64(time.Now().UnixMilli()),
				Payload: &gwpb.Packet_Samples{
					Samples: &gwpb.Samples{
						Samples: []*gwpb.Sample{
							{
								From:       from,
								NodeId:     uint64(sender),
								Timestamp:  uint64(time.Now().UnixMilli()),
								SensorType: 1,
								Value: &gwpb.Value{
									Value: &gwpb.Value_FloatVal{
										FloatVal: 3.14,
									},
								},
							},
						},
					},
				},
			}

			err := stream.Send(packet)
			if err != nil {
				log.Fatalf("can't send: %v", err)
			}

			time.Sleep(time.Millisecond * 200)
		}
	}()

	// Print incoming packets
	for {
		packet, err := stream.Recv()
		if err == io.EOF {
			log.Fatalf("got EOF")
		}
		if err != nil {
			log.Fatalf("can not receive %v", err)
		}
		log.Printf("packet: %s", packet.String())
	}
}
