.PHONY: build run clean client-install client-start client-build client-test client-test-coverage client-test-coverage-detail client-open-coverage server-test-coverage server-test-coverage-html server-open-coverage server-test-coverage-summary server-test-ginkgo server-test-ginkgo-verbose server-test-ginkgo-focus server-coverage-badge

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
	rm -rf client/coverage/
	rm -rf server/coverage/
	rm -f coverage.out

# Download dependencies
deps:
	go mod download

# Run tests
test:
	go test -v ./...

# Run server tests with coverage
server-test-coverage:
	go test -v -coverprofile=coverage.out ./server/...

# Generate HTML coverage report for server
server-test-coverage-html:
	go test -v -coverprofile=coverage.out ./server/...
	go tool cover -html=coverage.out -o server/coverage/coverage.html

# Show a summary of server code coverage
server-test-coverage-summary:
	go test -coverprofile=coverage.out ./server/...
	go tool cover -func=coverage.out

# Run server tests with Ginkgo and coverage
server-test-ginkgo:
	mkdir -p server/coverage
	cd server && $(shell go env GOPATH)/bin/ginkgo --cover --coverprofile=coverage.out ./...
	go tool cover -html=server/coverage.out -o server/coverage/coverage.html
	go tool cover -func=server/coverage.out

# Run server tests with Ginkgo, verbose output and coverage
server-test-ginkgo-verbose:
	mkdir -p server/coverage
	cd server && $(shell go env GOPATH)/bin/ginkgo --v --cover --coverprofile=coverage.out ./...
	go tool cover -html=server/coverage.out -o server/coverage/coverage.html
	go tool cover -func=server/coverage.out

# Run server tests with Ginkgo, focusing on specific tests or packages
# Usage: make server-test-ginkgo-focus FOCUS="TestName"
server-test-ginkgo-focus:
	mkdir -p server/coverage
	cd server && $(shell go env GOPATH)/bin/ginkgo --focus="$(FOCUS)" --cover --coverprofile=coverage.out ./...
	go tool cover -html=server/coverage.out -o server/coverage/coverage.html
	go tool cover -func=server/coverage.out

# Generate a coverage badge for the README
server-coverage-badge:
	mkdir -p server/coverage
	cd server && $(shell go env GOPATH)/bin/ginkgo --cover --coverprofile=coverage.out ./...
	go tool cover -func=server/coverage.out | grep total: | awk '{print $$3}' > server/coverage/coverage.txt
	@echo "Coverage badge data generated in server/coverage/coverage.txt"
	@echo "Add the following to your README.md:"
	@echo "![Coverage](https://img.shields.io/badge/coverage-$$(cat server/coverage/coverage.txt)-brightgreen)"

# Open the server coverage report in the default browser
server-open-coverage:
	open server/coverage/coverage.html

# Run client tests (excluding mock files by default)
client-test:
	cd client && npm test -- --testPathIgnorePatterns=mocks

# Run client tests with coverage report
client-test-coverage:
	cd client && npm test -- --testPathIgnorePatterns=mocks --coverage

# Run client tests with detailed coverage report
client-test-coverage-detail:
	cd client && npm test -- --testPathIgnorePatterns=mocks --coverage --coverageReporters="text" --coverageReporters="text-summary" --coverageReporters="html"

# Open the coverage report in the default browser
client-open-coverage:
	open client/coverage/index.html

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
