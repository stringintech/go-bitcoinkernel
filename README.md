# go-bitcoinkernel

A Go wrapper for Bitcoin Core's [`libbitcoinkernel`](https://github.com/bitcoin/bitcoin/pull/30595) library.

## Overview

This repository consists of:
- **Bitcoin Core Source**: Git subtree containing Bitcoin Core source code with libbitcoinkernel C API
- **CGO Bindings**: C wrapper functions that interface with the libbitcoinkernel C API
- **Go API**: Safe, idiomatic Go interfaces that manage memory and provide error handling

## Installation and Usage

Since this library includes native C++ dependencies that must be compiled from source, it cannot be installed directly via `go get` (at least for now). Follow these steps:

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

### Step 4: Use in Your Project

In your Go project directory, add a replace directive to point to your local copy:

```bash
# Initialize your Go module (if not already done)
go mod init your-project-name

# Add replace directive to use local go-bitcoinkernel
go mod edit -replace github.com/stringintech/go-bitcoinkernel=../path/to/go-bitcoinkernel

# Add the dependency
go get github.com/stringintech/go-bitcoinkernel/kernel
```

Your `go.mod` file should look like this:
```go
module your-project-name

go 1.23.3

require github.com/stringintech/go-bitcoinkernel/kernel v0.0.0-00010101000000-000000000000

replace github.com/stringintech/go-bitcoinkernel => ../path/to/go-bitcoinkernel
```

## Example Usage

```go
package main

import (
    "fmt"
    "log"
    
    "github.com/stringintech/go-bitcoinkernel/kernel"
)

func main() {
    // Create a new context with mainnet parameters
    ctx, err := kernel.NewContext(kernel.ChainParametersMainnet())
    if err != nil {
        log.Fatal(err)
    }
    defer ctx.Destroy()
    
    // Create a chainstate manager
    manager, err := kernel.NewChainstateManager(ctx)
    if err != nil {
        log.Fatal(err)
    }
    defer manager.Destroy()
    
    fmt.Println("Bitcoin kernel initialized successfully!")
}
``` 

## Important Notes

### Memory Management
The library handles memory management automatically through Go's finalizers, but it's highly recommended to explicitly call `Destroy()` methods when you're done with objects to free resources immediately.