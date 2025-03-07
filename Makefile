# Binary name
BINARY_NAME=TheDeeps

# Build directory
BUILD_DIR=bin

# Go build command
GO_BUILD=go build

# Server port
PORT=8080

# Ensure build directory exists
$(BUILD_DIR):
	mkdir -p $(BUILD_DIR)

# Build the server
build-server: $(BUILD_DIR)
	$(GO_BUILD) -o $(BUILD_DIR)/$(BINARY_NAME)

# Run the server
run-server: build-server
	$(BUILD_DIR)/$(BINARY_NAME)

# Build and run the server in one command
server: build-server run-server

# Run the server without rebuilding
start-server:
	$(BUILD_DIR)/$(BINARY_NAME)

# Install client dependencies
client-install:
	cd client && npm install

# Build the client
client-build: client-install
	cd client && npm run build

# Run the client in development mode
client-dev:
	cd client && npm run dev

# Build both client and server
build: build-server client-build

# Clean build artifacts
clean:
	rm -rf $(BUILD_DIR)
	rm -rf client/dist

# Run both client and server (in separate terminals)
# Note: This requires tmux or running in separate terminals
dev:
	@echo "Please run these commands in separate terminals:"
	@echo "make client-dev"
	@echo "make run-server"

.PHONY: build build-server run-server server start-server client-install client-build client-dev build clean dev 