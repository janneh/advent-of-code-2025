package main

import (
	"bufio"
	"fmt"
	"os"
	"strconv"
	"strings"
)

type Point struct {
	x, y int
}

func parseInput(filename string) ([]Point, error) {
	file, err := os.Open(filename)
	if err != nil {
		return nil, err
	}
	defer file.Close()

	var tiles []Point
	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		line := strings.TrimSpace(scanner.Text())
		if line == "" {
			continue
		}

		parts := strings.Split(line, ",")
		if len(parts) != 2 {
			continue
		}

		x, err := strconv.Atoi(strings.TrimSpace(parts[0]))
		if err != nil {
			return nil, err
		}
		y, err := strconv.Atoi(strings.TrimSpace(parts[1]))
		if err != nil {
			return nil, err
		}

		tiles = append(tiles, Point{x, y})
	}

	if err := scanner.Err(); err != nil {
		return nil, err
	}

	return tiles, nil
}

func abs(n int) int {
	if n < 0 {
		return -n
	}
	return n
}

func part1(tiles []Point) int {
	maxArea := 0

	// Try all pairs of tiles as opposite corners
	for i := range len(tiles) {
		for j := i + 1; j < len(tiles); j++ {
			// Include both corner tiles in the count
			width := abs(tiles[j].x-tiles[i].x) + 1
			height := abs(tiles[j].y-tiles[i].y) + 1
			area := width * height

			if area > maxArea {
				maxArea = area
			}
		}
	}

	return maxArea
}

func main() {
	tiles, err := parseInput("input.txt")
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error reading input: %v\n", err)
		os.Exit(1)
	}

	fmt.Printf("Part 1: %d\n", part1(tiles))
}
