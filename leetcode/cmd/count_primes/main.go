package main

import (
	"fmt"
	"time"

	"github.com/clyang82/leetcode"
)

func main() {
	// 1. Calculate number of primes up to 8^10
	limit := 1073741824 // 8^10
	fmt.Printf("Calculating the number of primes up to %d...\n", limit)

	start := time.Now()
	count := leetcode.CountPrimes(limit)
	elapsed := time.Since(start)

	fmt.Printf("Primes count: %d\n", count)
	fmt.Printf("Time elapsed: %v\n\n", elapsed)

	// 2. Get and print the first 100 prime numbers
	fmt.Println("Here are the first 100 prime numbers:")
	primes := leetcode.GetFirstNPrimes(100)
	
	// Format into nice rows of 10 for readability
	for i, p := range primes {
		fmt.Printf("%4d ", p)
		if (i+1)%10 == 0 {
			fmt.Println()
		}
	}
}
