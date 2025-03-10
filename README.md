# The Deeps

A roguelike dungeon crawler game written in Go with a React/TypeScript frontend.

## Description

The Deeps is a roguelike game where you explore procedurally generated dungeons, battle monsters, and collect treasure. The game features:

- Procedurally generated dungeons with multiple floors
- Character classes with unique abilities
- Inventory system with weapons, armor, and items
- Turn-based gameplay
- Save/load game functionality
- Class-specific styling and abilities

## Controls

- WASD: Move character
- G: Pick up item
- F: Attack
- I: Open inventory
- C: Open character profile
- Escape: Open main menu
- </>: Use stairs
- ?: Help screen

## Main Menu

Press Escape during gameplay to access the main menu, which includes:

- Return to Game: Continue playing
- Save Game: Save your current progress
- Load Game: Load a previously saved game
- New Game: Start a new character
- Quit to Main Menu: Return to the start screen

## Installation

### Server
1. Make sure you have Go installed on your system
2. Clone this repository
3. Install dependencies:

```bash
cd server
go mod tidy
```

4. Run the server:

```bash
go run main.go
```

### Client
1. Make sure you have Node.js installed
2. Install dependencies:

```bash
cd client
npm install
```

3. Run the client:

```bash
npm start
```

## Project Structure

- `server/`: Backend Go server
  - `game/`: Game state and WebSocket handling
  - `models/`: Data structures for game entities
  - `repositories/`: Data persistence
  - `handlers/`: HTTP request handlers
- `client/`: Frontend React application
  - `src/components/`: React components
  - `src/services/`: API and WebSocket services
  - `src/types/`: TypeScript type definitions
  - `src/constants/`: Game constants and configuration

## Documentation

- `docs/ai/`: Design requirements and specifications
- `docs/ui-todo.md`: UI implementation checklist
- `docs/character-todo.md`: Character system implementation checklist
- `docs/main-menu.md`: Main menu system documentation

## License

MIT 