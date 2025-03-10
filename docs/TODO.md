# The Deeps - TODO List

This document outlines the missing components and enhancements needed for the server to fully implement the requirements specified in the AI documents.

## Character System

- [x] **Inventory System**
  - [x] Add inventory field to Character struct
  - [x] Implement methods for adding/removing items from inventory
  - [x] Add weight limit based on strength attribute
  - [ ] Create UI components for inventory management

- [x] **Equipment System**
  - [x] Add equipment slots to Character struct (weapon, armor, etc.)
  - [x] Implement methods for equipping/unequipping items
  - [x] Add class and level requirements for equipment
  - [x] Calculate stat bonuses from equipped items
  - [ ] Implement Armor Class (AC) system for combat

- [ ] **Skill System**
  - [ ] Create Skills struct with various skills (Stealth, Perception, etc.)
  - [ ] Implement skill checks based on attributes
  - [ ] Add skill improvement through use or skill points
  - [ ] Implement class-specific skills or bonuses

- [ ] **Character Death and Respawn**
  - [ ] Implement logic for character death when HP reaches 0
  - [ ] Create respawn mechanics (at safe location, with penalties)
  - [ ] Add death statistics tracking

- [ ] **Character Customization**
  - [ ] Add visual appearance options
  - [ ] Implement background selection affecting starting skills/attributes
  - [ ] Create specialization paths within each class
  - [ ] Add talent or feat system

## Dungeon System

- [ ] **Enhanced Room Types**
  - [ ] Implement treasure rooms with special loot
  - [ ] Create boss rooms with powerful enemies
  - [ ] Add puzzle rooms requiring player interaction
  - [ ] Implement safe rooms with healing resources
  - [ ] Create shop rooms for trading

- [ ] **Room Properties**
  - [ ] Add environmental hazards (traps, lava, etc.)
  - [ ] Implement interactive objects (chests, levers, etc.)
  - [ ] Create hidden doors requiring discovery

- [ ] **Dungeon Themes**
  - [ ] Implement visual themes affecting appearance
  - [ ] Add theme-specific enemies and hazards
  - [ ] Create theme-appropriate loot

## Combat System

- [x] **Item Drop Mechanics**
  - [x] Implement loot tables for different mob types
  - [x] Scale drops based on mob variant and floor level
  - [ ] Add rare item drops for boss monsters

- [ ] **Enhanced Combat Mechanics**
  - [ ] Implement initiative system based on dexterity
  - [ ] Add status effects (poison, stun, etc.)
  - [ ] Create special abilities for different character classes
  - [ ] Implement area-of-effect attacks

- [ ] **Armor Class System**
  - [ ] Calculate AC based on armor and dexterity
  - [ ] Use AC in hit chance calculations
  - [ ] Implement armor penetration for certain attacks

## NPC System

- [ ] **NPC Interaction**
  - [ ] Create dialogue system with options
  - [ ] Implement attribute-based dialogue checks (e.g., charisma)
  - [ ] Add reputation system with different factions
  - [ ] Track relationships with key NPCs

- [ ] **Quest System**
  - [ ] Implement quest tracking
  - [ ] Create quest rewards (items, experience, reputation)
  - [ ] Add quest availability based on character attributes or choices

## Game State Management

- [ ] **Enhanced Save System**
  - [ ] Ensure all aspects of game state are saved
  - [ ] Implement auto-save functionality
  - [ ] Add save slots for multiple game states

- [ ] **Session Management**
  - [ ] Implement proper session handling for multiple players
  - [ ] Add authentication and authorization
  - [ ] Create admin tools for game management

## Technical Improvements

- [ ] **Performance Optimization**
  - [ ] Optimize map generation for large dungeons
  - [ ] Implement efficient visibility calculations
  - [ ] Add caching for frequently accessed data

- [ ] **Error Handling**
  - [ ] Improve error messages and logging
  - [ ] Implement graceful recovery from errors
  - [ ] Add validation for all user inputs

- [ ] **Testing**
  - [ ] Increase test coverage for all components
  - [ ] Add integration tests for complex interactions
  - [ ] Implement load testing for WebSocket connections

## Documentation

- [ ] **API Documentation**
  - [ ] Document all REST endpoints
  - [ ] Create WebSocket message reference
  - [ ] Add examples for common operations

- [ ] **Developer Guide**
  - [ ] Create setup instructions for new developers
  - [ ] Document code structure and architecture
  - [ ] Add contribution guidelines

- [ ] **User Guide**
  - [ ] Create player manual
  - [ ] Add tutorials for new players
  - [ ] Document game mechanics and systems 