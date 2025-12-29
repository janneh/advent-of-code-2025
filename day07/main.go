package main

import (
	"bufio"
	"fmt"
	"os"
)

type Beam struct {
	row, col int
}

func findStart(grid []string) (int, int) {
	for i, row := range grid {
		for j, ch := range row {
			if ch == 'S' {
				return i, j
			}
		}
	}
	return -1, -1
}

func part1(grid []string) int {
	rows := len(grid)
	if rows == 0 {
		return 0
	}
	cols := len(grid[0])

	startRow, startCol := findStart(grid)
	if startRow == -1 {
		return 0
	}

	splitCount := 0
	beams := []Beam{{startRow, startCol}}

	// Simulate beams moving downward
	for len(beams) > 0 {
		newBeams := []Beam{}
		splittersHit := make(map[[2]int]bool)

		for _, beam := range beams {
			// Move beam down one step
			nextRow := beam.row + 1

			// Check if beam exits the grid
			if nextRow >= rows {
				continue
			}

			// Check if beam hits a splitter
			if grid[nextRow][beam.col] == '^' {
				splittersHit[[2]int{nextRow, beam.col}] = true
				// Create two new beams at left and right of splitter
				if beam.col-1 >= 0 {
					newBeams = append(newBeams, Beam{nextRow, beam.col - 1})
				}
				if beam.col+1 < cols {
					newBeams = append(newBeams, Beam{nextRow, beam.col + 1})
				}
			} else {
				// Beam continues downward
				newBeams = append(newBeams, Beam{nextRow, beam.col})
			}
		}

		// Count unique splitters hit in this iteration
		splitCount += len(splittersHit)

		// Deduplicate beams at same position
		beamSet := make(map[Beam]bool)
		for _, beam := range newBeams {
			beamSet[beam] = true
		}
		beams = []Beam{}
		for beam := range beamSet {
			beams = append(beams, beam)
		}
	}

	return splitCount
}

func part2(grid []string) int {
	rows := len(grid)
	if rows == 0 {
		return 0
	}
	cols := len(grid[0])

	startRow, startCol := findStart(grid)
	if startRow == -1 {
		return 0
	}

	// Memoization for counting paths from each position
	memo := make(map[[2]int]int)

	var countPaths func(row, col int) int
	countPaths = func(row, col int) int {
		// Check bounds
		if col < 0 || col >= cols {
			return 0
		}

		// Check if we've exited the grid
		if row >= rows {
			return 1
		}

		// Check memo
		key := [2]int{row, col}
		if val, ok := memo[key]; ok {
			return val
		}

		// Check next row
		nextRow := row + 1
		if nextRow >= rows {
			// About to exit
			memo[key] = 1
			return 1
		}

		// Check if next position is a splitter
		var result int
		if grid[nextRow][col] == '^' {
			// Quantum split: particle takes both paths
			result = countPaths(nextRow, col-1) + countPaths(nextRow, col+1)
		} else {
			// Continue downward
			result = countPaths(nextRow, col)
		}

		memo[key] = result
		return result
	}

	return countPaths(startRow, startCol)
}

func main() {
	file, err := os.Open("input.txt")
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error opening file: %v\n", err)
		os.Exit(1)
	}
	defer file.Close()

	var grid []string
	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		grid = append(grid, scanner.Text())
	}

	if err := scanner.Err(); err != nil {
		fmt.Fprintf(os.Stderr, "Error reading file: %v\n", err)
		os.Exit(1)
	}

	fmt.Printf("Part 1: %d\n", part1(grid))
	fmt.Printf("Part 2: %d\n", part2(grid))
}
