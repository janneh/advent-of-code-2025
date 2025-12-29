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

func isInside(p Point, polygon []Point) bool {
	// Ray casting algorithm for point-in-polygon test
	count := 0
	n := len(polygon)

	for i := range n {
		j := (i + 1) % n
		p1, p2 := polygon[i], polygon[j]

		if (p1.y <= p.y && p2.y > p.y) || (p2.y <= p.y && p1.y > p.y) {
			// Edge crosses horizontal ray from p
			t := float64(p.y-p1.y) / float64(p2.y-p1.y)
			x := float64(p1.x) + t*float64(p2.x-p1.x)

			if float64(p.x) < x {
				count++
			}
		}
	}

	return count%2 == 1
}

func part2(tiles []Point) int {
	// Create a set of red tiles
	redTiles := make(map[Point]bool)
	for _, tile := range tiles {
		redTiles[tile] = true
	}

	// Pre-compute edge tiles for O(1) lookup
	edgeTiles := make(map[Point]bool)
	for i := range len(tiles) {
		next := (i + 1) % len(tiles)
		p1, p2 := tiles[i], tiles[next]

		if p1.x == p2.x {
			// Vertical line
			minY, maxY := min(p1.y, p2.y), max(p1.y, p2.y)
			for y := minY + 1; y < maxY; y++ {
				edgeTiles[Point{p1.x, y}] = true
			}
		} else if p1.y == p2.y {
			// Horizontal line
			minX, maxX := min(p1.x, p2.x), max(p1.x, p2.x)
			for x := minX + 1; x < maxX; x++ {
				edgeTiles[Point{x, p1.y}] = true
			}
		}
	}

	// Memoize point-in-polygon results for efficiency
	insideCache := make(map[Point]bool)

	isGreen := func(p Point) bool {
		if redTiles[p] || edgeTiles[p] {
			return true
		}

		// Check if inside polygon (with memoization)
		if val, ok := insideCache[p]; ok {
			return val
		}
		result := isInside(p, tiles)
		insideCache[p] = result
		return result
	}

	// Find the largest rectangle with red corners and all red/green tiles
	maxArea := 0

	for i := range len(tiles) {
		for j := i + 1; j < len(tiles); j++ {
			minX := min(tiles[i].x, tiles[j].x)
			maxX := max(tiles[i].x, tiles[j].x)
			minY := min(tiles[i].y, tiles[j].y)
			maxY := max(tiles[i].y, tiles[j].y)

			width := maxX - minX + 1
			height := maxY - minY + 1
			area := width * height

			// Skip rectangles that are too large to be all green
			// Increase threshold to see if we're missing larger rectangles
			if area > 5000000 {
				continue
			}

			// Check all tiles in the rectangle
			valid := true
			for x := minX; x <= maxX && valid; x++ {
				for y := minY; y <= maxY && valid; y++ {
					if !isGreen(Point{x, y}) {
						valid = false
					}
				}
			}

			if valid && area > maxArea {
				maxArea = area
				fmt.Fprintf(os.Stderr, "New max: %d (tiles %d,%d at (%d,%d)-(%d,%d) w=%d h=%d)\n",
					area, i, j, tiles[i].x, tiles[i].y, tiles[j].x, tiles[j].y, width, height)
			}
		}
		if i%100 == 0 {
			fmt.Fprintf(os.Stderr, "Progress: %d/%d, maxArea: %d\n", i, len(tiles), maxArea)
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
	fmt.Printf("Part 2: %d\n", part2(tiles))
}
