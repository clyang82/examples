package leetcode

import (
	"fmt"
)

// LengthOfLongestSubstring finds the length of the longest substring without repeating characters.
func LengthOfLongestSubstring(s string) int {
	charMap := make(map[rune]int)
	maxLength := 0
	start := 0

	for end, char := range s {
		if pos, ok := charMap[char]; ok && pos >= start {
			start = pos + 1
		}
		charMap[char] = end
		currentLength := end - start + 1
		if currentLength > maxLength {
			maxLength = currentLength
		}
	}

	return maxLength
}

func main() {
	testCases := []struct {
		s        string
		expected int
	}{
		{"abcabcbb", 3},
		{"bbbbb", 1},
		{"pwwkew", 3},
		{"", 0},
		{" ", 1},
		{"au", 2},
		{"dvdf", 3},
	}

	for _, tc := range testCases {
		result := LengthOfLongestSubstring(tc.s)
		fmt.Printf("Input: s = \"%s\"\n", tc.s)
		fmt.Printf("Output: %d\n", result)
		fmt.Printf("Expected: %d\n", tc.expected)
		if result == tc.expected {
			fmt.Println("Status: PASSED")
		} else {
			fmt.Println("Status: FAILED")
		}
		fmt.Println("--------------------")
	}
}
