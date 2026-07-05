package leetcode

func findSubstring(s string, words []string) []int {
	if len(s) == 0 || len(words) == 0 {
		return []int{}
	}

	wordLen := len(words[0])
	numWords := len(words)
	totalLen := wordLen * numWords
	n := len(s)

	if n < totalLen {
		return []int{}
	}

	wordCounts := make(map[string]int)
	for _, word := range words {
		wordCounts[word]++
	}

	var res []int

	// Iterate over each possible offset
	for i := 0; i < wordLen; i++ {
		left := i
		right := i
		currCounts := make(map[string]int)
		count := 0

		for right+wordLen <= n {
			w := s[right : right+wordLen]
			right += wordLen

			if _, ok := wordCounts[w]; ok {
				currCounts[w]++
				count++

				for currCounts[w] > wordCounts[w] {
					leftWord := s[left : left+wordLen]
					currCounts[leftWord]--
					count--
					left += wordLen
				}

				if count == numWords {
					res = append(res, left)
				}
			} else {
				currCounts = make(map[string]int)
				count = 0
				left = right
			}
		}
	}

	return res
}
