#!/bin/sh
set -e

# defaults
OS="${OS:-darwin}"
ARCH="${ARCH:-amd64}"
BINARY_NAME="characters-pipeline"
BUILD_DIR="build"
ENTRYPOINT="./cmd/main.go"

usage() {
  echo "Usage: $0 [build|test|clean]"
  echo "Environment variables:"
  echo "  OS   - target operating system (default: darwin)"
  echo "  ARCH - target architecture (default: amd64)"
}

build() {
  echo "Building for GOOS=$OS, GOARCH=$ARCH..."
  mkdir -p "$BUILD_DIR"
  GOOS="$OS" GOARCH="$ARCH" go build -o "$BUILD_DIR/${BINARY_NAME}-$OS-$ARCH" "$ENTRYPOINT"
  echo "Binary created at $BUILD_DIR/${BINARY_NAME}-$OS-$ARCH"
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
