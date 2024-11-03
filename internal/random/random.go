package random

import (
	"crypto/rand"
	"errors"
	"math/big"
)

// RandomInt generates a random integer between min (inclusive) and max (exclusive).
// It panics if max is less than min or if there is an error generating the random number.
//
// Parameters:
// - min: The minimum value (inclusive).
// - max: The maximum value (exclusive).
//
// Returns:
// - A random integer between min and max.
func Int(min, max int) int {
	// Calculate the difference between max and min
	diff := max - min

	// Ensure that max is greater than min
	if diff <= 0 {
		panic(errors.New("max must be greater than min"))
	}

	// Generate a random integer in the range [0, diff)
	nBig, err := rand.Int(rand.Reader, big.NewInt(int64(diff)))
	if err != nil {
		panic(err)
	}

	// Add the minimum value to the generated random number
	n := nBig.Int64() + int64(min)

	// Return the result as an integer
	return int(n)
}
