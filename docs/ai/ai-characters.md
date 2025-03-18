# TheDeeps Character System Documentation

## Overview

TheDeeps features a robust character system inspired by traditional role-playing games. Characters are defined by their class, attributes, skills, and equipment, providing players with diverse gameplay options and progression paths.

## Character Creation

### Basic Information

- **Name**: Player-chosen identifier, must be unique
- **Class**: Determines starting attributes, skills, and equipment
- **Appearance**: Visual representation based on class selection

### Character Classes

- **Warrior**: Melee combat specialist with high strength and constitution
- **Mage**: Spellcaster with high intelligence and arcane abilities
- **Rogue**: Stealthy character with high dexterity and critical hit potential
- **Cleric**: Divine spellcaster with healing abilities and high wisdom
- **Druid**: Nature-focused spellcaster with shapeshifting abilities
- **Warlock**: Pact-based spellcaster with eldritch abilities
- **Bard**: Charisma-based support class with music and inspiration abilities
- **Paladin**: Holy warrior with combat and divine magic abilities
- **Ranger**: Wilderness expert with tracking and archery skills
- **Monk**: Unarmed combat specialist with mystical abilities
- **Barbarian**: Rage-powered warrior with high damage resistance
- **Sorcerer**: Innate spellcaster with metamagic abilities

### Attribute System

- **Strength (STR)**: Affects melee damage, carrying capacity, and certain skill checks
- **Dexterity (DEX)**: Affects armor class, initiative, and certain skill checks
- **Constitution (CON)**: Affects hit points, stamina, and resistance to effects
- **Intelligence (INT)**: Affects spell power for arcane casters and certain skill checks
- **Wisdom (WIS)**: Affects spell power for divine casters and certain skill checks
- **Charisma (CHA)**: Affects social interactions and certain skill checks

### Attribute Allocation

- **Base Values**: Each class has predetermined base attribute values
- **Customization Points**: Players receive 5 points to distribute among attributes
- **Modifiers**: Attributes generate modifiers that affect various game mechanics
  - 3-4: -3 modifier
  - 5-6: -2 modifier
  - 7-8: -1 modifier
  - 9-12: 0 modifier
  - 13-14: +1 modifier
  - 15-16: +2 modifier
  - 17-18: +3 modifier
  - 19-20: +4 modifier

## Character Progression

### Experience and Leveling

- **Experience Points (XP)**: Gained from defeating enemies and completing objectives
- **Level Progression**: Characters advance through levels as they gain XP
- **Level Benefits**: Each level provides increased stats and potentially new abilities
- **Maximum Level**: Characters can reach level 20

### Health and Mana

- **Hit Points (HP)**: Represents character health, determined by class and CON modifier
- **Mana Points (MP)**: Represents magical energy, determined by class and primary spellcasting attribute
- **Regeneration**: HP and MP regenerate over time or through items and abilities

### Skills and Abilities

- **Skill System**: Characters have skills like Melee, Stealth, Perception, and Arcana
- **Skill Checks**: Skills are tested against difficulty values for success
- **Skill Advancement**: Skills improve through use and level advancement
- **Class Abilities**: Special powers and techniques unique to each class

## Equipment System

### Equipment Types

- **Weapons**: Melee and ranged options with varying damage and properties
- **Armor**: Body protection that affects armor class and defense
- **Accessories**: Items that provide special bonuses or abilities
- **Consumables**: One-use items like potions and scrolls

### Equipment Properties

- **Weight**: Affects carrying capacity and movement
- **Value**: Gold cost for buying and selling
- **Durability**: Some items may degrade with use
- **Magical Properties**: Special effects or bonuses
- **Class Restrictions**: Some items can only be used by certain classes

### Inventory Management

- **Carrying Capacity**: Limited by strength attribute
- **Weight System**: Items have weight values that accumulate
- **Organization**: Grid-based inventory interface
- **Equipment Slots**: Specific positions for equipped items (weapon, armor, etc.)

## Combat Mechanics

### Attack System

- **Hit Chance**: Calculated based on attacker's skill vs. defender's armor class
- **Damage Calculation**: Based on weapon, attributes, and critical hits
- **Critical Hits**: Chance to deal increased damage based on weapon and skills
- **Miss Chance**: Possibility of attacks failing to connect

### Defense System

- **Armor Class (AC)**: Determines difficulty to hit the character
- **Damage Reduction**: Some armor reduces incoming damage
- **Saving Throws**: Chances to resist or reduce certain effects
- **Evasion**: Some classes have abilities to avoid damage entirely

## Character State Management

- **Persistence**: Character data is stored on the server
- **Dungeon State**: Characters track their position and progress in dungeons
- **Multiple Characters**: Players can create up to 10 characters per account
- **Character Deletion**: Players can permanently remove characters they no longer want

## Technical Implementation

- **Data Structure**: Characters are represented as objects with properties and methods
- **Serialization**: Character data can be saved and loaded for game persistence
- **Client-Server Model**: Character logic runs on the server with state updates sent to clients
- **Validation**: Server-side validation prevents cheating and ensures data integrity

## Visual Representation

- **Class-Based Styling**: Each class has a distinct visual style and color scheme
- **Equipment Visualization**: Equipped items affect character appearance
- **Status Indicators**: Visual cues for character health, mana, and status effects
- **Animation**: Visual feedback for character actions and state changes
