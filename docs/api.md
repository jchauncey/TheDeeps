# The Deeps API Documentation

This document provides a comprehensive guide to all API endpoints available in The Deeps server.

## Table of Contents

- [Character Endpoints](#character-endpoints)
- [Dungeon Endpoints](#dungeon-endpoints)
- [Inventory Endpoints](#inventory-endpoints)
- [Combat Endpoints](#combat-endpoints)
- [WebSocket Endpoints](#websocket-endpoints)
- [Testing Endpoints](#testing-endpoints)

## Character Endpoints

### Get All Characters
- **URL**: `/characters`
- **Method**: `GET`
- **Description**: Returns a list of all characters.
- **Response**: Array of character objects.

### Create Character
- **URL**: `/characters`
- **Method**: `POST`
- **Description**: Creates a new character.
- **Request Body**:
  ```json
  {
    "name": "string",
    "class": "string",
    "attributes": {
      "strength": number,
      "dexterity": number,
      "constitution": number,
      "intelligence": number,
      "wisdom": number,
      "charisma": number
    }
  }
  ```
- **Response**: Created character object.

### Get Character
- **URL**: `/characters/{id}`
- **Method**: `GET`
- **Description**: Returns a specific character by ID.
- **URL Parameters**: `id` - Character ID.
- **Response**: Character object.

### Delete Character
- **URL**: `/characters/{id}`
- **Method**: `DELETE`
- **Description**: Deletes a character.
- **URL Parameters**: `id` - Character ID.
- **Response**: No content on success.

### Save Character
- **URL**: `/characters/{id}/save`
- **Method**: `POST`
- **Description**: Saves the current state of a character.
- **URL Parameters**: `id` - Character ID.
- **Request Body**:
  ```json
  {
    "position": { "x": number, "y": number },
    "currentHp": number,
    "currentMana": number,
    "gold": number,
    "experience": number,
    "currentFloor": number,
    "currentDungeon": "string"
  }
  ```
- **Response**: Updated character object.

### Get Character Floor
- **URL**: `/characters/{id}/floor`
- **Method**: `GET`
- **Description**: Returns the current floor number for a character.
- **URL Parameters**: `id` - Character ID.
- **Response**: 
  ```json
  {
    "floor": number
  }
  ```

## Dungeon Endpoints

### Get All Dungeons
- **URL**: `/dungeons`
- **Method**: `GET`
- **Description**: Returns a list of all dungeons.
- **Response**: Array of dungeon objects.

### Create Dungeon
- **URL**: `/dungeons`
- **Method**: `POST`
- **Description**: Creates a new dungeon.
- **Request Body**:
  ```json
  {
    "name": "string",
    "floors": number,
    "difficulty": "string",
    "seed": number (optional)
  }
  ```
- **Response**: Created dungeon object.

### Join Dungeon
- **URL**: `/dungeons/{id}/join`
- **Method**: `POST`
- **Description**: Makes a character join a dungeon.
- **URL Parameters**: `id` - Dungeon ID.
- **Request Body**:
  ```json
  {
    "characterId": "string"
  }
  ```
- **Response**: Success status with character's starting position.

### Get Floor
- **URL**: `/dungeons/{id}/floor/{level}`
- **Method**: `GET`
- **Description**: Returns a specific floor of a dungeon.
- **URL Parameters**: 
  - `id` - Dungeon ID.
  - `level` - Floor level.
- **Response**: Floor object with tiles, rooms, mobs, and items.

### Get Floor By Number (Alternative)
- **URL**: `/api/dungeons/{id}/floors/{floorNumber}`
- **Method**: `GET`
- **Description**: Alternative endpoint to get a specific floor of a dungeon.
- **URL Parameters**: 
  - `id` - Dungeon ID.
  - `floorNumber` - Floor number.
- **Response**: Floor object with tiles, rooms, mobs, and items.

## Inventory Endpoints

### Get Inventory
- **URL**: `/api/characters/{characterID}/inventory`
- **Method**: `GET`
- **Description**: Returns a character's inventory.
- **URL Parameters**: `characterID` - Character ID.
- **Response**: Array of item objects.

### Get Inventory Item
- **URL**: `/api/characters/{characterID}/inventory/{itemID}`
- **Method**: `GET`
- **Description**: Returns a specific item from a character's inventory.
- **URL Parameters**: 
  - `characterID` - Character ID.
  - `itemID` - Item ID.
- **Response**: Item object.

### Equip Item
- **URL**: `/api/characters/{characterID}/inventory/{itemID}/equip`
- **Method**: `POST`
- **Description**: Equips an item from a character's inventory.
- **URL Parameters**: 
  - `characterID` - Character ID.
  - `itemID` - Item ID.
- **Response**: Success status.

### Unequip Item
- **URL**: `/api/characters/{characterID}/inventory/{itemID}/unequip`
- **Method**: `POST`
- **Description**: Unequips an item from a character.
- **URL Parameters**: 
  - `characterID` - Character ID.
  - `itemID` - Item ID.
- **Response**: Success status.

### Use Item
- **URL**: `/api/characters/{characterID}/inventory/{itemID}/use`
- **Method**: `POST`
- **Description**: Uses an item from a character's inventory.
- **URL Parameters**: 
  - `characterID` - Character ID.
  - `itemID` - Item ID.
- **Response**: Success status.

### Get Equipment
- **URL**: `/api/characters/{characterID}/equipment`
- **Method**: `GET`
- **Description**: Returns a character's equipped items.
- **URL Parameters**: `characterID` - Character ID.
- **Response**: Equipment object with items equipped in different slots.

### Add Item To Inventory
- **URL**: `/api/characters/{characterID}/inventory/add`
- **Method**: `POST`
- **Description**: Adds an item to a character's inventory.
- **URL Parameters**: `characterID` - Character ID.
- **Request Body**: Item object.
- **Response**: Updated inventory.

### Get Character Weight
- **URL**: `/api/characters/{characterID}/weight`
- **Method**: `GET`
- **Description**: Returns a character's current weight (from inventory items).
- **URL Parameters**: `characterID` - Character ID.
- **Response**: 
  ```json
  {
    "currentWeight": number,
    "maxWeight": number
  }
  ```

### Get All Items
- **URL**: `/api/items`
- **Method**: `GET`
- **Description**: Returns all items in the repository.
- **Response**: Array of item objects.

### Generate Items
- **URL**: `/api/items/generate`
- **Method**: `POST`
- **Description**: Generates random items based on floor level.
- **Request Body**:
  ```json
  {
    "count": number,
    "floorLevel": number
  }
  ```
- **Response**: Array of generated item objects.

## Combat Endpoints

### Get Combat State
- **URL**: `/characters/{id}/combat`
- **Method**: `GET`
- **Description**: Returns the combat state for a character.
- **URL Parameters**: `id` - Character ID.
- **Response**: 
  ```json
  {
    "character": {Character Object},
    "nearbyMobs": {Object of Mob Objects},
    "inCombat": boolean
  }
  ```

## WebSocket Endpoints

### Combat WebSocket
- **URL**: `/ws/combat`
- **Method**: `WebSocket`
- **Description**: Handles real-time combat interactions.
- **Connection Parameters**: None
- **Client-to-Server Messages**:
  ```json
  {
    "action": "attack" | "useItem" | "flee",
    "characterId": "string",
    "mobId": "string" (for attack/flee),
    "itemId": "string" (for useItem)
  }
  ```
- **Server-to-Client Messages**:
  ```json
  {
    "action": "string",
    "success": boolean,
    "message": "string",
    "result": {
      "damageDealt": number,
      "damageTaken": number,
      "criticalHit": boolean,
      "killed": boolean,
      "expGained": number,
      "goldGained": number,
      "itemsDropped": [Item Objects]
    }
  }
  ```

### Game WebSocket
- **URL**: `/ws/game`
- **Method**: `WebSocket`
- **Description**: Handles real-time game state updates.
- **Connection Parameters**: `characterId` - Character ID (query parameter).
- **Client-to-Server Messages**:
  ```json
  {
    "type": "move" | "attack" | "pickup" | "useItem" | "dropItem" | "equipItem" | "unequipItem" | "ascend" | "descend",
    "characterId": "string",
    "direction": "up" | "down" | "left" | "right" (for move),
    "targetId": "string" (mob or item ID),
    "itemId": "string" (for item-related actions)
  }
  ```
- **Server-to-Client Messages**:
  ```json
  {
    "type": "updateMap" | "updatePlayer" | "updateMob" | "removeMob" | "addItem" | "removeItem" | "notification" | "floorUpdate" | "floorChange" | "error" | "initialState",
    "character": {Character Object},
    "floor": {Floor Object},
    "mob": {Mob Object},
    "item": {Item Object},
    "text": "string",
    "error": "string"
  }
  ```

## Testing Endpoints

### Generate Test Room
- **URL**: `/test/room`
- **Method**: `GET`
- **Description**: Generates a test room for client rendering testing.
- **Query Parameters**:
  - `type` - Room type (entrance, standard, treasure, boss, safe, shop). Default: entrance.
  - `width` - Room width (1-100). Default: 20.
  - `height` - Room height (1-100). Default: 20.
  - `roomWidth` - Inner room width. Default: 8 for entrance, 7 for others.
  - `roomHeight` - Inner room height. Default: 8 for entrance, 7 for others.
- **Response**: Floor object with a single room of the specified type. 