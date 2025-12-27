package main

import (
	"bufio"
	"fmt"
	"math/big"
	"os"
)

func maxJoltage(bank string) int {
	maxJolt := 0

	// Try all pairs of positions (i, j) where i < j
	for i := 0; i < len(bank); i++ {
		for j := i + 1; j < len(bank); j++ {
			// Form the two-digit number from digits at positions i and j
			d1 := int(bank[i] - '0')
			d2 := int(bank[j] - '0')
			joltage := d1*10 + d2

			if joltage > maxJolt {
				maxJolt = joltage
			}
		}
	}

	return maxJolt
}

func maxJoltagePart2(bank string) string {
	n := len(bank)
	k := 12

	if n < k {
		return bank
	}

	result := make([]byte, 0, k)
	lastPos := -1

	// Greedy selection: for each position, pick the largest digit
	// while ensuring we have enough digits remaining
	for i := 0; i < k; i++ {
		remaining := k - i
		maxDigit := byte('0')
		maxPos := lastPos + 1

		// Search for the maximum digit in the valid range
		// We can search from (lastPos+1) to (n-remaining)
		for j := lastPos + 1; j <= n-remaining; j++ {
			if bank[j] > maxDigit {
				maxDigit = bank[j]
				maxPos = j
			}
		}

		result = append(result, maxDigit)
		lastPos = maxPos
	}

	return string(result)
}

func part1(banks []string) int {
	total := 0
	for _, bank := range banks {
		total += maxJoltage(bank)
	}
	return total
}

func part2(banks []string) *big.Int {
	total := big.NewInt(0)
	for _, bank := range banks {
		joltageStr := maxJoltagePart2(bank)
		joltage := new(big.Int)
		joltage.SetString(joltageStr, 10)
		total.Add(total, joltage)
	}
	return total
}

func main() {
	file, err := os.Open("input.txt")
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error opening file: %v\n", err)
		os.Exit(1)
	}
	defer file.Close()

	var banks []string
	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		line := scanner.Text()
		if line != "" {
			banks = append(banks, line)
		}
	}

	if err := scanner.Err(); err != nil {
		fmt.Fprintf(os.Stderr, "Error reading file: %v\n", err)
		os.Exit(1)
	}

	fmt.Printf("Part 1: %d\n", part1(banks))
	fmt.Printf("Part 2: %s\n", part2(banks))
}
