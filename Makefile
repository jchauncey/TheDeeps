.PHONY: build run clean client-install client-start client-build client-test client-test-coverage client-test-coverage-detail client-open-coverage server-test-coverage server-test-coverage-html server-open-coverage server-test-coverage-summary server-test-ginkgo server-test-ginkgo-verbose server-test-ginkgo-focus server-coverage-badge server-test-ginkgo-coverage client-test-e2e client-test-e2e-ui client-test-e2e-headed client-test-e2e-debug client-test-e2e-with-server client-test-e2e-file-with-server client-test-e2e-headed-with-server client-test-e2e-file

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

# Run server tests with Ginkgo
server-test-ginkgo:
	cd server && $(shell go env GOPATH)/bin/ginkgo ./...

# Run server tests with Ginkgo and coverage
server-test-ginkgo-coverage:
	mkdir -p server/coverage
	cd server && $(shell go env GOPATH)/bin/ginkgo --cover --coverprofile=coverage.out ./...
	go tool cover -html=server/coverage.out -o server/coverage/coverage.html

# Run server tests with Ginkgo, verbose output and coverage
server-test-ginkgo-verbose:
	mkdir -p server/coverage
	cd server && $(shell go env GOPATH)/bin/ginkgo --v --cover --coverprofile=coverage.out ./...
	go tool cover -html=server/coverage.out -o server/coverage/coverage.html

# Run server tests with Ginkgo, focusing on specific tests or packages
# Usage: make server-test-ginkgo-focus FOCUS="TestName"
server-test-ginkgo-focus:
	mkdir -p server/coverage
	cd server && $(shell go env GOPATH)/bin/ginkgo --focus="$(FOCUS)" --cover --coverprofile=coverage.out ./...
	go tool cover -html=server/coverage.out -o server/coverage/coverage.html

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

# Client e2e test commands
client-test-e2e:
	cd client && npm run test:e2e

# Run e2e tests with UI mode
client-test-e2e-ui:
	cd client && npm run test:e2e:ui

# Run e2e tests in headed mode (shows browser)
client-test-e2e-headed:
	cd client && npm run test:e2e -- --headed

# Run e2e tests for a specific file
# Usage: make client-test-e2e-file FILE=movementDemo.spec.ts
client-test-e2e-file:
	cd client && npm run test:e2e -- $(FILE)

# Run e2e tests with debug mode
client-test-e2e-debug:
	cd client && npm run test:e2e -- --debug

# Start the development server and run e2e tests
# This target starts the server in the background, waits for it to be ready, runs the tests, and then stops the server
client-test-e2e-with-server:
	@echo "Starting development server..."
	cd client && npm run start & \
	SERVER_PID=$$! && \
	echo "Waiting for server to be ready..." && \
	sleep 15 && \
	echo "Running e2e tests..." && \
	cd client && npm run test:e2e && \
	echo "Stopping development server..." && \
	kill $$SERVER_PID || true

# Start the development server and run a specific e2e test file
# Usage: make client-test-e2e-file-with-server FILE=movementDemo.spec.ts
client-test-e2e-file-with-server:
	@echo "Starting development server..."
	cd client && npm run start & \
	SERVER_PID=$$! && \
	echo "Waiting for server to be ready..." && \
	sleep 15 && \
	echo "Running e2e tests for $(FILE)..." && \
	cd client && npm run test:e2e -- $(FILE) && \
	echo "Stopping development server..." && \
	kill $$SERVER_PID || true

# Start the development server and run e2e tests in headed mode
client-test-e2e-headed-with-server:
	@echo "Starting development server..."
	cd client && npm run start & \
	SERVER_PID=$$! && \
	echo "Waiting for server to be ready..." && \
	sleep 15 && \
	echo "Running e2e tests in headed mode..." && \
	cd client && npm run test:e2e -- --headed && \
	echo "Stopping development server..." && \
	kill $$SERVER_PID || true

# Run both server and client (requires tmux)
dev-full:
	tmux new-session -d -s thedeeps 'make run'
	tmux split-window -h -t thedeeps 'make client-start'
	tmux -2 attach-session -t thedeeps
