.PHONY: proto build build-client build-healthcheck build-all run-server test clean docker-build docker-run lint

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

# Build the healthcheck tool
build-healthcheck:
	mkdir -p bin
	go build -o bin/healthcheck cmd/healthcheck/main.go

# Build all binaries
build-all: build build-client build-healthcheck

# Run the server
run-server:
	./bin/server

# Run tests
test:
	go test -v ./...

# Run integration tests
integration-test:
	go test -v -tags=integration ./test/...

# Format code
fmt:
	go fmt ./...

# Lint code
lint:
	go vet ./...
	golangci-lint run

# Build Docker image
docker-build:
	docker build -t go-grpc-sqlite:latest .

# Run Docker container
docker-run:
	docker run -p 50051:50051 -p 9100:9100 --name go-grpc-sqlite go-grpc-sqlite:latest

# Clean build artifacts
clean:
	rm -rf bin/ 