# Name of the output binary
BINARY_NAME=fitanonymize

# Cross compilation targets
BUILD_DIR=build
TARGETS=linux window darwin
ARCH=amd64

# Go commands
GO=go
GOFMT=gofmt
GOFLAGS=-ldflags="-s -w"

.PHONY: all clean format build

# Default target
all: build

# Format the code
format:
	$(GOFMT) -w .

# Build the executable for the current platform
build:
	$(GO) build $(GOFLAGS) -o $(BUILD_DIR)/$(BINARY_NAME) main.go

# Cross compile for multiple targets
cross-compile:
	@for target in $(TARGETS); do \
		GOOS=$$target GOARCH=$(ARCH) $(GO) build $(GOFLAGS) -o $(BUILD_DIR)/$(BINARY_NAME)_$$target main.go; \
	done

# Cross compile specifically for Windows
compile-windows:
	GOOS=windows GOARCH=$(ARCH) $(GO) build $(GOFLAGS) -o $(BUILD_DIR)/$(BINARY_NAME).exe main.go

# Clean up build artifacts
clean:
	rm -rf $(BUILD_DIR)
