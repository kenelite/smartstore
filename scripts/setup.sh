#!/bin/bash

set -e

# Colors
RED='\033[0;31m'
GREEN='\033[0;32m'
YELLOW='\033[0;33m'
BLUE='\033[0;34m'
NC='\033[0m' # No Color

echo -e "${BLUE}==================================${NC}"
echo -e "${BLUE}SmartStore Setup Script${NC}"
echo -e "${BLUE}==================================${NC}\n"

# Check Go installation
echo -e "${BLUE}Checking Go installation...${NC}"
if ! command -v go &> /dev/null; then
    echo -e "${RED}✗ Go is not installed. Please install Go 1.23 or higher.${NC}"
    exit 1
fi

GO_VERSION=$(go version | awk '{print $3}' | sed 's/go//')
echo -e "${GREEN}✓ Go ${GO_VERSION} found${NC}\n"

# Check Docker installation (optional)
echo -e "${BLUE}Checking Docker installation...${NC}"
if command -v docker &> /dev/null; then
    echo -e "${GREEN}✓ Docker found${NC}"
    DOCKER_AVAILABLE=true
else
    echo -e "${YELLOW}⚠ Docker not found (optional)${NC}"
    DOCKER_AVAILABLE=false
fi
echo ""

# Download dependencies
echo -e "${BLUE}Downloading Go dependencies...${NC}"
make deps
echo -e "${GREEN}✓ Dependencies downloaded${NC}\n"

# Install development tools
echo -e "${BLUE}Installing development tools...${NC}"
read -p "Install development tools (golangci-lint, air, swag)? [Y/n] " -n 1 -r
echo
if [[ $REPLY =~ ^[Yy]$ ]] || [[ -z $REPLY ]]; then
    make install-tools
    echo -e "${GREEN}✓ Development tools installed${NC}"
else
    echo -e "${YELLOW}⚠ Skipped development tools installation${NC}"
fi
echo ""

# Setup Docker environment
if [ "$DOCKER_AVAILABLE" = true ]; then
    echo -e "${BLUE}Docker Setup${NC}"
    read -p "Start PostgreSQL, Redis, and MinIO with docker-compose? [Y/n] " -n 1 -r
    echo
    if [[ $REPLY =~ ^[Yy]$ ]] || [[ -z $REPLY ]]; then
        echo -e "${BLUE}Starting services...${NC}"
        make docker-compose-up
        echo -e "${GREEN}✓ Services started${NC}"
        echo -e "${YELLOW}Waiting for services to be ready...${NC}"
        sleep 5
        
        # Run database migrations
        echo -e "${BLUE}Running database migrations...${NC}"
        make db-migrate
        echo -e "${GREEN}✓ Database initialized${NC}"
    else
        echo -e "${YELLOW}⚠ Skipped Docker setup${NC}"
        echo -e "${YELLOW}Note: You'll need to setup PostgreSQL and Redis manually${NC}"
    fi
else
    echo -e "${YELLOW}Manual Setup Required:${NC}"
    echo -e "  1. Install and start PostgreSQL"
    echo -e "  2. Install and start Redis"
    echo -e "  3. Run: make db-migrate"
fi
echo ""

# Build the application
echo -e "${BLUE}Building application...${NC}"
make build
echo -e "${GREEN}✓ Application built successfully${NC}\n"

# Summary
echo -e "${GREEN}==================================${NC}"
echo -e "${GREEN}Setup Complete!${NC}"
echo -e "${GREEN}==================================${NC}\n"

echo -e "${BLUE}Next Steps:${NC}"
echo -e "  1. Edit config.yaml with your settings"
echo -e "  2. Run the application:"
echo -e "     ${YELLOW}make run${NC}     (run from source)"
echo -e "     ${YELLOW}make run-bin${NC} (run compiled binary)"
echo -e "     ${YELLOW}make dev${NC}     (run with hot-reload)"
echo ""

if [ "$DOCKER_AVAILABLE" = true ]; then
    echo -e "${BLUE}Docker Services:${NC}"
    echo -e "  • PostgreSQL: localhost:5432"
    echo -e "  • Redis: localhost:6379"
    echo -e "  • MinIO: localhost:9000 (API), localhost:9001 (Console)"
    echo -e "    Credentials: minioadmin / minioadmin"
    echo ""
fi

echo -e "${BLUE}Useful Commands:${NC}"
echo -e "  ${YELLOW}make help${NC}           - Show all available commands"
echo -e "  ${YELLOW}make test${NC}           - Run tests"
echo -e "  ${YELLOW}make lint${NC}           - Run linter"
echo -e "  ${YELLOW}make docker-compose-logs${NC} - View service logs"
echo ""

echo -e "${GREEN}Happy coding! 🚀${NC}\n"

