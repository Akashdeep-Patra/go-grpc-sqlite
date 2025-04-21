# gRPC User Service with SQLite

A production-ready Go microservice for user management using gRPC and SQLite for persistent storage, with REST API support via gRPC Gateway.

## Features

- User creation and retrieval via gRPC and REST API
- Persistent storage with SQLite
- Clean architecture with domain-driven design
- Structured logging with Zap
- Configuration management with Viper
- Metrics with Prometheus
- Health checking
- Middleware and interceptors
- Request tracing
- Rate limiting
- Graceful shutdown
- Docker support with security best practices
- OpenAPI/Swagger documentation
- REST API via gRPC Gateway

## Project Structure

```
/
├── .github/workflows/   # CI/CD configuration
├── api/                 # Protobuf/gRPC definitions and Swagger docs
├── bin/                 # Compiled binaries
├── cmd/                 # Main application entrypoints
│   ├── server/          # gRPC server
│   ├── client/          # gRPC client
│   ├── gateway/         # gRPC Gateway server
│   └── healthcheck/     # Health check tool
├── config/              # Configuration files
├── gen/                 # Generated code (protobuf)
├── internal/
│   ├── app/             # Core business logic (use cases)
│   ├── domain/          # Entities/interfaces
│   ├── handler/         # gRPC handlers
│   ├── repo/            # Data access
│   │   ├── memory/      # In-memory repository implementation
│   │   └── sqlite/      # SQLite repository implementation
├── pkg/                 # Reusable libraries
│   ├── config/          # Configuration utilities
│   ├── db/              # Database utilities
│   ├── logger/          # Logging utilities
│   ├── metrics/         # Metrics utilities
│   └── middleware/      # gRPC middleware
├── scripts/             # Helper scripts
├── test/                # Integration tests
├── third_party/         # Third-party proto definitions
├── user/                # Generated protobuf code
└── Makefile, Dockerfile, etc.
```

## Prerequisites

- Go 1.23+ (project uses Go 1.23.4)
- Protocol Buffers Compiler (protoc)
- protoc plugins:
  - protoc-gen-go
  - protoc-gen-go-grpc
  - protoc-gen-grpc-gateway
  - protoc-gen-openapiv2
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

3. Download third-party proto files and generate protobuf code
   ```
   make proto
   ```

4. Build the application
   ```
   make build-all
   ```

### Running the Server

Start the gRPC server:
```
make run-server
```

Start the REST API Gateway:
```
make run-gateway
```

### Using the Client

Create a new user:
```
./bin/client -create -name="John Doe" -email="john@example.com"
```

Retrieve a user by ID:
```
./bin/client -get="user-id"
```

### API Documentation

The API is documented using OpenAPI/Swagger. After starting the gateway server, access the Swagger UI at:

```
http://localhost:8080/swagger/
```

### SQLite Database

The SQLite database file is created at runtime. By default, it is stored at:

