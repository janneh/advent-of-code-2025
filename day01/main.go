package main

import (
	"bufio"
	"fmt"
	"os"
	"strconv"
)

func main() {
	file, err := os.Open("input.txt")
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error opening file: %v\n", err)
		os.Exit(1)
	}
	defer file.Close()

	position := 50
	count := 0

	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		line := scanner.Text()
		if len(line) < 2 {
			continue
		}

		direction := line[0]
		distance, err := strconv.Atoi(line[1:])
		if err != nil {
			fmt.Fprintf(os.Stderr, "Error parsing distance: %v\n", err)
			continue
		}

		switch direction {
		case 'L':
			// For left rotations, subtract and handle wraparound
			// Use modulo arithmetic with adjustment for negative numbers
			position = ((position-distance)%100 + 100) % 100
		case 'R':
			// For right rotations, add and handle wraparound
			position = (position + distance) % 100
		}

		if position == 0 {
			count++
		}
	}

	if err := scanner.Err(); err != nil {
		fmt.Fprintf(os.Stderr, "Error reading file: %v\n", err)
		os.Exit(1)
	}

	fmt.Printf("The password is: %d\n", count)
}
