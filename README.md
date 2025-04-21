# gRPC User Service with SQLite

A Go microservice for user management using gRPC and SQLite for persistent storage.

## Features

- User creation and retrieval via gRPC
- Persistent storage with SQLite
- Clean architecture with domain-driven design
- Proper error handling
- Graceful shutdown

## Project Structure

```
/
├── api/                # Protobuf/gRPC definitions
├── cmd/                # Main application entrypoints
│   ├── server/         # gRPC server
│   └── client/         # gRPC client for testing
├── internal/
│   ├── app/            # Core business logic (use cases)
│   ├── domain/         # Entities/interfaces
│   ├── handler/        # gRPC handlers
│   ├── repo/           # Data access
│   │   ├── memory/     # In-memory repository implementation
│   │   └── sqlite/     # SQLite repository implementation
│   └── service/        # External services
├── pkg/                # Reusable libraries
│   └── db/             # Database utilities
├── user/               # Generated protobuf code
└── Makefile, Dockerfile, etc.
```

## Prerequisites

- Go 1.18+ (tested with 1.23)
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

## Development

### Adding New Features

1. Add new definitions to the `api/user.proto` file
2. Regenerate the Go code with `make proto`
3. Implement the new functionality in the appropriate layers:
   - Add domain models/interfaces in `internal/domain`
   - Add business logic in `internal/app`
   - Add data access in `internal/repo/sqlite`
   - Add gRPC handlers in `internal/handler`

## License

MIT 