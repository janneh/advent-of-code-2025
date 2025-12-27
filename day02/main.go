package main

import (
	"fmt"
	"os"
	"strconv"
	"strings"
)

func isInvalidID(n int) bool {
	s := strconv.Itoa(n)
	// Must have even length to be splittable into two equal parts
	if len(s)%2 != 0 {
		return false
	}

	mid := len(s) / 2
	left := s[:mid]
	right := s[mid:]

	// Check if both halves are equal
	return left == right
}

func parseRanges(input string) ([]struct{ start, end int }, error) {
	// Remove any whitespace and newlines
	input = strings.TrimSpace(input)

	// Split by comma
	parts := strings.Split(input, ",")

	var ranges []struct{ start, end int }
	for _, part := range parts {
		part = strings.TrimSpace(part)
		if part == "" {
			continue
		}

		// Split by dash
		rangeParts := strings.Split(part, "-")
		if len(rangeParts) != 2 {
			return nil, fmt.Errorf("invalid range format: %s", part)
		}

		start, err := strconv.Atoi(strings.TrimSpace(rangeParts[0]))
		if err != nil {
			return nil, fmt.Errorf("invalid start number: %s", rangeParts[0])
		}

		end, err := strconv.Atoi(strings.TrimSpace(rangeParts[1]))
		if err != nil {
			return nil, fmt.Errorf("invalid end number: %s", rangeParts[1])
		}

		ranges = append(ranges, struct{ start, end int }{start, end})
	}

	return ranges, nil
}

func isInvalidIDPart2(n int) bool {
	s := strconv.Itoa(n)
	length := len(s)

	// Try all possible pattern lengths
	for patternLen := 1; patternLen <= length/2; patternLen++ {
		// Pattern length must divide the total length evenly
		if length%patternLen != 0 {
			continue
		}

		pattern := s[:patternLen]
		repetitions := length / patternLen

		// Check if repeating the pattern creates the entire string
		if repetitions >= 2 {
			repeated := strings.Repeat(pattern, repetitions)
			if repeated == s {
				return true
			}
		}
	}

	return false
}

func part1(input string) int {
	ranges, err := parseRanges(input)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error parsing ranges: %v\n", err)
		return 0
	}

	sum := 0
	for _, r := range ranges {
		for id := r.start; id <= r.end; id++ {
			if isInvalidID(id) {
				sum += id
			}
		}
	}

	return sum
}

func part2(input string) int {
	ranges, err := parseRanges(input)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error parsing ranges: %v\n", err)
		return 0
	}

	sum := 0
	for _, r := range ranges {
		for id := r.start; id <= r.end; id++ {
			if isInvalidIDPart2(id) {
				sum += id
			}
		}
	}

	return sum
}

func main() {
	data, err := os.ReadFile("input.txt")
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error reading input: %v\n", err)
		os.Exit(1)
	}

	input := string(data)
	fmt.Printf("Part 1: %d\n", part1(input))
	fmt.Printf("Part 2: %d\n", part2(input))
}
