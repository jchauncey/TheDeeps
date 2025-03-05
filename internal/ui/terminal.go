package ui

import (
	"fmt"
	"os"
	"os/signal"
	"syscall"

	"golang.org/x/term"
)

// TerminalSize represents the width and height of the terminal
type TerminalSize struct {
	Width  int
	Height int
}

// Current terminal size
var CurrentSize TerminalSize

// GetTerminalSize returns the current terminal width and height
func GetTerminalSize() (TerminalSize, error) {
	width, height, err := term.GetSize(int(os.Stdout.Fd()))
	if err != nil {
		return TerminalSize{}, fmt.Errorf("failed to get terminal size: %w", err)
	}
	return TerminalSize{Width: width, Height: height}, nil
}

// InitTerminalSize initializes the terminal size and sets up a resize handler
func InitTerminalSize() error {
	// Get initial terminal size
	size, err := GetTerminalSize()
	if err != nil {
		return err
	}
	CurrentSize = size

	// Update the UI dimensions based on terminal size
	UpdateUIDimensions(size)

	// Set up resize handler
	go handleResize()

	return nil
}

// handleResize listens for terminal resize signals and updates the UI dimensions
func handleResize() {
	// Create a channel to receive signals
	ch := make(chan os.Signal, 1)
	signal.Notify(ch, syscall.SIGWINCH)

	for {
		// Wait for a SIGWINCH signal
		<-ch

		// Get the new terminal size
		size, err := GetTerminalSize()
		if err != nil {
			fmt.Fprintf(os.Stderr, "Error getting terminal size: %v\n", err)
			continue
		}

		// Update the current size
		CurrentSize = size

		// Update the UI dimensions
		UpdateUIDimensions(size)

		// Debug output
		fmt.Printf("Terminal resized to %dx%d\n", size.Width, size.Height)
	}
}

// UpdateUIDimensions updates the UI dimensions based on the terminal size
func UpdateUIDimensions(size TerminalSize) {
	// Set minimum sizes to ensure the UI doesn't break
	width := max(size.Width, 80)
	height := max(size.Height, 24)

	// Update window dimensions
	WindowWidth = width
	WindowHeight = height

	// Calculate message window height (proportional to screen size)
	MessageWindowHeight = max(10, height/4) // At least 10 lines or 1/4 of screen height

	// Available height for game area (excluding message window and borders)
	gameAreaHeight := height - MessageWindowHeight - 3 // 3 for borders and controls

	// Update map view dimensions (left panel)
	MapViewWidth = int(float64(width) * 0.65) // 65% of width
	MapViewHeight = gameAreaHeight            // Use all available height for game area

	// Update character panel dimensions (right panel)
	CharPanelWidth = width - MapViewWidth - 2 // 2 for separator
	CharPanelHeight = MapViewHeight

	// Update message log width
	MessageLogWidth = width

	// Debug output
	fmt.Printf("UI dimensions updated: Window(%d,%d), Map(%d,%d), Char(%d,%d), Msg(%d,%d)\n",
		WindowWidth, WindowHeight, MapViewWidth, MapViewHeight,
		CharPanelWidth, CharPanelHeight, MessageLogWidth, MessageWindowHeight)
}

// max returns the larger of two integers
func max(a, b int) int {
	if a > b {
		return a
	}
	return b
}

// min returns the smaller of two integers
func min(a, b int) int {
	if a < b {
		return a
	}
	return b
}

// HideCursor hides the terminal cursor
func HideCursor() {
	fmt.Print("\033[?25l")
}

// ShowCursor shows the terminal cursor
func ShowCursor() {
	fmt.Print("\033[?25h")
}
