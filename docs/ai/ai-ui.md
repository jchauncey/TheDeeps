# TheDeeps UI

## UI Requirements

1. **Technical Foundation**
   - ✅ React and TypeScript front-end architecture
   - ✅ WebSocket connection to backend API server
   - ✅ Responsive design with minimum resolution support of 1280x720
   - ✅ Cross-browser compatibility (Chrome, Firefox, Safari, Edge)
   - Mobile-friendly layout with touch controls where appropriate
   - Optimized asset loading for faster initial render

2. **General UI Principles**
   - ✅ Map window anchored to the left of the browser window
   - ✅ Simplified character status panel anchored to the right
   - Main browser window should avoid heights that add scroll bars
   - ✅ Consistent color scheme and typography throughout the application
   - Accessible design with proper contrast ratios and screen reader support
   - Dark mode support with automatic detection of system preferences
   - ✅ Consistent padding and margins throughout the interface
   - ✅ Text colors should have sufficient contrast with background colors
   - ✅ Full-screen layouts for character selection and dungeon selection screens

3. **Notification System**
   - ✅ Toast banners should be clickable to dismiss
   - Toast banners should automatically disappear after 5 seconds unless clicked
   - ✅ Different notification types (info, warning, error, success) with appropriate styling
   - Notification queue system for handling multiple notifications
   - Critical notifications should remain until dismissed
   - ✅ Notifications should not obstruct gameplay elements
   - ✅ Floor transition notifications should display the floor number being entered

4. **Input Controls**
   - ✅ Modal window that displays all available hotkeys
   - ✅ Avoid hotkeys that interfere with the movement keys (WASD)
   - ✅ Do not use arrow keys as movement keys
   - ✅ Avoid F1-F12 keys as hotkeys
   - ✅ Common actions should have single key hotkeys
   - Complex hotkeys (Ctrl+key) reserved for less frequently used actions
   - Customizable keybindings with reset to default option
   - ✅ Mouse controls for all actions with appropriate tooltips
   - ✅ 'C' key hotkey to open character profile modal
   - ✅ Specific hotkeys for ascending ('u' key) and descending ('d' key) stairs

5. **Loading Screen**
   - ✅ Display The Deeps logo (stored at client/public/logo.png)
   - ✅ New Game and Load Game buttons positioned side by side
   - ✅ New Game button navigates to Character Creation screen
   - ✅ Load Game button (to be implemented in future)
   - Loading progress indicator for initial asset loading
   - ✅ Background artwork that sets the game's tone
   - Version number and developer credits
   - Tips or lore snippets during loading

6. **Character Creation Screen**
   - ✅ User-friendly name input field
   - Random Name generator with thematic options
   - ✅ Class selection dropdown featuring all D&D-based classes
   - ✅ Visual representation of selected character class
   - ✅ Attributes allocation area with intuitive controls
   - ✅ Automatic attribute allocation based on class with 5 remaining points for user customization
   - ✅ Clear indication of primary attributes for selected class
   - ✅ Display of modifier values for each attribute
   - ✅ Preview of starting equipment and abilities
   - ✅ Confirmation button with validation for required fields
   - ✅ Proper API integration to avoid duplicate character creation

7. **Character Selection Screen**
   - ✅ Grid layout displaying all available characters
   - ✅ Character cards showing name, class, and class-colored badge
   - ✅ Selection mechanism by clicking on character cards
   - ✅ Delete button on each character card with confirmation modal
   - ✅ Full-screen layout with fixed positioning
   - ✅ Maximum of 10 characters allowed per account
   - ✅ Create New Character button when below character limit
   - ✅ Back button to return to title screen
   - ✅ Select Character button to proceed with selected character
   - ✅ Proper error handling and user feedback for character deletion

8. **Dungeon Selection Screen**
   - ✅ Full-screen layout with fixed positioning
   - ✅ List of available dungeons with details (name, floors, creation date, player count)
   - ✅ Dungeon creation form with name and floor count inputs
   - ✅ Selection mechanism for choosing a dungeon
   - ✅ Join Selected Dungeon button to enter the dungeon
   - ✅ Back button to return to character selection
   - ✅ Consistent styling with other screens
   - ✅ REST API integration for dungeon creation and joining
   - ✅ Loading indicators for async operations
   - ✅ Error handling and user feedback

9. **Character Status Panel**
   - ✅ Anchored to the right of the browser window
   - ✅ Simplified design showing only essential information
   - ✅ Character name and level at the top
   - ✅ Health bar with current/maximum values
   - ✅ Mana bar with current/maximum values
   - Experience bar with progress to next level
   - ✅ Profile button with class-specific icon
   - ✅ Class-specific styling and colors
   - ✅ Responsive design that works well on smaller screens
   - ✅ Current floor indicator showing which floor the player is on

10. **Character Profile Modal**
    - ✅ Accessible via 'C' hotkey or Profile button in status panel
    - ✅ Overlays the game map when opened
    - ✅ Detailed character information in a scrollable modal
    - ✅ Class-specific icon displayed in a square container
    - ✅ Thematic icons for each character class (warrior, mage, rogue, etc.)
    - ✅ Character name, level, and class prominently displayed
    - ✅ Text colors with sufficient contrast against background colors
    - ✅ Attributes section with modifier values
    - Class abilities section with visual indicators
    - Equipment section showing equipped items
    - Inventory section with grid-based layout
    - Gold and potion counters with appropriate icons
    - ✅ Close button to return to gameplay
    - ✅ Display of current floor and total dungeon depth

11. **Dungeon Window**
    - ✅ Renders map in a nethack or rogue-style ASCII/tile-based format
    - ✅ Character represented by a stylized @ symbol with class-specific colors
    - ✅ No fog of war implementation
    - ✅ Clear visual distinction between different terrain types
    - ✅ Animated transitions between map areas
    - ✅ Visual indicators for interactive objects (doors, chests, etc.)
    - ✅ Enemy representations with distinct symbols/colors
    - Turn counter or game clock display
    - ✅ Distinct visual representation for up stairs (< symbol) and down stairs (> symbol)
    - ✅ Smooth transition animations when changing floors

12. **Combat Interface**
    - ✅ Visual feedback for attacks and damage
    - ✅ Health/mana bars for player and visible enemies
    - Combat log showing recent actions
    - Quick-access bar for common combat abilities
    - Initiative order display for turn-based combat
    - ✅ Target selection mechanism
    - Status effect icons with duration indicators
    - Critical hit and miss animations/indicators

13. **Menu System**
    - ✅ Easily accessible main menu via hotkey or button
    - Settings menu with audio, video, and gameplay options
    - ✅ Save/Load game functionality
    - ✅ Help section with game mechanics explanation and hotkey reference
    - ✅ Exit game confirmation dialog
    - ✅ Controls reference section explaining floor navigation mechanics

14. **Floor Navigation**
    - ✅ Visual indication when player is standing on stairs
    - ✅ Tooltip or prompt showing available actions when standing on stairs
    - ✅ Feedback message when attempting to ascend from the top floor or descend from the bottom floor
    - ✅ Floor transition animation when moving between floors
    - ✅ Floor number indicator that updates when changing floors
    - ✅ Action buttons for ascending and descending when standing on appropriate stairs
