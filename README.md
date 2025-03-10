# The Deeps

A roguelike dungeon crawler game written in Go with a React/TypeScript frontend.

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
- Node.js and npm (for the client)

### Building and Running the Server
```bash
# Build the server
make build-server

# Run the server
make run-server
```

### Building and Running the Client
```bash
# Install client dependencies
make client-install

# Build the client
make client-build

# Run the client in development mode
make client-dev
```

## License

MIT 