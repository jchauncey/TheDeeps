# TheDeeps UI

This is the general UI definition -

- React and Typescript front end that uses websockets to connect to a backend API server
- Recieves the map from the api server
- Map window is anchored to the left of the browser window
- The main browser window should avoid heights that add scroll bar
- Toast banners should be clickable to dismiss
- Toast banners should only stay for 5 seconds unless clicked
- It should have a modal window that displays the hotkeys
- Avoid hotkeys that interfere with the movement keys wasd
- Do not use the arrow keys as movement keys
- Avoid the F1-F12 keys as hotkeys
- Common actions should have single key hotkeys
- Using complex hotkeys such as ctrl+g is fine but should be reserved for less used actions

It should have a loading screen with the following elements - 
- It should have The Deeps logo whichi s stored at client/public/logo.png
- It has a New Game and Load Game button which are side by side 
- Clicking New Game takes you to the Character creation screen
- The Load Game button isnt implemented yet

Character Creation Screen - 
- Allow a user to specify a name
- Also provide a Random Name generator
- Provide a class selection drop down based on the base classes from D&D
- Provide a attributes allocation area
- Attributes are automatically allocated based on the class selected with 5 remaining points that the user can spend
- Show which attributes are the primary for the class selected
- Show the modifier value for attributes given their current value

Character Profile Window - 
- Character profile window is anchored to the right of the browser window.
- The character profile window may have a scroll bar
- It should display health, mana, experience, and AC in the first section
- It should display current attirbutes and modifier values
- It should display buffs for attributes and health/mana/ac
- It should list any primary skills the character has
- It should show how much gold the character currently has
- It should have a counter for how many potions the character is currently holding

Dungeon Window - 
- The dungeon window should render the map like a nethack or rogue style game. 
- The character icon should be a stylized @ where each class is a unique color
- It should not have a fog of war
