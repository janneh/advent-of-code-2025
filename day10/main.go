package main

import (
	"bufio"
	"fmt"
	"os"
	"strconv"
	"strings"
)

type Machine struct {
	target  []int
	buttons [][]int
}

func parseLine(line string) Machine {
	// Extract target pattern [.##.]
	start := strings.Index(line, "[")
	end := strings.Index(line, "]")
	pattern := line[start+1 : end]

	target := make([]int, len(pattern))
	for i, ch := range pattern {
		if ch == '#' {
			target[i] = 1
		}
	}

	// Extract buttons (0,1,3) (2,3) etc.
	rest := line[end+1:]
	var buttons [][]int

	for {
		start := strings.Index(rest, "(")
		if start == -1 {
			break
		}
		end := strings.Index(rest, ")")
		if end == -1 {
			break
		}

		buttonStr := rest[start+1 : end]
		parts := strings.Split(buttonStr, ",")
		button := make([]int, 0)
		for _, p := range parts {
			num, err := strconv.Atoi(strings.TrimSpace(p))
			if err == nil {
				button = append(button, num)
			}
		}
		buttons = append(buttons, button)
		rest = rest[end+1:]
	}

	return Machine{target: target, buttons: buttons}
}

// Brute force search for minimum button presses
func solveMachine(machine Machine) int {
	numLights := len(machine.target)
	numButtons := len(machine.buttons)

	// For reasonable number of buttons, try all combinations
	if numButtons > 25 {
		// Use heuristic for very large problems
		return solveMachineHeuristic(machine)
	}

	minPresses := numButtons + 1

	// Try all 2^n combinations of button presses
	for mask := 0; mask < (1 << numButtons); mask++ {
		// Simulate pressing buttons according to mask
		lights := make([]int, numLights)

		for buttonIdx := 0; buttonIdx < numButtons; buttonIdx++ {
			if (mask & (1 << buttonIdx)) != 0 {
				// Press this button
				for _, lightIdx := range machine.buttons[buttonIdx] {
					if lightIdx < numLights {
						lights[lightIdx] ^= 1
					}
				}
			}
		}

		// Check if we achieved target
		match := true
		for i := 0; i < numLights; i++ {
			if lights[i] != machine.target[i] {
				match = false
				break
			}
		}

		if match {
			// Count number of button presses
			presses := 0
			for buttonIdx := 0; buttonIdx < numButtons; buttonIdx++ {
				if (mask & (1 << buttonIdx)) != 0 {
					presses++
				}
			}
			if presses < minPresses {
				minPresses = presses
			}
		}
	}

	if minPresses == numButtons+1 {
		return -1 // No solution found
	}
	return minPresses
}

// Heuristic solver for larger problems using Gaussian elimination
func solveMachineHeuristic(machine Machine) int {
	numLights := len(machine.target)
	numButtons := len(machine.buttons)

	// Build augmented matrix [A | b]
	matrix := make([][]int, numLights)
	for i := range matrix {
		matrix[i] = make([]int, numButtons+1)
		matrix[i][numButtons] = machine.target[i]
	}

	// Fill matrix A
	for buttonIdx, button := range machine.buttons {
		for _, lightIdx := range button {
			if lightIdx < numLights {
				matrix[lightIdx][buttonIdx] = 1
			}
		}
	}

	// Gaussian elimination in GF(2)
	pivotRow := 0
	for col := 0; col < numButtons && pivotRow < numLights; col++ {
		foundPivot := false
		for row := pivotRow; row < numLights; row++ {
			if matrix[row][col] == 1 {
				matrix[pivotRow], matrix[row] = matrix[row], matrix[pivotRow]
				foundPivot = true
				break
			}
		}

		if !foundPivot {
			continue
		}

		for row := 0; row < numLights; row++ {
			if row != pivotRow && matrix[row][col] == 1 {
				for c := 0; c <= numButtons; c++ {
					matrix[row][c] ^= matrix[pivotRow][c]
				}
			}
		}
		pivotRow++
	}

	// Check for inconsistency
	for row := 0; row < numLights; row++ {
		allZero := true
		for col := 0; col < numButtons; col++ {
			if matrix[row][col] != 0 {
				allZero = false
				break
			}
		}
		if allZero && matrix[row][numButtons] == 1 {
			return -1
		}
	}

	// Back-substitution
	solution := make([]int, numButtons)
	pivotCols := make(map[int]int)
	for row := 0; row < numLights; row++ {
		for col := 0; col < numButtons; col++ {
			if matrix[row][col] == 1 {
				pivotCols[col] = row
				break
			}
		}
	}

	for col := numButtons - 1; col >= 0; col-- {
		if row, isPivot := pivotCols[col]; isPivot {
			val := matrix[row][numButtons]
			for c := col + 1; c < numButtons; c++ {
				val ^= (matrix[row][c] * solution[c])
			}
			solution[col] = val
		}
	}

	count := 0
	for _, v := range solution {
		count += v
	}

	return count
}

func main() {
	file, err := os.Open("input.txt")
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error opening file: %v\n", err)
		os.Exit(1)
	}
	defer file.Close()

	totalPresses := 0
	scanner := bufio.NewScanner(file)
	lineNum := 0

	for scanner.Scan() {
		lineNum++
		line := strings.TrimSpace(scanner.Text())
		if line == "" {
			continue
		}

		machine := parseLine(line)
		presses := solveMachine(machine)

		if presses == -1 {
			fmt.Fprintf(os.Stderr, "No solution for machine %d\n", lineNum)
			continue
		}

		totalPresses += presses
	}

	if err := scanner.Err(); err != nil {
		fmt.Fprintf(os.Stderr, "Error reading file: %v\n", err)
		os.Exit(1)
	}

	fmt.Printf("Part 1: %d\n", totalPresses)
}
