# TheDeeps UI - TODO List

This document lists the missing UI elements that need to be implemented to fully match the requirements specified in `docs/ai/ai-ui.md`.

## Character Creation Screen
- [ ] Implement a Random Name generator button that generates fantasy-appropriate character names

## Character Profile Window
- [ ] Add a gold counter to display how much gold the character currently has
- [ ] Add a potion counter to display how many potions the character is currently holding
- [ ] Enhance the inventory section to show gold, potions, and other items
- [ ] Ensure experience points are properly tracked and updated during gameplay

## Dungeon Window
- [ ] Ensure the character icon is a stylized @ symbol
- [ ] Implement unique colors for each character class's @ symbol

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

This TODO list should be updated as items are completed or new requirements are identified. 