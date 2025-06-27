# go-bitcoinkernel

A Go wrapper for Bitcoin Core's [`libbitcoinkernel`](https://github.com/bitcoin/bitcoin/pull/30595) library.

## ⚠️ Work in Progress

This library is currently under active development. Check [TODO.md](./TODO.md) for current status and missing wrappers.

## Overview

This repository consists of:

- **Bitcoin Core Source**: [Git subtree](./depend/bitcoin) containing Bitcoin Core source code with `libbitcoinkernel` C
  API
- **Kernel Package**: Safe, idiomatic Go interfaces with integrated CGO bindings that manage memory and provide error handling
- **Utils Package**: Helper functions and utilities built on the kernel package wrappers for common operations

## Installation and Usage

Since this library includes native C++ dependencies that must be compiled from source, it cannot be installed directly
via `go get` (at least for now). Follow these steps:

### Step 1: Clone the Repository

```bash
git clone https://github.com/stringintech/go-bitcoinkernel.git
cd go-bitcoinkernel
```

### Step 2: Build the Native Library

```bash
make build-kernel
```

This command will:

- Configure Bitcoin Core's CMake build system
- Build only the `libbitcoinkernel` library
- Use parallel compilation for faster builds

Refer to Bitcoin Core's build documentation to for the minimum requirements to compile `libbitcoinkernel` from source:
([Unix](./depend/bitcoin/doc/build-unix.md),
[macOS](./depend/bitcoin/doc/build-osx.md),
[Windows](./depend/bitcoin/doc/build-windows.md))

### Step 3: Run Tests

```bash
make test
```

This ensures that both the native library and Go bindings are working correctly.

The tests also include examples demonstrating how to use different components. For example, see:
- [`chainstate_manager_test.go`](./kernel/chainstate_manager_test.go)
- [`logger_test.go`](./utils/logger_test.go)

### Step 4: Use in Your Project

In your Go project directory, add a replace directive to point to your local copy:

```bash
# Initialize your Go module (if not already done)
go mod init your-project-name

# Add replace directive to use local go-bitcoinkernel
go mod edit -replace github.com/stringintech/go-bitcoinkernel=/path/to/go-bitcoinkernel

# Add the dependency
go get github.com/stringintech/go-bitcoinkernel/kernel
```

Your `go.mod` file should look like this:

```go
module your-project-name

go 1.23.3

require github.com/stringintech/go-bitcoinkernel/kernel v0.0.0-00010101000000-000000000000

replace github.com/stringintech/go-bitcoinkernel => /path/to/go-bitcoinkernel
```

## Important Notes

### Memory Management

The library handles memory management automatically through Go's finalizers, but it's highly recommended to explicitly
call `Destroy()` methods when you're done with objects to free resources immediately.

### Runtime Dependencies

Your Go application will have a runtime dependency on the shared `libbitcoinkernel` library produced by `make build-kernel` in `/path/to/go-bitcoinkernel/depend/bitcoin/build`. Do not delete or move these built library files as your application needs them to run.