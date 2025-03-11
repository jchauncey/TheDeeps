.PHONY: build run clean client-install client-start client-build

# Build the server
build:
	go build -o bin/server ./server

# Run the server
run:
	cd server && go run .

# Clean build artifacts
clean:
	rm -rf bin/
	rm -rf client/build/

# Download dependencies
deps:
	go mod download

# Run tests
test:
	go test -v ./...

# Build and run the server
dev: build
	./bin/server

# Client commands
client-install:
	cd client && npm install

client-start:
	cd client && npm start

client-build:
	cd client && npm run build

# Run both server and client (requires tmux)
dev-full:
	tmux new-session -d -s thedeeps 'make run'
	tmux split-window -h -t thedeeps 'make client-start'
	tmux -2 attach-session -t thedeeps
