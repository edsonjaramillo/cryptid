package flags

import (
	"fmt"

	"github.com/atotto/clipboard"
	"github.com/urfave/cli/v2"
)

var NoClipboardFlag = &cli.BoolFlag{
	Name:  "no-clipboard",
	Usage: "Do not copy the password to the clipboard",
	Value: false,
}

// ClipboardPrinter copies the provided text to the clipboard if the status is false.
// If the clipboard operation fails, it prints an error message and panics.
//
// Parameters:
// - status: A boolean flag indicating whether to skip the clipboard operation.
// - text: The text to be copied to the clipboard.
func ClipboardPrinter(status bool, text string) {
	// Check if the clipboard operation should be performed
	if !status {
		// Attempt to write the text to the clipboard
		err := clipboard.WriteAll(text)
		if err != nil {
			// Print an error message if the clipboard operation fails
			println("Failed to copy text to clipboard")
			// Panic to indicate a critical failure
			panic(err)
		}
	}
}

var QuietFlag = &cli.BoolFlag{
	Name:    "quiet",
	Aliases: []string{"q"},
	Usage:   "Do not print the password to the console",
	Value:   false,
}

// QuietPrinter prints the provided text to the console if the status is false.
// This function is useful for suppressing console output when the user specifies a flag.
// If the status is true, the text is not printed to the console.
// This function is intended to be used in conjunction with the QuietFlag.
func QuietPrinter(status bool, text string) {
	if !status {
		fmt.Println(text)
	}
}
