package jwt_test

import (
	"testing"

	"github.com/edsonjaramillo/crpytid/internal/commands/password"
)

func TestGenerateSecret(t *testing.T) {
	testCases := []struct {
		length       int
		expectError  bool
		expectSecret bool
	}{
		{length: 0, expectError: true, expectSecret: false},
		{length: 16, expectError: false, expectSecret: true},
		{length: 32, expectError: false, expectSecret: true},
		{length: 64, expectError: false, expectSecret: true},
		{length: -1, expectError: true, expectSecret: false},
	}

	for _, tc := range testCases {
		secret := password.Generate(tc.length, true, true)

		if tc.expectError && secret != "" {
			t.Errorf("Expected error, got nil")
		}

		if !tc.expectError && secret == "" {
			t.Errorf("Expected secret, got nil")
		}

		if !tc.expectSecret && len(secret) != 0 {
			t.Errorf("Expected secret length to be 0, got %d", len(secret))
		}

	}
}
