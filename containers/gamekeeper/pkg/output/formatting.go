package output

import (
	"fmt"
)

// ANSI color codes
const (
	Reset  = "\033[0m"
	Bold   = "\033[1m"
	Dim    = "\033[2m"
	Red    = "\033[0;31m"
	Green  = "\033[0;32m"
	Yellow = "\033[1;33m"
	Cyan   = "\033[0;36m"
)

// Section prints a formatted section header
func Section(title string) {
	fmt.Println()
	fmt.Printf("%s%s═══════════════════════════════════════════════════════════%s\n", Bold, Cyan, Reset)
	fmt.Printf("%s%s  %s%s\n", Bold, Cyan, title, Reset)
	fmt.Printf("%s%s═══════════════════════════════════════════════════════════%s\n", Bold, Cyan, Reset)
}

// Step prints a step message (without newline, expects Success/Error to follow)
func Step(message string) {
	fmt.Printf("  %s→%s %s: ", Cyan, Reset, message)
}

// Success prints a success indicator
func Success() {
	fmt.Printf("%s✓%s\n", Green, Reset)
}

// SuccessWithMessage prints a success indicator with a message
func SuccessWithMessage(message string) {
	fmt.Printf("%s%s%s\n", Green, message, Reset)
}

// Error prints an error message
func Error(message string) {
	fmt.Printf("%s✗ %s%s\n", Red, message, Reset)
}

// Info prints an info message
func Info(message string) {
	fmt.Printf("  %sℹ%s  %s\n", Cyan, Reset, message)
}

// Warning prints a warning message
func Warning(message string) {
	fmt.Printf("  %s⚠%s  %s\n", Yellow, Reset, message)
}

// Launch prints a launch message
func Launch(message string) {
	fmt.Println()
	fmt.Printf("%s%s🚀 %s%s\n", Bold, Green, message, Reset)
	fmt.Println()
}
