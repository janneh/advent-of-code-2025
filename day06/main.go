package main

import (
	"bufio"
	"fmt"
	"os"
	"strconv"
	"strings"
)

func parseWorksheet(lines []string) [][]string {
	if len(lines) == 0 {
		return nil
	}

	// Find max width and pad all lines
	maxWidth := 0
	for _, line := range lines {
		if len(line) > maxWidth {
			maxWidth = len(line)
		}
	}

	// Pad lines to same width
	paddedLines := make([]string, len(lines))
	for i, line := range lines {
		paddedLines[i] = line + strings.Repeat(" ", maxWidth-len(line))
	}

	// Identify problem boundaries
	// A column is a separator if it's all spaces
	problems := [][]string{}
	inProblem := false
	problemStart := 0

	for col := 0; col <= maxWidth; col++ {
		isSeparator := true
		if col < maxWidth {
			for row := 0; row < len(paddedLines); row++ {
				if paddedLines[row][col] != ' ' {
					isSeparator = false
					break
				}
			}
		}

		if isSeparator || col == maxWidth {
			if inProblem {
				// End of problem, extract it
				problem := make([]string, len(paddedLines))
				for row := 0; row < len(paddedLines); row++ {
					problem[row] = paddedLines[row][problemStart:col]
				}
				problems = append(problems, problem)
				inProblem = false
			}
		} else {
			if !inProblem {
				// Start of new problem
				problemStart = col
				inProblem = true
			}
		}
	}

	return problems
}

func solveProblem(problem []string) int {
	// Last line contains the operator
	operatorLine := strings.TrimSpace(problem[len(problem)-1])
	var operator rune
	for _, ch := range operatorLine {
		if ch == '*' || ch == '+' {
			operator = ch
			break
		}
	}

	// Extract numbers from all lines except the last
	numbers := []int{}
	for i := 0; i < len(problem)-1; i++ {
		numStr := strings.TrimSpace(problem[i])
		if numStr != "" {
			num, _ := strconv.Atoi(numStr)
			numbers = append(numbers, num)
		}
	}

	// Evaluate based on operator
	if len(numbers) == 0 {
		return 0
	}

	result := numbers[0]
	for i := 1; i < len(numbers); i++ {
		if operator == '*' {
			result *= numbers[i]
		} else {
			result += numbers[i]
		}
	}

	return result
}

func part1(lines []string) int {
	problems := parseWorksheet(lines)
	grandTotal := 0

	for _, problem := range problems {
		answer := solveProblem(problem)
		grandTotal += answer
	}

	return grandTotal
}

func main() {
	file, err := os.Open("input.txt")
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error opening file: %v\n", err)
		os.Exit(1)
	}
	defer file.Close()

	var lines []string
	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		lines = append(lines, scanner.Text())
	}

	if err := scanner.Err(); err != nil {
		fmt.Fprintf(os.Stderr, "Error reading file: %v\n", err)
		os.Exit(1)
	}

	fmt.Printf("Part 1: %d\n", part1(lines))
}
