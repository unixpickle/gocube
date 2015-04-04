package gocube

var factorials = []int{1, 1, 2, 6, 24, 120, 720, 5040, 40320, 362880, 3628800,
	39916800, 479001600}

var pascalsTriangle = [][]int{
	[]int{1},
	[]int{1, 1},
	[]int{1, 2, 1},
	[]int{1, 3, 3, 1},
	[]int{1, 4, 6, 4, 1},
	[]int{1, 5, 10, 10, 5, 1},
	[]int{1, 6, 15, 20, 15, 6, 1},
	[]int{1, 7, 21, 35, 35, 21, 7, 1},
	[]int{1, 8, 28, 56, 70, 56, 28, 8, 1},
	[]int{1, 9, 36, 84, 126, 126, 84, 36, 9, 1},
	[]int{1, 10, 45, 120, 210, 252, 210, 120, 45, 10, 1},
	[]int{1, 11, 55, 165, 330, 462, 462, 330, 165, 55, 11, 1},
	[]int{1, 12, 66, 220, 495, 792, 924, 792, 495, 220, 66, 12, 1},
}

func allPermutations(size int) [][]int {
	if size == 0 {
		return [][]int{[]int{}}
	} else if size == 1 {
		return [][]int{[]int{0}}
	}
	
	result := make([][]int, 0, factorial(size))
	subPermutations := allPermutations(size - 1)

	// For every starting element, go through every sub permutation and generate
	// a new permutation
	for start := 0; start < size; start++ {
		for _, subPerm := range subPermutations {
			perm := append([]int{start}, subPerm...)
			// Increment values which are >= start
			for j := 1; j < len(perm); j++ {
				if perm[j] >= start {
					perm[j]++
				}
			}
			result = append(result, perm)
		}
	}

	return result
}

func choose(a, b int) int {
	if a < 13 {
		return pascalsTriangle[a][b]
	}
	res := 1
	for i := 0; i < b; i++ {
		res *= a
		a -= 1
	}
	return res / factorial(b)
}

func encodeChoice(choice []bool) int {
	trueCount := 0
	for _, x := range choice {
		if x {
			trueCount++
		}
	}
	return encodeExplicitChoice(choice, 0, trueCount)
}

func encodeExplicitChoice(choice []bool, start, numTrue int) int {
	if len(choice)-start <= 1 || numTrue == 0 {
		return 0
	} else if numTrue == 1 {
		for i := start; i < len(choice); i++ {
			if choice[i] {
				return i - start
			}
		}
	}

	numMissed := 0
	for i := start; i < len(choice); i++ {
		if choice[i] {
			subValue := encodeExplicitChoice(choice, i+1, numTrue-1)
			return subValue + numMissed
		}
		numMissed += choose(len(choice)-(i+1), numTrue-1)
	}

	panic("internal inconsistency in encodeExplicitChoice")
	return -1
}

func encodePermutation(perm []int) int {
	c := make([]int, len(perm))
	copy(c, perm)
	return encodePermutationInPlace(c)
}

func encodePermutationInPlace(perm []int) int {
	if len(perm) == 0 {
		return 0
	}

	count := len(perm) - 1
	result := 0
	factorial := factorial(count)

	for i := 0; i < count; i++ {
		current := perm[i]

		// Add the element to the result.
		result += factorial * current
		factorial /= count - i

		// Shift all the elements which were above the current element.
		for j := i + 1; j < count; j++ {
			if perm[j] > current {
				perm[j]--
			}
		}
	}

	return result
}

func factorial(n int) int {
	if n >= len(factorials) {
		return n * factorial(n-1)
	}
	return factorials[n]
}
