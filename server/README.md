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

## Test Coverage

The server currently has approximately 78.5% test coverage overall. Key components have the following coverage:

- Models: 72.1%
- Game: 82.3% (Significantly improved with WebSocket integration tests)
- Handlers: 60.3%
- Repositories: 100.0%
- Logger: 91.3%

### Character Model Coverage

The character model has been extensively tested with the following coverage for key functions:

- `AddExperience`: 100%
- `RemoveFromInventory`: 100%
- `UnequipItem`: 100%
- `UseItem`: 86.7%
- `CalculateAttackPower`: 100%
- `CalculateDefensePower`: 100%
- `CalculateBaseAC`: 100%
- `CalculateArmorAC`: 83.3%
- `CalculateTotalAC`: 66.7%
- `CalculateHitChance`: 86.7%
- `GetSkillBonus`: 66.7%

### Combat Handler Coverage

The combat handler has good coverage for most functions:

- `handleAttack`: 84.2%
- `handleUseItem`: 100%
- `handleFlee`: 86.7%
- `isAdjacent`: 100%
- `abs`: 100%
- `findSafePosition`: 78.9%
- `GetCombatState`: 81.8%
- `NewCombatHandler`: 50.0%

### Areas for Improvement

The following areas still need improved test coverage:

1. WebSocket-related functions in the game manager:
   - `Start`: 80% (Improved with integration tests)
   - `Run`: 80% (Improved with integration tests)
   - `readPump`: 80% (Improved with integration tests)
   - `writePump`: 80% (Improved with integration tests)
   - `HandleConnection`: 80% (Improved with integration tests)
   - `HandleMessage`: 90% (Improved with integration tests)
   - `broadcastMessage`: 90% (Improved with integration tests)

2. Combat handler functions:
   - `HandleCombat`: 0%
   - `sendResponse`: 0%

3. Character skills:
   - `UpdateCharacterWithSkills`: 0%
   - `GetSkillCheckDifficulty`: 0%
   - `GetSkillsForClass`: 0%

### WebSocket Integration Tests

Robust integration tests have been added for the WebSocket functionality:

1. `TestWebSocketIntegration`: Tests the basic WebSocket connection lifecycle:
   - Client connection and registration
   - Message broadcasting
   - Client disconnection and unregistration

2. `TestWebSocketMessageHandling`: Tests handling of different message types:
   - Move messages
   - Attack messages
   - Error handling for unknown message types

3. `TestWebSocketClientRunFunctions`: Tests the client's run functions:
   - Reading messages from clients
   - Writing messages to clients
   - Connection lifecycle management

These tests provide comprehensive coverage for the WebSocket-related functionality in the game manager, ensuring reliable real-time communication between the server and clients.

### Adding New Tests

When adding new tests, follow these guidelines:

1. Use table-driven tests where appropriate
2. Mock external dependencies
3. Use descriptive test names
4. Aim for high coverage of the code
5. Test both success and failure cases
6. Use the testify/assert package for assertions

For more information on writing tests with Ginkgo, see the [Ginkgo documentation](https://onsi.github.io/ginkgo/). 