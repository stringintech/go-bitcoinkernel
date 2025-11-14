# Makefile for go-bitcoinkernel

.PHONY: all build-kernel build test clean help

all: build-kernel test

build-kernel:
	cd depend/bitcoin && \
	cmake -B build \
		-DCMAKE_BUILD_TYPE=RelWithDebInfo \
		-DBUILD_SHARED_LIBS=OFF \
		-DBUILD_KERNEL_LIB=ON \
		-DBUILD_KERNEL_TEST=OFF \
		-DBUILD_TESTS=OFF \
		-DBUILD_TX=OFF \
		-DBUILD_WALLET_TOOL=OFF \
		-DENABLE_WALLET=OFF \
		-DENABLE_EXTERNAL_SIGNER=OFF \
		-DBUILD_UTIL=OFF \
		-DBUILD_BITCOIN_BIN=OFF \
		-DBUILD_DAEMON=OFF \
		-DBUILD_UTIL_CHAINSTATE=OFF \
		-DBUILD_CLI=OFF \
		-DENABLE_IPC=OFF \
		-DCMAKE_INSTALL_LIBDIR=lib \
		-DCMAKE_INSTALL_PREFIX=$$PWD/install && \
	cmake --build build --config RelWithDebInfo --parallel$(if $(NUM_JOBS),=$(NUM_JOBS)) && \
	cmake --install build --config RelWithDebInfo

build:
	go build ./...

test:
	go test -v ./...

clean:
	rm -rf depend/bitcoin/build
	rm -rf depend/bitcoin/install
	go clean ./...
	go clean -testcache

lint:
	golangci-lint run ./...

deps:
	go install github.com/golangci/golangci-lint/v2/cmd/golangci-lint@v2.1.6

update-kernel:
	git subtree pull --prefix=depend/bitcoin https://github.com/bitcoin/bitcoin.git master --squash

help:
	@echo "Available targets:"
	@echo "  all			- Build kernel library and run tests (default)"
	@echo "  build-kernel		- Build Bitcoin kernel library"
	@echo "  build			- Compile Go code"
	@echo "  test        		- Run Go tests"
	@echo "  clean       		- Clean build artifacts"
	@echo "  lint        		- Lint Go code"
	@echo "  deps        		- Install development dependencies"
	@echo "  update-kernel  	- Update Bitcoin dependency using git subtree"
	@echo "  help        		- Show this help message"