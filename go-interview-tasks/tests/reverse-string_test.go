package main

import (
	"interviews/go-interview-tasks/strings"
	"reflect"
	"testing"
)

func TestReverseString(t *testing.T) {
	tests := []struct {
		name     string
		input    []byte
		expected []byte
	}{
		{
			name:     "Empty string",
			input:    []byte{},
			expected: []byte{},
		},
		{
			name:     "Single character",
			input:    []byte{'a'},
			expected: []byte{'a'},
		},
		{
			name:     "Simple string",
			input:    []byte{'h', 'e', 'l', 'l', 'o'},
			expected: []byte{'o', 'l', 'l', 'e', 'h'},
		},
		{
			name:     "String with spaces",
			input:    []byte{'h', 'e', 'l', 'l', 'o', ' ', 'w', 'o', 'r', 'l', 'd'},
			expected: []byte{'d', 'l', 'r', 'o', 'w', ' ', 'o', 'l', 'l', 'e', 'h'},
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			strings.ReverseString(test.input)
			if !reflect.DeepEqual(test.input, test.expected) {
				t.Errorf("Expected %v, got %v", test.expected, test.input)
			}
		})
	}
}
