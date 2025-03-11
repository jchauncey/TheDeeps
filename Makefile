.PHONY: build run clean client-install client-start client-build client-test client-test-coverage client-test-coverage-detail client-open-coverage

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

# Download dependencies
deps:
	go mod download

# Run tests
test:
	go test -v ./...

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
