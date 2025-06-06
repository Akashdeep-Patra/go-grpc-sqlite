.PHONY: proto build build-client build-healthcheck build-gateway build-all run-server run-gateway test integration-test clean docker-build docker-run lint docs swagger-ui fmt ci-local ci-local-lint ci-local-test ci-local-build ci-local-docker ci-local-clean ci-local-help

# Generate protobuf files
proto:
	./scripts/download_protos.sh
	mkdir -p api/swagger gen/go
	protoc --proto_path=. \
		--proto_path=./third_party/proto \
		--go_out=./gen/go --go-grpc_out=./gen/go \
		--grpc-gateway_out=logtostderr=true:./gen/go \
		--openapiv2_out=logtostderr=true:./api/swagger \
		api/*.proto

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

# Build the API gateway
build-gateway:
	mkdir -p bin
	./scripts/download_swagger_ui.sh
	go build -o bin/gateway cmd/gateway/main.go

# Build all binaries
build-all: build build-client build-healthcheck build-gateway

# Run the server
run-server:
	./bin/server

# Run the gateway server
run-gateway:
	./bin/gateway

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
	docker run -p 50051:50051 -p 9100:9100 -p 8080:8080 --name go-grpc-sqlite go-grpc-sqlite:latest

# Generate OpenAPI documentation
docs:
	./scripts/download_protos.sh
	mkdir -p api/swagger
	protoc --proto_path=. \
		--proto_path=./third_party/proto \
		--openapiv2_out=logtostderr=true:./api/swagger \
		api/*.proto

# Run Swagger UI for API documentation
swagger-ui:
	@echo "Access Swagger UI at http://localhost:8080/swagger/"
	./bin/gateway

# Clean build artifacts
clean:
	rm -rf bin/ gen/

# Run all GitHub Actions locally
ci-local:
	./.github/actions-runner/run-actions.sh

# Run linting job from GitHub Actions locally
ci-local-lint:
	./.github/actions-runner/run-actions.sh lint

# Run test job from GitHub Actions locally
ci-local-test:
	./.github/actions-runner/run-actions.sh test

# Run build job from GitHub Actions locally
ci-local-build:
	./.github/actions-runner/run-actions.sh build

# Run Docker build job from GitHub Actions locally
ci-local-docker:
	./.github/actions-runner/run-actions.sh docker-build

# Clean GitHub Actions local runner resources
ci-local-clean:
	./.github/actions-runner/run-actions.sh clean

# Show help for GitHub Actions local runner
ci-local-help:
	./.github/actions-runner/run-actions.sh help 