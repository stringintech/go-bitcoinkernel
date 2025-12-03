# Conformance Handler

This binary implements the JSON protocol required by the [kernel-bindings-spec](https://github.com/stringintech/kernel-bindings-spec) conformance testing framework.

## Purpose

The conformance handler acts as a bridge between the test runner and the Go Bitcoin Kernel bindings. It:

- Reads test requests from stdin (JSON protocol)
- Executes operations using the Go binding API
- Returns responses to stdout (JSON protocol)

## Testing

This handler is designed to work with the conformance test suite. The easiest way to run tests is using the Makefile:

```bash
# Run conformance tests (builds handler and downloads test runner automatically)
make test

# Or manually build and run
make build
make download-tests
./.conformance-tests/runner --handler ./handler
```

The test suite is automatically downloaded for your platform (darwin_arm64, darwin_amd64, linux_amd64, or linux_arm64).

## Pinned Test Version

This handler is compatible with:
- Test Suite Version: `0.0.3-alpha.3`
- Test Repository: [stringintech/kernel-bindings-tests](https://github.com/stringintech/kernel-bindings-tests)