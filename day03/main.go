package main

import (
	"bufio"
	"fmt"
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

func part1(banks []string) int {
	total := 0
	for _, bank := range banks {
		total += maxJoltage(bank)
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
}
