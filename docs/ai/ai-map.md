# AI Map
## Map Generation Requirements

1. **Map Structure**
   - Maps must have a minimum of 5 rooms and a maximum of 20 rooms per floor
   - Each floor should have at least one staircase to the next floor
   - The final floor must contain a boss room
   - Maps should have a logical layout with no disconnected areas
   - Hallways should connect rooms in a sensible manner
   - ✅ Each floor must have both up and down staircases (except for the first and last floors)

2. **Room Types**
   - Standard rooms (empty or with basic enemies)
   - Treasure rooms (containing loot)
   - Boss rooms (containing a powerful enemy)
   - Puzzle rooms (requiring player to solve a puzzle)
   - Safe rooms (no enemies, may contain healing resources)
   - Shop rooms (where players can trade resources)

3. **Room Properties**
   - Each room must have a unique identifier
   - Rooms should have variable sizes (small, medium, large)
   - Rooms can contain environmental hazards (traps, lava, etc.)
   - Rooms may contain interactive objects (chests, levers, etc.)
   - Rooms should have appropriate enemy spawns based on difficulty level

4. **Doors and Connections**
   - Doors can be standard, locked, or hidden
   - Locked doors require specific keys or conditions to open
   - Hidden doors require player discovery (through skills or items)
   - Some doors may be one-way only
   - Doors should visually match their environment

5. **Difficulty Scaling**
   - Lower floors should be easier than deeper floors
   - Enemy density and strength should increase with depth
   - Loot quality should improve with depth
   - Environmental hazards should become more dangerous with depth
   - Puzzles should become more complex with depth

6. **Procedural Generation Parameters**
   - Maps should be generated with a seed value for reproducibility
   - Generation should account for desired difficulty level
   - Maps should have themes that affect visual appearance and enemy types
   - Generation should ensure all areas are accessible
   - Special rooms (boss, treasure) should be appropriately placed
   - ✅ Staircases should be placed in accessible locations on each floor

7. **Navigation and Exploration**
   - Unexplored areas should be hidden until discovered
   - Map should track discovered rooms, doors, and stairs
   - Map should be fully visible when character profile modal is open
   - Player position should be clearly indicated on the map
   - ✅ Staircases should be clearly visible and distinguishable from other map elements

8. **Technical Requirements**
   - Map data structure should be serializable for saving/loading
   - Generation should be efficient and complete in under 5 seconds
   - Map should support dynamic modifications during gameplay
   - Collision detection system for walls and obstacles
   - Pathfinding system for AI navigation
   - Map rendering should be optimized for performance
   - Map should maintain visibility when UI modals are displayed
   - ✅ Map should track player positions across multiple floors

9. **UI Integration**
   - Map should be the primary focus of the game interface
   - Map should occupy the majority of the screen space
   - Map should be responsive to window size changes
   - Map should maintain proper aspect ratio
   - Map should have clear visual indicators for player, enemies, and objects
   - Map should be accessible via keyboard navigation
   - Map should support hotkeys for common actions (movement, interaction)
   - ✅ Map should update smoothly when transitioning between floors

10. **Floor Navigation**
    - ✅ Players can move between floors using staircases
    - ✅ Up staircases allow players to move to the previous floor
    - ✅ Down staircases allow players to move to the next floor
    - ✅ Players are positioned at the corresponding staircase when changing floors
    - ✅ The first floor only has down staircases
    - ✅ The last floor only has up staircases
    - ✅ Middle floors have both up and down staircases
    - ✅ Floor transitions maintain player state (health, inventory, etc.)