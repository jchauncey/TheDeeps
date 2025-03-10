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

## Character selection screen
- Displays all available characters in a grid layout
- Each character card shows the character's name, class, and a class-colored badge
- Characters can be selected by clicking on their card
- A delete button on each character card allows for character deletion
- Deletion requires confirmation via a modal to prevent accidental deletions
- The screen uses fixed positioning to fill the entire browser window
- Maximum of 10 characters allowed per account

## Dungeon selection screen
- Accessible after selecting a character
- Displays a list of available dungeons with details (name, floors, creation date, player count)
- Allows creation of new dungeons with custom name and number of floors
- Uses REST API calls instead of WebSockets for dungeon creation and joining
- Dungeons can be selected and joined with the selected character
- The screen uses fixed positioning to fill the entire browser window
- Consistent styling with the character selection screen

## REST API Endpoints

### Character Management
- `GET /characters` - Retrieves a list of all characters for the current user
- `GET /characters/{id}` - Retrieves a specific character by ID
- `POST /characters` - Creates a new character
- `DELETE /characters/{id}` - Deletes a character by ID

### Dungeon Management
- `GET /dungeons` - Retrieves a list of available dungeons
- `POST /dungeons` - Creates a new dungeon
- `POST /dungeons/{id}/join` - Joins a dungeon with a specified character
- `GET /dungeons/{id}/floor/{level}` - Retrieves data for a specific floor in a dungeon

### Game State
- `POST /characters/{id}/save` - Saves the current state of a character
- `GET /characters/{id}/floor` - Retrieves the current floor data for a character

## Implementation Details
- Character creation now properly avoids duplicate API calls by having the CharacterCreation component pass data to the parent component
- Character deletion includes proper error handling and user feedback
- Dungeon creation and joining use REST API calls instead of WebSockets
- All screens use consistent styling and layout to fill the entire browser window
