package main

import (
	"interviews/go-interview-tasks/strings"
	"testing"
)

func TestRemoveDuplicates(t *testing.T) {
	tests := []struct {
		name     string
		input    string
		expected string
	}{
		{
			name:     "Example 1: abbaca -> ca",
			input:    "abbaca",
			expected: "ca",
		},
		{
			name:     "Example 2: azxxzy -> ay",
			input:    "azxxzy",
			expected: "ay",
		},
		{
			name:     "Empty string",
			input:    "",
			expected: "",
		},
		{
			name:     "Single character",
			input:    "a",
			expected: "a",
		},
		{
			name:     "No duplicates",
			input:    "abc",
			expected: "abc",
		},
		{
			name:     "All duplicates",
			input:    "aaaa",
			expected: "",
		},
		{
			name:     "Alternating duplicates",
			input:    "abab",
			expected: "abab",
		},
		{
			name:     "Multiple consecutive duplicates",
			input:    "aaabbb",
			expected: "ab",
		},
		{
			name:     "Complex case with multiple removals",
			input:    "aabccba",
			expected: "a",
		},
		{
			name:     "Long string with duplicates",
			input:    "abcdefghhgfedcba",
			expected: "",
		},
		{
			name:     "Partial removal at end",
			input:    "abcdeed",
			expected: "abc",
		},
		{
			name:     "Partial removal at beginning",
			input:    "aabcde",
			expected: "bcde",
		},
		{
			name:     "Mixed duplicates",
			input:    "aabccb",
			expected: "",
		},
		{
			name:     "Single pair in middle",
			input:    "abcdeedfgh",
			expected: "abcfgh",
		},
		{
			name:     "Multiple pairs scattered",
			input:    "aabbccddee",
			expected: "",
		},
		{
			name:     "No adjacent duplicates",
			input:    "abcdef",
			expected: "abcdef",
		},
		{
			name:     "Duplicates at boundaries",
			input:    "aabcc",
			expected: "b",
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			result := strings.RemoveDuplicates(test.input)
			if result != test.expected {
				t.Errorf("RemoveDuplicates(%q) = %q; expected %q",
					test.input, result, test.expected)
			}
		})
	}
}

// Benchmark test for performance testing
func BenchmarkRemoveDuplicates(b *testing.B) {
	testString := "aabbccddeeaabbccddeeaabbccddee"
	for i := 0; i < b.N; i++ {
		strings.RemoveDuplicates(testString)
	}
}

// Test with very long string to check performance
func TestRemoveDuplicatesLongString(t *testing.T) {
	// Create a long string with many duplicates
	longString := ""
	for i := 0; i < 1000; i++ {
		longString += "ab"
	}

	result := strings.RemoveDuplicates(longString)
	if result != longString {
		t.Errorf("Expected unchanged string for alternating pattern, got %q", result)
	}
}

// Test with string that becomes empty after all removals
func TestRemoveDuplicatesBecomesEmpty(t *testing.T) {
	testCases := []string{
		"aa",
		"aaaa",
		"aabbcc",
	}

	for _, input := range testCases {
		result := strings.RemoveDuplicates(input)
		if result != "" {
			t.Errorf("RemoveDuplicates(%q) = %q; expected empty string", input, result)
		}
	}
}
