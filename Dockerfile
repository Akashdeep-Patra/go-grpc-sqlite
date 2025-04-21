# Build stage
FROM golang:1.21-alpine AS builder

# Install git for fetching dependencies
RUN apk --no-cache add git

# Set working directory
WORKDIR /app

# Copy go.mod and go.sum files
COPY go.mod go.sum ./

# Download dependencies
RUN go mod download

# Copy the source code
COPY . .

# Build the application with security flags
RUN CGO_ENABLED=1 \
    GOOS=linux \
    GOARCH=amd64 \
    go build -a -ldflags="-w -s -extldflags=-static" \
    -o /app/bin/server /app/cmd/server/main.go

# Final stage
FROM alpine:latest

# Add necessary packages
RUN apk --no-cache add ca-certificates tzdata sqlite && \
    update-ca-certificates

# Create a non-root user and group
RUN addgroup -S appgroup && \
    adduser -S appuser -G appgroup

# Create app directories
RUN mkdir -p /app/data /app/config && \
    chown -R appuser:appgroup /app

# Set working directory
WORKDIR /app

# Copy binary from builder stage
COPY --from=builder --chown=appuser:appgroup /app/bin/server /app/server

# Copy config files
COPY --from=builder --chown=appuser:appgroup /app/config /app/config

# Switch to non-root user
USER appuser

# Expose ports for gRPC and metrics
EXPOSE 50051 9100

# Set healthcheck
HEALTHCHECK --interval=30s --timeout=5s --start-period=5s --retries=3 \
    CMD [ "/app/server", "-health" ] || exit 1

# Run the application
ENTRYPOINT ["/app/server"] 