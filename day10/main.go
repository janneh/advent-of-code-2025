package main

import (
	"bufio"
	"fmt"
	"os"
	"strconv"
	"strings"
)

type Machine struct {
	target   []int
	buttons  [][]int
	joltages []int
}

func parseLine(line string) Machine {
	start := strings.Index(line, "[")
	end := strings.Index(line, "]")
	pattern := line[start+1 : end]

	target := make([]int, len(pattern))
	for i, ch := range pattern {
		if ch == '#' {
			target[i] = 1
		}
	}

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

	var joltages []int
	joltStart := strings.Index(line, "{")
	joltEnd := strings.Index(line, "}")
	if joltStart != -1 && joltEnd != -1 {
		joltStr := line[joltStart+1 : joltEnd]
		parts := strings.Split(joltStr, ",")
		for _, p := range parts {
			num, err := strconv.Atoi(strings.TrimSpace(p))
			if err == nil {
				joltages = append(joltages, num)
			}
		}
	}

	return Machine{target: target, buttons: buttons, joltages: joltages}
}

func solvePart1(machine Machine) int {
	numLights := len(machine.target)
	numButtons := len(machine.buttons)

	if numButtons > 25 {
		return -1
	}

	minPresses := numButtons + 1

	for mask := 0; mask < (1 << numButtons); mask++ {
		lights := make([]int, numLights)

		for buttonIdx := range numButtons {
			if (mask & (1 << buttonIdx)) != 0 {
				for _, lightIdx := range machine.buttons[buttonIdx] {
					if lightIdx < numLights {
						lights[lightIdx] ^= 1
					}
				}
			}
		}

		match := true
		for i := 0; i < numLights; i++ {
			if lights[i] != machine.target[i] {
				match = false
				break
			}
		}

		if match {
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
		return -1
	}
	return minPresses
}

func solvePart2(machine Machine) int {
	numCounters := len(machine.joltages)
	numButtons := len(machine.buttons)

	if numCounters == 0 || numButtons == 0 {
		return 0
	}

	A := make([][]float64, numCounters)
	for i := range A {
		A[i] = make([]float64, numButtons)
	}

	for buttonIdx, button := range machine.buttons {
		for _, counterIdx := range button {
			if counterIdx < numCounters {
				A[counterIdx][buttonIdx] = 1
			}
		}
	}

	b := make([]float64, numCounters)
	for i, val := range machine.joltages {
		b[i] = float64(val)
	}

	return recursiveDivideConquer(A, b, numButtons)
}

func recursiveDivideConquer(A [][]float64, b []float64, numButtons int) int {
	target := make([]int, len(b))
	for i, v := range b {
		target[i] = int(v)
	}
	memo := make(map[string]int)
	result := solveRecursive(A, target, numButtons, memo)
	return result
}

func solveRecursive(A [][]float64, target []int, numButtons int, memo map[string]int) int {
	allZero := true
	for _, v := range target {
		if v != 0 {
			allZero = false
			break
		}
	}
	if allZero {
		return 0
	}

	key := arrayToString(target)
	if val, ok := memo[key]; ok {
		return val
	}

	minPresses := int(1e9)
	numCounters := len(target)

	parityTarget := make([]int, numCounters)
	for i, v := range target {
		if v%2 == 1 {
			parityTarget[i] = 1
		}
	}

	if numButtons > 25 {
		memo[key] = -1
		return -1
	}

	maxMask := 1 << numButtons
	for mask := range maxMask {
		parity := make([]int, numCounters)
		effect := make([]int, numCounters)
		pressCount := 0

		for btn := range numButtons {
			if (mask & (1 << btn)) != 0 {
				pressCount++
				for ctr := range numCounters {
					if A[ctr][btn] > 0.5 {
						parity[ctr] ^= 1 // XOR for parity
						effect[ctr]++    // Actual effect
					}
				}
			}
		}

		parityMatch := true
		for i := range parityTarget {
			if parity[i] != parityTarget[i] {
				parityMatch = false
				break
			}
		}

		if !parityMatch {
			continue
		}

		remaining := make([]int, numCounters)
		allEven := true
		anyNegative := false

		for i := range target {
			remaining[i] = target[i] - effect[i]
			if remaining[i] < 0 {
				anyNegative = true
				break
			}
			if remaining[i]%2 != 0 {
				allEven = false
				break
			}
		}

		if !allEven || anyNegative {
			continue
		}

		halfRemaining := make([]int, numCounters)
		for i := range remaining {
			halfRemaining[i] = remaining[i] / 2
		}

		recursivePresses := solveRecursive(A, halfRemaining, numButtons, memo)
		if recursivePresses != -1 {
			// Formula: k + 2 * f((b-effect)/2)
			totalPresses := pressCount + 2*recursivePresses
			if totalPresses < minPresses {
				minPresses = totalPresses
			}
		}
	}

	result := -1
	if minPresses < int(1e9) {
		result = minPresses
	}

	memo[key] = result
	return result
}

func arrayToString(arr []int) string {
	var sb strings.Builder
	for i, v := range arr {
		if i > 0 {
			sb.WriteString(",")
		}
		fmt.Fprintf(&sb, "%d", v)
	}
	return sb.String()
}

func main() {
	file, err := os.Open("input.txt")
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error: %v\n", err)
		os.Exit(1)
	}
	defer file.Close()

	totalPart1 := 0
	totalPart2 := 0
	scanner := bufio.NewScanner(file)
	machineNum := 0
	skipped := 0

	for scanner.Scan() {
		line := strings.TrimSpace(scanner.Text())
		if line == "" {
			continue
		}

		machine := parseLine(line)
		machineNum++

		presses1 := solvePart1(machine)
		if presses1 != -1 {
			totalPart1 += presses1
		}

		presses2 := solvePart2(machine)
		if presses2 == -1 {
			fmt.Printf("Machine %d: SKIPPED (no solution found) - %d buttons, %d counters\n",
				machineNum, len(machine.buttons), len(machine.joltages))
			skipped++
		} else {
			totalPart2 += presses2
		}
	}

	fmt.Printf("\nSkipped %d machines out of %d\n", skipped, machineNum)

	if err := scanner.Err(); err != nil {
		fmt.Fprintf(os.Stderr, "Error: %v\n", err)
		os.Exit(1)
	}

	fmt.Printf("Part 1: %d\n", totalPart1)
	fmt.Printf("Part 2: %d\n", totalPart2)
}
