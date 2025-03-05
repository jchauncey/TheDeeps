# Binary name
BINARY_NAME=TheDeeps

# Build directory
BUILD_DIR=bin

# Go build command
GO_BUILD=go build

# Ensure build directory exists
$(BUILD_DIR):
	mkdir -p $(BUILD_DIR)

# Build the application
build: $(BUILD_DIR)
	$(GO_BUILD) -o $(BUILD_DIR)/$(BINARY_NAME)

# Clean build artifacts
clean:
	rm -rf $(BUILD_DIR)

.PHONY: build clean 