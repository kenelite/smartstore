#!/bin/bash

set -e

# Colors
RED='\033[0;31m'
GREEN='\033[0;32m'
YELLOW='\033[0;33m'
BLUE='\033[0;34m'
NC='\033[0m' # No Color

echo -e "${BLUE}==================================${NC}"
echo -e "${BLUE}SmartStore Test Suite${NC}"
echo -e "${BLUE}==================================${NC}\n"

# Parse command line arguments
VERBOSE=false
COVERAGE=false
BENCH=false
SHORT=false

while [[ $# -gt 0 ]]; do
    case $1 in
        -v|--verbose)
            VERBOSE=true
            shift
            ;;
        -c|--coverage)
            COVERAGE=true
            shift
            ;;
        -b|--bench)
            BENCH=true
            shift
            ;;
        -s|--short)
            SHORT=true
            shift
            ;;
        -h|--help)
            echo "Usage: $0 [OPTIONS]"
            echo ""
            echo "Options:"
            echo "  -v, --verbose    Verbose output"
            echo "  -c, --coverage   Generate coverage report"
            echo "  -b, --bench      Run benchmarks"
            echo "  -s, --short      Run short tests only"
            echo "  -h, --help       Show this help message"
            exit 0
            ;;
        *)
            echo -e "${RED}Unknown option: $1${NC}"
            exit 1
            ;;
    esac
done

# Format code
echo -e "${BLUE}1. Formatting code...${NC}"
make fmt
echo -e "${GREEN}✓ Code formatted${NC}\n"

# Run go vet
echo -e "${BLUE}2. Running go vet...${NC}"
make vet
echo -e "${GREEN}✓ go vet passed${NC}\n"

# Run linter (if available)
echo -e "${BLUE}3. Running linter...${NC}"
if make lint 2>/dev/null; then
    echo -e "${GREEN}✓ Linter passed${NC}\n"
else
    echo -e "${YELLOW}⚠ Linter not available or found issues${NC}\n"
fi

# Run tests
if [ "$SHORT" = true ]; then
    echo -e "${BLUE}4. Running short tests...${NC}"
    make test-short
elif [ "$COVERAGE" = true ]; then
    echo -e "${BLUE}4. Running tests with coverage...${NC}"
    make test-coverage
    echo -e "${GREEN}✓ Coverage report generated: coverage.html${NC}\n"
elif [ "$VERBOSE" = true ]; then
    echo -e "${BLUE}4. Running tests (verbose)...${NC}"
    go test -v -race ./...
else
    echo -e "${BLUE}4. Running tests...${NC}"
    make test
fi
echo -e "${GREEN}✓ Tests passed${NC}\n"

# Run benchmarks
if [ "$BENCH" = true ]; then
    echo -e "${BLUE}5. Running benchmarks...${NC}"
    make test-bench
    echo -e "${GREEN}✓ Benchmarks completed${NC}\n"
fi

# Summary
echo -e "${GREEN}==================================${NC}"
echo -e "${GREEN}All checks passed! ✓${NC}"
echo -e "${GREEN}==================================${NC}\n"

