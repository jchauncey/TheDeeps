# TheDeeps UI - TODO List

This document lists the missing UI elements that need to be implemented to fully match the requirements specified in `docs/ai/ai-ui.md`.

## Character Creation Screen
- [ ] Implement a Random Name generator button that generates fantasy-appropriate character names

## Character Profile Window
- [x] Add a gold counter to display how much gold the character currently has
- [ ] Add a potion counter to display how many potions the character is currently holding
- [ ] Enhance the inventory section to show gold, potions, and other items
- [ ] Ensure experience points are properly tracked and updated during gameplay

## Dungeon Window
- [x] Ensure the character icon is a stylized @ symbol
- [x] Implement unique colors for each character class's @ symbol

## Menu System
- [x] Implement main menu accessible via Escape key
- [x] Add Save Game functionality
- [x] Add Load Game option
- [x] Add New Game option
- [x] Add Return to Game option
- [x] Add Quit to Main Menu option
- [ ] Implement Settings menu with audio, video, and gameplay options

## Floor Navigation
- [x] Visual indication when player is standing on stairs
- [x] Tooltip or prompt showing available actions when standing on stairs
- [x] Feedback message when attempting to ascend from the top floor or descend from the bottom floor
- [x] Floor transition animation when moving between floors
- [x] Floor number indicator that updates when changing floors
- [x] Action buttons for ascending and descending when standing on appropriate stairs

## General UI Improvements
- [ ] Verify that toast banners automatically disappear after 5 seconds
- [ ] Ensure the main browser window avoids heights that add scroll bars
- [ ] Make sure the character profile window can scroll if content exceeds the window height

## Implementation Notes
- The inventory system should track:
  - Gold amount (numeric value)
  - Potion count (numeric value)
  - Other equipment and items

- The character icon system should:
  - Use @ as the base symbol
  - Apply different colors based on character class
  - Stand out clearly against the dungeon background

- The menu system now includes:
  - Main menu accessible via Escape key
  - Save game functionality that stores character state
  - Load game option to retrieve saved characters
  - New game option to start a fresh character
  - Return to game option to continue playing
  - Quit option to return to the start screen

This TODO list should be updated as items are completed or new requirements are identified. 