package ui

import (
	"fmt"
	"strings"

	"github.com/jchauncey/TheDeeps/internal/character"
	mapgen "github.com/jchauncey/TheDeeps/internal/map"
)

// Renderer handles all UI rendering
type Renderer struct {
	// Reference to game state
	Player        *character.Player
	Dungeon       *mapgen.Dungeon
	CurrentFloor  int
	ExploredTiles map[string]bool
	MessageLog    *MessageLog

	// Character creation data
	CreationName string
}

// NewRenderer creates a new UI renderer
func NewRenderer(player *character.Player, dungeon *mapgen.Dungeon) *Renderer {
	return &Renderer{
		Player:        player,
		Dungeon:       dungeon,
		CurrentFloor:  0,
		ExploredTiles: make(map[string]bool),
		MessageLog:    NewMessageLog(100),
	}
}

// GetTileKey returns a unique key for a tile position
func GetTileKey(x, y, floor int) string {
	return fmt.Sprintf("%d:%d:%d", x, y, floor)
}

// UpdateExploredTiles marks tiles around the player as explored
func (r *Renderer) UpdateExploredTiles() {
	// Mark tiles as explored in a radius around the player
	viewRadius := 5
	for y := r.Player.Y - viewRadius; y <= r.Player.Y+viewRadius; y++ {
		for x := r.Player.X - viewRadius; x <= r.Player.X+viewRadius; x++ {
			// Check if coordinates are within the map
			if y >= 0 && y < r.Dungeon.Floors[r.CurrentFloor].Level.Height &&
				x >= 0 && x < r.Dungeon.Floors[r.CurrentFloor].Level.Width {
				r.ExploredTiles[GetTileKey(x, y, r.CurrentFloor)] = true
			}
		}
	}
}

// AddMessage adds a message to the message log (backward compatibility)
func (r *Renderer) AddMessage(format string, args ...interface{}) {
	r.MessageLog.AddMessage(format, args...)
}

// AddGameMessage adds a game message to the message log
func (r *Renderer) AddGameMessage(format string, args ...interface{}) {
	r.MessageLog.AddGameMessage(format, args...)
}

// AddSystemMessage adds a system message to the message log
func (r *Renderer) AddSystemMessage(format string, args ...interface{}) {
	r.MessageLog.AddSystemMessage(format, args...)
}

// RenderGame renders the current game state
func (r *Renderer) RenderGame() {
	// Clear the screen
	fmt.Print("\033[H\033[2J")

	// Check if terminal size has changed
	size, err := GetTerminalSize()
	if err == nil && (size.Width != CurrentSize.Width || size.Height != CurrentSize.Height) {
		// Update dimensions if terminal size has changed
		CurrentSize = size
		UpdateUIDimensions(size)
	}

	// Debug output
	fmt.Printf("RenderGame called - Player: %v, Dungeon: %v\n",
		r.Player != nil, r.Dungeon != nil)

	// Check if we have a player and dungeon
	if r.Player == nil || r.Dungeon == nil {
		// We're in character creation or main menu
		fmt.Println("Rendering character creation screen")
		r.renderCharacterCreation()
		return
	}

	fmt.Printf("Rendering game screen - Player at (%d, %d), Floor %d\n",
		r.Player.X, r.Player.Y, r.CurrentFloor)

	// Update explored tiles
	r.UpdateExploredTiles()

	// Get current floor
	currentFloor := r.Dungeon.Floors[r.CurrentFloor]

	// Calculate visible map region (centered on player)
	startX := r.Player.X - MapViewWidth/2
	endX := startX + MapViewWidth
	startY := r.Player.Y - MapViewHeight/2
	endY := startY + MapViewHeight

	// Ensure map view stays within map boundaries
	if startX < 0 {
		startX = 0
		endX = MapViewWidth
	}
	if endX > currentFloor.Level.Width {
		endX = currentFloor.Level.Width
		startX = endX - MapViewWidth
	}
	if startX < 0 {
		startX = 0
	}

	if startY < 0 {
		startY = 0
		endY = MapViewHeight
	}
	if endY > currentFloor.Level.Height {
		endY = currentFloor.Level.Height
		startY = endY - MapViewHeight
	}
	if startY < 0 {
		startY = 0
	}

	// Draw the game area (map and character panel)
	r.renderGameArea(currentFloor, startX, startY, endX, endY)

	// Draw the message window
	r.renderMessageWindow()

	// Draw controls
	fmt.Println(strings.Repeat("─", WindowWidth))
	fmt.Printf("%sControls:%s [↑↓←→] Move  [e] Use Stairs  [a] Special Ability  [i] Inventory  [q] Quit\n",
		ColorGreen, ColorReset)
}

