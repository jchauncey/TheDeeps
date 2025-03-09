# AI API

This sets the requirements for the server API

## Core Requirements

- ✅ Should have a generate Map endpoint that will generate a complete dungeon and return it to the client.
- ✅ Should have a current floor endpoint that will return the floor of the current character. This endpoint should allow for mutiple characters to be on multiple floors.
- ✅ Should have an endpoint that tracks the movement of characters and validates their position on the map.
- Should have an endpoint for combat mechanics related to player vs monsters
- ✅ Should support floor navigation via ascend and descend actions, allowing players to move between dungeon floors

## Implementation Guidelines

- ✅ Real-time game mechanics (movement, combat, actions) should be implemented using WebSockets for optimal performance and responsiveness.
- ✅ Non-real-time operations (character creation/retrieval, initial data loading) should be implemented as REST endpoints.
- ✅ The server should maintain a hybrid approach: WebSockets for dynamic gameplay and REST for stateless operations.
- ✅ WebSocket connections should handle reconnection attempts and message queuing for reliability.
- ✅ REST endpoints should be provided for initial game setup and as fallbacks when WebSocket connections fail.
- ✅ Map generation should be triggered from the client side, not automatically by the server. This allows users to either start a new dungeon or join an existing one.

## Multi-Character Support

- ✅ The server should support multiple characters being on different floors simultaneously.
- ✅ Each character should have its own game state, including current floor, position, and visibility.
- ✅ The server should efficiently broadcast updates only to relevant clients based on their character's location.
- ✅ The server should provide a mechanism for users to discover and join existing dungeons with other players.
- ✅ Characters can independently navigate between floors using stairs, with their positions tracked separately

## Performance Considerations

- ✅ Map generation should be optimized to minimize latency when creating new dungeons.
- ✅ The server should implement efficient visibility calculations to reduce bandwidth usage.
- Combat calculations should be performed server-side to prevent cheating.
- ✅ The server should handle concurrent connections from multiple clients without significant performance degradation.

## Dungeon Management

- ✅ Clients should be able to request the creation of a new dungeon instance.
- ✅ Clients should be able to request a list of active dungeons they can join.
- ✅ When joining an existing dungeon, the client should receive the current state of that dungeon.
- ✅ Players joining an existing dungeon should always start at the first floor, regardless of where other players are located.
- ✅ Dungeons should persist even when no players are present, but may be cleaned up after a configurable period of inactivity.
- ✅ Each dungeon should have a unique identifier that clients can use to join specific instances.

## Floor Navigation

- ✅ Players can ascend to previous floors by standing on up stairs and using the "ascend" action
- ✅ Players can descend to deeper floors by standing on down stairs and using the "descend" action
- ✅ When ascending or descending, players are positioned at the corresponding stairs on the destination floor
- ✅ The server validates that players are standing on the appropriate stairs before allowing floor transitions
- ✅ The server updates the player's floor and position in the dungeon instance when changing floors
- ✅ The server sends updated floor data to the client after a successful floor transition
- ✅ The server prevents players from ascending above the first floor or descending below the last floor