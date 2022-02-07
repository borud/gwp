# Gateway Protocol

This repository contains the Protobuffer and gRPC definitions for the
Gateway Protocol.

## Using from Go

In order to use the types defined in this module you only have to include
the following in your Go program.

```go
import gwp "go.buf.build/grpc/go/borud/gwp/gwpb/v1"
```

This will use the buf.build proxy and serve you a ready-made library with the generated code.

## Documentation

You can find the API documentation at <https://buf.build/borud/gwp/docs>.

## Modifying the protocol

If you modify the protocol **please** make sure that you are not introducing any breaking changes.  You can check for breaking changes by running:

```shell
make breaking
```
If it doesn't list any errors you should be okay.  If it lists any errors you
should discuss the change with @borud before continuing.
