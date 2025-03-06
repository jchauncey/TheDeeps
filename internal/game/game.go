package game

import (
	"encoding/json"
	"fmt"

	"github.com/jchauncey/TheDeeps/internal/character"
	"github.com/jchauncey/TheDeeps/internal/input"
	mapgen "github.com/jchauncey/TheDeeps/internal/map"
	"github.com/jchauncey/TheDeeps/internal/ui"
)

// GameState represents the current state of the game
type GameState int

const (
	StateMainMenu GameState = iota
	StateCharacterCreation
	StatePlaying
	StateInventory
	StateGameOver
	StateQuit
)

// Game represents the main game structure
type Game struct {
	Player       *character.Player
	Dungeon      *mapgen.Dungeon
	CurrentFloor int
	Renderer     *ui.Renderer
	State        GameState
	Running      bool

	// Character creation fields
	TempName      string
	TempClass     character.ClassType
	CreationStage int // 0 = class selection, 1 = name entry
}

// GameStateResponse represents the game state sent to the client
type GameStateResponse struct {
	Player struct {
		X, Y  int    `json:"x,y"`
		HP    int    `json:"hp"`
		MaxHP int    `json:"maxHp"`
		Name  string `json:"name"`
		Class string `json:"class"`
	} `json:"player"`
	Map struct {
		Tiles  [][]string `json:"tiles"`
		Width  int        `json:"width"`
		Height int        `json:"height"`
	} `json:"map"`
	Messages []string  `json:"messages"`
	State    GameState `json:"state"`
}

// NewGame creates a new game instance
func NewGame() *Game {
	// Create renderer
	renderer := ui.NewRenderer(nil, nil)

	// Add welcome message
	renderer.AddSystemMessage("Welcome to The Deeps! Create your character to begin.")

	// Return new game
	return &Game{
		CurrentFloor:  0,
		Renderer:      renderer,
		State:         StateCharacterCreation,
		Running:       true,
		TempClass:     character.Warrior, // Default selection
		CreationStage: 0,                 // Start with class selection
	}
}

// GetState returns the current game state in a format suitable for the client
func (g *Game) GetState() []byte {
	response := GameStateResponse{}

	if g.Player != nil {
		response.Player.X = g.Player.X
		response.Player.Y = g.Player.Y
		response.Player.HP = g.Player.HP
		response.Player.MaxHP = g.Player.MaxHP
		response.Player.Name = g.Player.Name
		response.Player.Class = string(g.TempClass)
	}

	if g.Dungeon != nil && g.CurrentFloor >= 0 && g.CurrentFloor < len(g.Dungeon.Floors) {
		currentFloor := g.Dungeon.Floors[g.CurrentFloor]
		response.Map.Width = len(currentFloor.Level.Tiles[0])
		response.Map.Height = len(currentFloor.Level.Tiles)
		response.Map.Tiles = make([][]string, response.Map.Height)
		for y := range currentFloor.Level.Tiles {
			response.Map.Tiles[y] = make([]string, response.Map.Width)
			for x := range currentFloor.Level.Tiles[y] {
				tile := currentFloor.Level.Tiles[y][x]
				switch tile.Type {
				case mapgen.TileWall:
					response.Map.Tiles[y][x] = "#"
				case mapgen.TileFloor:
					response.Map.Tiles[y][x] = "."
				case mapgen.TileEntrance:
					response.Map.Tiles[y][x] = "<"
				case mapgen.TileExit:
					response.Map.Tiles[y][x] = ">"
				case mapgen.TileHallway:
					response.Map.Tiles[y][x] = "."
				case mapgen.TilePillar:
					response.Map.Tiles[y][x] = "O"
				case mapgen.TileWater:
					response.Map.Tiles[y][x] = "~"
				case mapgen.TileRubble:
					response.Map.Tiles[y][x] = ","
				default:
					response.Map.Tiles[y][x] = " "
				}
			}
		}
	}

	response.State = g.State

	data, err := json.Marshal(response)
	if err != nil {
		return []byte("{}")
	}
	return data
}

