# Makefile for go-bitcoinkernel

.PHONY: all build-kernel test clean help

all: build-kernel test

build-kernel:
	cd depend/bitcoin && \
	cmake -B build -DBUILD_SHARED_LIBS=ON -DBUILD_KERNEL_LIB=ON -DBUILD_BENCH=OFF -DBUILD_CLI=OFF -DBUILD_DAEMON=OFF -DBUILD_FOR_FUZZING=OFF -DBUILD_FUZZ_BINARY=OFF -DBUILD_GUI=OFF -DBUILD_KERNEL_TEST=OFF -DBUILD_TESTS=OFF -DBUILD_TX=OFF -DBUILD_UTIL=OFF -DBUILD_UTIL_CHAINSTATE=OFF -DBUILD_WALLET_TOOL=OFF -DENABLE_WALLET=OFF && \
	cmake --build build --target bitcoinkernel -j $(shell nproc 2>/dev/null || echo 4)

test: build-kernel
	go test -v ./...

clean:
	rm -rf depend/bitcoin/build
	go clean ./...

lint:
	golangci-lint run ./...

deps:
	go install github.com/golangci/golangci-lint/cmd/golangci-lint@latest

help:
	@echo "Available targets:"
	@echo "  all		- Build kernel library and run tests (default)"
	@echo "  build-kernel	- Build Bitcoin kernel library"
	@echo "  test        	- Run Go tests"
	@echo "  clean       	- Clean build artifacts"
	@echo "  lint        	- Lint Go code"
	@echo "  deps        	- Install development dependencies"
	@echo "  help        	- Show this help message"