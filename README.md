# GWP - Gateway Protocol

This is just an exercise in prototyping a simple gateway protocol.  The parts of the code you care about are really in the `proto` directory.

If you want to look at how this ends up looking in Go code you can have a look at the `pkg/server/server.go` and `cmd/gateway/main.go` files.

## Before you start

It is assumed that you have Go 1.17 or newer installed.

The project uses generated types that are generated from protobuffer definitions.  We don't check in the generated code, so be aware that
until you have run code generation, editors will complain about missing
types.

To install the tooling run:

```shell
make dep-install
```

Then run the code-generation with

```shell
make gen
```

The code will be regenerated on each build when you just run `make`.  If you add or rename types you may have to do a `make clean` before running `make` to get rid of superfloous files.

## What problems are we trying to solve?

### Round 1

- We are assuming that we have a set of gateways.

- Behind these gateways are devices that are connected to the gateway
  in some local network.

- We wish to be able to communicate with these devices through the gateways.

- We don't know how the gateway communicates with the devices - that depends
  entirely on the nature of the networking technology used between the
  gateways and the devices.

- Devices behind gateways need some form of persistent address and the gateway
  must be able to maintan the mapping between the persistent address and the
  (possibly ephemeral) actual address.

### Round 2

- Given a set of gateways, behind which there is a set of devices, we
  need a mechanism to determine which gateway(s) can relay messages to
  a given device.

- For meshed network topologies we must have support for representing the
  topology so the application can be aware of current topology and compute changes in topology.

## gRPC interface and binaries

There is a `server` and a `gateway` binary that gets built.  These are just examples of the server and the gateway end of the protocol which we can evolve into testing tools.

We currently use gRPC as the transport, but don't get too hung up on that.  We can also use UDP, HTTP, Websockets and whatever other transports to transport the packet instances.  We use a gRPC interface just as a convenience.

*Right now we have a gRPC interface with just one function: `Connect`.  Since we want to be able to run this across whatever connectivity layer we just assume that the transport transports `gwpb.Packet` instances.*
