# The Deeps

A roguelike dungeon crawler game written in Go with a React/TypeScript frontend.

![Coverage](https://img.shields.io/badge/coverage-55.7%25-brightgreen)

## Description

The Deeps is a roguelike game where you explore procedurally generated dungeons, battle monsters, and collect treasure. The game features:

- Procedurally generated dungeons with multiple floors
- Character classes with unique abilities
- Turn-based gameplay
- WebSocket-based real-time communication
- Floor navigation with stairs
- Multiple characters and dungeons
- Combat system with attacks, critical hits, and fleeing

## Features

### Character System
- 12 character classes (Warrior, Mage, Rogue, Cleric, Druid, Warlock, Bard, Paladin, Ranger, Monk, Barbarian, Sorcerer)
- D&D-style attributes (Strength, Dexterity, Constitution, Intelligence, Wisdom, Charisma)
- Experience and leveling system
- Character creation and management

### Dungeon System
- Procedurally generated floors with rooms and corridors
- Multiple room types (standard, treasure, boss, etc.)
- Stairs for navigating between floors
- Mob placement based on floor difficulty

### Mob System
- Various mob types with different stats and abilities
- Mob variants (easy, normal, hard, boss)
- Scaling difficulty based on floor level

### Item System
- Different item types (weapons, armor, potions, etc.)
- Item generation based on floor level
- Gold as currency

### Combat System
- Turn-based combat with attacks and counterattacks
- Critical hit system based on character attributes
- Damage calculation using character and mob stats
- Experience and gold rewards for defeating mobs
- Flee mechanics with success chance based on attributes
- Item usage during combat (potions, etc.)

## Project Structure

- `server/`: Backend Go server
  - `models/`: Data structures for game entities
  - `repositories/`: Data persistence
  - `handlers/`: HTTP request handlers
  - `game/`: Game state and WebSocket handling
- `client/`: Frontend React/TypeScript application
  - `src/components/`: Reusable UI components
  - `src/pages/`: Page components for different routes
  - `src/services/`: API services for communicating with the server
  - `src/types/`: TypeScript type definitions

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

## WebSocket Messages

### Game WebSocket (Client to Server)
- `move`: Move the character
- `pickup`: Pick up an item
- `ascend`: Go up stairs
- `descend`: Go down stairs
- `useItem`: Use an item
- `dropItem`: Drop an item
- `equipItem`: Equip an item
- `unequipItem`: Unequip an item

### Game WebSocket (Server to Client)
- `updateMap`: Update the map
- `updatePlayer`: Update the player
- `updateMob`: Update a mob
- `removeMob`: Remove a mob
- `addItem`: Add an item
- `removeItem`: Remove an item
- `notification`: Show a notification
- `floorChange`: Change the floor
- `error`: Show an error

### Combat WebSocket (Client to Server)
- `attack`: Attack a mob
- `useItem`: Use an item during combat
- `flee`: Attempt to flee from combat

### Combat WebSocket (Server to Client)
- `attack`: Result of an attack
- `useItem`: Result of using an item
- `flee`: Result of a flee attempt
- `error`: Combat error message

## Getting Started

### Prerequisites
- Go 1.21 or higher
- Node.js 14+ and npm 6+ (for the client)

### Building and Running the Server
```bash
# Build the server
make build

# Run the server
make run

# Or run directly from the server directory
cd server
go run .
```

### Building and Running the Client
```bash
# Navigate to the client directory
cd client

# Install client dependencies
npm install

# Run the client in development mode
npm start

# Or use the run script
./run.sh
```

## Testing

### Server Tests
The server code is tested using Go's standard testing package and the Ginkgo testing framework.

```bash
# Run server tests with coverage
make server-test-coverage

# Run server tests with Ginkgo and coverage
make server-test-ginkgo

# Run server tests with Ginkgo, verbose output and coverage
make server-test-ginkgo-verbose

# Run specific tests with Ginkgo
make server-test-ginkgo-focus FOCUS="TestName"

# Open the coverage report in your browser
make server-open-coverage
```

Current test coverage is approximately 55.7% of statements. See the [server README](server/README.md) for more details on testing.

### Client Tests
The client code is tested using Jest and React Testing Library.

```bash
# Run client tests
make client-test

# Run client tests with coverage
make client-test-coverage

# Run client tests with detailed coverage
make client-test-coverage-detail

# Open the client coverage report
make client-open-coverage
```

## Troubleshooting

### Client Issues
If you encounter the error `ENOENT: no such file or directory, uv_cwd` when running `npm run dev`, try using `npm start` instead.

If you encounter dependency conflicts during installation, you can use:
```bash
npm install --legacy-peer-deps
```

### Server Issues
If you encounter an error like `undefined: NewServer` when running `go run main.go`, use `go run .` instead to compile all files in the package.

## Client Features

### Character Selection
- View all your characters in a grid layout
- See character details including name, class, level, HP, mana, and XP
- Delete characters you no longer want
- Create new characters (up to a maximum of 10)

### Character Creation
- Choose from 12 different character classes
- Each class has unique attributes and abilities
- Simple creation process with class descriptions

## License

MIT 