// StartGame initializes the game after character creation
func (g *Game) StartGame() {
	// Debug output
	fmt.Printf("Starting game with character: %s (class: %s)\n",
		g.TempName, character.GetClassName(g.TempClass))

	// Create player with chosen name and class
	g.Player = character.NewPlayer(0, 0, g.TempClass)
	g.Player.Name = g.TempName
	fmt.Printf("Player created: %s at position (%d, %d)\n", g.Player.Name, g.Player.X, g.Player.Y)

	// Create dungeon
	g.Dungeon = mapgen.NewDungeon(5) // 5 floors
	fmt.Printf("Dungeon created with %d floors\n", len(g.Dungeon.Floors))

	// Update renderer with player and dungeon
	g.Renderer.Player = g.Player
	g.Renderer.Dungeon = g.Dungeon
	fmt.Println("Renderer updated with player and dungeon")

	// Place player at entrance of first floor
	entrance := g.Dungeon.Floors[0].Entrance
	g.Player.X = entrance.X
	g.Player.Y = entrance.Y
	fmt.Printf("Player placed at entrance: (%d, %d)\n", g.Player.X, g.Player.Y)

	// Add game start messages
	g.Renderer.AddSystemMessage("Game initialized successfully")
	g.Renderer.AddGameMessage("Welcome, %s the %s! Your adventure begins!",
		g.TempName, character.GetClassName(g.TempClass))

	// Change state to playing
	g.State = StatePlaying

	// Debug output
	fmt.Println("Game state changed to Playing")
}

// Run starts the game loop
func (g *Game) Run() {
	// Main game loop
	for g.Running {
		// Debug output for current game state
		fmt.Printf("Current game state: %d\n", g.State)

		// Render the game
		g.Renderer.RenderGame()

		// Handle input based on current state
		switch g.State {
		case StateCharacterCreation:
			fmt.Println("Handling character creation input")
			g.handleCharacterCreationInput()
		case StatePlaying:
			fmt.Println("Handling playing input")
			g.handlePlayingInput()
		case StateInventory:
			fmt.Println("Handling inventory input")
			g.handleInventoryInput()
		case StateGameOver:
			fmt.Println("Handling game over input")
			g.handleGameOverInput()
		case StateMainMenu:
			fmt.Println("Handling main menu input")
			g.handleMainMenuInput()
		case StateQuit:
			fmt.Println("Quitting game")
			g.Running = false
		}
	}
}

// handleCharacterCreationInput handles input during character creation
func (g *Game) handleCharacterCreationInput() {
	// Get input
	key := input.GetSingleKey()
	fmt.Printf("Character creation input: '%s', stage: %d\n", key, g.CreationStage)

	// Name entry input received
	if g.CreationStage == 1 {
		fmt.Printf("Name entry input received: '%s'\n", key)

		// Handle backspace
		if key == "backspace" && len(g.TempName) > 0 {
			g.TempName = g.TempName[:len(g.TempName)-1]
			g.Renderer.CreationName = g.TempName
			return
		}

		// Handle escape (go back to class selection)
		if key == "escape" {
			g.CreationStage = 0
			g.TempName = ""
			g.Renderer.CreationName = ""
			g.Renderer.AddSystemMessage("Returned to class selection")
			return
		}

		// Handle enter (confirm name)
		if key == "enter" || key == "" {
			fmt.Println("Empty key detected (Enter pressed)")

			// Validate name
			if g.TempName == "" {
				g.Renderer.AddSystemMessage("Please enter a name for your character.")
				return
			}

			fmt.Printf("Starting game with name: %s\n", g.TempName)
			g.StartGame()
			return
		}

		// Add character to name if it's a valid input and not too long
		if len(key) == 1 && len(g.TempName) < 20 {
			g.TempName += key
			g.Renderer.CreationName = g.TempName
		}

		return
	}

	// Class selection
	switch key {
	case "1": // Warrior
		g.TempClass = character.Warrior
		g.CreationStage = 1
		g.Renderer.AddGameMessage("Selected Warrior: Strong and resilient fighters.")
		g.Renderer.AddSystemMessage("Enter your character's name:")
	case "2": // Wizard
		g.TempClass = character.Wizard
		g.CreationStage = 1
		g.Renderer.AddGameMessage("Selected Wizard: Powerful wielders of arcane magic.")
		g.Renderer.AddSystemMessage("Enter your character's name:")
	case "3": // Rogue
		g.TempClass = character.Rogue
		g.CreationStage = 1
		g.Renderer.AddGameMessage("Selected Rogue: Quick and stealthy adventurers.")
		g.Renderer.AddSystemMessage("Enter your character's name:")
	case "4": // Ranger
		g.TempClass = character.Ranger
		g.CreationStage = 1
		g.Renderer.AddGameMessage("Selected Ranger: Skilled hunters with ranged attacks.")
		g.Renderer.AddSystemMessage("Enter your character's name:")
	case "5": // Cleric
		g.TempClass = character.Cleric
		g.CreationStage = 1
		g.Renderer.AddGameMessage("Selected Cleric: Divine healers with protective abilities.")
		g.Renderer.AddSystemMessage("Enter your character's name:")
	case "escape": // Quit
		g.Running = false
	default:
		// If we're at the class selection stage, remind the player
		if g.CreationStage == 0 {
			g.Renderer.AddSystemMessage("Choose your character class:")
		}
	}
}

