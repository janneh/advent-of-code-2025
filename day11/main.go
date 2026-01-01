package main

import (
	"bufio"
	"fmt"
	"os"
	"strings"
)

func main() {
	file, err := os.Open("input.txt")
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error: %v\n", err)
		os.Exit(1)
	}
	defer file.Close()

	// Build adjacency list (directed graph)
	graph := make(map[string][]string)
	scanner := bufio.NewScanner(file)

	for scanner.Scan() {
		line := strings.TrimSpace(scanner.Text())
		if line == "" {
			continue
		}

		parts := strings.Split(line, ": ")
		if len(parts) != 2 {
			continue
		}

		from := parts[0]
		outputs := strings.Fields(parts[1])
		graph[from] = outputs
	}

	if err := scanner.Err(); err != nil {
		fmt.Fprintf(os.Stderr, "Error: %v\n", err)
		os.Exit(1)
	}

	// Count all paths from "you" to "out"
	visited := make(map[string]bool)
	paths := countPaths(graph, "you", "out", visited)

	fmt.Printf("Part 1: %d\n", paths)
}

// DFS to count all paths from current to target
func countPaths(graph map[string][]string, current, target string, visited map[string]bool) int {
	// Base case: reached the target
	if current == target {
		return 1
	}

	// Avoid cycles: check if we've visited this node in the current path
	if visited[current] {
		return 0
	}

	// Mark as visited for this path
	visited[current] = true
	defer func() {
		// Unmark when backtracking to allow other paths through this node
		visited[current] = false
	}()

	// Explore all neighbors and sum up paths
	totalPaths := 0
	for _, neighbor := range graph[current] {
		totalPaths += countPaths(graph, neighbor, target, visited)
	}

	return totalPaths
}
