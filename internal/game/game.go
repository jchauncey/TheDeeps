package game

import (
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

// NewGame creates a new game instance
func NewGame() *Game {
	// Create renderer without player or dungeon initially
	renderer := ui.NewRenderer(nil, nil)

	// Add welcome message
	renderer.AddMessage("Welcome to The Deeps! Create your character to begin.")

	return &Game{
		Player:        nil, // Will be created after character creation
		Dungeon:       nil, // Will be created after character creation
		CurrentFloor:  0,
		Renderer:      renderer,
		State:         StateCharacterCreation,
		Running:       true,
		TempClass:     character.Warrior, // Default selection
		CreationStage: 0,                 // Start with class selection
	}
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

	// Add game start message
	g.Renderer.AddMessage(fmt.Sprintf("Welcome, %s the %s! Your adventure begins!",
		g.TempName, character.GetClassName(g.TempClass)))

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
	if g.CreationStage == 0 {
		// Class selection stage
		key := input.GetSingleKey()
		fmt.Printf("Class selection key: '%s'\n", key)

		switch key {
		case "1", "w", "W":
			g.TempClass = character.Warrior
			g.Renderer.AddMessage("Selected Warrior: Strong and resilient fighters.")
			// Automatically advance to name entry
			g.CreationStage = 1
			g.Renderer.AddMessage("Enter your character's name:")
			g.TempName = ""              // Initialize empty name
			g.Renderer.CreationName = "" // Clear the creation name
			// Force redraw to show name entry screen
			g.Renderer.RenderGame()
		case "2", "r", "R":
			g.TempClass = character.Rogue
			g.Renderer.AddMessage("Selected Rogue: Quick and stealthy adventurers.")
			// Automatically advance to name entry
			g.CreationStage = 1
			g.Renderer.AddMessage("Enter your character's name:")
			g.TempName = ""              // Initialize empty name
			g.Renderer.CreationName = "" // Clear the creation name
			// Force redraw to show name entry screen
			g.Renderer.RenderGame()
		case "3", "m", "M":
			g.TempClass = character.Wizard
			g.Renderer.AddMessage("Selected Wizard: Powerful wielders of arcane magic.")
			// Automatically advance to name entry
			g.CreationStage = 1
			g.Renderer.AddMessage("Enter your character's name:")
			g.TempName = ""              // Initialize empty name
			g.Renderer.CreationName = "" // Clear the creation name
			// Force redraw to show name entry screen
			g.Renderer.RenderGame()
		case "q", "Q", "escape", "\x1b":
			g.State = StateQuit
		}
	} else if g.CreationStage == 1 {
		// Name entry stage - use GetSingleKey for better control
		key := input.GetSingleKey()

		// Debug output
		fmt.Printf("Name entry key: '%s'\n", key)

		switch key {
		case "enter":
			// Enter key pressed
			fmt.Println("Enter key detected")
			if len(g.TempName) > 0 {
				// Name entered, start the game
				fmt.Printf("Starting game with name: %s\n", g.TempName)
				g.StartGame()
				// Force a redraw after starting the game
				g.Renderer.RenderGame()
			} else {
				g.Renderer.AddMessage("Please enter a name for your character.")
			}
		case "backspace":
			// Remove last character
			if len(g.TempName) > 0 {
				g.TempName = g.TempName[:len(g.TempName)-1]
				g.Renderer.CreationName = g.TempName // Update the creation name
			}
		case "escape":
			// Go back to class selection
			g.CreationStage = 0
			g.Renderer.AddMessage("Choose your character class:")
		default:
			// Add character to name if it's a single printable character
			if len(key) == 1 && key[0] >= 32 && key[0] <= 126 {
				if len(g.TempName) < 20 { // Limit name length
					g.TempName += key
					g.Renderer.CreationName = g.TempName // Update the creation name
				}
			}
		}
	}
}

// handlePlayingInput handles input during normal gameplay
func (g *Game) handlePlayingInput() {
	key := input.GetSingleKey()

	// Movement
	newX, newY := g.Player.X, g.Player.Y

	switch key {
	case "up", "w":
		newY--
	case "down", "s":
		newY++
	case "left", "a":
		newX--
	case "right", "d":
		newX++
	case "i":
		g.State = StateInventory
		g.Renderer.AddMessage("Opened inventory")
		return
	case "q":
		g.State = StateQuit
		return
	}

	// Check if movement is valid
	if g.isValidMove(newX, newY) {
		g.Player.X, g.Player.Y = newX, newY

		// Check for special tiles
		g.checkSpecialTiles()
	}
}

// handleInventoryInput handles input in the inventory screen
func (g *Game) handleInventoryInput() {
	key := input.GetSingleKey()

	switch key {
	case "i", "escape":
		g.State = StatePlaying
		g.Renderer.AddMessage("Closed inventory")
	case "q":
		g.State = StateQuit
	}
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

// checkSpecialTiles checks for special tiles at the player's position
func (g *Game) checkSpecialTiles() {
	// Get current floor
	floor := g.Dungeon.Floors[g.CurrentFloor]

	// Get tile at player position
	tile := floor.Level.Tiles[g.Player.Y][g.Player.X]

	// Check for exit
	if tile.Type == mapgen.TileExit {
		if g.CurrentFloor < len(g.Dungeon.Floors)-1 {
			g.CurrentFloor++
			g.Renderer.CurrentFloor = g.CurrentFloor

			// Place player at entrance of next floor
			entrance := g.Dungeon.Floors[g.CurrentFloor].Entrance
			g.Player.X = entrance.X
			g.Player.Y = entrance.Y

			g.Renderer.AddMessage("You descend deeper into the dungeon...")
		} else {
			// Player has reached the end of the dungeon
			g.Renderer.AddMessage("Congratulations! You have reached the end of the dungeon!")
			g.State = StateGameOver
		}
	}
}