// handlePlayingInput handles input during gameplay
func (g *Game) handlePlayingInput() {
	// Get input
	key := input.GetSingleKey()

	// Debug output
	fmt.Printf("Playing input: '%s'\n", key)

	// Handle input
	switch key {
	case "up", "w", "W":
		g.movePlayer(0, -1)
	case "down", "s", "S":
		g.movePlayer(0, 1)
	case "left", "h", "H":
		g.movePlayer(-1, 0)
	case "right", "l", "L":
		g.movePlayer(1, 0)
	case "e", "E":
		g.useStairs()
	case "i", "I":
		g.openInventory()
	case "a", "A":
		g.useSpecialAbility()
	case "q", "Q", "escape":
		g.Running = false
	}
}

// openInventory opens the inventory screen
func (g *Game) openInventory() {
	g.State = StateInventory
	g.Renderer.AddSystemMessage("Opened inventory")
}

// handleInventoryInput handles input in the inventory screen
func (g *Game) handleInventoryInput() {
	// Get input
	key := input.GetSingleKey()

	// Debug output
	fmt.Printf("Inventory input: '%s'\n", key)

	// Handle input
	switch key {
	case "i", "I", "escape":
		g.closeInventory()
	}
}

// closeInventory closes the inventory screen
func (g *Game) closeInventory() {
	g.State = StatePlaying
	g.Renderer.AddSystemMessage("Closed inventory")
}

// handleGameOverInput handles input on the game over screen
func (g *Game) handleGameOverInput() {
	key := input.GetSingleKey()

	switch key {
	case "r":
		// Reset game
		*g = *NewGame()
	case "q":
		g.State = StateQuit
	}
}

// handleMainMenuInput handles input on the main menu
func (g *Game) handleMainMenuInput() {
	key := input.GetSingleKey()

	switch key {
	case "n":
		g.State = StatePlaying
	case "q":
		g.State = StateQuit
	}
}

// isValidMove checks if a move to the given coordinates is valid
func (g *Game) isValidMove(x, y int) bool {
	// Get current floor
	floor := g.Dungeon.Floors[g.CurrentFloor]

	// Check if coordinates are within the map
	if y < 0 || y >= floor.Level.Height || x < 0 || x >= floor.Level.Width {
		return false
	}

	// Check if the tile is walkable
	tile := floor.Level.Tiles[y][x]
	return tile.Walkable
}

// checkSpecialTiles checks for special tiles at the player's position and shows appropriate messages
func (g *Game) checkSpecialTiles() {
	// Get current floor
	floor := g.Dungeon.Floors[g.CurrentFloor]

	// Get tile at player position
	tile := floor.Level.Tiles[g.Player.Y][g.Player.X]

	// Debug output
	fmt.Printf("Player at position (%d, %d) on floor %d, tile type: %d\n",
		g.Player.X, g.Player.Y, g.CurrentFloor, tile.Type)

	// Print debug information about tiles around the player
	g.debugPrintTileTypes()

	// Check for exit
	if tile.Type == mapgen.TileExit {
		fmt.Println("Player is on an exit tile")

		// Add a message prompting the player to press 'e' to use the exit
		if g.CurrentFloor < len(g.Dungeon.Floors)-1 {
			g.Renderer.AddMessage("Press 'e' to descend deeper into the dungeon...")
		} else {
			g.Renderer.AddMessage("Press 'e' to complete your adventure!")
		}
		return // Return early to avoid showing multiple messages
	}

	// Check for entrance (to go up a level)
	if tile.Type == mapgen.TileEntrance && g.CurrentFloor > 0 {
		fmt.Println("Player is on an entrance tile")

		// Add a message prompting the player to press 'e' to use the entrance
		g.Renderer.AddMessage("Press 'e' to climb back up to the previous floor...")
	}
}

