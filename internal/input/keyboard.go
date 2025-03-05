package input

import (
	"bufio"
	"fmt"
	"os"
	"strings"

	"golang.org/x/term"
)

// GetSingleKey reads a single keypress from the terminal
func GetSingleKey() string {
	// Save the current state of the terminal
	oldState, err := term.MakeRaw(int(os.Stdin.Fd()))
	if err != nil {
		fmt.Println("Error setting terminal to raw mode:", err)
		return ""
	}
	defer term.Restore(int(os.Stdin.Fd()), oldState)

	// Read a single byte
	b := make([]byte, 3)
	n, err := os.Stdin.Read(b)
	if err != nil {
		fmt.Println("Error reading from stdin:", err)
		return ""
	}

	// Handle special keys
	if n == 3 && b[0] == 27 && b[1] == 91 {
		// Arrow keys
		switch b[2] {
		case 65:
			return "up"
		case 66:
			return "down"
		case 67:
			return "right"
		case 68:
			return "left"
		}
	} else if n == 1 {
		// Regular keys
		switch b[0] {
		case 27:
			return "escape"
		case 13:
			return "enter"
		case 127:
			return "backspace"
		default:
			return string(b[0])
		}
	}

	return ""
}

// GetKey reads a line of input from the terminal
func GetKey() string {
	reader := bufio.NewReader(os.Stdin)
	input, err := reader.ReadString('\n')

	// Debug output
	fmt.Printf("Raw input received: %q, error: %v\n", input, err)

	// Trim newline characters
	input = strings.TrimSuffix(input, "\n")
	input = strings.TrimSuffix(input, "\r")

	// Debug output
	fmt.Printf("Trimmed input: %q\n", input)

	return input
}
