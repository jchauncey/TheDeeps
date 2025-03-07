# TheDeeps UI

## UI Requirements

1. **Technical Foundation**
   - React and TypeScript front-end architecture
   - WebSocket connection to backend API server
   - Responsive design with minimum resolution support of 1280x720
   - Cross-browser compatibility (Chrome, Firefox, Safari, Edge)
   - Mobile-friendly layout with touch controls where appropriate
   - Optimized asset loading for faster initial render

2. **General UI Principles**
   - Map window anchored to the left of the browser window
   - Character profile window anchored to the right
   - Main browser window should avoid heights that add scroll bars
   - Consistent color scheme and typography throughout the application
   - Accessible design with proper contrast ratios and screen reader support
   - Dark mode support with automatic detection of system preferences
   - Consistent padding and margins throughout the interface

3. **Notification System**
   - Toast banners should be clickable to dismiss
   - Toast banners should automatically disappear after 5 seconds unless clicked
   - Different notification types (info, warning, error, success) with appropriate styling
   - Notification queue system for handling multiple notifications
   - Critical notifications should remain until dismissed
   - Notifications should not obstruct gameplay elements

4. **Input Controls**
   - Modal window that displays all available hotkeys
   - Avoid hotkeys that interfere with the movement keys (WASD)
   - Do not use arrow keys as movement keys
   - Avoid F1-F12 keys as hotkeys
   - Common actions should have single key hotkeys
   - Complex hotkeys (Ctrl+key) reserved for less frequently used actions
   - Customizable keybindings with reset to default option
   - Mouse controls for all actions with appropriate tooltips

5. **Loading Screen**
   - Display The Deeps logo (stored at client/public/logo.png)
   - New Game and Load Game buttons positioned side by side
   - New Game button navigates to Character Creation screen
   - Load Game button (to be implemented in future)
   - Loading progress indicator for initial asset loading
   - Background artwork that sets the game's tone
   - Version number and developer credits
   - Tips or lore snippets during loading

6. **Character Creation Screen**
   - User-friendly name input field
   - Random Name generator with thematic options
   - Class selection dropdown featuring all D&D-based classes
   - Visual representation of selected character class
   - Attributes allocation area with intuitive controls
   - Automatic attribute allocation based on class with 5 remaining points for user customization
   - Clear indication of primary attributes for selected class
   - Display of modifier values for each attribute
   - Preview of starting equipment and abilities
   - Confirmation button with validation for required fields

7. **Character Profile Window**
   - Anchored to the right of the browser window
   - May include a scroll bar for additional information
   - Character portrait or icon at the top
   - First section displays health, mana, experience, and AC with visual indicators
   - Current attributes and modifier values clearly displayed
   - Visual indicators for buffs affecting attributes and stats
   - List of primary skills with proficiency levels
   - Gold counter with appropriate icon
   - Potion counter with visual representation
   - Equipment section showing equipped items
   - Inventory section with grid-based layout
   - Character level and progress to next level

8. **Dungeon Window**
   - Renders map in a nethack or rogue-style ASCII/tile-based format
   - Character represented by a stylized @ symbol with class-specific colors
   - No fog of war implementation
   - Clear visual distinction between different terrain types
   - Animated transitions between map areas
   - Zoom functionality for map view
   - Mini-map in corner for larger dungeons
   - Visual indicators for interactive objects (doors, chests, etc.)
   - Enemy representations with distinct symbols/colors
   - Turn counter or game clock display

9. **Combat Interface**
   - Visual feedback for attacks and damage
   - Health/mana bars for player and visible enemies
   - Combat log showing recent actions
   - Quick-access bar for common combat abilities
   - Initiative order display for turn-based combat
   - Target selection mechanism
   - Status effect icons with duration indicators
   - Critical hit and miss animations/indicators

10. **Menu System**
    - Easily accessible main menu via hotkey or button
    - Settings menu with audio, video, and gameplay options
    - Save/Load game functionality
    - Character sheet with detailed statistics
    - Quest log with active and completed quests
    - Help section with game mechanics explanation
    - Exit game confirmation dialog
