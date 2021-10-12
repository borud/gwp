package main

import (
	"context"
	"encoding/json"
	"io"
	"log"
	"math/rand"
	"os"
	"time"

	"github.com/borud/gwp/pkg/gwpb"
	"github.com/jessevdk/go-flags"
	"google.golang.org/grpc"
	"google.golang.org/protobuf/proto"
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
		log.Fatalf("open stream error %v", err)
	}

	// Sender
	go func() {
		for {
			sender := rand.Intn(100000-1000) + 1000
			from := &gwpb.Address{Addr: &gwpb.Address_B32{B32: uint32(sender)}}

			packet := &gwpb.Packet{
				Payload: &gwpb.Packet_Samples{
					Samples: &gwpb.Samples{
						Samples: []*gwpb.Sample{
							{
								From:      from,
								Timestamp: uint64(time.Now().UnixMilli()),
								Type:      1,
								Value: &gwpb.Value{
									Value: &gwpb.Value_FloatVal{FloatVal: 3.14},
								},
							},
							{
								From:      from,
								Timestamp: uint64(time.Now().UnixMilli()),
								Type:      99,
								Value: &gwpb.Value{
									Value: &gwpb.Value_Int32Val{Int32Val: 314},
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

			time.Sleep(time.Millisecond * 2000)
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
		b, _ := proto.Marshal(packet)

		json, _ := json.MarshalIndent(packet, "", "  ")

		log.Printf("GW [%4d] : %s\n\t%s", len(b), packet.String(), json)
	}
}
