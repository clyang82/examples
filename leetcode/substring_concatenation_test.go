package leetcode

import (
	"reflect"
	"sort"
	"testing"
)

func TestFindSubstring(t *testing.T) {
	tests := []struct {
		s        string
		words    []string
		expected []int
	}{
		{
			s:        "barfoothefoobarman",
			words:    []string{"foo", "bar"},
			expected: []int{0, 9},
		},
		{
			s:        "wordgoodgoodgoodbestword",
			words:    []string{"word", "good", "best", "word"},
			expected: []int{},
		},
		{
			s:        "barfoofoobarthefoobarman",
			words:    []string{"bar", "foo", "the"},
			expected: []int{6, 9, 12},
		},
		{
			s:        "wordgoodgoodgoodbestword",
			words:    []string{"word", "good", "best", "good"},
			expected: []int{8},
		},
		{
			s:        "lingmindraboofooowingdingbarwingmonkeypoundcake",
			words:    []string{"fooo", "barw", "wing", "ding", "wing"},
			expected: []int{13},
		},
		{
			s:        "aaaaaaaaaaaaaa",
			words:    []string{"aa", "aa"},
			expected: []int{0, 1, 2, 3, 4, 5, 6, 7, 8, 9, 10},
		},
	}

	for _, tt := range tests {
		actual := findSubstring(tt.s, tt.words)
		sort.Ints(actual)
		sort.Ints(tt.expected)
		if !reflect.DeepEqual(actual, tt.expected) {
			t.Errorf("findSubstring(%q, %v) = %v; expected %v", tt.s, tt.words, actual, tt.expected)
		}
	}
}
