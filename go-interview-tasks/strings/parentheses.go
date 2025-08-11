package strings

/*
Given a string s containing just the characters '(', ')', '{', '}', '[' and ']', determine if the input string is valid.

An input string is valid if:

Open brackets must be closed by the same type of brackets.
Open brackets must be closed in the correct order.
Every close bracket has a corresponding open bracket of the same type.

Constraints:

1 <= s.length <= 104
s consists of parentheses only '()[]{}'.
*/

func IsValidParentheses(s string) bool {
	stack := []rune{}

	for _, char := range s {
		switch char {
		case '(', '[', '{':
			stack = append(stack, char)
		case ')', ']', '}':
			if len(stack) == 0 {
				return false
			}
			top := stack[len(stack)-1]
			if (char == ')' && top != '(') || (char == ']' && top != '[') || (char == '}' && top != '{') {
				return false
			}
			stack = stack[:len(stack)-1]
		}
	}
	return len(stack) == 0
}
