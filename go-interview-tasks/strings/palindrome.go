package strings

import (
	"unicode"
)

/*
https://leetcode.com/problems/valid-palindrome/
A phrase is a palindrome if, after converting all uppercase letters into lowercase letters and removing all non-alphanumeric characters, it reads the same forward and backward. Alphanumeric characters include letters and numbers.

Given a string s, return true if it is a palindrome, or false otherwise.



Example 1:

Input: s = "A man, a plan, a canal: Panama"
Output: true
Explanation: "amanaplanacanalpanama" is a palindrome.
Example 2:

Input: s = "race a car"
Output: false
Explanation: "raceacar" is not a palindrome.
Example 3:

Input: s = " "
Output: true
Explanation: s is an empty string "" after removing non-alphanumeric characters.
Since an empty string reads the same forward and backward, it is a palindrome.


Constraints:

1 <= s.length <= 2 * 105
s consists only of printable ASCII characters.
*/

// isAlphanumeric проверяет, является ли символ буквой или цифрой
func isAlphanumeric(r rune) bool {
	return unicode.IsLetter(r) || unicode.IsDigit(r)
}

// toLowerRune преобразует символ в нижний регистр
func toLowerRune(r rune) rune {
	return unicode.ToLower(r)
}

func IsPalindrome(s string) bool {
	if len(s) == 0 {
		return true
	}

	left, right := 0, len(s)-1

	for left < right {
		// Пропускаем неалфавитно-цифровые символы слева
		for left < right && !isAlphanumeric(rune(s[left])) {
			left++
		}

		// Пропускаем неалфавитно-цифровые символы справа
		for left < right && !isAlphanumeric(rune(s[right])) {
			right--
		}

		// Если указатели встретились, значит палиндром
		if left >= right {
			break
		}

		// Сравниваем символы в нижнем регистре
		if toLowerRune(rune(s[left])) != toLowerRune(rune(s[right])) {
			return false
		}

		left++
		right--
	}

	return true
}
