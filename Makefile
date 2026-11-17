.PHONY: all proto build build-server build-client run-server run-client clean tidy help

BINARY_DIR ?= bin
SERVER_BIN ?= $(BINARY_DIR)/server
CLIENT_BIN ?= $(BINARY_DIR)/client

all: proto tidy build

help:
	@echo "Available commands:"
	@echo "  make proto        - Generate protobuf and gRPC Go stubs"
	@echo "  make build        - Build both server and client binaries into $(BINARY_DIR)/"
	@echo "  make run-server   - Run the server directly"
	@echo "  make run-client   - Run the client directly"
	@echo "  make tidy         - Run go mod tidy"
	@echo "  make clean        - Remove generated binaries"

proto:
	protoc --go_out=. --go_opt=paths=source_relative --go-grpc_out=. --go-grpc_opt=paths=source_relative proto/service.proto

build: build-server build-client

build-server:
	go build -o $(SERVER_BIN) ./server

build-client:
	go build -o $(CLIENT_BIN) ./client

run-server:
	go run ./server

run-client:
	go run ./client

tidy:
	go mod tidy

clean:
	go clean
	rm -rf $(BINARY_DIR) *.exe
