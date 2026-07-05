package leetcode

import (
	"math"
	"math/bits"
)

// CountPrimes returns the number of prime numbers between 1 and n (inclusive).
// It uses an optimized, odd-only, bit-packed Sieve of Eratosthenes.
func CountPrimes(n int) int {
	if n < 2 {
		return 0
	}
	if n == 2 {
		return 1
	}

	// We only sieve odd numbers starting from 3.
	// Find the largest odd number <= n.
	maxOdd := n
	if maxOdd%2 == 0 {
		maxOdd--
	}
	if maxOdd < 3 {
		return 1 // only the prime 2
	}

	// The number of odd numbers in the range [3, maxOdd]
	numOdd := (maxOdd-3)/2 + 1

	// Pack 64 odd numbers into each uint64 word.
	// 0 bit means prime, 1 bit means composite.
	sieve := make([]uint64, (numOdd+63)/64)

	limit := int(math.Sqrt(float64(n)))
	for p := 3; p <= limit; p += 2 {
		pIdx := (p - 3) / 2
		word := pIdx / 64
		bit := pIdx % 64
		if (sieve[word] & (1 << bit)) == 0 {
			// p is prime. Mark its odd multiples starting from p * p.
			p2 := p * 2
			for j := p * p; j <= n; j += p2 {
				jIdx := (j - 3) / 2
				jWord := jIdx / 64
				jBit := jIdx % 64
				sieve[jWord] |= (1 << jBit)
			}
		}
	}

	// Mask out the unused bits in the last uint64 word so they are marked as 1 (composite).
	if rem := numOdd % 64; rem != 0 {
		mask := ^uint64(0) << rem
		sieve[len(sieve)-1] |= mask
	}

	// Count the composites (1-bits)
	composites := 0
	for _, word := range sieve {
		composites += bits.OnesCount64(word)
	}

	// Total primes = (total odd positions) - (composites) + 1 (for the prime 2)
	totalPositions := len(sieve) * 64
	return totalPositions - composites + 1
}

// GetFirstNPrimes returns a slice containing the first n prime numbers.
func GetFirstNPrimes(n int) []int {
	if n <= 0 {
		return nil
	}
	primes := make([]int, 0, n)
	primes = append(primes, 2)
	if n == 1 {
		return primes
	}

	candidate := 3
	for len(primes) < n {
		isPrime := true
		limit := int(math.Sqrt(float64(candidate)))
		for _, p := range primes {
			if p > limit {
				break
			}
			if candidate%p == 0 {
				isPrime = false
				break
			}
		}
		if isPrime {
			primes = append(primes, candidate)
		}
		candidate += 2
	}
	return primes
}

