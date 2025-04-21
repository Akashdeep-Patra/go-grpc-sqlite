# Contributing to go-grpc-sqlite

Thank you for your interest in contributing to this project! This document outlines the process for contributing to this project and provides instructions for testing your changes locally.

## Development Process

1. Fork the repository
2. Create a feature branch (`git checkout -b feature/your-feature-name`)
3. Make your changes
4. Run tests locally (see below)
5. Commit your changes (`git commit -am 'Add some feature'`)
6. Push to the branch (`git push origin feature/your-feature-name`)
7. Create a new Pull Request

## Testing Locally

### Running Tests

You can run the standard Go tests with:

```bash
make test
```

For integration tests:

```bash
make integration-test
```

### Testing GitHub Actions Locally

This project includes a local CI/CD testing environment that allows you to test your changes against the same GitHub Actions workflows that run in the cloud. This helps catch issues before pushing to GitHub.

#### Prerequisites

- Docker
- Docker Compose

#### Running Local CI/CD Tests

To run all CI jobs locally:

```bash
make ci-local
```

To run specific CI jobs:

```bash
make ci-local-lint    # Run linting job
make ci-local-test    # Run tests job
make ci-local-build   # Run build job
make ci-local-docker  # Run Docker build job
```

To clean up local CI resources:

```bash
make ci-local-clean
```

For help:

```bash
make ci-local-help
```

#### How It Works

Each job in the GitHub Actions workflow is implemented as a Docker Compose service in `.github/actions-runner/docker-compose.yml`. The `run-actions.sh` script orchestrates these services and provides a summary of results.

## Code Style

Please follow these code style guidelines:

- Run `go fmt ./...` before committing
- Use meaningful variable and function names
- Write clear comments
- Follow the project's existing code style

## Documentation

When adding new features, please update the documentation:

- Add/modify usage examples in the README
- Document new functions and types with Go doc comments
- Update API documentation when changing gRPC services

## License

By contributing, you agree that your contributions will be licensed under the project's license. 