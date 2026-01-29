# Service Registry

A simple service registry API for tracking and monitoring services.

## Overview

This service provides a REST API for registering, discovering, and health-checking services.
It's designed to be a lightweight service registry for development and testing environments.

## Prerequisites

- Go 1.22 or later
- golangci-lint (for linting)

## Getting Started

1. Clone the repository:

   ```bash
   git clone <repository-url>
   cd service-registry
   ```

2. Copy the environment file:

   ```bash
   cp .env.example .env
   ```

3. Run the server:

   ```bash
   make run
   ```

   The server will start on `http://localhost:8080` by default.

## Available Commands

Run `make help` to see all available commands:

```
build   Build the server binary
run     Run the server locally
test    Run all tests
lint    Run linter (requires golangci-lint)
fmt     Format code
clean   Remove build artifacts
tidy    Tidy go modules
```

## API Endpoints

| Method | Endpoint | Description |
|--------|----------|-------------|
| GET | `/health` | Health check for the registry |
| GET | `/services` | List all registered services |
| POST | `/services` | Register a new service |
| GET | `/services/{name}` | Get a specific service |
| PUT | `/services/{name}` | Update a service |
| DELETE | `/services/{name}` | Deregister a service |
| GET | `/services/{name}/health` | Check health of a service |

## Example Usage

Register a service:

```bash
curl -X POST http://localhost:8080/services \
  -H "Content-Type: application/json" \
  -d '{
    "name": "my-api",
    "endpoint": "http://localhost:3000",
    "health_check_url": "http://localhost:3000/health",
    "description": "My API service"
  }'
```

List all services:

```bash
curl http://localhost:8080/services
```

Check a service's health:

```bash
curl http://localhost:8080/services/my-api/health
```

## Configuration

The server can be configured using environment variables:

| Variable | Description | Default |
|----------|-------------|---------|
| `PORT` | Server port | `8080` |

## Contributing

Please read [CONTRIBUTING.md](CONTRIBUTING.md) for guidelines on how to contribute to this project.

## License

Internal use only.
