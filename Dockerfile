# Build stage
FROM golang:1.23-alpine AS builder

# Install git for fetching dependencies and other tools
RUN apk --no-cache add git curl protobuf protobuf-dev gcc musl-dev

# Set Go toolchain to auto
ENV GOTOOLCHAIN=auto

# Set working directory
WORKDIR /app

# Copy scripts first for better caching
COPY scripts/ scripts/
RUN chmod +x scripts/*.sh && ls -la scripts/

# Install protoc plugins with specific versions
RUN go install google.golang.org/protobuf/cmd/protoc-gen-go@v1.31.0 && \
    go install google.golang.org/grpc/cmd/protoc-gen-go-grpc@v1.3.0 && \
    go install github.com/grpc-ecosystem/grpc-gateway/v2/protoc-gen-grpc-gateway@v2.16.0 && \
    go install github.com/grpc-ecosystem/grpc-gateway/v2/protoc-gen-openapiv2@v2.16.0

# Copy go.mod and go.sum files
COPY go.mod go.sum ./

# Download dependencies
RUN go mod download

# Copy the source code
COPY . .

# Create third_party directory manually
RUN mkdir -p third_party/proto/google/api third_party/proto/protoc-gen-openapiv2/options

# Download proto files manually
RUN curl -sSL "https://raw.githubusercontent.com/googleapis/googleapis/master/google/api/annotations.proto" \
    -o "third_party/proto/google/api/annotations.proto" && \
    curl -sSL "https://raw.githubusercontent.com/googleapis/googleapis/master/google/api/http.proto" \
    -o "third_party/proto/google/api/http.proto" && \
    curl -sSL "https://raw.githubusercontent.com/googleapis/googleapis/master/google/api/field_behavior.proto" \
    -o "third_party/proto/google/api/field_behavior.proto" && \
    curl -sSL "https://raw.githubusercontent.com/grpc-ecosystem/grpc-gateway/main/protoc-gen-openapiv2/options/annotations.proto" \
    -o "third_party/proto/protoc-gen-openapiv2/options/annotations.proto" && \
    curl -sSL "https://raw.githubusercontent.com/grpc-ecosystem/grpc-gateway/main/protoc-gen-openapiv2/options/openapiv2.proto" \
    -o "third_party/proto/protoc-gen-openapiv2/options/openapiv2.proto"

# Generate protobuf files
RUN mkdir -p api/swagger gen/go && \
    protoc --proto_path=. \
    --proto_path=./third_party/proto \
    --go_out=./gen/go --go-grpc_out=./gen/go \
    --grpc-gateway_out=logtostderr=true:./gen/go \
    --openapiv2_out=logtostderr=true:./api/swagger \
    api/*.proto

# Create swagger-ui directory manually
RUN mkdir -p cmd/gateway/swagger-ui

# Build the server and gateway with CGO enabled
RUN mkdir -p bin && \
    CGO_ENABLED=1 go build -o bin/server cmd/server/main.go && \
    CGO_ENABLED=1 go build -o bin/gateway cmd/gateway/main.go

# Final stage
FROM alpine:latest

# Add necessary packages
RUN apk --no-cache add ca-certificates tzdata sqlite && \
    update-ca-certificates

# Create a non-root user and group
RUN addgroup -S appgroup && \
    adduser -S appuser -G appgroup

# Create app directories
RUN mkdir -p /app/data /app/config /app/api/swagger && \
    chown -R appuser:appgroup /app

# Set working directory
WORKDIR /app

# Copy binaries from builder stage
COPY --from=builder --chown=appuser:appgroup /app/bin/server /app/server
COPY --from=builder --chown=appuser:appgroup /app/bin/gateway /app/gateway

# Copy config files
COPY --from=builder --chown=appuser:appgroup /app/config /app/config

# Copy swagger files
COPY --from=builder --chown=appuser:appgroup /app/api/swagger /app/api/swagger

# Create entrypoint script directly
RUN printf '#!/bin/sh\n\n# Start the gRPC server in the background\n/app/server &\nSERVER_PID=$!\n\n# Wait a bit for the server to start\nsleep 2\n\n# Start the gateway server\n/app/gateway --grpc-server-endpoint=localhost:50051 &\nGATEWAY_PID=$!\n\n# Set up traps for graceful shutdown\ntrap '\''kill -TERM $SERVER_PID $GATEWAY_PID; wait $SERVER_PID $GATEWAY_PID'\'' TERM INT\n\n# Wait for the processes to terminate\nwait $SERVER_PID $GATEWAY_PID\n' > /app/entrypoint.sh && \
    chmod +x /app/entrypoint.sh && \
    chown appuser:appgroup /app/entrypoint.sh && \
    ls -la /app

# Switch to non-root user
USER appuser

# Expose ports for gRPC, metrics, and HTTP gateway
EXPOSE 50051 9100 8080

# Set healthcheck
HEALTHCHECK --interval=30s --timeout=5s --start-period=5s --retries=3 \
    CMD [ "/app/server", "-health" ] || exit 1

ENTRYPOINT ["/app/entrypoint.sh"] 