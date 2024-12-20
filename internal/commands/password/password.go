package password

import (
	_ "embed"
	"fmt"
	"os"
	"strings"

	"github.com/edsonjaramillo/crpytid/internal/flags"
	"github.com/edsonjaramillo/crpytid/internal/random"
	"github.com/urfave/cli/v2"
)

const (
	wordsFilePath = "internal/commands/password/words.txt"
)

//go:embed words.txt
var wordsFileContent string

// Main Password command
var PasswordCommand = &cli.Command{
	Name:        "password",
	Usage:       "Generate passwords",
	Subcommands: []*cli.Command{complexSubcommand, passphraseSubcommand},
}

// Subcommands
var complexSubcommand = &cli.Command{
	Name:   "complex",
	Usage:  "Generate a complex password. Ex: n48h@3fj!2f",
	Flags:  []cli.Flag{complexLengthFlag, noNumbersFlag, noSymbolsFlag, flags.NoClipboardFlag, flags.QuietFlag},
	Action: complexAction,
}

var passphraseSubcommand = &cli.Command{
	Name:   "passphrase",
	Usage:  "Generate a passphrase. Ex: Apple-Banana9-Orange$",
	Flags:  []cli.Flag{numberOfWordsFlag, passphraseSeparatorFlag, flags.NoClipboardFlag, flags.QuietFlag},
	Action: passphraseAction,
}

// Actions

func complexAction(cCtx *cli.Context) error {
	length := cCtx.Int("length")
	noNumbers := cCtx.Bool("no-numbers")
	noSymbols := cCtx.Bool("no-symbols")
	noClipboard := cCtx.Bool("no-clipboard")
	quiet := cCtx.Bool("quiet")

	passwordGenerated := GenerateRandom(length, noNumbers, noSymbols)

	flags.QuietPrinter(quiet, passwordGenerated)
	flags.ClipboardPrinter(noClipboard, passwordGenerated)

	return nil
}

func passphraseAction(cCtx *cli.Context) error {
	numberOfWords := cCtx.Int("count")
	separator := cCtx.String("separator")
	noClipboard := cCtx.Bool("no-clipboard")
	quiet := cCtx.Bool("quiet")

	passphraseGenerated := generatePassphrase(numberOfWords, separator)

	flags.QuietPrinter(quiet, passphraseGenerated)
	flags.ClipboardPrinter(noClipboard, passphraseGenerated)

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

var numberOfWordsFlag = &cli.IntFlag{
	Name:     "count",
	Aliases:  []string{"c"},
	Usage:    "Number of words in the passphrase",
	Category: "Passphrase Options",
	Value:    4,
}

var passphraseSeparatorFlag = &cli.StringFlag{
	Name:     "separator",
	Aliases:  []string{"s"},
	Usage:    "Separator between words",
	Category: "Passphrase Options",
	Value:    "-",
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

func generatePassphrase(numberOfWordsNeeded int, separator string) string {
	wordWithNumberIdx := random.Int(0, numberOfWordsNeeded-1)

	words, numberOfWordsFromList := getWords()
	maxWordListIdx := numberOfWordsFromList - 1

	passphrase := ""

	for i := 0; i < numberOfWordsNeeded; i++ {
		wordIndex := random.Int(0, maxWordListIdx)
		word := words[wordIndex]

		if i == wordWithNumberIdx {
			randomNumber := random.Int(0, 100)
			word += fmt.Sprintf("%d", randomNumber)
		}

		passphrase += word

		// if it is not the last word, add the separator
		if i != numberOfWordsNeeded-1 {
			passphrase += separator
		}

	}

	// add a special character at the end
	randomSymbol := string(symbols[random.Int(0, len(symbols)-1)])
	passphrase += randomSymbol

	return passphrase
}

func getWords() ([]string, int) {
	// Use cached content if available
	if wordsFileContent != "" {
		words := strings.Split(wordsFileContent, "\n")
		return words, len(words)
	}

	// Read file content
	content, err := os.ReadFile(wordsFilePath)
	if err != nil {
		return nil, 0
	}

	// Split into lines and clean up
	words := make([]string, 0)
	for _, word := range strings.Split(string(content), "\n") {
		word = strings.TrimSpace(word)
		if word != "" {
			words = append(words, word)
		}
	}

	if len(words) == 0 {
		return nil, 0
	}

	return words, len(words)
}
