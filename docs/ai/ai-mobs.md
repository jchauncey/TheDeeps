# TheDeeps Mob System Documentation

## Overview

TheDeeps features a diverse ecosystem of monsters (mobs) that populate the dungeon floors. The mob system is designed to provide balanced challenges that scale with dungeon depth and difficulty settings.

## Mob Types

### Basic Enemies

- **Skeletons**: Basic undead enemies, common on early floors
- **Goblins**: Quick, weak humanoids that appear in groups
- **Ratmen**: Agile rodent-humanoids with moderate damage
- **Orcs**: Stronger humanoid enemies with higher health
- **Oozes**: Amorphous creatures with special resistances

### Advanced Enemies

- **Trolls**: High health regenerating enemies
- **Wraiths**: Incorporeal undead with special abilities
- **Ogres**: Large, high-damage brutes
- **Drakes**: Dragon-like creatures with elemental attacks

### Elite Enemies

- **Liches**: Powerful undead spellcasters
- **Elementals**: Embodiments of elemental forces
- **Dragons**: Ultimate enemies with devastating attacks

### Special Entities

- **Shopkeepers**: Non-hostile NPCs that offer items for purchase
- **Boss Variants**: Unique, powerful versions of standard enemies

## Mob Properties

### Core Attributes

- **Health**: Determines how much damage a mob can take
- **Attack Power**: Determines base damage output
- **Defense**: Reduces incoming damage
- **Speed**: Affects turn order and movement capabilities
- **Level**: Scales with floor depth and affects overall power

### Variant System

- **Easy Variant**: Lower stats, simpler behavior, less valuable drops
- **Normal Variant**: Standard stats and behavior
- **Hard Variant**: Enhanced stats, more complex behavior, better drops
- **Boss Variant**: Significantly enhanced stats, unique abilities, valuable drops

### Behavior Patterns

- **Aggressive**: Actively pursues and attacks the player
- **Territorial**: Attacks when player enters its territory
- **Passive**: Only attacks when provoked
- **Special**: Unique behaviors for specific mob types

## Mob Generation

### Placement Rules

- **Room-Based**: Mobs are assigned to specific rooms based on type
- **Exclusion Zones**: No mobs in entrance rooms or safe rooms
- **Density Control**: Number of mobs scales with room size and floor depth
- **Boss Placement**: Boss mobs only appear in boss rooms
- **Shop Placement**: Shopkeepers only appear in shop rooms

### Scaling Mechanics

- **Floor Depth**: Higher floors spawn stronger mob types and variants
- **Difficulty Setting**: Affects the distribution of mob variants
  - Easy: More easy variants, fewer hard variants
  - Normal: Balanced distribution of variants
  - Hard: More hard variants, fewer easy variants
- **Room Type**: Affects the quantity and quality of mobs
  - Standard Rooms: Normal distribution
  - Treasure Rooms: Fewer but stronger mobs
  - Boss Rooms: Single powerful boss mob

## Combat Mechanics

### Attack System

- **Hit Chance**: Calculated based on mob's attack rating vs. player's defense
- **Damage Calculation**: Based on attack power, modified by equipment and abilities
- **Critical Hits**: Chance to deal increased damage
- **Miss Chance**: Possibility of attacks failing to connect

### Defense System

- **Armor Class**: Determines difficulty to hit the mob
- **Damage Reduction**: Decreases incoming damage based on defense rating
- **Resistances**: Some mobs have specific damage type resistances
- **Vulnerabilities**: Some mobs take increased damage from certain sources

## Reward System

### Experience Points

- **Base Value**: Each mob type has a base XP value
- **Level Scaling**: XP scales with mob level
- **Variant Bonus**: Hard and boss variants provide additional XP
- **Level Difference**: XP is adjusted based on level difference between player and mob

### Loot Drops

- **Gold**: Currency dropped based on mob type, level, and variant
- **Items**: Equipment, consumables, and special items
- **Drop Rates**: Probability-based system for determining drops
- **Quality Scaling**: Higher level mobs drop better quality items

## Technical Implementation

- **Data Structure**: Mobs are represented as objects with properties and methods
- **Serialization**: Mob data can be saved and loaded for game persistence
- **Client-Server Model**: Mob logic runs on the server with state updates sent to clients
- **Performance Optimization**: Efficient algorithms for mob behavior and combat calculations

## Visual Representation

- **Symbol-Based**: Each mob type has a distinct symbol in the tile-based display
- **Color Coding**: Different mob types and variants have specific colors
- **Status Indicators**: Visual cues for mob health and status effects
- **Animation**: Visual feedback for mob actions and state changes