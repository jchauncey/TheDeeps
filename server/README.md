# TheDeeps Server

This is the server component of TheDeeps, a roguelike dungeon crawler game.

## Running the Server

You can run the server using one of the following methods:

### Method 1: Using the run script
```bash
./run.sh
```

### Method 2: Using Go directly
```bash
go run .
```

### Method 3: From the project root using Make
```bash
cd ..
make run
```

## Server Structure

- `models/`: Data structures for game entities
- `repositories/`: Data persistence
- `handlers/`: HTTP request handlers
- `game/`: Game state and WebSocket handling
- `utils/`: Utility functions
- `log/`: Logging system

## API Endpoints

### Character Endpoints
- `GET /characters`: Get all characters
- `GET /characters/{id}`: Get a character by ID
- `POST /characters`: Create a new character
- `DELETE /characters/{id}`: Delete a character
- `POST /characters/{id}/save`: Save a character's state
- `GET /characters/{id}/floor`: Get a character's current floor
- `GET /characters/{id}/combat`: Get a character's combat state

### Dungeon Endpoints
- `GET /dungeons`: Get all dungeons
- `POST /dungeons`: Create a new dungeon
- `POST /dungeons/{id}/join`: Join a dungeon with a character
- `GET /dungeons/{id}/floor/{level}`: Get a specific floor of a dungeon

### WebSocket Endpoints
- `/ws/game?characterId={id}`: Connect to the game with a character
- `/ws/combat`: Connect to the combat system

## Development

The server is built with Go and uses:
- Gorilla Mux for routing
- Gorilla WebSocket for WebSocket connections
- Custom logging system
- In-memory repositories (can be extended to use a database)

## Running Tests with Ginkgo

[Ginkgo](https://onsi.github.io/ginkgo/) is a BDD-style testing framework for Go. This project uses Ginkgo for running tests with coverage reporting.

### Prerequisites

1. Install the Ginkgo CLI:

```bash
go install github.com/onsi/ginkgo/v2/ginkgo@latest
```

2. Make sure the Ginkgo binary is in your PATH:

```bash
export PATH=$PATH:$(go env GOPATH)/bin
```

### Running Tests

The project includes several make targets for running tests with Ginkgo:

#### Run all server tests with Ginkgo and coverage

```bash
make server-test-ginkgo
```

This will:
- Run all tests in the server directory
- Generate a coverage report
- Create an HTML coverage report
- Display a summary of the coverage

#### Run tests with verbose output

```bash
make server-test-ginkgo-verbose
```

This is useful for debugging as it provides more detailed output.

#### Run specific tests

```bash
make server-test-ginkgo-focus FOCUS="TestName"
```

Replace `TestName` with the name of the test you want to run. This is useful when you want to focus on a specific test or group of tests.

### Viewing Coverage Reports

To open the HTML coverage report in your default browser:

```bash
make server-open-coverage
```

## Test Coverage Summary

The current test coverage is approximately 68.5% of statements (improved from 63.9%). Areas with high coverage include:

- Repository implementations (100%)
- Logger functionality (91.3%)
- Character and skill models (62.6%)
- Map generator (92.1%)
- Character handler (92.3%)
- Game manager (71.4%)
- Combat handler (most functions now have good coverage)

Key improvements in the combat handler coverage:
- `handleAttack`: 84.2% coverage
- `handleUseItem`: 100% coverage
- `handleFlee`: 86.7% coverage
- `isAdjacent`: 100% coverage
- `abs`: 100% coverage
- `findSafePosition`: 78.9% coverage
- `GetCombatState`: 81.8% coverage

Areas that still need improved test coverage:

- WebSocket-related functions in game manager (`Run`, `readPump`, `writePump`, `HandleConnection`): 0%
- `Start` function in game manager: 0%
- `HandleCombat` in combat handler: 0% (requires WebSocket connection)
- `sendResponse` in combat handler: 0% (requires WebSocket connection)

## Adding New Tests

When adding new tests, follow these guidelines:

1. Place test files in the same package as the code being tested
2. Name test files with the `_test.go` suffix
3. Use descriptive test names that explain what is being tested
4. Use table-driven tests where appropriate to test multiple scenarios
5. Mock external dependencies to isolate the code being tested

For more information on writing tests with Ginkgo, see the [Ginkgo documentation](https://onsi.github.io/ginkgo/). 