package main

import (
	"bufio"
	"fmt"
	"os"
	"strconv"
	"strings"
)

type Shape [][]bool

type Region struct {
	width, height int
	counts        []int
}

func main() {
	shapes, regions := parseInput("input.txt")

	allVariants := make([][]Shape, len(shapes))
	for i, shape := range shapes {
		allVariants[i] = generateVariants(shape)
	}

	validRegions := 0
	for _, region := range regions {
		if canFitPresents(region, allVariants) {
			validRegions++
		}
	}

	fmt.Printf("Part 1: %d\n", validRegions)
}

func parseInput(filename string) ([]Shape, []Region) {
	file, err := os.Open(filename)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error: %v\n", err)
		os.Exit(1)
	}
	defer file.Close()

	scanner := bufio.NewScanner(file)
	shapes := []Shape{}
	regions := []Region{}

	for scanner.Scan() {
		line := strings.TrimSpace(scanner.Text())
		if line == "" {
			continue
		}

		if strings.HasSuffix(line, ":") {
			shape := [][]bool{}
			for scanner.Scan() {
				line := strings.TrimSpace(scanner.Text())
				if line == "" {
					break
				}
				if strings.Contains(line, "x") && strings.Contains(line, ":") {
					regions = append(regions, parseRegion(line))
					goto parseRegions
				}
				row := []bool{}
				for _, c := range line {
					row = append(row, c == '#')
				}
				shape = append(shape, row)
			}
			if len(shape) > 0 {
				shapes = append(shapes, shape)
			}
		} else if strings.Contains(line, "x") && strings.Contains(line, ":") {
			regions = append(regions, parseRegion(line))
		}
	}

parseRegions:
	for scanner.Scan() {
		line := strings.TrimSpace(scanner.Text())
		if line == "" {
			continue
		}
		if strings.Contains(line, "x") && strings.Contains(line, ":") {
			regions = append(regions, parseRegion(line))
		}
	}

	return shapes, regions
}

func parseRegion(line string) Region {
	parts := strings.Split(line, ": ")
	if len(parts) != 2 {
		return Region{}
	}

	dims := strings.Split(parts[0], "x")
	width, _ := strconv.Atoi(dims[0])
	height, _ := strconv.Atoi(dims[1])

	countStrs := strings.Fields(parts[1])
	counts := []int{}
	for _, cs := range countStrs {
		c, _ := strconv.Atoi(cs)
		counts = append(counts, c)
	}

	return Region{width, height, counts}
}

func generateVariants(shape Shape) []Shape {
	variants := []Shape{}
	current := copyShape(shape)

	for r := 0; r < 4; r++ {
		variants = append(variants, copyShape(current))
		flipped := flipShape(current)
		variants = append(variants, flipped)
		current = rotateShape(current)
	}

	return removeDuplicateShapes(variants)
}

func rotateShape(shape Shape) Shape {
	if len(shape) == 0 || len(shape[0]) == 0 {
		return shape
	}
	rows := len(shape)
	cols := len(shape[0])
	rotated := make(Shape, cols)
	for i := range rotated {
		rotated[i] = make([]bool, rows)
	}

	for r := range rows {
		for c := 0; c < cols; c++ {
			rotated[c][rows-1-r] = shape[r][c]
		}
	}
	return rotated
}

func flipShape(shape Shape) Shape {
	flipped := make(Shape, len(shape))
	for i := range shape {
		flipped[i] = make([]bool, len(shape[i]))
		for j := range shape[i] {
			flipped[i][len(shape[i])-1-j] = shape[i][j]
		}
	}
	return flipped
}

func copyShape(shape Shape) Shape {
	copied := make(Shape, len(shape))
	for i := range shape {
		copied[i] = make([]bool, len(shape[i]))
		copy(copied[i], shape[i])
	}
	return copied
}

func removeDuplicateShapes(shapes []Shape) []Shape {
	unique := []Shape{}
	for _, s := range shapes {
		isDup := false
		for _, u := range unique {
			if shapesEqual(s, u) {
				isDup = true
				break
			}
		}
		if !isDup {
			unique = append(unique, s)
		}
	}
	return unique
}

func shapesEqual(a, b Shape) bool {
	if len(a) != len(b) {
		return false
	}
	for i := range a {
		if len(a[i]) != len(b[i]) {
			return false
		}
		for j := range a[i] {
			if a[i][j] != b[i][j] {
				return false
			}
		}
	}
	return true
}

func canFitPresents(region Region, allVariants [][]Shape) bool {
	grid := make([][]int, region.height)
	for i := range grid {
		grid[i] = make([]int, region.width)
	}

	type Present struct {
		shapeIdx int
		size     int
	}
	presents := []Present{}
	for shapeIdx, count := range region.counts {
		if count > 0 && shapeIdx < len(allVariants) && len(allVariants[shapeIdx]) > 0 {
			size := 0
			for _, row := range allVariants[shapeIdx][0] {
				for _, cell := range row {
					if cell {
						size++
					}
				}
			}
			for i := 0; i < count; i++ {
				presents = append(presents, Present{shapeIdx, size})
			}
		}
	}

	for i := 0; i < len(presents)-1; i++ {
		for j := i + 1; j < len(presents); j++ {
			if presents[j].size > presents[i].size {
				presents[i], presents[j] = presents[j], presents[i]
			}
		}
	}

	shapeIndices := make([]int, len(presents))
	for i, p := range presents {
		shapeIndices[i] = p.shapeIdx
	}

	return tryPlace(grid, shapeIndices, allVariants, 0)
}

func tryPlace(grid [][]int, presents []int, allVariants [][]Shape, presentIdx int) bool {
	if presentIdx >= len(presents) {
		return true
	}

	remainingArea := 0
	for i := range grid {
		for j := range grid[i] {
			if grid[i][j] == 0 {
				remainingArea++
			}
		}
	}

	neededArea := 0
	for i := presentIdx; i < len(presents); i++ {
		shapeIdx := presents[i]
		if shapeIdx < len(allVariants) && len(allVariants[shapeIdx]) > 0 {
			for _, row := range allVariants[shapeIdx][0] {
				for _, cell := range row {
					if cell {
						neededArea++
					}
				}
			}
		}
	}

	if neededArea > remainingArea {
		return false
	}

	shapeIdx := presents[presentIdx]
	variants := allVariants[shapeIdx]

	for _, variant := range variants {
		if len(variant) == 0 || len(variant[0]) == 0 {
			continue
		}
		for r := 0; r <= len(grid)-len(variant); r++ {
			for c := 0; c <= len(grid[0])-len(variant[0]); c++ {
				if canPlace(grid, variant, r, c) {
					place(grid, variant, r, c, presentIdx+1)
					if tryPlace(grid, presents, allVariants, presentIdx+1) {
						return true
					}
					unplace(grid, variant, r, c)
				}
			}
		}
	}

	return false
}

func canPlace(grid [][]int, shape Shape, row, col int) bool {
	for r := range shape {
		for c := range shape[r] {
			if shape[r][c] && grid[row+r][col+c] != 0 {
				return false
			}
		}
	}
	return true
}

func place(grid [][]int, shape Shape, row, col int, id int) {
	for r := range shape {
		for c := range shape[r] {
			if shape[r][c] {
				grid[row+r][col+c] = id
			}
		}
	}
}

func unplace(grid [][]int, shape Shape, row, col int) {
	for r := range shape {
		for c := range shape[r] {
			if shape[r][c] {
				grid[row+r][col+c] = 0
			}
		}
	}
}
