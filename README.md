# The Deeps

A roguelike dungeon crawler game written in Go.

## Description

The Deeps is a terminal-based roguelike game where you explore procedurally generated dungeons, battle monsters, and collect treasure. The game features:

- Procedurally generated dungeons with multiple floors
- Character classes with unique abilities
- Inventory system with weapons, armor, and items
- Fog of war exploration system
- Turn-based gameplay

## Controls

- Arrow keys: Move character
- A: Use special ability
- I: Open inventory
- Q: Quit game

## Installation

1. Make sure you have Go installed on your system
2. Clone this repository
3. Install dependencies:

```bash
go mod tidy
```

4. Run the game:

```bash
go run main.go
```

## Project Structure

- `main.go`: Entry point for the game
- `internal/`: Internal packages
  - `character/`: Character-related code (player, classes, inventory)
  - `game/`: Game state and loop management
  - `input/`: Input handling
  - `map/`: Dungeon generation and map management
  - `ui/`: User interface rendering

## License

MIT 