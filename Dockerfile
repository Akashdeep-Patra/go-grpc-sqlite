# Build stage
FROM golang:1.21-alpine AS builder

# Install git for fetching dependencies
RUN apk --no-cache add git curl

# Set working directory
WORKDIR /app

# Copy scripts first for better caching
COPY scripts/ scripts/
RUN chmod +x scripts/*.sh

# Download proto files and Swagger UI
RUN ./scripts/download_protos.sh && ./scripts/download_swagger_ui.sh

# Copy go.mod and go.sum files
COPY go.mod go.sum ./

# Download dependencies
RUN go mod download

# Copy the source code
COPY . .

# Generate protobuf files
RUN mkdir -p api/swagger && \
    protoc --proto_path=. \
    --proto_path=./third_party/proto \
    --go_out=. --go-grpc_out=. \
    --grpc-gateway_out=logtostderr=true:. \
    --openapiv2_out=logtostderr=true:./api/swagger \
    api/*.proto

# Build the server and gateway with security flags
RUN CGO_ENABLED=1 \
    GOOS=linux \
    GOARCH=amd64 \
    go build -a -ldflags="-w -s -extldflags=-static" \
    -o /app/bin/server /app/cmd/server/main.go

RUN CGO_ENABLED=1 \
    GOOS=linux \
    GOARCH=amd64 \
    go build -a -ldflags="-w -s -extldflags=-static" \
    -o /app/bin/gateway /app/cmd/gateway/main.go

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

# Switch to non-root user
USER appuser

# Expose ports for gRPC, metrics, and HTTP gateway
EXPOSE 50051 9100 8080

# Set healthcheck
HEALTHCHECK --interval=30s --timeout=5s --start-period=5s --retries=3 \
    CMD [ "/app/server", "-health" ] || exit 1

# Run the application - use an entrypoint script that starts both servers
COPY --from=builder --chown=appuser:appgroup /app/scripts/entrypoint.sh /app/entrypoint.sh
RUN chmod +x /app/entrypoint.sh

ENTRYPOINT ["/app/entrypoint.sh"] 