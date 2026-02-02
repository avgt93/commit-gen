#!/bin/bash

# Build release script for commit-gen
# Creates cross-platform binaries

set -e

VERSION="0.1.0"
PROJECT_NAME="commit-gen"
BUILD_DIR="dist"

# Colors for output
RED='\033[0;31m'
GREEN='\033[0;32m'
BLUE='\033[0;34m'
NC='\033[0m' # No Color

echo -e "${BLUE}Building commit-gen v${VERSION}...${NC}\n"

# Clean old builds
rm -rf "$BUILD_DIR"
mkdir -p "$BUILD_DIR"

# Define build targets
declare -A TARGETS=(
    [linux/amd64]="commit-gen-linux-amd64"
    [linux/arm64]="commit-gen-linux-arm64"
    [darwin/amd64]="commit-gen-darwin-amd64"
    [darwin/arm64]="commit-gen-darwin-arm64"
    [windows/amd64]="commit-gen-windows-amd64.exe"
)

# Build each target
for platform in "${!TARGETS[@]}"; do
    IFS='/' read -r GOOS GOARCH <<< "$platform"
    BINARY="${TARGETS[$platform]}"
    
    echo -e "${BLUE}Building for $platform...${NC}"
    
    GOOS="$GOOS" GOARCH="$GOARCH" go build \
        -ldflags "-X main.version=$VERSION" \
        -o "$BUILD_DIR/$BINARY" \
        ./cmd/commit-gen
    
    if [ -f "$BUILD_DIR/$BINARY" ]; then
        SIZE=$(du -h "$BUILD_DIR/$BINARY" | cut -f1)
        echo -e "${GREEN}✓ Built: $BINARY ($SIZE)${NC}\n"
    else
        echo -e "${RED}✗ Failed to build: $BINARY${NC}\n"
        exit 1
    fi
done

echo -e "${GREEN}All builds completed successfully!${NC}\n"
echo "Release files in: $BUILD_DIR/"
ls -lh "$BUILD_DIR"
