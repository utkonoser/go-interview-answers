package main

import (
	"interviews/go-interview-tasks/strings"
	"testing"
)

func TestRemoveStars(t *testing.T) {
	tests := []struct {
		name     string
		input    string
		expected string
	}{
		{
			name:     "Example 1: leet**cod*e",
			input:    "leet**cod*e",
			expected: "lecoe",
		},
		{
			name:     "Example 2: erase*****",
			input:    "erase*****",
			expected: "",
		},
		{
			name:     "Empty string",
			input:    "",
			expected: "",
		},
		{
			name:     "No stars",
			input:    "hello",
			expected: "hello",
		},
		{
			name:     "Single star at beginning",
			input:    "*hello",
			expected: "hello",
		},
		{
			name:     "Single star at end",
			input:    "hello*",
			expected: "hell",
		},
		{
			name:     "Multiple consecutive stars",
			input:    "abc***",
			expected: "",
		},
		{
			name:     "Stars with single character between",
			input:    "a*b*c*d*",
			expected: "",
		},
		{
			name:     "Complex pattern 1",
			input:    "ab*c*d*e*",
			expected: "a",
		},
		{
			name:     "Complex pattern 2",
			input:    "abc*def*ghi*",
			expected: "abdegh",
		},
		{
			name:     "Stars at beginning and end",
			input:    "*hello*",
			expected: "hell",
		},
		{
			name:     "Alternating stars and letters",
			input:    "a*b*c*d*e*",
			expected: "",
		},
		{
			name:     "Long string with stars",
			input:    "abcdefghijklmnopqrstuvwxyz*",
			expected: "abcdefghijklmnopqrstuvwxy",
		},
		{
			name:     "String with only stars",
			input:    "*****",
			expected: "",
		},
		{
			name:     "Single character with star",
			input:    "a*",
			expected: "",
		},
		{
			name:     "Star followed by characters",
			input:    "*abc",
			expected: "abc",
		},
		{
			name:     "Characters between stars",
			input:    "a*b*c",
			expected: "c",
		},
		{
			name:     "Multiple stars in middle",
			input:    "hello***world",
			expected: "heworld",
		},
		{
			name:     "Stars with spaces",
			input:    "hello *world*",
			expected: "helloworl",
		},
		{
			name:     "Unicode characters",
			input:    "привет*мир*",
			expected: "привеми",
		},
		{
			name:     "Numbers and stars",
			input:    "123*456*789*",
			expected: "124578",
		},
		{
			name:     "Mixed characters and stars",
			input:    "a1b*2c*3d*",
			expected: "a123",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := strings.RemoveStars(tt.input)
			if result != tt.expected {
				t.Errorf("RemoveStars(%q) = %q, want %q", tt.input, result, tt.expected)
			}
		})
	}
}

func TestRemoveStarsEdgeCases(t *testing.T) {
	t.Run("Very long string", func(t *testing.T) {
		// Создаем длинную строку для тестирования производительности
		longStr := ""
		for range 1000 {
			longStr += "a"
		}
		longStr += "*"

		result := strings.RemoveStars(longStr)
		expected := longStr[:len(longStr)-2] // убираем последний символ и звезду

		if result != expected {
			t.Errorf("RemoveStars(long string) = %q, want %q", result, expected)
		}
	})

	t.Run("String with many stars", func(t *testing.T) {
		input := "a*b*c*d*e*f*g*h*i*j*k*l*m*n*o*p*q*r*s*t*u*v*w*x*y*z*"
		result := strings.RemoveStars(input)
		expected := ""

		if result != expected {
			t.Errorf("RemoveStars(%q) = %q, want %q", input, result, expected)
		}
	})

	t.Run("String with stars at every other position", func(t *testing.T) {
		input := "a*b*c*d*e*f*g*h*i*j*"
		result := strings.RemoveStars(input)
		expected := ""

		if result != expected {
			t.Errorf("RemoveStars(%q) = %q, want %q", input, result, expected)
		}
	})
}

func BenchmarkRemoveStars(b *testing.B) {
	testString := "hello*world*this*is*a*test*string*with*many*stars*"

	for b.Loop() {
		strings.RemoveStars(testString)
	}
}

func BenchmarkRemoveStarsLong(b *testing.B) {
	// Создаем длинную строку для бенчмарка
	longStr := ""
	for range 10000 {
		longStr += "a"
	}
	longStr += "*"

	for b.Loop() {
		strings.RemoveStars(longStr)
	}
}
