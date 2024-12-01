package password

import (
	"github.com/edsonjaramillo/crpytid/internal/flags"
	"github.com/edsonjaramillo/crpytid/internal/random"
	"github.com/urfave/cli/v2"
)

// Main Password command
var PasswordCommand = &cli.Command{
	Name:        "password",
	Usage:       "Generate passwords",
	Subcommands: []*cli.Command{complexSubcommand},
}

// Subcommands
var complexSubcommand = &cli.Command{
	Name:   "complex",
	Usage:  "Generate a complex password. Ex: n48h@3fj!2f",
	Flags:  []cli.Flag{complexLengthFlag, noNumbersFlag, noSymbolsFlag, flags.NoClipboardFlag, flags.NoConsoleFlag},
	Action: complexAction,
}

// Actions

func complexAction(cCtx *cli.Context) error {
	length := cCtx.Int("length")
	noNumbers := cCtx.Bool("no-numbers")
	noSymbols := cCtx.Bool("no-symbols")
	noClipboard := cCtx.Bool("no-clipboard")
	noConsole := cCtx.Bool("no-console")

	passwordGenerated := GenerateRandom(length, noNumbers, noSymbols)

	flags.NoConsolePrinter(noConsole, passwordGenerated)
	flags.ClipboardPrinter(noClipboard, passwordGenerated)

	return nil
}

// Flags

var complexLengthFlag = &cli.IntFlag{
	Name:     "length",
	Aliases:  []string{"l"},
	Usage:    "Length of the password",
	Category: "Complex Password Options",
	Value:    14,
}

var noNumbersFlag = &cli.BoolFlag{
	Name:    "no-numbers",
	Aliases: []string{"no-digits"},
	Usage:   "Exclude numbers",
	Value:   false,
}

var noSymbolsFlag = &cli.BoolFlag{
	Name:    "no-symbols",
	Aliases: []string{"no-specials"},
	Usage:   "Exclude special characters",
	Value:   false,
}

// Functionality

const alphabet = "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ"
const numbers = "0123456789"
const symbols = "!#$%&()*+,-.:;<=>?@[]^_{|}~"

// Generate creates a random password of the specified length using the provided character sets.
// It ensures that at least one number and one special character are included if specified.
//
// Parameters:
// - length: The desired length of the password.
// - withAlphabet: Whether to include alphabetic characters (A-Z, a-z).
// - noNumbers: Whether to include numeric characters (0-9).
// - noSymbols: Whether to include special characters: !#$%&()*+,-.:;<=>?@[]^_{|}~
//
// Returns:
// - A randomly generated password as a string.
func GenerateRandom(length int, noNumbers bool, noSymbols bool) string {

	charset := alphabet

	if !noNumbers {
		charset += numbers
	}

	if !noSymbols {
		charset += symbols
	}

	password := make([]byte, 0, length)
	charsetLength := len(charset)

	for i := 0; i < length; i++ {
		randomIndex := random.Int(0, charsetLength-1)
		password = append(password, charset[randomIndex])
	}

	if !noNumbers {
		randomIndex := random.Int(0, len(numbers)-1)
		password[random.Int(0, length)] = numbers[randomIndex]
	}

	if !noSymbols {
		randomIndex := random.Int(0, len(symbols)-1)
		password[random.Int(0, length)] = symbols[randomIndex]
	}

	return string(password)
}