```
data/sqlite.db
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

Generate API documentation:
```
make docs
```

View API documentation with Swagger UI:
```
make swagger-ui
```

### Testing GitHub Actions Locally

You can test GitHub Actions workflow locally before pushing to GitHub:

Run all CI jobs locally:
```
make ci-local
```

Run specific CI jobs:
```
make ci-local-lint    # Run linting job
make ci-local-test    # Run tests job
make ci-local-build   # Run build job
make ci-local-docker  # Run Docker build job
```

Clean up local CI resources:
```
make ci-local-clean
```

Show help for local CI:
```
make ci-local-help
```

### Configuration

The service can be configured through:
1. Configuration file (`config/config.yaml`)
2. Environment variables 

Key configuration options:
- `APP_SERVER_PORT`: gRPC server port (default: 50051)
- `APP_HTTP_PORT`: Gateway HTTP server port (default: 8080)
- `APP_DB_PATH`: SQLite database path
- `APP_ENVIRONMENT`: Environment (development/production)

## Developer Setup and Workflow

### First-Time Setup

1. **Install Required Tools**
   ```bash
   # Install Go (if not already installed)
   brew install go  # macOS
   # For other platforms see: https://golang.org/doc/install

   # Install Protocol Buffers compiler
   brew install protobuf  # macOS
   # apt-get install protobuf-compiler  # Debian/Ubuntu

   # Install required Go tools
   go install google.golang.org/protobuf/cmd/protoc-gen-go@latest
   go install google.golang.org/grpc/cmd/protoc-gen-go-grpc@latest
   go install github.com/grpc-ecosystem/grpc-gateway/v2/protoc-gen-grpc-gateway@latest
   go install github.com/grpc-ecosystem/grpc-gateway/v2/protoc-gen-openapiv2@latest
   ```

2. **Clone and Set Up the Repository**
   ```bash
   git clone https://github.com/Akashdeep-Patra/go-grpc-sqlite.git
   cd go-grpc-sqlite
   go mod download
   make proto  # Downloads required proto files and generates code
   ```

3. **VSCode Setup (Recommended)**
   - Install the Go extension for VSCode
   - Install the Protocol Buffers extension for .proto file support
   - Configure settings.json with:
     ```json
     {
       "go.lintTool": "golangci-lint",
       "go.formatTool": "goimports",
       "editor.formatOnSave": true
     }
     ```

### Development Workflow

1. **Making Changes to API**
   - Edit proto files in the `api/` directory
   - Regenerate code:
     ```bash
     make proto
     ```
   - This will update generated files in `gen/go/` directory

2. **Implementation Workflow**
   - **Domain Layer**: Add/modify entity definitions in `internal/domain/`
   - **Repository Layer**: Implement data access in `internal/repo/`
   - **Service Layer**: Implement business logic in `internal/app/`
   - **Handler Layer**: Connect gRPC endpoints to services in `internal/handler/`

3. **Test-Driven Development**
   - Write tests before implementation
   - Run tests frequently:
     ```bash
     go test ./... # Run all tests
     go test ./internal/app/... # Test specific package
     ```

4. **Local Development Loop**
   ```bash
   # Make changes to code
   make build  # Build the server
   make run-server  # Run in one terminal

   # In another terminal
   make build-gateway  # Build the gateway
   make run-gateway  # Run the REST API gateway

   # Test with the client
   ./bin/client -create -name="Test User" -email="test@example.com"
   ```

5. **API Documentation**
   - After making changes to proto files, regenerate docs:
     ```bash
     make docs
     ```
   - View updated documentation:
     ```bash
     make swagger-ui
     ```
   - Access Swagger UI at http://localhost:8080/swagger/

### Branch Strategy

1. **Main Branch Structure**
   - `main`: Production-ready code
   - `develop`: Integration branch for feature work

2. **Feature Development**
   - Create feature branches from `develop`:
     ```bash
     git checkout -b feature/user-authentication develop
     ```
   - Make small, focused commits
   - Write meaningful commit messages

3. **Pull Request Workflow**
   - Push feature branch to GitHub
   - Create PR against `develop` branch
   - Ensure tests pass and code meets standards
   - Request code review

### Versioning and Releases

1. **Semantic Versioning**
   - Follow [SemVer](https://semver.org/) for version numbers
   - MAJOR.MINOR.PATCH (e.g., 1.2.3)
   - Increment accordingly based on backward compatibility

2. **Release Process**
   - Merge `develop` into `main` for releases
   - Tag releases with version number:
     ```bash
     git tag -a v1.0.0 -m "Version 1.0.0"
     git push origin v1.0.0
     ```

### Troubleshooting Common Issues

1. **Proto Generation Issues**
   - Ensure all required tools are installed and in PATH
   - Check that third-party proto files are downloaded:
     ```bash
     ./scripts/download_protos.sh
     ```

2. **Import Path Problems**
   - If you see import errors, ensure the module path in go.mod matches your repository
   - Update import paths across the codebase if module name changes

3. **Swagger UI Not Loading**
   - Check that the gateway server is running
   - Verify Swagger JSON file is generated at `api/swagger/api/user.swagger.json`
   - Check browser console for CORS or other errors

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
- Recovery interceptor for panic handling
- Rate limiting interceptor

### REST API Gateway
- HTTP/JSON API via gRPC Gateway
- OpenAPI/Swagger documentation
- Automatic translation between HTTP/JSON and gRPC

## Known Issues

- There are linter errors in the protobuf imports that need to be resolved with proper third-party proto imports

## License

MIT 