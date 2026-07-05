package leetcode

import (
	"fmt"
	"math"
	"testing"
)

func TestMyAtoi(t *testing.T) {
	tests := []struct {
		s        string
		expected int
	}{
		{"42", 42},
		{"   -42", -42},
		{"4193 with words", 4193},
		{"words and 987", 0},
		{"-91283472332", math.MinInt32},
		{"+1", 1},
		{"  +0 123", 0},
		{"2147483647", 2147483647},
		{"2147483648", 2147483647},
		{"-2147483648", -2147483648},
		{"-2147483649", -2147483648},
		{"", 0},
		{" ", 0},
		{"+-12", 0},
		{"00000-42a123", 0},
	}

	for _, tt := range tests {
		actual := myAtoi(tt.s)
		if actual != tt.expected {
			t.Errorf("myAtoi(%q) = %d; expected %d", tt.s, actual, tt.expected)
		}
	}
}

func TestXxx(t *testing.T) {
	s := "12abcd"
	fmt.Println(s[0])
	fmt.Println(s[0] - '0')
	t.Fatal("fail")
}
