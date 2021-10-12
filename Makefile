ifeq ($(GOPATH),)
GOPATH := $(HOME)/go
endif

all: gen test lint vet build

build: server gateway

server:
	@cd cmd/$@ && go build -o ../../bin/$@ -tags osusergo,netgo
	
gateway:
	@cd cmd/$@ && go build -o ../../bin/$@ -tags osusergo,netgo

test:
	@go test ./...

vet:
	@go vet ./...

lint:
	@revive ./...

clean:
	@rm -rf pkg/gwpb
	
gen:
	@rm -rf pkg/gwpb
	@buf generate

proto-lint:
	@buf lint

count:
	@echo "Linecounts excluding generated and third party code"
	@gocloc --not-match-d='gwpb' .

dep-install:
	go install google.golang.org/protobuf/cmd/protoc-gen-go@v1.27.1
	go install google.golang.org/grpc/cmd/protoc-gen-go-grpc@v1.1.0
	go install github.com/grpc-ecosystem/grpc-gateway/v2/protoc-gen-grpc-gateway@v2.5.0
	go install github.com/grpc-ecosystem/grpc-gateway/v2/protoc-gen-openapiv2@v2.5.0
	go install github.com/bufbuild/buf/cmd/buf@v0.51.1
	go install github.com/mgechev/revive@latest
