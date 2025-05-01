package commands_test

import (
	"testing"

	"github.com/edsonjaramillo/hyde/backend/internal/commands"
	"github.com/stretchr/testify/assert"
)

func TestAddEncExtension(t *testing.T) {
	tests := []struct {
		name     string
		input    string
		expected string
	}{
		{
			name:     "File with .enc extension in folder",
			input:    "./hello/example.txt",
			expected: "example.txt.enc",
		},
		{
			name:     "File with .enc extension in absolute folder",
			input:    "/hello/example.txt",
			expected: "example.txt.enc",
		},
		{
			name:     "File without .enc extension",
			input:    "example.txt",
			expected: "example.txt.enc",
		},
		{
			name:     "File with .enc extension",
			input:    "example.txt.enc",
			expected: "example.txt.enc",
		},
		{
			name:     "File with no extension",
			input:    "example",
			expected: "example.enc",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := commands.AddEncExtension(tt.input)
			assert.Equal(t, tt.expected, result)
		})
	}
}

func TestRemoveEncExtension(t *testing.T) {
	tests := []struct {
		name     string
		input    string
		expected string
	}{
		{
			name:     "File with .enc extension in relative folder",
			input:    "./hello/example.txt.enc",
			expected: "example.txt",
		},
		{
			name:     "File with .enc extension in absolute folder",
			input:    "/hello/example.txt.enc",
			expected: "example.txt",
		},
		{
			name:     "File with .enc extension",
			input:    "example.txt.enc",
			expected: "example.txt",
		},
		{
			name:     "File without .enc extension",
			input:    "example.txt",
			expected: "example.txt",
		},
		{
			name:     "File with no extension",
			input:    "example",
			expected: "example",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := commands.RemoveEncExtension(tt.input)
			assert.Equal(t, tt.expected, result)
		})
	}
}
