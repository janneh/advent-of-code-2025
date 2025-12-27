package main

import (
	"bufio"
	"fmt"
	"os"
	"strconv"
)

type Rotation struct {
	direction byte
	distance  int
}

func parseRotations(filename string) ([]Rotation, error) {
	file, err := os.Open(filename)
	if err != nil {
		return nil, err
	}
	defer file.Close()

	var rotations []Rotation
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

		rotations = append(rotations, Rotation{direction, distance})
	}

	if err := scanner.Err(); err != nil {
		return nil, err
	}

	return rotations, nil
}

func part1(rotations []Rotation) int {
	position := 50
	count := 0

	for _, rot := range rotations {
		switch rot.direction {
		case 'L':
			position = ((position-rot.distance)%100 + 100) % 100
		case 'R':
			position = (position + rot.distance) % 100
		}

		if position == 0 {
			count++
		}
	}

	return count
}

func countZeros(start, distance int, goingRight bool) int {
	if goingRight {
		// Going right, we hit 0 every time we complete a cycle of 100
		return (start + distance) / 100
	}
	// Going left
	if start == 0 {
		// Starting at 0, we hit it again at 100, 200, etc.
		return distance / 100
	}
	// Starting at non-zero, we hit 0 at distances: start, start+100, start+200, ...
	if distance >= start {
		return 1 + (distance-start)/100
	}
	return 0
}

func part2(rotations []Rotation) int {
	position := 50
	count := 0

	for _, rot := range rotations {
		switch rot.direction {
		case 'L':
			count += countZeros(position, rot.distance, false)
			position = ((position-rot.distance)%100 + 100) % 100
		case 'R':
			count += countZeros(position, rot.distance, true)
			position = (position + rot.distance) % 100
		}
	}

	return count
}

func main() {
	rotations, err := parseRotations("input.txt")
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error reading input: %v\n", err)
		os.Exit(1)
	}

	fmt.Printf("Part 1: %d\n", part1(rotations))
	fmt.Printf("Part 2: %d\n", part2(rotations))
}
