.PHONY: build clean
GO=CGO_ENABLED=1 GO111MODULE=on go

BINARIES=my-plugin/rpc/rpc my-plugin/grpc/grpc my-service/my-service
.PHONY: $(BINARIES)

build: $(BINARIES)

my-plugin/rpc/rpc:
	cd ./my-plugin/rpc; $(GO) build

my-plugin/grpc/grpc:
	cd ./my-plugin/grpc; $(GO) build

my-service/my-service:
	cd ./my-service; $(GO) build

clean:
	rm -f $(BINARIES)

run-grpc:
	cd ./my-service; ./my-service

run-rpc:
	cd ./my-service; ./my-service rpc