// handleLevelTransition handles the player's transition between dungeon levels
func (g *Game) handleLevelTransition() {
	// Get current floor
	floor := g.Dungeon.Floors[g.CurrentFloor]

	// Get tile at player position
	tile := floor.Level.Tiles[g.Player.Y][g.Player.X]

	fmt.Printf("Attempting level transition. Player at (%d, %d) on floor %d, tile type: %d\n",
		g.Player.X, g.Player.Y, g.CurrentFloor, tile.Type)

	// Print debug information about tiles around the player
	g.debugPrintTileTypes()

	// Handle exit (going down)
	if tile.Type == mapgen.TileExit {
		fmt.Println("Confirmed player is on an exit tile, attempting to go down")
		if g.CurrentFloor < len(g.Dungeon.Floors)-1 {
			g.CurrentFloor++
			g.Renderer.CurrentFloor = g.CurrentFloor
			fmt.Printf("Moving to floor %d\n", g.CurrentFloor)

			// Place player near the entrance of the next floor
			entrance := g.Dungeon.Floors[g.CurrentFloor].Entrance
			fmt.Printf("Entrance on new floor is at (%d, %d)\n", entrance.X, entrance.Y)
			g.placePlayerNearPosition(entrance)
			fmt.Printf("Placed player at (%d, %d)\n", g.Player.X, g.Player.Y)

			// Print debug information about the new floor
			g.debugPrintTileTypes()

			g.Renderer.AddMessage("You descend deeper into the dungeon...")
		} else {
			// Player has reached the end of the dungeon
			g.Renderer.AddMessage("Congratulations! You have reached the end of the dungeon!")
			g.State = StateGameOver
		}
		return // Add return to prevent checking entrance condition
	}

	// Handle entrance (going up)
	if tile.Type == mapgen.TileEntrance && g.CurrentFloor > 0 {
		fmt.Println("Confirmed player is on an entrance tile, attempting to go up")
		g.CurrentFloor--
		g.Renderer.CurrentFloor = g.CurrentFloor
		fmt.Printf("Moving to floor %d\n", g.CurrentFloor)

		// Place player near the exit of the previous floor
		exit := g.Dungeon.Floors[g.CurrentFloor].Exit
		fmt.Printf("Exit on previous floor is at (%d, %d)\n", exit.X, exit.Y)
		g.placePlayerNearPosition(exit)
		fmt.Printf("Placed player at (%d, %d)\n", g.Player.X, g.Player.Y)

		// Print debug information about the new floor
		g.debugPrintTileTypes()

		g.Renderer.AddMessage("You climb back up to the previous floor...")
		return
	}

	// If we get here, the player is not on a special tile
	fmt.Println("Player is not on a level transition tile")
	g.Renderer.AddMessage("You need to find stairs to change floors.")
}

