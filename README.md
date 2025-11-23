# SmartStore

SmartStore is an intelligent object storage gateway that provides unified access to multiple cloud storage backends (S3, GCS) with intelligent routing and metadata management.

## Features

- 🚀 **Multi-Backend Support**: Seamlessly integrate with AWS S3, Google Cloud Storage, and other S3-compatible services
- 🎯 **Smart Routing**: Intelligent request routing based on regions and storage classes
- 💾 **Metadata Management**: Centralized metadata storage with PostgreSQL
- ⚡ **High Performance**: Redis caching for improved response times
- 🔄 **RESTful API**: Clean HTTP API for object operations

## Prerequisites

- Go 1.23 or higher
- PostgreSQL 16+ (for metadata storage)
- Redis 7+ (for caching)
- Docker & Docker Compose (optional, for local development)

## Quick Start

### 1. Clone and Setup

```bash
git clone https://github.com/kenelite/smartstore.git
cd smartstore
```

### 2. Start Dependencies

Using Docker Compose:

```bash
make docker-compose-up
```

This starts PostgreSQL, Redis, and MinIO (S3-compatible storage) in containers.

### 3. Configure

Edit `config.yaml` to match your environment and cloud storage credentials.

### 4. Initialize Database

```bash
make db-migrate
```

### 5. Run the Application

```bash
make run
```

Or with hot-reload for development:

```bash
make dev
```

## Makefile Commands

### Development

```bash
make deps              # Download Go dependencies
make tidy              # Tidy Go modules
make fmt               # Format Go code
make vet               # Run go vet
make lint              # Run golangci-lint
make install-tools     # Install development tools (golangci-lint, air, etc.)
```

### Build

```bash
make build             # Build the application
make build-linux       # Build for Linux (amd64)
make build-darwin      # Build for macOS (amd64/arm64)
make build-all         # Build for all platforms
make install           # Install binary to $GOPATH/bin
```

### Run

```bash
make run               # Run from source
make run-bin           # Build and run binary
make dev               # Run with hot-reload (requires air)
```

### Test

```bash
make test              # Run all tests with race detector
make test-short        # Run short tests
make test-coverage     # Run tests and generate coverage report
make test-bench        # Run benchmark tests
```

### Database

```bash
make db-create         # Create database
make db-drop           # Drop database
make db-migrate        # Run migrations
make db-reset          # Drop, create, and migrate
make db-shell          # Connect to database shell
```

Database configuration via environment variables:

```bash
DB_HOST=localhost \
DB_PORT=5432 \
DB_USER=smartstore \
DB_PASSWORD=smartstore \
DB_NAME=smartstore \
make db-migrate
```

### Docker

```bash
make docker-build           # Build Docker image
make docker-run             # Run Docker container
make docker-compose-up      # Start all services
make docker-compose-down    # Stop all services
make docker-compose-logs    # View logs
```

### Utilities

```bash
make clean             # Clean build artifacts
make clean-all         # Deep clean (including caches)
make info              # Display project information
make check-deps        # Check for outdated dependencies
make release-check     # Run pre-release checks
make release-build     # Build release binaries
make help              # Display all available commands
```

## Project Structure

```
smartstore/
├── cmd/
│   └── gateway/          # Application entry point
├── internal/
│   ├── api/
│   │   └── http/         # HTTP handlers
│   ├── app/              # Application initialization
│   ├── cache/            # Redis cache implementation
│   ├── config/           # Configuration management
│   ├── metadata/         # Metadata repository (memory, SQL)
│   └── storage/
│       ├── objectstore/  # Storage adapters (S3, GCS)
│       └── smart/        # Smart routing logic
├── db/
│   └── schema.sql        # Database schema
├── config.yaml           # Configuration file
├── Makefile              # Build automation
├── Dockerfile            # Container image definition
├── docker-compose.yaml   # Local development environment
└── README.md
```

## Configuration

The `config.yaml` file contains all application settings:

- HTTP server configuration
- Database connection settings
- Redis cache settings
- Storage backend configurations (S3, GCS)
- Routing rules

See `config.yaml` for detailed configuration options.

## API Usage

### Upload Object

```bash
curl -X PUT http://localhost:8080/v1/buckets/my-bucket/objects/my-file.txt \
  -H "Content-Type: text/plain" \
  --data-binary @my-file.txt
```

### Download Object

```bash
curl -X GET http://localhost:8080/v1/buckets/my-bucket/objects/my-file.txt \
  -o my-file.txt
```

### Delete Object

```bash
curl -X DELETE http://localhost:8080/v1/buckets/my-bucket/objects/my-file.txt
```

## Development

### Hot Reload with Air

```bash
make install-tools  # Install air
make dev           # Start with hot-reload
```

### Code Quality

```bash
make fmt lint vet test  # Format, lint, vet, and test
```

### Pre-commit Checks

```bash
make release-check  # Runs fmt, vet, and test
```

## Docker Development

Run the entire stack locally:

```bash
make docker-compose-up
```

This starts:
- **PostgreSQL** on port 5432
- **Redis** on port 6379
- **MinIO** on ports 9000 (API) and 9001 (Console)

Access MinIO Console at http://localhost:9001 (minioadmin/minioadmin)

## Production Deployment

### Build Release

```bash
make release-build
```

This creates optimized binaries for Linux and macOS in the `bin/` directory.

### Docker

```bash
make docker-build
docker push your-registry/smartstore:latest
```

## Environment Variables

Override configuration with environment variables:

```bash
export DB_HOST=your-db-host
export DB_PORT=5432
export DB_USER=your-db-user
export DB_PASSWORD=your-db-password
export DB_NAME=smartstore

make run
```

## Contributing

1. Fork the repository
2. Create your feature branch (`git checkout -b feature/amazing-feature`)
3. Make your changes
4. Run tests (`make test`)
5. Run linting (`make lint`)
6. Commit your changes (`git commit -m 'Add amazing feature'`)
7. Push to the branch (`git push origin feature/amazing-feature`)
8. Open a Pull Request

## License

[TBU]

## Support

For issues and questions, please open an issue on GitHub.
