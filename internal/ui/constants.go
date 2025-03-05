package ui

// ANSI color codes
const (
	ColorReset  = "\033[0m"
	ColorRed    = "\033[31m"
	ColorGreen  = "\033[32m"
	ColorYellow = "\033[33m"
	ColorBlue   = "\033[34m"
	ColorPurple = "\033[35m"
	ColorCyan   = "\033[36m"
	ColorWhite  = "\033[37m"
	ColorGray   = "\033[90m"

	// Background colors
	BgBlack = "\033[40m"
	BgRed   = "\033[41m"
	BgGreen = "\033[42m"
	BgBlue  = "\033[44m"

	// Message log dimensions
	MessageLogHeight = 5
)

// Window dimensions - now variables that can be updated
var (
	// Default window dimensions
	WindowWidth  = 80
	WindowHeight = 40

	// Map view dimensions (left panel)
	MapViewWidth  = 50
	MapViewHeight = 30

	// Character panel dimensions (right panel)
	CharPanelWidth  = 28 // 80 - 50 - 2 (separator)
	CharPanelHeight = 30
)
