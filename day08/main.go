package main

import (
	"bufio"
	"fmt"
	"math"
	"os"
	"sort"
	"strconv"
	"strings"
)

type Point struct {
	x, y, z int
}

type Edge struct {
	i, j int
	dist float64
}

// Union-Find data structure for tracking connected components
type UnionFind struct {
	parent []int
	size   []int
}

func NewUnionFind(n int) *UnionFind {
	uf := &UnionFind{
		parent: make([]int, n),
		size:   make([]int, n),
	}
	for i := range n {
		uf.parent[i] = i
		uf.size[i] = 1
	}
	return uf
}

func (uf *UnionFind) Find(x int) int {
	if uf.parent[x] != x {
		uf.parent[x] = uf.Find(uf.parent[x]) // path compression
	}
	return uf.parent[x]
}

func (uf *UnionFind) Union(x, y int) bool {
	rootX := uf.Find(x)
	rootY := uf.Find(y)

	if rootX == rootY {
		return false // already in same component
	}

	// Union by size
	if uf.size[rootX] < uf.size[rootY] {
		uf.parent[rootX] = rootY
		uf.size[rootY] += uf.size[rootX]
	} else {
		uf.parent[rootY] = rootX
		uf.size[rootX] += uf.size[rootY]
	}

	return true
}

func (uf *UnionFind) GetComponentSizes() []int {
	sizeMap := make(map[int]int)
	for i := range len(uf.parent) {
		root := uf.Find(i)
		sizeMap[root] = uf.size[root]
	}

	sizes := []int{}
	for _, size := range sizeMap {
		sizes = append(sizes, size)
	}
	return sizes
}

func distance(p1, p2 Point) float64 {
	dx := float64(p1.x - p2.x)
	dy := float64(p1.y - p2.y)
	dz := float64(p1.z - p2.z)
	return math.Sqrt(dx*dx + dy*dy + dz*dz)
}

func parsePoint(line string) (Point, error) {
	parts := strings.Split(line, ",")
	if len(parts) != 3 {
		return Point{}, fmt.Errorf("invalid point format")
	}

	x, err := strconv.Atoi(strings.TrimSpace(parts[0]))
	if err != nil {
		return Point{}, err
	}
	y, err := strconv.Atoi(strings.TrimSpace(parts[1]))
	if err != nil {
		return Point{}, err
	}
	z, err := strconv.Atoi(strings.TrimSpace(parts[2]))
	if err != nil {
		return Point{}, err
	}

	return Point{x, y, z}, nil
}

func part1(points []Point, numConnections int) int {
	n := len(points)

	// Calculate all pairwise distances
	var edges []Edge
	for i := range n {
		for j := i + 1; j < n; j++ {
			dist := distance(points[i], points[j])
			edges = append(edges, Edge{i, j, dist})
		}
	}

	// Sort edges by distance (Kruskal's algorithm)
	sort.Slice(edges, func(i, j int) bool {
		return edges[i].dist < edges[j].dist
	})

	// Use Union-Find to connect closest pairs
	// Count attempts, not just successful connections
	uf := NewUnionFind(n)

	for i := 0; i < numConnections && i < len(edges); i++ {
		uf.Union(edges[i].i, edges[i].j)
	}

	// Get component sizes
	sizes := uf.GetComponentSizes()

	// Sort sizes in descending order
	sort.Slice(sizes, func(i, j int) bool {
		return sizes[i] > sizes[j]
	})

	// Multiply the three largest
	if len(sizes) >= 3 {
		return sizes[0] * sizes[1] * sizes[2]
	}

	return 0
}

func part2(points []Point) int {
	n := len(points)

	// Calculate all pairwise distances
	var edges []Edge
	for i := range n {
		for j := i + 1; j < n; j++ {
			dist := distance(points[i], points[j])
			edges = append(edges, Edge{i, j, dist})
		}
	}

	// Sort edges by distance
	sort.Slice(edges, func(i, j int) bool {
		return edges[i].dist < edges[j].dist
	})

	// Use Union-Find to connect boxes until all are in one circuit
	uf := NewUnionFind(n)
	successfulConnections := 0
	var lastEdge Edge

	for _, edge := range edges {
		if uf.Union(edge.i, edge.j) {
			successfulConnections++
			lastEdge = edge

			// All boxes in one circuit when we've made n-1 successful connections
			if successfulConnections == n-1 {
				break
			}
		}
	}

	// Return product of X coordinates of last two boxes connected
	return points[lastEdge.i].x * points[lastEdge.j].x
}

func main() {
	file, err := os.Open("input.txt")
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error opening file: %v\n", err)
		os.Exit(1)
	}
	defer file.Close()

	var points []Point
	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		line := strings.TrimSpace(scanner.Text())
		if line == "" {
			continue
		}

		point, err := parsePoint(line)
		if err != nil {
			fmt.Fprintf(os.Stderr, "Error parsing point: %v\n", err)
			continue
		}
		points = append(points, point)
	}

	if err := scanner.Err(); err != nil {
		fmt.Fprintf(os.Stderr, "Error reading file: %v\n", err)
		os.Exit(1)
	}

	fmt.Printf("Part 1: %d\n", part1(points, 1000))
	fmt.Printf("Part 2: %d\n", part2(points))
}
