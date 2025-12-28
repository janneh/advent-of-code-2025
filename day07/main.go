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
}
