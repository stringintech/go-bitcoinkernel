# go-bitcoinkernel

[![Go Reference](https://pkg.go.dev/badge/github.com/stringintech/go-bitcoinkernel.svg)](https://pkg.go.dev/github.com/stringintech/go-bitcoinkernel@v0.1.0-alpha.2)
[![Go Report Card](https://goreportcard.com/badge/github.com/stringintech/go-bitcoinkernel)](https://goreportcard.com/report/github.com/stringintech/go-bitcoinkernel)
[![License](https://img.shields.io/github/license/stringintech/go-bitcoinkernel)](https://github.com/stringintech/go-bitcoinkernel/blob/main/LICENSE)
[![GitHub Actions](https://github.com/stringintech/go-bitcoinkernel/actions/workflows/ci.yml/badge.svg)](https://github.com/stringintech/go-bitcoinkernel/actions)
[![GitHub release](https://img.shields.io/github/v/release/stringintech/go-bitcoinkernel?include_prereleases)](https://github.com/stringintech/go-bitcoinkernel/releases)

A Go wrapper for Bitcoin Core's [`libbitcoinkernel`](https://github.com/bitcoin/bitcoin/pull/30595) library.

## ⚠️ Work in Progress

This library is experimental and not production-ready. The underlying C API is subject to change, which may cause breaking changes in this wrapper. The wrapper itself may also change based on user feedback, which is welcome and appreciated. Feel free to open issues to help make this wrapper more useful for everyone.

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

This command will configure Bitcoin Core's CMake build system and statically compile the `libbitcoinkernel` library. Building requires CMake, a C++ compiler, Boost library, and Make. For platform-specific setup instructions, consult Bitcoin Core's build documentation:
([Unix](./depend/bitcoin/doc/build-unix.md),
[macOS](./depend/bitcoin/doc/build-osx.md),
[Windows](./depend/bitcoin/doc/build-windows.md))

**Additional notes:**
- **For Windows:** MinGW toolchain is required. CGo uses GCC on Windows, so static linking requires the library to be compiled with GCC/MinGW rather than MSVC for ABI compatibility.
- **For Nix users:** Run `nix develop` to automatically get all the build and runtime dependencies, then run `make build-kernel` and continue with the remaining steps in the Nix dev shell.

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

The library handles memory management automatically through Go's finalizers (see [common.go](./kernel/common.go)), but it's highly recommended to explicitly
call `Destroy()` methods when you're done with owned objects to free resources immediately.

### Error Handling

The library uses structured error types for better error handling (see [errors.go](./kernel/errors.go)).