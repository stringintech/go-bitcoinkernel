# Conformance Handler Makefile

# Test suite configuration
TEST_VERSION := 0.0.2
TEST_REPO := stringintech/kernel-bindings-tests
TEST_DIR := .conformance-tests

# Platform detection
UNAME_S := $(shell uname -s)
UNAME_M := $(shell uname -m)

ifeq ($(UNAME_S),Darwin)
	ifeq ($(UNAME_M),arm64)
		PLATFORM := darwin_arm64
	else
		PLATFORM := darwin_amd64
	endif
else ifeq ($(UNAME_S),Linux)
	ifeq ($(UNAME_M),x86_64)
		PLATFORM := linux_amd64
	else ifeq ($(UNAME_M),aarch64)
		PLATFORM := linux_arm64
	else
		PLATFORM := linux_amd64
	endif
else
	$(error Unsupported platform: $(UNAME_S) $(UNAME_M))
endif

# Binary names
TEST_RUNNER := $(TEST_DIR)/runner
HANDLER_BIN := handler

.PHONY: all build download-tests test clean help

all: build test

help:
	@echo "Conformance Handler Makefile"
	@echo ""
	@echo "Targets:"
	@echo "  build           - Build the conformance handler binary"
	@echo "  download-tests  - Download the test suite for your platform"
	@echo "  test            - Run conformance tests against the handler"
	@echo "  clean           - Remove built binaries and downloaded tests"
	@echo "  help            - Show this help message"
	@echo ""
	@echo "Configuration:"
	@echo "  Test Version: $(TEST_VERSION)"
	@echo "  Platform:     $(PLATFORM)"

build:
	@echo "Building conformance handler..."
	go build -o $(HANDLER_BIN) .

download-tests:
	@echo "Downloading test suite $(TEST_VERSION) for $(PLATFORM)..."
	@mkdir -p $(TEST_DIR)
	$(eval DOWNLOAD_URL := https://github.com/$(TEST_REPO)/releases/download/v$(TEST_VERSION)/kernel-bindings-tests_$(TEST_VERSION)_$(PLATFORM).tar.gz)
	@echo "URL: $(DOWNLOAD_URL)"
	@curl -L -o $(TEST_DIR)/test-runner.tar.gz "$(DOWNLOAD_URL)"
	@echo "Extracting test runner..."
	@tar -xzf $(TEST_DIR)/test-runner.tar.gz -C $(TEST_DIR)
	@chmod +x $(TEST_RUNNER)
	@rm $(TEST_DIR)/test-runner.tar.gz
	@echo "Test runner downloaded to $(TEST_RUNNER)"

test: build
	@if [ ! -f "$(TEST_RUNNER)" ]; then \
		echo "Test runner not found. Downloading..."; \
		$(MAKE) download-tests; \
	fi
	@echo "Running conformance tests..."
	$(TEST_RUNNER) --handler ./$(HANDLER_BIN)

clean:
	@echo "Cleaning up..."
	rm -f $(HANDLER_BIN)
	rm -rf $(TEST_DIR)