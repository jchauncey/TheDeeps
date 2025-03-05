package main

import (
	"fmt"

	"github.com/jchauncey/TheDeeps/internal/game"
)

func main() {
	fmt.Println("Welcome to The Deeps!")
	fmt.Println("Loading game...")

	// Create and run the game
	g := game.NewGame()
	g.Run()

	fmt.Println("Thanks for playing The Deeps!")
}
