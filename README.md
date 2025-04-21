# gRPC User Service with SQLite

A production-ready Go microservice for user management using gRPC and SQLite for persistent storage.

## Features

- User creation and retrieval via gRPC
- Persistent storage with SQLite
- Clean architecture with domain-driven design
- Structured logging with Zap
- Configuration management with Viper
- Metrics with Prometheus
- Health checking
- Middleware and interceptors
- Request tracing
- Authentication and authorization
- Rate limiting
- Graceful shutdown
- Docker support with security best practices
- CI/CD with GitHub Actions

## Project Structure

```
/
├── .github/workflows/   # CI/CD configuration
├── api/                 # Protobuf/gRPC definitions
├── cmd/                 # Main application entrypoints
│   ├── server/          # gRPC server
│   ├── client/          # gRPC client for testing
│   └── healthcheck/     # Health check tool
├── config/              # Configuration files
├── internal/
│   ├── app/             # Core business logic (use cases)
│   ├── domain/          # Entities/interfaces
│   ├── handler/         # gRPC handlers
│   ├── repo/            # Data access
│   │   ├── memory/      # In-memory repository implementation
│   │   └── sqlite/      # SQLite repository implementation
│   └── service/         # External services
├── pkg/                 # Reusable libraries
│   ├── config/          # Configuration utilities
│   ├── db/              # Database utilities
│   ├── logger/          # Logging utilities
│   ├── metrics/         # Metrics utilities
│   └── middleware/      # gRPC middleware
├── test/                # Integration tests
├── user/                # Generated protobuf code
└── Makefile, Dockerfile, etc.
```

## Prerequisites

- Go 1.18+ (tested with 1.21)
- Protocol Buffers Compiler (protoc)
- SQLite3

## Getting Started

### Installation

1. Clone the repository
   ```
   git clone https://github.com/akashdeep-patra/go-grpc-sqlite.git
   cd go-grpc-sqlite
   ```

2. Install dependencies
   ```
   go mod download
   ```

3. Generate protobuf files
   ```
   make proto
   ```

4. Build the application
   ```
   make build-all
   ```

### Running the Server

```
make run-server
```

### Using the Client

Create a new user:
```
./bin/client --create --name="John Doe" --email="john@example.com"
```

Retrieve a user by ID:
```
./bin/client --get=<user-id>
```

### Running with Docker

Build the Docker image:
```
make docker-build
```

Run the Docker container:
```
make docker-run
```

## Development

### Building and Testing

Build all binaries:
```
make build-all
```

Run unit tests:
```
make test
```

Run integration tests:
```
make integration-test
```

Format code:
```
make fmt
```

Lint code:
```
make lint
```

### Configuration

The service can be configured through:
1. Configuration file (`config/config.yaml`)
2. Environment variables (e.g., `APP_SERVER_PORT=50051`)

## Production Features

### Logging
- Structured logging with Zap
- Different log formats based on environment (development/production)
- Request context preservation through interceptors

### Metrics
- Prometheus metrics for request counts, durations, and errors
- Metrics server exposed on port 9100
- Endpoint: `/metrics`

### Health Checking
- Implementation of gRPC Health Checking Protocol
- Health check tool for Docker health checks
- Service status management during startup and shutdown

### Middleware
- Logging interceptor
- Authentication interceptor (placeholder for your auth implementation)
- Recovery interceptor for panic handling
- Rate limiting interceptor

### Security
- Non-root user in Docker
- Static binaries with security flags
- Authorization middleware
- HTTPS support (add your certificates)

## License

MIT 