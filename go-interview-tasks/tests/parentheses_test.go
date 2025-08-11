package main

import (
	"interviews/go-interview-tasks/strings"
	"testing"
)

func TestIsValidParentheses(t *testing.T) {
	tests := []struct {
		name     string
		input    string
		expected bool
	}{
		{
			name:     "Empty string",
			input:    "",
			expected: true,
		},
		{
			name:     "Simple parentheses",
			input:    "()",
			expected: true,
		},
		{
			name:     "Simple brackets",
			input:    "[]",
			expected: true,
		},
		{
			name:     "Simple braces",
			input:    "{}",
			expected: true,
		},
		{
			name:     "Mixed valid parentheses",
			input:    "()[]{}",
			expected: true,
		},
		{
			name:     "Nested valid parentheses",
			input:    "({[]})",
			expected: true,
		},
		{
			name:     "Complex valid case",
			input:    "{[()()]}",
			expected: true,
		},
		{
			name:     "Single opening parenthesis",
			input:    "(",
			expected: false,
		},
		{
			name:     "Single closing parenthesis",
			input:    ")",
			expected: false,
		},
		{
			name:     "Mismatched parentheses",
			input:    "(]",
			expected: false,
		},
		{
			name:     "Mismatched brackets",
			input:    "[}",
			expected: false,
		},
		{
			name:     "Wrong order",
			input:    "([)]",
			expected: false,
		},
		{
			name:     "Multiple opening without closing",
			input:    "(((",
			expected: false,
		},
		{
			name:     "Multiple closing without opening",
			input:    ")))",
			expected: false,
		},
		{
			name:     "Mixed invalid case",
			input:    "{[()()]",
			expected: false,
		},
		{
			name:     "Long valid string",
			input:    repeatString("()", 50),
			expected: true,
		},
		{
			name:     "Long invalid string",
			input:    repeatString("(", 100),
			expected: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := strings.IsValidParentheses(tt.input)
			if result != tt.expected {
				t.Errorf("isValidParentheses(%q) = %v, want %v", tt.input, result, tt.expected)
			}
		})
	}
}

// Benchmark test for performance testing
func BenchmarkIsValidParentheses(b *testing.B) {
	testString := "{[()()]}{[()()]}{[()()]}"
	for b.Loop() {
		strings.IsValidParentheses(testString)
	}
}

// Helper function to repeat a string n times
func repeatString(s string, n int) string {
	result := ""
	for range n {
		result += s
	}
	return result
}
