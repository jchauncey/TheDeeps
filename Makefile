.PHONY: build run clean

# Build the server
build:
	go build -o bin/server ./server

# Run the server
run:
	go run ./server

# Clean build artifacts
clean:
	rm -rf bin/

# Download dependencies
deps:
	go mod download

# Run tests
test:
	go test -v ./...

# Build and run the server
dev: build
	./bin/server
