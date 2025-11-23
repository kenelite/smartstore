#!/bin/bash

# Colors
RED='\033[0;31m'
GREEN='\033[0;32m'
YELLOW='\033[0;33m'
BLUE='\033[0;34m'
NC='\033[0m' # No Color

echo -e "${BLUE}==================================${NC}"
echo -e "${BLUE}SmartStore Setup Verification${NC}"
echo -e "${BLUE}==================================${NC}\n"

ERRORS=0
WARNINGS=0

# Function to check file exists
check_file() {
    if [ -f "$1" ]; then
        echo -e "${GREEN}✓${NC} Found: $1"
    else
        echo -e "${RED}✗${NC} Missing: $1"
        ((ERRORS++))
    fi
}

# Function to check directory exists
check_dir() {
    if [ -d "$1" ]; then
        echo -e "${GREEN}✓${NC} Found: $1/"
    else
        echo -e "${RED}✗${NC} Missing: $1/"
        ((ERRORS++))
    fi
}

# Function to check command
check_command() {
    if command -v $1 &> /dev/null; then
        VERSION=$($1 version 2>&1 | head -1 || echo "installed")
        echo -e "${GREEN}✓${NC} $1: ${VERSION}"
    else
        echo -e "${YELLOW}⚠${NC}  $1: not found (optional)"
        ((WARNINGS++))
    fi
}

echo -e "${BLUE}1. Checking required files...${NC}\n"

# Core files
check_file "Makefile"
check_file "go.mod"
check_file "go.sum"
check_file "config.yaml"
check_file "README.md"

echo ""
echo -e "${BLUE}2. Checking Docker files...${NC}\n"

check_file "Dockerfile"
check_file "docker-compose.yaml"
check_file ".dockerignore"

echo ""
echo -e "${BLUE}3. Checking configuration files...${NC}\n"

check_file ".air.toml"
check_file ".golangci.yml"
check_file ".gitignore"

echo ""
echo -e "${BLUE}4. Checking documentation...${NC}\n"

check_file "QUICK_START.md"
check_file "MAKEFILE_GUIDE.md"
check_file "CONTRIBUTING.md"
check_file "PROJECT_SETUP_SUMMARY.md"

echo ""
echo -e "${BLUE}5. Checking scripts...${NC}\n"

check_file "scripts/setup.sh"
check_file "scripts/test.sh"
check_file "scripts/clean.sh"
check_file "scripts/deploy.sh"
check_file "scripts/verify-setup.sh"

echo ""
echo -e "${BLUE}6. Checking CI/CD workflows...${NC}\n"

check_dir ".github/workflows"
check_file ".github/workflows/ci.yml"
check_file ".github/workflows/release.yml"

echo ""
echo -e "${BLUE}7. Checking project structure...${NC}\n"

check_dir "cmd/gateway"
check_dir "internal"
check_dir "db"
check_file "db/schema.sql"

echo ""
echo -e "${BLUE}8. Checking required commands...${NC}\n"

check_command "go"
check_command "make"
check_command "git"
check_command "docker"
check_command "docker-compose"

echo ""
echo -e "${BLUE}9. Testing Makefile commands...${NC}\n"

# Test Makefile syntax
if make -n help > /dev/null 2>&1; then
    echo -e "${GREEN}✓${NC} Makefile syntax is valid"
else
    echo -e "${RED}✗${NC} Makefile has syntax errors"
    ((ERRORS++))
fi

# Count available targets
TARGET_COUNT=$(make -qp | awk -F':' '/^[a-zA-Z0-9][^$#\/\t=]*:([^=]|$)/ {split($1,A,/ /);for(i in A)print A[i]}' | grep -v '^\.PHONY' | sort -u | wc -l)
echo -e "${GREEN}✓${NC} Found $TARGET_COUNT Makefile targets"

echo ""
echo -e "${BLUE}10. Checking Go dependencies...${NC}\n"

if go mod verify > /dev/null 2>&1; then
    echo -e "${GREEN}✓${NC} Go modules verified"
else
    echo -e "${YELLOW}⚠${NC}  Go modules need verification (run: make deps)"
    ((WARNINGS++))
fi

# Summary
echo ""
echo -e "${BLUE}==================================${NC}"
echo -e "${BLUE}Verification Summary${NC}"
echo -e "${BLUE}==================================${NC}\n"

if [ $ERRORS -eq 0 ]; then
    echo -e "${GREEN}✓ All critical checks passed!${NC}"
else
    echo -e "${RED}✗ Found $ERRORS error(s)${NC}"
fi

if [ $WARNINGS -gt 0 ]; then
    echo -e "${YELLOW}⚠ Found $WARNINGS warning(s) (optional tools)${NC}"
fi

echo ""

if [ $ERRORS -eq 0 ]; then
    echo -e "${GREEN}Your SmartStore project is ready! 🚀${NC}\n"
    echo -e "${BLUE}Quick start:${NC}"
    echo -e "  1. Run: ${YELLOW}make help${NC} to see all commands"
    echo -e "  2. Run: ${YELLOW}make setup${NC} for automated setup"
    echo -e "  3. Run: ${YELLOW}make dev${NC} to start development"
    echo ""
    exit 0
else
    echo -e "${RED}Please fix the errors above before proceeding.${NC}\n"
    exit 1
fi

