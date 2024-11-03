package random_test

import (
	"testing"

	"github.com/edsonjaramillo/crpytid/internal/random"
)

// TestIntValidRange tests the Int function with valid ranges.
func TestIntValidRange(t *testing.T) {
	for i := 0; i < 100; i++ {
		result := random.Int(1, 10)
		if result < 1 || result >= 10 {
			t.Errorf("Expected result to be in range [1, 10), got %d", result)
		}
	}
}

// TestIntMinEqualsMax tests the Int function when min equals max.
func TestIntMinEqualsMax(t *testing.T) {
	defer func() {
		if r := recover(); r == nil {
			t.Errorf("Expected panic when min equals max, but did not panic")
		}
	}()
	random.Int(5, 5)
}

// TestIntMinGreaterThanMax tests the Int function when min is greater than max.
func TestIntMinGreaterThanMax(t *testing.T) {
	defer func() {
		if r := recover(); r == nil {
			t.Errorf("Expected panic when min is greater than max, but did not panic")
		}
	}()
	random.Int(10, 5)
}
