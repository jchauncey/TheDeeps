package main

import (
	"fmt"
	"os"
	"os/signal"
	"syscall"

	"github.com/jchauncey/TheDeeps/internal/game"
	"github.com/jchauncey/TheDeeps/internal/ui"
)

func main() {
	fmt.Println("Welcome to The Deeps!")
	fmt.Println("Loading game...")

	// Set up signal handling for clean exit
	setupSignalHandling()

	// Initialize terminal size and set up resize handler
	if err := ui.InitTerminalSize(); err != nil {
		fmt.Fprintf(os.Stderr, "Error initializing terminal: %v\n", err)
		os.Exit(1)
	}

	// Hide cursor during gameplay
	ui.HideCursor()

	// Create and run the game
	g := game.NewGame()
	g.Run()

	// Restore terminal state
	ui.ShowCursor()
	fmt.Print("\033[0m") // Reset colors

	fmt.Println("Thanks for playing The Deeps!")
}

// setupSignalHandling ensures the terminal is restored on unexpected exit
func setupSignalHandling() {
	c := make(chan os.Signal, 1)
	signal.Notify(c, os.Interrupt, syscall.SIGTERM)

	go func() {
		<-c
		// Restore terminal
		ui.ShowCursor()
		fmt.Print("\033[0m") // Reset colors
		fmt.Println("\nGame terminated.")
		os.Exit(0)
	}()
}
