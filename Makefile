ifeq ($(GOPATH),)
GOPATH := $(HOME)/go
endif

all: gen

clean:
	@rm -rf pkg/gwpb
	
gen:
	@buf generate

proto-lint:
	@buf lint

dep-install:
	go get google.golang.org/protobuf/cmd/protoc-gen-go@v1.27.1
	go get google.golang.org/grpc/cmd/protoc-gen-go-grpc@v1.1.0
	go get github.com/grpc-ecosystem/grpc-gateway/v2/protoc-gen-grpc-gateway@v2.5.0
	go get github.com/grpc-ecosystem/grpc-gateway/v2/protoc-gen-openapiv2@v2.5.0
	go get github.com/bufbuild/buf/cmd/buf@v0.51.1
	go get github.com/mgechev/revive@latest