// placePlayerNearPosition places the player in a valid position near the given position
func (g *Game) placePlayerNearPosition(pos mapgen.Position) {
	// Define possible offsets, prioritizing cardinal directions first
	cardinalOffsets := []struct{ dx, dy int }{
		{0, 1}, {1, 0}, {0, -1}, {-1, 0}, // Cardinal directions
	}

	diagonalOffsets := []struct{ dx, dy int }{
		{1, 1}, {1, -1}, {-1, 1}, {-1, -1}, // Diagonals
	}

	fmt.Printf("Trying to place player near position (%d, %d)\n", pos.X, pos.Y)

	// First try cardinal directions (more intuitive for player movement)
	for _, offset := range cardinalOffsets {
		newX := pos.X + offset.dx
		newY := pos.Y + offset.dy

		fmt.Printf("Checking cardinal position (%d, %d): ", newX, newY)
		if g.isValidMove(newX, newY) {
			fmt.Println("valid")
			g.Player.X = newX
			g.Player.Y = newY
			return
		} else {
			fmt.Println("invalid")
		}
	}

	// If no valid cardinal position, try diagonals
	for _, offset := range diagonalOffsets {
		newX := pos.X + offset.dx
		newY := pos.Y + offset.dy

		fmt.Printf("Checking diagonal position (%d, %d): ", newX, newY)
		if g.isValidMove(newX, newY) {
			fmt.Println("valid")
			g.Player.X = newX
			g.Player.Y = newY
			return
		} else {
			fmt.Println("invalid")
		}
	}

	// If no valid position found, try positions 2 spaces away
	for dy := -2; dy <= 2; dy++ {
		for dx := -2; dx <= 2; dx++ {
			// Skip already checked positions and the center
			if (dx == 0 && dy == 0) || (abs(dx) == 1 && abs(dy) == 1) || (abs(dx) == 1 && dy == 0) || (dx == 0 && abs(dy) == 1) {
				continue
			}

			newX := pos.X + dx
			newY := pos.Y + dy

			fmt.Printf("Checking extended position (%d, %d): ", newX, newY)
			if g.isValidMove(newX, newY) {
				fmt.Println("valid")
				g.Player.X = newX
				g.Player.Y = newY
				return
			} else {
				fmt.Println("invalid")
			}
		}
	}

	// If no valid position found, use the position itself as a fallback
	fmt.Println("No valid position found, using the position itself")
	g.Player.X = pos.X
	g.Player.Y = pos.Y
}

// abs returns the absolute value of an integer
func abs(x int) int {
	if x < 0 {
		return -x
	}
	return x
}

// debugPrintTileTypes prints the tile types around the player's position
func (g *Game) debugPrintTileTypes() {
	// Get current floor
	floor := g.Dungeon.Floors[g.CurrentFloor]

	fmt.Printf("Tile types around player on floor %d:\n", g.CurrentFloor)

	// Print a 5x5 grid around the player
	for y := g.Player.Y - 2; y <= g.Player.Y+2; y++ {
		for x := g.Player.X - 2; x <= g.Player.X+2; x++ {
			if y >= 0 && y < floor.Level.Height && x >= 0 && x < floor.Level.Width {
				tile := floor.Level.Tiles[y][x]
				if x == g.Player.X && y == g.Player.Y {
					fmt.Printf("[%d]", tile.Type) // Player position
				} else {
					fmt.Printf(" %d ", tile.Type)
				}
			} else {
				fmt.Printf(" X ") // Out of bounds
			}
		}
		fmt.Println()
	}

	// Print legend
	fmt.Println("Legend: 0=Wall, 1=Floor, 2=Entrance, 3=Exit, 4=Hallway, 5=Pillar, 6=Water, 7=Rubble")

	// Print entrance and exit positions
	fmt.Printf("Entrance: (%d, %d), Exit: (%d, %d)\n",
		floor.Entrance.X, floor.Entrance.Y, floor.Exit.X, floor.Exit.Y)
}

// useStairs uses stairs at the player's position
func (g *Game) useStairs() {
	// Get current floor
	currentFloor := g.Dungeon.Floors[g.CurrentFloor]

	// Get tile at player position
	tile := currentFloor.Level.Tiles[g.Player.Y][g.Player.X]

	// Check if player is on stairs
	if tile.Type == mapgen.TileExit {
		// Check if this is the last floor
		if g.CurrentFloor == len(g.Dungeon.Floors)-1 {
			g.Renderer.AddGameMessage("Press 'e' to complete your adventure!")
		} else {
			g.Renderer.AddGameMessage("Press 'e' to descend deeper into the dungeon...")
		}
		g.descendStairs()
	} else if tile.Type == mapgen.TileEntrance && g.CurrentFloor > 0 {
		g.Renderer.AddGameMessage("Press 'e' to climb back up to the previous floor...")
		g.ascendStairs()
	} else {
		g.Renderer.AddSystemMessage("You need to find stairs to change floors.")
	}
}

