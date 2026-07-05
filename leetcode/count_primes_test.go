package leetcode

import (
	"testing"
)

func TestCountPrimes(t *testing.T) {
	tests := []struct {
		n        int
		expected int
	}{
		{n: -5, expected: 0},
		{n: 0, expected: 0},
		{n: 1, expected: 0},
		{n: 2, expected: 1},
		{n: 3, expected: 2},
		{n: 4, expected: 2},
		{n: 5, expected: 3},
		{n: 10, expected: 4}, // 2, 3, 5, 7
		{n: 100, expected: 25},
		{n: 1000, expected: 168},
		{n: 10000, expected: 1229},
		{n: 100000, expected: 9592},
		{n: 1000000, expected: 78498},
	}

	for _, tc := range tests {
		actual := CountPrimes(tc.n)
		if actual != tc.expected {
			t.Errorf("CountPrimes(%d) = %d; expected %d", tc.n, actual, tc.expected)
		}
	}
}

func TestCountPrimesLarge(t *testing.T) {
	if testing.Short() {
		t.Skip("skipping large test in short mode")
	}

	limit := 1073741824 // 8^10
	expected := 54400028
	actual := CountPrimes(limit)
	if actual != expected {
		t.Errorf("CountPrimes(%d) = %d; expected %d", limit, actual, expected)
	}
}

func BenchmarkCountPrimesLarge(b *testing.B) {
	limit := 1073741824 // 8^10
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		CountPrimes(limit)
	}
}

func TestGetFirstNPrimes(t *testing.T) {
	// First 10 primes: 2, 3, 5, 7, 11, 13, 17, 19, 23, 29
	expected := []int{2, 3, 5, 7, 11, 13, 17, 19, 23, 29}
	actual := GetFirstNPrimes(10)

	if len(actual) != len(expected) {
		t.Fatalf("expected %d primes, got %d", len(expected), len(actual))
	}

	for i, v := range expected {
		if actual[i] != v {
			t.Errorf("at index %d: expected %d, got %d", i, v, actual[i])
		}
	}

	// Test negative/zero
	if GetFirstNPrimes(0) != nil {
		t.Error("expected nil for 0 primes")
	}
	if GetFirstNPrimes(-5) != nil {
		t.Error("expected nil for negative primes")
	}
}

