# The Deeps - Server

This is the server component for The Deeps, a roguelike dungeon crawler game.

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