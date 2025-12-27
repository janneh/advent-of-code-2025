package main

import (
	"bufio"
	"fmt"
	"os"
)

func countAccessible(grid []string) int {
	count := 0
	rows := len(grid)
	if rows == 0 {
		return 0
	}
	cols := len(grid[0])

	// Directions for 8 adjacent cells
	dirs := [][]int{
		{-1, -1},
		{-1, 0},
		{-1, 1},
		{0, -1},
		{0, 1},
		{1, -1},
		{1, 0},
		{1, 1},
	}

	for i := range rows {
		for j := range cols {
			if grid[i][j] == '@' {
				// Count adjacent paper rolls
				adjacent := 0
				for _, dir := range dirs {
					ni, nj := i+dir[0], j+dir[1]
					if ni >= 0 && ni < rows && nj >= 0 && nj < cols {
						if grid[ni][nj] == '@' {
							adjacent++
						}
					}
				}

				// Accessible if fewer than 4 adjacent rolls
				if adjacent < 4 {
					count++
				}
			}
		}
	}

	return count
}

func part1(grid []string) int {
	return countAccessible(grid)
}

func part2(grid []string) int {
	rows := len(grid)
	if rows == 0 {
		return 0
	}
	cols := len(grid[0])

	// Create a mutable 2D grid
	mutableGrid := make([][]byte, rows)
	for i := range rows {
		mutableGrid[i] = []byte(grid[i])
	}

	totalRemoved := 0

	// Directions for 8 adjacent cells
	dirs := [][]int{
		{-1, -1}, {-1, 0}, {-1, 1},
		{0, -1}, {0, 1},
		{1, -1}, {1, 0}, {1, 1},
	}

	// Keep removing until no more accessible rolls
	for {
		// Find all accessible rolls
		accessible := [][2]int{}

		for i := range rows {
			for j := range cols {
				if mutableGrid[i][j] == '@' {
					// Count adjacent paper rolls
					adjacent := 0
					for _, dir := range dirs {
						ni, nj := i+dir[0], j+dir[1]
						if ni >= 0 && ni < rows && nj >= 0 && nj < cols {
							if mutableGrid[ni][nj] == '@' {
								adjacent++
							}
						}
					}

					if adjacent < 4 {
						accessible = append(accessible, [2]int{i, j})
					}
				}
			}
		}

		// If no accessible rolls, stop
		if len(accessible) == 0 {
			break
		}

		// Remove all accessible rolls
		for _, pos := range accessible {
			mutableGrid[pos[0]][pos[1]] = '.'
			totalRemoved++
		}
	}

	return totalRemoved
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
		line := scanner.Text()
		if line != "" {
			grid = append(grid, line)
		}
	}

	if err := scanner.Err(); err != nil {
		fmt.Fprintf(os.Stderr, "Error reading file: %v\n", err)
		os.Exit(1)
	}

	fmt.Printf("Part 1: %d\n", part1(grid))
	fmt.Printf("Part 2: %d\n", part2(grid))
}
