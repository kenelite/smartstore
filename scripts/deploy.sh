#!/bin/bash

set -e

# Colors
RED='\033[0;31m'
GREEN='\033[0;32m'
YELLOW='\033[0;33m'
BLUE='\033[0;34m'
NC='\033[0m' # No Color

echo -e "${BLUE}==================================${NC}"
echo -e "${BLUE}SmartStore Deployment Script${NC}"
echo -e "${BLUE}==================================${NC}\n"

# Check if we're on the right branch
CURRENT_BRANCH=$(git rev-parse --abbrev-ref HEAD)
echo -e "${BLUE}Current branch: ${YELLOW}${CURRENT_BRANCH}${NC}\n"

# Check for uncommitted changes
if [[ -n $(git status -s) ]]; then
    echo -e "${RED}✗ You have uncommitted changes!${NC}"
    git status -s
    echo ""
    read -p "Continue anyway? [y/N] " -n 1 -r
    echo
    if [[ ! $REPLY =~ ^[Yy]$ ]]; then
        exit 1
    fi
fi

# Get version tag
echo -e "${BLUE}Enter version tag (e.g., v1.0.0):${NC}"
read -r VERSION

if [[ -z "$VERSION" ]]; then
    echo -e "${RED}✗ Version tag is required${NC}"
    exit 1
fi

# Validate version format
if [[ ! $VERSION =~ ^v[0-9]+\.[0-9]+\.[0-9]+$ ]]; then
    echo -e "${YELLOW}⚠ Version should follow semver format (e.g., v1.0.0)${NC}"
    read -p "Continue anyway? [y/N] " -n 1 -r
    echo
    if [[ ! $REPLY =~ ^[Yy]$ ]]; then
        exit 1
    fi
fi

echo ""
echo -e "${BLUE}Deployment checklist:${NC}"
echo -e "  Version: ${GREEN}${VERSION}${NC}"
echo -e "  Branch: ${GREEN}${CURRENT_BRANCH}${NC}"
echo ""
read -p "Proceed with deployment? [Y/n] " -n 1 -r
echo
if [[ ! $REPLY =~ ^[Yy]$ ]] && [[ ! -z $REPLY ]]; then
    echo -e "${YELLOW}Deployment cancelled${NC}"
    exit 0
fi
echo ""

# Run tests
echo -e "${BLUE}1. Running tests...${NC}"
make release-check
echo -e "${GREEN}✓ All tests passed${NC}\n"

# Build release binaries
echo -e "${BLUE}2. Building release binaries...${NC}"
make release-build
echo -e "${GREEN}✓ Release binaries built${NC}\n"

# Build Docker image
echo -e "${BLUE}3. Building Docker image...${NC}"
make docker-build
docker tag smartstore-gateway:latest smartstore-gateway:${VERSION}
echo -e "${GREEN}✓ Docker image built and tagged${NC}\n"

# Create git tag
echo -e "${BLUE}4. Creating git tag...${NC}"
git tag -a ${VERSION} -m "Release ${VERSION}"
echo -e "${GREEN}✓ Git tag created: ${VERSION}${NC}\n"

# Summary
echo -e "${GREEN}==================================${NC}"
echo -e "${GREEN}Deployment preparation complete!${NC}"
echo -e "${GREEN}==================================${NC}\n"

echo -e "${BLUE}Next steps:${NC}"
echo -e "  1. Push the tag: ${YELLOW}git push origin ${VERSION}${NC}"
echo -e "  2. Push Docker image: ${YELLOW}docker push your-registry/smartstore-gateway:${VERSION}${NC}"
echo -e "  3. Deploy to your infrastructure"
echo ""

echo -e "${YELLOW}Note: This script only prepares the deployment.${NC}"
echo -e "${YELLOW}You need to manually push the tag and deploy the image.${NC}\n"

