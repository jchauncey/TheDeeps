package main

import (
	"github.com/jchauncey/TheDeeps/server/log"
)

func main() {
	// Set the log level to Debug to show all messages
	log.SetLevel(log.DebugLevel)

	// Ensure colors are enabled
	log.SetUseColors(true)

	// Log messages at different levels
	log.Debug("This is a DEBUG message - should be BLUE")
	log.Info("This is an INFO message - should be GREEN")
	log.Warn("This is a WARN message - should be YELLOW")
	log.Error("This is an ERROR message - should be RED")

	// Don't use Fatal in this demo as it would exit the program
	log.Info("Notice how each log level has its own distinct color")

	// Show an example with formatted output
	log.Info("You can also use %s in your log messages with %d or more arguments", "formatting", 2)

	// Show an example with caller information
	log.Debug("This message includes the file and line number where it was called")
}
