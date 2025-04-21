.PHONY: proto build build-client run-server test clean

# Generate protobuf files
proto:
	protoc --go_out=. --go-grpc_out=. api/*.proto

# Build the server
build:
	mkdir -p bin
	go build -o bin/server cmd/server/main.go

# Build the client
build-client:
	mkdir -p bin
	go build -o bin/client cmd/client/main.go

# Build all
build-all: build build-client

# Run the server
run-server:
	./bin/server

# Run tests
test:
	go test -v ./...

# Clean build artifacts
clean:
	rm -rf bin/ 