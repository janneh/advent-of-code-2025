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

	// Part 1: Count all paths from "you" to "out"
	visited := make(map[string]bool)
	paths := countPaths(graph, "you", "out", visited)
	fmt.Printf("Part 1: %d\n", paths)

	// Part 2: Count paths from "svr" to "out" that visit both "dac" and "fft"
	visited2 := make(map[string]bool)
	memo := make(map[State]int)
	pathsWithRequired := countPathsWithRequiredMemo(graph, "svr", "out", visited2, false, false, memo)
	fmt.Printf("Part 2: %d\n", pathsWithRequired)
}

type State struct {
	node    string
	seenDAC bool
	seenFFT bool
}

func countPaths(graph map[string][]string, current, target string, visited map[string]bool) int {
	if current == target {
		return 1
	}

	if visited[current] {
		return 0
	}

	visited[current] = true
	defer func() {
		visited[current] = false
	}()

	totalPaths := 0
	for _, neighbor := range graph[current] {
		totalPaths += countPaths(graph, neighbor, target, visited)
	}

	return totalPaths
}

func countPathsWithRequiredMemo(graph map[string][]string, current, target string,
	visited map[string]bool, seenDAC, seenFFT bool, memo map[State]int) int {

	if current == "dac" {
		seenDAC = true
	}

	if current == "fft" {
		seenFFT = true
	}

	if current == target {
		if seenDAC && seenFFT {
			return 1
		}
		return 0
	}

	if visited[current] {
		return 0
	}

	state := State{current, seenDAC, seenFFT}
	if val, ok := memo[state]; ok {
		return val
	}

	visited[current] = true
	defer func() {
		visited[current] = false
	}()

	totalPaths := 0
	for _, neighbor := range graph[current] {
		totalPaths += countPathsWithRequiredMemo(graph, neighbor, target, visited, seenDAC, seenFFT, memo)
	}

	memo[state] = totalPaths

	return totalPaths
}
