package main

import (
	"interviews/go-interview-tasks/strings"
	"testing"
)

func TestIsPalindrome(t *testing.T) {
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
			name:     "Single character",
			input:    "a",
			expected: true,
		},
		{
			name:     "Simple palindrome",
			input:    "racecar",
			expected: true,
		},
		{
			name:     "Palindrome with spaces and punctuation",
			input:    "A man, a plan, a canal: Panama",
			expected: true,
		},
		{
			name:     "Not a palindrome",
			input:    "race a car",
			expected: false,
		},
		{
			name:     "Palindrome with mixed case",
			input:    "RaCeCaR",
			expected: true,
		},
		{
			name:     "Palindrome with numbers",
			input:    "A1b2b1a",
			expected: true,
		},
		{
			name:     "Not palindrome with numbers",
			input:    "A1b2c3",
			expected: false,
		},
		{
			name:     "Only spaces and punctuation",
			input:    "  ,.!?  ",
			expected: true,
		},
		{
			name:     "Single alphanumeric with spaces",
			input:    "  a  ",
			expected: true,
		},
		{
			name:     "Long palindrome",
			input:    "Never odd or even",
			expected: true,
		},
	}

	passed := 0
	total := len(tests)

	for _, test := range tests {
		result := strings.IsPalindrome(test.input)
		if result == test.expected {
			passed++
			println("✅ PASS:", test.name)
		} else {
			println("❌ FAIL:", test.name)
			println("   Input:", test.input)
			println("   Expected:", test.expected)
			println("   Got:", result)
		}
	}

	println("\n📊 Results:", passed, "/", total, "tests passed")
	if passed == total {
		println("🎉 All tests passed!")
	} else {
		println("⚠️  Some tests failed!")
	}
}
