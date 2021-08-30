package main

import (
	"log"
	"net"
	"os"

	"github.com/jessevdk/go-flags"
	"google.golang.org/grpc"

	"github.com/borud/gwp/pkg/gwpb"
	"github.com/borud/gwp/pkg/server"
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
	listen, err := net.Listen("tcp", opt.GRPCAddr)
	if err != nil {
		log.Fatalf("failed to listen: %v", err)
	}

	s := grpc.NewServer()
	gwpb.RegisterGatewaysServer(s, server.Server{})

	log.Printf("server listening to %s", listen.Addr().String())
	err = s.Serve(listen)
	if err != nil {
		log.Fatalf("failed to serve: %v", err)
	}
}
