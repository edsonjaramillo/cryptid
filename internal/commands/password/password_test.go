package password_test

import (
	"testing"
	"unicode"

	"github.com/edsonjaramillo/crpytid/internal/commands/password"
)

func TestGenerate(t *testing.T) {
	testCases := []struct {
		name        string
		length      int
		withNumbers bool
		withSpecial bool
	}{
		{"OnlyAlphabet", 12, false, false},
		{"AlphabetNumbers", 12, true, false},
		{"AlphabetSpecials", 12, false, true},
		{"AllCharsets", 12, true, true},
		{"OnlyNumbers", 12, true, false},
		{"OnlySpecials", 12, false, true},
		{"NoCharsets", 12, false, false},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			password := password.Generate(tc.length, tc.withNumbers, tc.withSpecial)
			if len(password) != tc.length {
				t.Errorf("Expected password length %d, got %d", tc.length, len(password))
			}

			var hasNumber, hasSpecial bool
			for _, char := range password {
				switch {
				case unicode.IsDigit(char):
					hasNumber = true
				case isSpecialCharacter(char):
					hasSpecial = true
				}
			}

			if tc.withNumbers && !hasNumber {
				t.Errorf("Password should contain numeric characters")
			}
			if tc.withSpecial && !hasSpecial {
				t.Errorf("Password should contain special characters")
			}
			if !tc.withNumbers && hasNumber {
				t.Errorf("Password should not contain numeric characters")
			}
			if !tc.withSpecial && hasSpecial {
				t.Errorf("Password should not contain special characters")
			}
		})
	}
}

func isSpecialCharacter(c rune) bool {
	symbols := "!*-.@_"
	for _, s := range symbols {
		if c == s {
			return true
		}
	}
	return false
}
