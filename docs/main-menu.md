# TheDeeps Main Menu System

This document describes the main menu system implemented in TheDeeps game.

## Overview

The main menu system provides players with essential game management options, including saving and loading games, starting new games, and returning to the current game. The menu is accessible via the Escape key during gameplay and is designed to be intuitive and visually consistent with the rest of the game's UI.

## Features

### Menu Access
- The main menu can be accessed by pressing the **Escape** key during gameplay
- If the help modal is open, pressing Escape will close it first
- Pressing Escape while the menu is open will close it

### Menu Options

The main menu includes the following options:

1. **Return to Game**
   - Closes the menu and returns to the current game session
   - Highlighted as the primary action with a filled button style

2. **Save Game**
   - Saves the current character state to the server
   - Includes position, health, mana, and other attributes
   - Displays a success/failure toast notification
   - Uses the character's name as the identifier

3. **Load Game**
   - Provides access to previously saved games
   - Currently links to the existing load game functionality
   - Will be enhanced to show a list of saved characters

4. **New Game**
   - Allows starting a new game with a new character
   - Takes the player to the character creation screen

5. **Quit to Main Menu**
   - Returns to the start screen
   - Styled in red to indicate it's a potentially destructive action

### Visual Design

- Modal overlay with a semi-transparent backdrop
- Purple accent color matching the game's theme
- Consistent button styling with the rest of the UI
- Character information displayed at the bottom
- Clear visual hierarchy with the primary action highlighted

## Implementation Details

### Client-Side Components

- `MainMenu.tsx`: The main component that renders the menu UI
- `GameControls.tsx`: Updated to handle Escape key and show/hide the menu
- `App.tsx`: Passes character data and navigation functions to GameControls

### Server-Side Components

- `SaveGame` function in `server.go`: Handles saving the character state
- Character repository: Stores and retrieves character data

### Data Flow

1. Player presses Escape to open the menu
2. When Save Game is clicked:
   - Client sends a save_game action to the server
   - Server retrieves the character associated with the connection
   - Server updates the character's timestamp and saves to repository
   - Server sends a success/failure response
   - Client displays a toast notification with the result

## Future Enhancements

- Settings menu with audio, video, and gameplay options
- More detailed save game information (timestamp, level, location)
- Multiple save slots per character
- Auto-save functionality at key points
- Confirmation dialog for potentially destructive actions
- Visual preview of saved games

## Usage

To access the main menu:
1. Press the Escape key during gameplay
2. Select the desired option
3. For Save Game, wait for the success notification
4. To return to the game, select "Return to Game" or press Escape again 