# TheDeeps Character System - TODO List

This document lists the character system features that need to be implemented to fully match the requirements specified in `docs/ai/ai-characters.md`.

## Character Data Structure Enhancements

- [ ] Update the `CharacterData` interface to include:
  - [ ] Current and maximum HP fields
  - [ ] Current and maximum mana fields
  - [ ] Current experience and experience needed for next level
  - [ ] Level field
  - [ ] Equipped items (weapon, armor, shield, accessories)
  - [ ] Inventory items (including gold and potions)

## Equipment System

- [ ] Create an `Equipment` interface with:
  - [ ] Weapon slot
  - [ ] Armor slot
  - [ ] Shield slot
  - [ ] Accessory slot(s)

- [ ] Implement equipment restrictions based on class:
  - [ ] Only allow equipping items the character's class is proficient with
  - [ ] Display error messages when attempting to equip incompatible items

- [ ] Create item data structures:
  - [ ] Weapons with damage values and special properties
  - [ ] Armor with AC bonus values and special properties
  - [ ] Shields with AC bonus values
  - [ ] Accessories with various bonuses

- [ ] Update AC calculation to account for equipped armor and shields

## Inventory System

- [ ] Create an `Inventory` interface with:
  - [ ] Gold amount
  - [ ] Potion count
  - [ ] Array of other items

- [ ] Implement inventory management functions:
  - [ ] Add item to inventory
  - [ ] Remove item from inventory
  - [ ] Use consumable items (potions, scrolls)
  - [ ] Equip/unequip items

## Character Progression

- [ ] Implement experience gain system:
  - [ ] Award experience for defeating enemies
  - [ ] Award experience for completing objectives

- [ ] Implement leveling system:
  - [ ] Define experience thresholds for each level
  - [ ] Increase character stats on level up
  - [ ] Unlock new abilities on level up

## Combat Stats

- [ ] Implement derived combat stats:
  - [ ] Attack bonus based on attributes and proficiencies
  - [ ] Damage calculation based on weapon and attributes
  - [ ] Defense calculation based on armor, shield, and attributes

## UI Updates

- [ ] Update the GameStatus component to display:
  - [ ] Current equipment with stats
  - [ ] Inventory items
  - [ ] Gold amount
  - [ ] Potion count

- [ ] Add equipment management UI:
  - [ ] Equip/unequip buttons
  - [ ] Item comparison tooltips

- [ ] Add inventory management UI:
  - [ ] Use item buttons
  - [ ] Drop item buttons

## Backend Integration

- [x] Implement basic character save functionality
- [x] Add server-side save game endpoint
- [ ] Implement full character data persistence
- [ ] Implement server-side validation for equipment and inventory changes
- [ ] Synchronize character stats between client and server

## Character State Management

- [x] Implement main menu with save/load options
- [x] Add save game functionality that preserves character state
- [ ] Implement load game functionality to restore character state
- [ ] Add game state tracking (position, health, inventory)
- [ ] Implement auto-save feature at key points (level changes, etc.)

## Implementation Notes

- The equipment system should respect D&D rules regarding armor types and weapon proficiencies
- AC calculation should follow D&D formula: 10 + DEX modifier + armor bonus + shield bonus
- Special class features like Unarmored Defense for Monks and Barbarians should be preserved
- The inventory system should have reasonable limits based on character strength or other factors
- The save game system now stores character data in the server's character repository
- Character state includes position, health, mana, and other attributes

This TODO list should be updated as items are completed or new requirements are identified. 