// renderGameArea renders the map and character panel
func (r *Renderer) renderGameArea(currentFloor mapgen.Floor, startX, startY, endX, endY int) {
	// Draw the top border
	fmt.Println(strings.Repeat("─", WindowWidth))

	// Draw the map and character panel
	for y := 0; y < MapViewHeight; y++ {
		// Draw map row
		for x := 0; x < MapViewWidth; x++ {
			mapX := startX + x
			mapY := startY + y

			// Check if coordinates are within the map
			if mapY >= 0 && mapY < currentFloor.Level.Height &&
				mapX >= 0 && mapX < currentFloor.Level.Width {

				// Check if tile has been explored
				tileKey := GetTileKey(mapX, mapY, r.CurrentFloor)
				if !r.ExploredTiles[tileKey] {
					fmt.Print(" ") // Unexplored tile
					continue
				}

				// Draw player
				if mapX == r.Player.X && mapY == r.Player.Y {
					fmt.Print(ColorYellow + string(r.Player.Symbol) + ColorReset)
				} else {
					// Draw map tile
					tile := currentFloor.Level.Tiles[mapY][mapX]
					switch tile.Type {
					case mapgen.TileWall:
						fmt.Print(ColorGray + string(tile.Symbol) + ColorReset)
					case mapgen.TileFloor:
						fmt.Print(ColorWhite + string(tile.Symbol) + ColorReset)
					case mapgen.TileEntrance:
						fmt.Print(ColorGreen + string(tile.Symbol) + ColorReset)
					case mapgen.TileExit:
						fmt.Print(ColorRed + string(tile.Symbol) + ColorReset)
					case mapgen.TileHallway:
						fmt.Print(ColorWhite + string(tile.Symbol) + ColorReset)
					case mapgen.TilePillar:
						fmt.Print(ColorCyan + string(tile.Symbol) + ColorReset)
					case mapgen.TileWater:
						fmt.Print(ColorBlue + string(tile.Symbol) + ColorReset)
					case mapgen.TileRubble:
						fmt.Print(ColorPurple + string(tile.Symbol) + ColorReset)
					default:
						fmt.Print(" ") // Unknown tile
					}
				}
			} else {
				fmt.Print(" ") // Out of bounds
			}
		}

		// Draw separator
		fmt.Print("│")

		// Draw character panel (right panel)
		if y < CharPanelHeight {
			switch y {
			case 0:
				// Character name and class
				fmt.Printf(" %s%s%s the %s",
					ColorYellow, r.Player.Name, ColorReset, r.Player.Class.Name)
			case 1:
				// Divider
				fmt.Printf(" %s", strings.Repeat("─", CharPanelWidth-2))
			case 3:
				// Health label
				fmt.Printf(" HP: %s%d/%d%s",
					ColorRed, r.Player.HP, r.Player.MaxHP, ColorReset)
			case 4:
				// Health bar
				barWidth := CharPanelWidth - 4
				healthPercentage := float64(r.Player.HP) / float64(r.Player.MaxHP)
				filledWidth := int(healthPercentage * float64(barWidth))

				// Ensure at least one character for non-zero health
				if r.Player.HP > 0 && filledWidth == 0 {
					filledWidth = 1
				}

				// Choose color based on health percentage
				healthColor := ColorRed
				if healthPercentage > 0.7 {
					healthColor = ColorGreen
				} else if healthPercentage > 0.3 {
					healthColor = ColorYellow
				}

				// Draw the bar
				fmt.Print(" [")
				fmt.Print(healthColor)
				fmt.Print(strings.Repeat("█", filledWidth))
				fmt.Print(ColorReset)
				fmt.Print(strings.Repeat("░", barWidth-filledWidth))
				fmt.Print("]")
			case 5:
				// Mana label
				fmt.Printf(" MP: %s%d/%d%s",
					ColorBlue, r.Player.Mana, r.Player.MaxMana, ColorReset)
			case 6:
				// Mana bar
				barWidth := CharPanelWidth - 4
				manaPercentage := float64(r.Player.Mana) / float64(r.Player.MaxMana)
				filledWidth := int(manaPercentage * float64(barWidth))

				// Ensure at least one character for non-zero mana
				if r.Player.Mana > 0 && filledWidth == 0 {
					filledWidth = 1
				}

				// Draw the bar
				fmt.Print(" [")
				fmt.Print(ColorBlue)
				fmt.Print(strings.Repeat("█", filledWidth))
				fmt.Print(ColorReset)
				fmt.Print(strings.Repeat("░", barWidth-filledWidth))
				fmt.Print("]")
			case 8:
				// Divider
				fmt.Printf(" %s", strings.Repeat("─", CharPanelWidth-2))
			case 9:
				// Stats header
				fmt.Printf(" %sSTATS%s", ColorCyan, ColorReset)
			case 10:
				// Strength
				fmt.Printf(" STR: %s%d%s (%+d)",
					ColorWhite, r.Player.Strength, ColorReset, (r.Player.Strength-10)/2)
			case 11:
				// Wisdom
				fmt.Printf(" WIS: %s%d%s (%+d)",
					ColorWhite, r.Player.Wisdom, ColorReset, (r.Player.Wisdom-10)/2)
			case 12:
				// Constitution
				fmt.Printf(" CON: %s%d%s (%+d)",
					ColorWhite, r.Player.Constitution, ColorReset, (r.Player.Constitution-10)/2)
			case 13:
				// Dexterity
				fmt.Printf(" DEX: %s%d%s (%+d)",
					ColorWhite, r.Player.Dexterity, ColorReset, (r.Player.Dexterity-10)/2)
			case 14:
				// Charisma
				fmt.Printf(" CHA: %s%d%s (%+d)",
					ColorWhite, r.Player.Charisma, ColorReset, (r.Player.Charisma-10)/2)
			case 15:
				// Divider
				fmt.Printf(" %s", strings.Repeat("─", CharPanelWidth-2))
			case 16:
				// Equipment header
				fmt.Printf(" %sEQUIPMENT%s", ColorCyan, ColorReset)
			case 17:
				// Weapon
				fmt.Printf(" Weapon: %s (%+d)",
					r.Player.Weapon.Name, r.Player.GetDamageValue())
			case 18:
				// Armor
				fmt.Printf(" Armor: %s (%+d)",
					r.Player.Armor.Name, r.Player.GetArmorValue())
			case 19:
				// Divider
				fmt.Printf(" %s", strings.Repeat("─", CharPanelWidth-2))
			case 20:
				// Location header
				fmt.Printf(" %sLOCATION%s", ColorCyan, ColorReset)
			case 21:
				// Floor
				fmt.Printf(" Floor: %s%d%s",
					ColorYellow, r.CurrentFloor+1, ColorReset)
			case 22:
				// Position
				fmt.Printf(" Position: (%d, %d)",
					r.Player.X, r.Player.Y)
			case 23:
				// Divider
				fmt.Printf(" %s", strings.Repeat("─", CharPanelWidth-2))
			case 24:
				// Inventory header
				fmt.Printf(" %sINVENTORY%s", ColorCyan, ColorReset)
			case 25:
				// Gold
				fmt.Printf(" Gold: %s%d%s",
					ColorYellow, r.Player.Gold, ColorReset)
			case 26:
				// Items
				fmt.Printf(" Items: %d/%d",
					len(r.Player.Bag.Items), r.Player.Bag.Capacity)
			case 27:
				// Divider
				fmt.Printf(" %s", strings.Repeat("─", CharPanelWidth-2))
			case 28:
				// Special ability header
				fmt.Printf(" %sSPECIAL ABILITY%s", ColorCyan, ColorReset)
			case 29:
				// Ability name
				fmt.Printf(" %s", r.Player.Class.SpecialAbility)
			}
		}
		fmt.Println()
	}

	// Draw the bottom border of the game area
	fmt.Println(strings.Repeat("─", WindowWidth))
}

