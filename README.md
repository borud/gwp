# GWP - Gateway Protocol

This is just an exercise in prototyping a simple gateway protocol.  The parts of the code you care about are really in the `proto` directory.

If you want to look at how this ends up looking in Go code you can have a look at the `pkg/server/server.go` and `cmd/gateway/main.go` files.

## Addressing

The addressing within this protocol deals with the address space of network the devices connected to the gateway reside on.  The IP addresses of the gateways belong in the layer above this protocol.

## gRPC interface

Right now we have a gRPC interface with just one function: `Connect`.  Since we want to be able to run this across whatever connectivity layer we just assume that the transport transports `gwpb.Packet` instances.  So we might as well just do that using gRPC for test purposes.  In other implementations we probably want to use UDP packets.
