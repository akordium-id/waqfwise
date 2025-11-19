#!/bin/bash

# Build script for WaqfWise Community Edition
# Licensed under AGPL v3

set -e

echo "Building WaqfWise Community Edition..."

# Set build tags for community edition
BUILD_TAGS="community"

# Build metadata
VERSION=$(git describe --tags --always --dirty 2>/dev/null || echo "dev")
BUILD_TIME=$(date -u '+%Y-%m-%d_%H:%M:%S')
GIT_COMMIT=$(git rev-parse --short HEAD 2>/dev/null || echo "unknown")

# Build flags
LDFLAGS="-X main.version=${VERSION} -X main.buildTime=${BUILD_TIME} -X main.gitCommit=${GIT_COMMIT}"

# Output binary path
OUTPUT="./bin/waqfwise-community"

# Create bin directory if it doesn't exist
mkdir -p ./bin

# Build the community binary
echo "Building with tags: ${BUILD_TAGS}"
go build -tags="${BUILD_TAGS}" -ldflags="${LDFLAGS}" -o "${OUTPUT}" ./cmd/waqfwise-community

echo "Build complete: ${OUTPUT}"
echo "Version: ${VERSION}"
echo "Build time: ${BUILD_TIME}"
echo "Git commit: ${GIT_COMMIT}"
