package leetcode

import (
	"reflect"
	"testing"
)

func TestUpdateMatrix(t *testing.T) {
	tests := []struct {
		mat      [][]int
		expected [][]int
	}{
		{
			mat: [][]int{
				{0, 0, 0},
				{0, 1, 0},
				{0, 0, 0},
			},
			expected: [][]int{
				{0, 0, 0},
				{0, 1, 0},
				{0, 0, 0},
			},
		},
		{
			mat: [][]int{
				{0, 0, 0},
				{0, 1, 0},
				{1, 1, 1},
			},
			expected: [][]int{
				{0, 0, 0},
				{0, 1, 0},
				{1, 2, 1},
			},
		},
	}

	for _, tt := range tests {
		actual := updateMatrix(tt.mat)
		if !reflect.DeepEqual(actual, tt.expected) {
			t.Errorf("updateMatrix(%v) = %v; expected %v", tt.mat, actual, tt.expected)
		}
	}
}
