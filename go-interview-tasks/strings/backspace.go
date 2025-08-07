package strings

/*
Given two strings s and t, return true if they are equal when both are typed into empty text editors. '#' means a backspace character.

Note that after backspacing an empty text, the text will continue empty.



Example 1:

Input: s = "ab#c", t = "ad#c"
Output: true
Explanation: Both s and t become "ac".
Example 2:

Input: s = "ab##", t = "c#d#"
Output: true
Explanation: Both s and t become "".
Example 3:

Input: s = "a#c", t = "b"
Output: false
Explanation: s becomes "c" while t becomes "b".


Constraints:

1 <= s.length, t.length <= 200
s and t only contain lowercase letters and '#' characters.
*/

func BackspaceCompare(s string, t string) bool {
	runesS := make([]rune, 0, len(s))
	runesT := make([]rune, 0, len(t))

	for _, item := range s {
		if item == '#' {
			if len(runesS) > 0 {
				runesS = runesS[:len(runesS)-1]
			}
		} else {
			runesS = append(runesS, item)
		}
	}

	for _, item := range t {
		if item == '#' {
			if len(runesT) > 0 {
				runesT = runesT[:len(runesT)-1]
			}
		} else {
			runesT = append(runesT, item)
		}
	}

	return string(runesS) == string(runesT)
}
