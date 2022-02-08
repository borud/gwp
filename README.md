# Gateway Protocol

This repository contains the Protobuffer and gRPC definitions for the
Gateway Protocol.

## Prerequisites

### Golang

It is assumed that you have a relatively recent version of Go installed.  Preferably 1.17 or newer.  If you don't have this, please visit
<https://go.dev/doc/install> and follow the instructions.

### buf

If you run "make dep" and you have Go installed properly (with $GOPATH/bin in your execution path), this will take care of installing buf for you.  This will currently install the latest version at the time you run it.  We may change this later to refer to some specific version.

### Make

Since we use Make to build this project it is preferable if you have make installed, but in a pinch you can look at the `Makefile` and figure out how to run the commands.

## Using from Go

In order to use the types defined in this module in your programs you only have to include the following import in your Go program.

```go
import gw "go.buf.build/borud/grpc-gateway4/borud/gwp/gwpb/v1"
```

...and run "go mod tidy" in order to trigger download of dependencies.

This will use the buf.build proxy and serve you a ready-made module with the generated code.

## Documentation

You can find the API documentation at <https://buf.build/borud/gwp/docs>.

## Modifying the protocol

If you modify the protocol **please** make sure that you are not introducing any breaking changes.  You can check for breaking changes by running:

```shell
make breaking
```

This will check against the last main branch and see if you have introduced any changes that break backwards compatibility.  If this command doesn't return any output or list any errors, you should be okay. If it lists any errors you
should discuss the change with @borud before continuing.