// descendStairs moves the player down to the next floor
func (g *Game) descendStairs() {
	// Check if player is on exit stairs
	currentFloor := g.Dungeon.Floors[g.CurrentFloor]
	playerTile := currentFloor.Level.Tiles[g.Player.Y][g.Player.X]

	if playerTile.Type != mapgen.TileExit {
		return
	}

	// Check if this is the last floor
	if g.CurrentFloor == len(g.Dungeon.Floors)-1 {
		// Player has completed the dungeon!
		g.Renderer.AddGameMessage("Congratulations! You have reached the end of the dungeon!")
		g.State = StateGameOver
		return
	}

	// Move to next floor
	g.CurrentFloor++
	g.Renderer.CurrentFloor = g.CurrentFloor

	// Place player at entrance of new floor
	entrance := g.Dungeon.Floors[g.CurrentFloor].Entrance
	g.Player.X = entrance.X
	g.Player.Y = entrance.Y

	// Add message
	g.Renderer.AddGameMessage("You descend deeper into the dungeon...")
}

// ascendStairs moves the player up to the previous floor
func (g *Game) ascendStairs() {
	// Check if player is on entrance stairs
	currentFloor := g.Dungeon.Floors[g.CurrentFloor]
	playerTile := currentFloor.Level.Tiles[g.Player.Y][g.Player.X]

	if playerTile.Type != mapgen.TileEntrance || g.CurrentFloor == 0 {
		return
	}

	// Move to previous floor
	g.CurrentFloor--
	g.Renderer.CurrentFloor = g.CurrentFloor

	// Place player at exit of previous floor
	exit := g.Dungeon.Floors[g.CurrentFloor].Exit
	g.Player.X = exit.X
	g.Player.Y = exit.Y

	// Add message
	g.Renderer.AddGameMessage("You climb back up to the previous floor...")
}

// movePlayer attempts to move the player in the specified direction
func (g *Game) movePlayer(dx, dy int) {
	// Calculate new position
	newX := g.Player.X + dx
	newY := g.Player.Y + dy

	// Get current floor
	currentFloor := g.Dungeon.Floors[g.CurrentFloor]

	// Check if the new position is within bounds
	if newX < 0 || newX >= currentFloor.Level.Width ||
		newY < 0 || newY >= currentFloor.Level.Height {
		return
	}

	// Check if the new position is walkable
	tile := currentFloor.Level.Tiles[newY][newX]
	if tile.Type == mapgen.TileWall || tile.Type == mapgen.TilePillar {
		return
	}

	// Move player
	g.Player.X = newX
	g.Player.Y = newY

	// Regenerate mana each turn
	g.Player.RegenerateMana()

	// Debug output
	fmt.Printf("Player moved to (%d, %d)\n", g.Player.X, g.Player.Y)

	// Check for special tiles
	if tile.Type == mapgen.TileEntrance {
		if g.CurrentFloor > 0 {
			g.Renderer.AddGameMessage("You found stairs leading up.")
		} else {
			g.Renderer.AddGameMessage("This is the entrance to the dungeon.")
		}
	} else if tile.Type == mapgen.TileExit {
		if g.CurrentFloor < len(g.Dungeon.Floors)-1 {
			g.Renderer.AddGameMessage("You found stairs leading down.")
		} else {
			g.Renderer.AddGameMessage("You found the exit of the dungeon!")
		}
	}
}

// useSpecialAbility uses the player's special ability
func (g *Game) useSpecialAbility() {
	// Check if player exists
	if g.Player == nil {
		return
	}

	// Use the special ability
	message, success := g.Player.UseSpecialAbility()

	// Add appropriate message
	if success {
		g.Renderer.AddGameMessage(message)
	} else {
		g.Renderer.AddSystemMessage(message)
	}

	// Regenerate mana each turn (this would normally be in a turn processing function)
	g.Player.RegenerateMana()
}

// MovePlayer moves the player by the given delta if possible
func (g *Game) MovePlayer(dx, dy int) bool {
	if g.State != StatePlaying || g.Player == nil {
		return false
	}

	newX := g.Player.X + dx
	newY := g.Player.Y + dy

	if g.isValidMove(newX, newY) {
		g.Player.X = newX
		g.Player.Y = newY
		g.checkSpecialTiles()
		return true
	}

	return false
}
