# AI Refactor

This is the requirements for refactoring the current state

- Prefer rest calls for static initalization instead of websockets
- modularize the code so that it can be tested
- When the game starts it should fetch a list of your characters from the server
- Create character should send the character data to the sever
- Once you create a character it should take you to the character selection screen
- From here you can choose a character and then choose an existing dungeon or to create a new dungeon
- These steps should be done without using websockets
- Once the character is in the dungeon the game can utilize websockets to maintain game state.

## Character creation workflow 
- Loading screen should display the deeps logo
- Fetch existing characters and list them on the right for selection
- Provide a create character button if no character exists or they are below the 10 character limit
