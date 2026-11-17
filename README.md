# Go gRPC Example

A simple, self-documenting Go implementation of gRPC demonstrating all 4 communication patterns with error handling.

## Communication Patterns

1. **Normal (Unary) RPC**: Standard request-response. Client sends a message; server returns a single response string.
2. **Server Streaming RPC**: Client sends a single request; server streams responses back one by one with a delay.
3. **Client Streaming RPC**: Client streams messages incrementally; server prints each message received and returns a summary response upon completion.
4. **Bi-Directional Streaming RPC**: Client sends messages incrementally; server prints each message and concurrently streams responses back with a delay.

Each method supports triggering errors dynamically via the `error` parameter in the request.

---

## Project Structure

```text
.
├── client/
│   └── main.go           # Interactive console client (5 menu options)
├── server/
│   └── main.go           # gRPC server implementation
├── proto/
│   ├── service.proto     # Protobuf service & message definitions
│   ├── service.pb.go     # Generated Go Protobuf structs
│   └── service_grpc.pb.go# Generated Go gRPC interfaces and stubs
├── .gitignore            # Go gitignore rules
├── go.mod                # Go module definition
├── go.sum                # Go checksums
├── Makefile              # Build automation
└── README.md             # Project documentation
```

---

## Prerequisites

- [Go](https://go.dev/dl/) (1.25 or newer)
- [Protocol Buffer Compiler (`protoc`)](https://github.com/protocolbuffers/protobuf/releases)
- Go Protobuf & gRPC plugins:
  ```bash
  go install google.golang.org/protobuf/cmd/protoc-gen-go@latest
  go install google.golang.org/grpc/cmd/protoc-gen-go-grpc@latest
  ```

---

## Quick Start

### 1. Compile Protobuf (Optional if stubs already exist)

Using `make`:
```bash
make proto
```

Or using `protoc` directly:
```bash
protoc --go_out=. --go_opt=paths=source_relative --go-grpc_out=. --go-grpc_opt=paths=source_relative proto/service.proto
```

### 2. Run the Server

Open a terminal and start the gRPC server:

```bash
# Using make
make run-server

# Or with go run
go run ./server
```

The server will listen on port `:50051`.

### 3. Run the Client

Open a second terminal and start the interactive client:

```bash
# Using make
make run-client

# Or with go run
go run ./client
```

---

## Interactive Client Menu

When running the client, you will see the interactive console prompt:

```text
=== Select gRPC Request ===
1. Normal
2. Server Stream
3. Client Stream
4. Bi Directional
5. Exit
Select option (1-5):
```

- **Option 1 (Normal)**: Enter a message and choose whether to trigger an error (`y/n`).
- **Option 2 (Server Stream)**: Enter a message; server streams 3 chunks with a 500ms delay.
- **Option 3 (Client Stream)**: Enter comma-separated messages; client streams them one-by-one to the server.
- **Option 4 (Bi-Directional)**: Enter comma-separated messages; client streams messages while simultaneously receiving streamed replies from the server.
- **Option 5 (Exit)**: Closes the connection and exits.

---

## Building Binaries

Build standalone executables for both server and client into the `bin/` directory:

```bash
# Using make
make build

# Or manually
go build -o bin/server ./server
go build -o bin/client ./client
```

To clean up compiled binaries:
```bash
make clean
```
