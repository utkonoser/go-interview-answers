package main

import (
	"interviews/go-interview-tasks/strings"
	"testing"
)

func TestBackspaceCompare(t *testing.T) {
	tests := []struct {
		name     string
		s        string
		t        string
		expected bool
	}{
		{
			name:     "Example 1: ab#c vs ad#c",
			s:        "ab#c",
			t:        "ad#c",
			expected: true,
		},
		{
			name:     "Example 2: ab## vs c#d#",
			s:        "ab##",
			t:        "c#d#",
			expected: true,
		},
		{
			name:     "Example 3: a#c vs b",
			s:        "a#c",
			t:        "b",
			expected: false,
		},
		{
			name:     "Empty strings",
			s:        "",
			t:        "",
			expected: true,
		},
		{
			name:     "One empty string",
			s:        "a",
			t:        "",
			expected: false,
		},
		{
			name:     "Both become empty",
			s:        "a#",
			t:        "b#",
			expected: true,
		},
		{
			name:     "Multiple backspaces on empty",
			s:        "###",
			t:        "####",
			expected: true,
		},
		{
			name:     "No backspaces, equal",
			s:        "abc",
			t:        "abc",
			expected: true,
		},
		{
			name:     "No backspaces, different",
			s:        "abc",
			t:        "def",
			expected: false,
		},
		{
			name:     "Complex backspace sequence",
			s:        "ab##c#d#",
			t:        "##",
			expected: true,
		},
		{
			name:     "Backspace at beginning",
			s:        "#abc",
			t:        "abc",
			expected: true,
		},
		{
			name:     "Different lengths after processing",
			s:        "a##c",
			t:        "#a#c",
			expected: true,
		},
		{
			name:     "Long equal strings",
			s:        "a#b#c#d#e#f#",
			t:        "######",
			expected: true,
		},
		{
			name:     "Partial backspace",
			s:        "abc#def",
			t:        "abdef",
			expected: true,
		},
		{
			name:     "All backspaces vs empty",
			s:        "a#b#c#",
			t:        "",
			expected: true,
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			result := strings.BackspaceCompare(test.s, test.t)
			if result != test.expected {
				t.Errorf("BackspaceCompare(%q, %q) = %v; expected %v",
					test.s, test.t, result, test.expected)
			}
		})
	}
}
