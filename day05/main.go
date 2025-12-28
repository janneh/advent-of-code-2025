package main

import (
	"bufio"
	"fmt"
	"os"
	"strconv"
	"strings"
)

type Range struct {
	start, end int
}

func parseInput(filename string) ([]Range, []int, error) {
	file, err := os.Open(filename)
	if err != nil {
		return nil, nil, err
	}
	defer file.Close()

	var ranges []Range
	var ids []int
	parsingRanges := true

	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		line := strings.TrimSpace(scanner.Text())

		if line == "" {
			parsingRanges = false
			continue
		}

		if parsingRanges {
			// Parse range like "3-5"
			parts := strings.Split(line, "-")
			if len(parts) == 2 {
				start, err1 := strconv.Atoi(parts[0])
				end, err2 := strconv.Atoi(parts[1])
				if err1 == nil && err2 == nil {
					ranges = append(ranges, Range{start, end})
				}
			}
		} else {
			// Parse single ID
			id, err := strconv.Atoi(line)
			if err == nil {
				ids = append(ids, id)
			}
		}
	}

	if err := scanner.Err(); err != nil {
		return nil, nil, err
	}

	return ranges, ids, nil
}

func isFresh(id int, ranges []Range) bool {
	for _, r := range ranges {
		if id >= r.start && id <= r.end {
			return true
		}
	}
	return false
}

func part1(ranges []Range, ids []int) int {
	count := 0
	for _, id := range ids {
		if isFresh(id, ranges) {
			count++
		}
	}
	return count
}

func main() {
	ranges, ids, err := parseInput("input.txt")
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error reading input: %v\n", err)
		os.Exit(1)
	}

	fmt.Printf("Part 1: %d\n", part1(ranges, ids))
}
