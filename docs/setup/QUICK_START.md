# SmartStore Quick Start

Get up and running with SmartStore in minutes! 🚀

## Prerequisites

- Go 1.23+
- Docker & Docker Compose (optional but recommended)
- Make

## Quick Setup (5 minutes)

### Option 1: Automated Setup (Recommended)

```bash
# Clone the repository
git clone https://github.com/kenelite/smartstore.git
cd smartstore

# Run automated setup
make setup

# The setup script will:
# ✓ Download dependencies
# ✓ Install development tools
# ✓ Start Docker services (PostgreSQL, Redis, MinIO)
# ✓ Initialize database
# ✓ Build the application
```

### Option 2: Manual Setup

```bash
# 1. Download dependencies
make deps

# 2. Start services
make docker-compose-up

# 3. Initialize database
make db-migrate

# 4. Build
make build
```

## Running the Application

### Development Mode (with hot-reload)

```bash
make dev
```

### Production Mode

```bash
make run
```

### Using Binary

```bash
make run-bin
```

## Essential Commands

### Daily Development

```bash
make dev                    # Start with hot-reload
make test                   # Run tests
make fmt                    # Format code
make lint                   # Run linter
```

### Common Tasks

```bash
make help                   # Show all commands
make info                   # Show project info
make docker-compose-logs    # View service logs
make db-shell              # Connect to database
make clean                 # Clean build artifacts
```

## Testing Your Setup

### 1. Check if services are running

```bash
docker ps
```

You should see:
- smartstore-postgres
- smartstore-redis
- smartstore-minio

### 2. Test the API (after starting the app)

```bash
# Upload a file
curl -X PUT http://localhost:8080/v1/buckets/test-bucket/objects/test.txt \
  -H "Content-Type: text/plain" \
  --data "Hello SmartStore!"

# Download the file
curl http://localhost:8080/v1/buckets/test-bucket/objects/test.txt
```

## Configuration

Edit `config.yaml` to customize:
- Server port
- Database connection
- Storage backends
- Cache settings

## Service Endpoints

### Application
- API: http://localhost:8080

### Supporting Services
- PostgreSQL: localhost:5432
- Redis: localhost:6379
- MinIO API: http://localhost:9000
- MinIO Console: http://localhost:9001
  - Username: `minioadmin`
  - Password: `minioadmin`

## Troubleshooting

### Port already in use

```bash
# Stop existing services
make docker-compose-down

# Check what's using the port
lsof -i :8080  # or :5432, :6379, etc.
```

### Database connection failed

```bash
# Check if PostgreSQL is running
docker ps | grep postgres

# Restart services
make docker-compose-down
make docker-compose-up

# Wait a few seconds, then migrate
make db-migrate
```

### Build fails

```bash
# Clean and rebuild
make clean
make deps
make build
```

### Tests fail

```bash
# Reset test database
DB_NAME=smartstore_test make db-reset
make test
```

## Development Workflow

```bash
# 1. Start services (once)
make docker-compose-up

# 2. Start development with hot-reload
make dev

# 3. In another terminal, make changes and test
make test

# 4. Before committing
make fmt lint test
```

## Building for Production

```bash
# Run pre-release checks
make release-check

# Build for all platforms
make release-build

# Build Docker image
make docker-build
```

## Next Steps

1. **Read the full README**: `README.md`
2. **Explore Makefile commands**: `make help`
3. **Check contribution guidelines**: `CONTRIBUTING.md`
4. **Detailed Makefile guide**: `MAKEFILE_GUIDE.md`

## Need Help?

- `make help` - List all available commands
- Check logs: `make docker-compose-logs`
- Open an issue on GitHub
- Read the documentation in `README.md`

## Common Patterns

### Reset Everything

```bash
make clean-docker  # Stop and clean Docker
make clean        # Clean build artifacts
make setup        # Start fresh
```

### Update Dependencies

```bash
make check-deps   # Check for updates
make tidy        # Clean up go.mod
make deps        # Download dependencies
```

### Full Test Suite

```bash
make fmt          # Format
make vet          # Static analysis
make lint         # Linting
make test         # Tests
# Or all at once:
make release-check
```

## Useful Scripts

Located in `scripts/` directory:

- `setup.sh` - Automated setup
- `test.sh` - Comprehensive testing
- `clean.sh` - Advanced cleanup
- `deploy.sh` - Deployment preparation

## Tips & Tricks

1. **Alias for help**: Add `alias mh='make help'` to your shell rc file
2. **Watch logs**: Run `make docker-compose-logs` in a separate terminal
3. **Quick test**: Use `make test-short` for faster feedback
4. **Hot reload**: Always use `make dev` during development
5. **Clean often**: Run `make clean` when switching branches

---

**🎉 You're ready to go!**

Start the application with `make dev` and begin developing!

For detailed documentation, see `README.md` and `MAKEFILE_GUIDE.md`.