// renderMessageWindow renders the message window with system and game messages
func (r *Renderer) renderMessageWindow() {
	// Draw message window header
	fmt.Printf("%sMESSAGE LOG%s\n", ColorCyan, ColorReset)

	// Get recent messages
	messageCount := min(MessageWindowHeight-2, len(r.MessageLog.Messages))
	recentMessages := r.MessageLog.GetRecentMessages(messageCount)

	// Display messages with appropriate colors
	for i := 0; i < messageCount; i++ {
		msg := recentMessages[i]

		// Format timestamp
		timestamp := fmt.Sprintf("[%s]", msg.Timestamp)

		// Format prefix based on message type
		prefix := ""
		if msg.Type == MessageTypeSystem {
			prefix = ColorYellow + "[System]" + ColorReset
		} else {
			prefix = ColorGreen + "[Game]" + ColorReset
		}

		// Print the formatted message
		fmt.Printf("%s %s %s\n", timestamp, prefix, msg.Text)
	}

	// Fill remaining lines if needed
	for i := messageCount; i < MessageWindowHeight-2; i++ {
		fmt.Println()
	}
}

// renderCharacterCreation renders the character creation screen
func (r *Renderer) renderCharacterCreation() {
	// Get terminal size
	width := WindowWidth

	// Calculate padding for title
	title := "CHARACTER CREATION"
	padding := (width - len(title)) / 2
	if padding < 0 {
		padding = 0
	}

	// Draw title
	fmt.Print(ColorCyan)
	fmt.Print(strings.Repeat("=", width))
	fmt.Print("\n")
	fmt.Print(strings.Repeat(" ", padding))
	fmt.Print(title)
	fmt.Print("\n")
	fmt.Print(strings.Repeat("=", width))
	fmt.Print(ColorReset)
	fmt.Print("\n\n")

	// Check if we're in name entry stage
	inNameEntry := false
	if r.MessageLog != nil && len(r.MessageLog.Messages) > 0 {
		lastMsg := r.MessageLog.Messages[len(r.MessageLog.Messages)-1]
		inNameEntry = strings.Contains(lastMsg.Text, "Enter your character's name")
	}
	fmt.Printf("In name entry stage: %v\n", inNameEntry)

	// If in name entry, show a simplified screen
	if inNameEntry {
		fmt.Print(ColorCyan + "Enter your character's name: " + ColorReset)

		// Display the current name with a cursor
		nameField := ""
		if r.CreationName != "" {
			nameField = r.CreationName
		}

		// Add cursor and padding
		nameField += "_"
		fmt.Print(nameField)

		// Add some instructions
		fmt.Print("\n\n")
		fmt.Println("Type your character's name and press Enter when done.")
		fmt.Println("Press Escape to go back to class selection.")
		fmt.Print("\n")
	} else {
		// Draw class selection
		fmt.Print(ColorCyan + "Choose your character class:" + ColorReset + "\n\n")

		// Display class options
		fmt.Printf("  %s[1] Warrior%s: Strong fighters with high health and damage\n", ColorGreen, ColorReset)
		fmt.Printf("      Bonus: +3 Strength, +2 Constitution\n")
		fmt.Printf("      Special: Berserk - Deal extra damage when low on health\n\n")

		fmt.Printf("  %s[2] Wizard%s: Masters of arcane magic with powerful spells\n", ColorBlue, ColorReset)
		fmt.Printf("      Bonus: +3 Wisdom, +2 Intelligence\n")
		fmt.Printf("      Special: Fireball - Deal area damage to enemies\n\n")

		fmt.Printf("  %s[3] Rogue%s: Stealthy adventurers with high evasion\n", ColorPurple, ColorReset)
		fmt.Printf("      Bonus: +3 Dexterity, +2 Charisma\n")
		fmt.Printf("      Special: Backstab - Deal critical damage from stealth\n\n")

		fmt.Printf("  %s[4] Ranger%s: Skilled hunters with ranged attacks\n", ColorGreen, ColorReset)
		fmt.Printf("      Bonus: +3 Dexterity, +2 Wisdom\n")
		fmt.Printf("      Special: Eagle Eye - Increased vision range\n\n")

		fmt.Printf("  %s[5] Cleric%s: Divine healers with protective abilities\n", ColorYellow, ColorReset)
		fmt.Printf("      Bonus: +3 Wisdom, +2 Charisma\n")
		fmt.Printf("      Special: Heal - Restore health points\n\n")
	}

	// Draw message log
	fmt.Print(ColorYellow)
	fmt.Print(strings.Repeat("-", width))
	fmt.Print(ColorReset + "\n")

	if r.MessageLog != nil {
		// Show the last few messages that fit in the window
		messages := r.MessageLog.Messages
		maxMessages := min(len(messages), WindowHeight/4) // Use at most 1/4 of the window height

		for i := len(messages) - maxMessages; i < len(messages); i++ {
			if i >= 0 {
				msg := messages[i]
				// Format the message
				formattedMsg := r.MessageLog.FormatMessage(msg)

				// Truncate message if it's too long for the window
				if len(formattedMsg) > width {
					formattedMsg = formattedMsg[:width-3] + "..."
				}
				fmt.Println(formattedMsg)
			}
		}
	}
}
