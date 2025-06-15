# Makefile for go-bitcoinkernel

.PHONY: all build-kernel test clean help

# Default target
all: build-kernel test

# Build the Bitcoin kernel library
build-kernel:
	@echo "Building Bitcoin kernel library..."
	cd depend/bitcoin && \
	cmake -B build -DBUILD_SHARED_LIBS=ON -DBUILD_KERNEL_LIB=ON -DBUILD_BENCH=OFF -DBUILD_CLI=OFF -DBUILD_DAEMON=OFF -DBUILD_FOR_FUZZING=OFF -DBUILD_FUZZ_BINARY=OFF -DBUILD_GUI=OFF -DBUILD_KERNEL_TEST=OFF -DBUILD_TESTS=OFF -DBUILD_TX=OFF -DBUILD_UTIL=OFF -DBUILD_UTIL_CHAINSTATE=OFF -DBUILD_WALLET_TOOL=OFF -DENABLE_WALLET=OFF && \
	cmake --build build --target bitcoinkernel -j $(shell nproc 2>/dev/null || echo 4)

# Test the Go bindings
test: build-kernel
	@echo "Running Go tests..."
	cd kernel && go test -v

# Build Go package (compilation check)
build: build-kernel
	@echo "Building Go package..."
	cd kernel && go build -v

# Clean build artifacts
clean:
	@echo "Cleaning build artifacts..."
	rm -rf depend/bitcoin/build
	cd kernel && go clean

# Lint Go code (requires golangci-lint)
lint:
	@echo "Linting Go code..."
	golangci-lint run ./...

# Install development dependencies
deps:
	@echo "Installing development dependencies..."
	go install github.com/golangci/golangci-lint/cmd/golangci-lint@latest

# Show help
help:
	@echo "Available targets:"
	@echo "  all		- Build kernel library and run tests (default)"
	@echo "  build-kernel	- Build Bitcoin kernel library"
	@echo "  test        	- Run Go tests"
	@echo "  build       	- Build Go package (compilation check)"
	@echo "  clean       	- Clean build artifacts"
	@echo "  lint        	- Lint Go code"
	@echo "  deps        	- Install development dependencies"
	@echo "  help        	- Show this help message"