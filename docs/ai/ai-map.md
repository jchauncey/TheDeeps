# TheDeeps Map System Documentation

## Overview

TheDeeps features a procedurally generated dungeon system with multiple floors, room types, and navigation mechanics. The map generation system creates varied and interesting environments while ensuring gameplay balance and proper progression.

## Map Structure

### Dungeon Organization

- **Multi-Floor Design**: Dungeons consist of multiple floors (levels) with increasing difficulty
- **Floor Structure**: Each floor contains a grid of tiles representing walls, floors, and special features
- **Room-Based Layout**: Floors are organized into distinct rooms connected by corridors
- **Stair System**: Floors are connected by staircases allowing vertical movement through the dungeon

### Tile System

- **Tile Types**:
  - Wall: Non-walkable boundary tiles
  - Floor: Basic walkable tiles
  - Up Stairs: Connection to previous floor (except on first floor)
  - Down Stairs: Connection to next floor (except on final floor)
  - Special tiles: Shop counters, traps, etc.
- **Tile Properties**:
  - Walkability: Determines if entities can move onto the tile
  - Visibility: Controls if the tile is visible to the player
  - Explored status: Tracks if the player has seen the tile
  - Room ID: Associates the tile with a specific room
  - Entity references: Links to characters, mobs, or items on the tile

## Room Generation

### Room Types

- **Entrance Room**: Starting point on the first floor, always contains down stairs
- **Safe Room**: First room on floors 2+, always contains up stairs, never has mobs
- **Standard Room**: Common room type with basic encounters
- **Treasure Room**: Contains higher quality and quantity of items
- **Shop Room**: Contains a shopkeeper and purchasable items
- **Boss Room**: Present on the final floor, contains a boss enemy

### Room Properties

- **Size**: Variable dimensions (minimum 5x5, maximum 10x10)
- **Position**: Placed to avoid overlapping with other rooms
- **Type**: Determines the room's purpose and contents
- **Exploration Status**: Tracks if the player has discovered the room

### Room Distribution

- **First Floor**:
  - Always has an entrance room near the center
  - Always has a shop room
  - Additional standard rooms based on floor size
- **Middle Floors**:
  - Always has a safe room with up stairs
  - Mix of standard, treasure, and shop rooms
  - Room count increases with floor depth
- **Final Floor**:
  - Always has a safe room with up stairs
  - Always has a boss room
  - Additional standard and treasure rooms

## Corridor System

- **Connection Algorithm**: Rooms are connected with L-shaped corridors
- **Pathfinding**: Ensures all rooms are accessible
- **Corridor Properties**: Corridors are walkable and can contain items or mobs

## Stair Placement

- **Up Stairs**:
  - Located in safe rooms on floors 2+
  - Positioned in the center of the room
  - Marked as explored when the floor is generated
- **Down Stairs**:
  - Located in entrance rooms on floor 1
  - Located in the last room on other floors
  - Positioned in the center of the room on floor 1
  - Randomly positioned in other rooms

## Entity Placement

### Mob Placement

- **Room-Based**: Mobs are placed in appropriate rooms based on type
- **Exclusion Zones**: No mobs in entrance rooms or safe rooms
- **Density Control**: Number of mobs scales with room size and floor depth
- **Variant Distribution**: Mix of easy, normal, and hard variants based on difficulty setting
- **Boss Placement**: Boss mobs only appear in boss rooms

### Item Placement

- **Room-Based**: Items are distributed throughout rooms with higher concentrations in treasure rooms
- **Type Distribution**: Item types and quality scale with floor depth
- **Placement Rules**: Items are only placed on walkable tiles not occupied by other entities

## Floor Navigation

- **Player Positioning**:
  - When descending, players are positioned near up stairs in the safe room
  - When ascending, players are positioned near down stairs
  - Player positions are always on walkable tiles, avoiding stairs and other entities
- **Floor Transitions**:
  - Maintain player state across floors (inventory, health, etc.)
  - Update visibility and exploration status
  - Trigger appropriate notifications

## Map Generation Parameters

- **Seed-Based**: Maps can be generated deterministically using a seed value
- **Difficulty Levels**: Easy, normal, and hard difficulties affect mob distribution
- **Floor Scaling**: Deeper floors have more rooms, stronger mobs, and better loot
- **Room Count**: 5-20 rooms per floor, scaling with floor depth

## Technical Implementation

- **Data Structure**: Maps are represented as 2D arrays of tile objects
- **Serialization**: Map data can be saved and loaded for game persistence
- **Client-Server Model**: Maps are generated on the server and transmitted to clients
- **Optimization**: Generation algorithms are designed for efficiency and performance

## Visualization

- **Rendering**: Maps are displayed as a grid of colored tiles
- **Entity Representation**: Characters, mobs, and items have distinct visual styles
- **Room Differentiation**: Different room types have subtle visual distinctions
- **Navigation Aids**: Stairs are clearly marked with directional indicators