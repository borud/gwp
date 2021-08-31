# GWP - Gateway Protocol

This is just an exercise in prototyping a simple gateway protocol.  The parts of the code you care about are really in the `proto` directory.

If you want to look at how this ends up looking in Go code you can have a look at the `pkg/server/server.go` and `cmd/gateway/main.go` files.

## Addressing

The addressing within this protocol deals with the address space of network the devices connected to the gateway reside on.  The IP addresses of the gateways belong in the layer above this protocol.

## gRPC interface and binaries

There is a `server` and a `gateway` binary that gets built.  These are just examples of the server and the gateway end of the protocol which we can evolve into testing tools.

We currently use gRPC as the transport, but don't get too hung up on that.  We can also use UDP, HTTP, Websockets and whatever other transports to transport the packet instances.  We use a gRPC interface just as a convenience.

*Right now we have a gRPC interface with just one function: `Connect`.  Since we want to be able to run this across whatever connectivity layer we just assume that the transport transports `gwpb.Packet` instances.*
