#!/bin/sh
set -e

# defaults resolved from the host system
OS="${OS:-$(uname | tr '[:upper:]' '[:lower:]')}"
ARCH="${ARCH:-$(uname -m)}"

# normalize ARCH
case "$ARCH" in
  x86_64) ARCH="amd64" ;;
  aarch64) ARCH="arm64" ;;
esac

BINARY_NAME="characters-pipeline"
BUILD_DIR="build"
ENTRYPOINT="./cmd/main.go"

usage() {
  echo "Usage: $0 [build|test|clean]"
  echo "Environment variables:"
  echo "  OS   - target operating system (default: auto-detected)"
  echo "  ARCH - target architecture (default: auto-detected)"
}

build() {
  echo "Building for GOOS=$OS, GOARCH=$ARCH..."
  mkdir -p "$BUILD_DIR"
  go mod tidy && GOOS="$OS" GOARCH="$ARCH" go build -o "$BUILD_DIR/${BINARY_NAME}" "$ENTRYPOINT"
  echo "Binary created at $BUILD_DIR/${BINARY_NAME}"
}

test() {
  echo "Running tests..."
  go mod tidy && go test -v ./pkg/...
}

clean() {
  echo "Cleaning build artifacts..."
  rm -rf "$BUILD_DIR"
}

cmd="$1"

case "$cmd" in
  build)
    build
    ;;
  test)
    test
    ;;
  clean)
    clean
    ;;
  *)
    usage
    exit 1
    ;;
esac
