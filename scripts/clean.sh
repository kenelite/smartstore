#!/bin/bash

# Colors
RED='\033[0;31m'
GREEN='\033[0;32m'
YELLOW='\033[0;33m'
BLUE='\033[0;34m'
NC='\033[0m' # No Color

echo -e "${BLUE}==================================${NC}"
echo -e "${BLUE}SmartStore Cleanup${NC}"
echo -e "${BLUE}==================================${NC}\n"

# Parse command line arguments
DEEP=false
DOCKER=false

while [[ $# -gt 0 ]]; do
    case $1 in
        --deep)
            DEEP=true
            shift
            ;;
        --docker)
            DOCKER=true
            shift
            ;;
        -h|--help)
            echo "Usage: $0 [OPTIONS]"
            echo ""
            echo "Options:"
            echo "  --deep     Deep clean (includes caches)"
            echo "  --docker   Stop and remove Docker containers"
            echo "  -h, --help Show this help message"
            exit 0
            ;;
        *)
            echo -e "${RED}Unknown option: $1${NC}"
            exit 1
            ;;
    esac
done

# Stop Docker services
if [ "$DOCKER" = true ]; then
    echo -e "${BLUE}Stopping Docker services...${NC}"
    if command -v docker-compose &> /dev/null || command -v docker &> /dev/null; then
        make docker-compose-down 2>/dev/null || true
        echo -e "${GREEN}✓ Docker services stopped${NC}\n"
    else
        echo -e "${YELLOW}⚠ Docker not available${NC}\n"
    fi
fi

# Clean build artifacts
echo -e "${BLUE}Cleaning build artifacts...${NC}"
make clean
echo -e "${GREEN}✓ Build artifacts cleaned${NC}\n"

# Deep clean
if [ "$DEEP" = true ]; then
    echo -e "${YELLOW}Performing deep clean...${NC}"
    make clean-all
    
    # Remove additional files
    echo -e "${BLUE}Removing additional files...${NC}"
    rm -f *.log
    rm -rf tmp/
    rm -rf .air/
    
    echo -e "${GREEN}✓ Deep clean completed${NC}\n"
fi

# Summary
echo -e "${GREEN}==================================${NC}"
echo -e "${GREEN}Cleanup completed!${NC}"
echo -e "${GREEN}==================================${NC}\n"

if [ "$DEEP" = false ]; then
    echo -e "${YELLOW}Tip: Use --deep for a more thorough cleanup${NC}"
fi

if [ "$DOCKER" = false ]; then
    echo -e "${YELLOW}Tip: Use --docker to stop Docker services${NC}"
fi
echo ""

