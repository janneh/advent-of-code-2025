package main

import (
	"bufio"
	"fmt"
	"os"
	"sort"
	"strconv"
	"strings"
)

type Point struct {
	x, y int
}

type Edge struct {
	x1, y1, x2, y2 int
}

type RectCandidate struct {
	i, j, area int
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

func isOnEdge(p Point, tiles []Point) bool {
	for i := 0; i < len(tiles); i++ {
		next := (i + 1) % len(tiles)
		from := tiles[i]
		to := tiles[next]

		// Check if p is on the line segment from -> to
		if from.x == to.x && p.x == from.x {
			// Same column
			minY, maxY := from.y, to.y
			if minY > maxY {
				minY, maxY = maxY, minY
			}
			if p.y >= minY && p.y <= maxY {
				return true
			}
		} else if from.y == to.y && p.y == from.y {
			// Same row
			minX, maxX := from.x, to.x
			if minX > maxX {
				minX, maxX = maxX, minX
			}
			if p.x >= minX && p.x <= maxX {
				return true
			}
		}
	}
	return false
}

func part2(tiles []Point) int {
	redTiles := make(map[Point]bool)
	for _, tile := range tiles {
		redTiles[tile] = true
	}

	edges := make([]Edge, 0, len(tiles))
	for i := range len(tiles) {
		next := (i + 1) % len(tiles)
		edges = append(edges, Edge{tiles[i].x, tiles[i].y, tiles[next].x, tiles[next].y})
	}

	intersectsEdge := func(minX, maxX, minY, maxY int, edge Edge) bool {
		edgeMinX, edgeMaxX := edge.x1, edge.x2
		if edgeMinX > edgeMaxX {
			edgeMinX, edgeMaxX = edgeMaxX, edgeMinX
		}
		edgeMinY, edgeMaxY := edge.y1, edge.y2
		if edgeMinY > edgeMaxY {
			edgeMinY, edgeMaxY = edgeMaxY, edgeMinY
		}
		return minX < edgeMaxX && maxX > edgeMinX && minY < edgeMaxY && maxY > edgeMinY
	}

	var candidates []RectCandidate
	areaLimit := part1(tiles) // Use theoretical max

	for i := range len(tiles) {
		for j := i + 1; j < len(tiles); j++ {
			p1 := tiles[i]
			p2 := tiles[j]

			if p1.x == p2.x || p1.y == p2.y {
				continue
			}

			width := abs(p2.x-p1.x) + 1
			height := abs(p2.y-p1.y) + 1
			area := width * height

			if area > areaLimit {
				continue
			}

			candidates = append(candidates, RectCandidate{i, j, area})
		}
	}

	fmt.Fprintf(os.Stderr, "Generated %d candidates, now sorting...\n", len(candidates))

	sort.Slice(candidates, func(a, b int) bool {
		return candidates[a].area > candidates[b].area
	})

	fmt.Fprintf(os.Stderr, "Sorting complete, starting validation...\n")

	insideCache := make(map[Point]bool)

	isValidTileCached := func(p Point) bool {
		if redTiles[p] {
			return true
		}
		if result, ok := insideCache[p]; ok {
			return result
		}
		valid := isOnEdge(p, tiles) || isInside(p, tiles)
		insideCache[p] = valid
		return valid
	}

	maxArea := 0
	checked := 0

	for _, cand := range candidates {
		if cand.area <= maxArea {
			break
		}

		checked++
		if checked%10000 == 0 {
			fmt.Fprintf(os.Stderr, "Checked %d/%d candidates, maxArea: %d, current: %d\n",
				checked, len(candidates), maxArea, cand.area)
		}

		p1 := tiles[cand.i]
		p2 := tiles[cand.j]

		minX, maxX := p1.x, p2.x
		if minX > maxX {
			minX, maxX = maxX, minX
		}
		minY, maxY := p1.y, p2.y
		if minY > maxY {
			minY, maxY = maxY, minY
		}

		if !isValidTileCached(Point{minX, minY}) || !isValidTileCached(Point{maxX, maxY}) ||
			!isValidTileCached(Point{minX, maxY}) || !isValidTileCached(Point{maxX, minY}) {
			continue
		}

		hasIntersection := false
		for _, edge := range edges {
			if intersectsEdge(minX, maxX, minY, maxY, edge) {
				hasIntersection = true
				break
			}
		}

		if !hasIntersection {
			maxArea = cand.area
			continue
		}

		allValid := true
		for x := minX; x <= maxX && allValid; x++ {
			if !isValidTileCached(Point{x, minY}) || (minY != maxY && !isValidTileCached(Point{x, maxY})) {
				allValid = false
			}
		}

		if allValid {
			for y := minY + 1; y < maxY && allValid; y++ {
				if !isValidTileCached(Point{minX, y}) || (minX != maxX && !isValidTileCached(Point{maxX, y})) {
					allValid = false
				}
			}
		}

		if allValid && maxX-minX > 1 && maxY-minY > 1 {
			for x := minX + 1; x < maxX && allValid; x++ {
				for y := minY + 1; y < maxY && allValid; y++ {
					if !isValidTileCached(Point{x, y}) {
						allValid = false
					}
				}
			}
		}

		if allValid {
			maxArea = cand.area
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
