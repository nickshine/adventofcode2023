package main

import (
	"fmt"
	"os"
	"strings"
)

func parseInput() []string {
	raw, _ := os.ReadFile(os.Args[1])
	return strings.Split(strings.Trim(string(raw), "\n"), "\n\n")
}

func parsePatterns(strPatterns []string) [][]string {

	var patterns [][]string

	for _, p := range strPatterns {
		patterns = append(patterns, strings.Split(p, "\n"))
	}

	return patterns
}

func rotate(pattern []string) []string {

	transposed := make([][]byte, len(pattern[0]))

	for x := 0; x < len(pattern[0]); x++ {
		transposed[x] = make([]byte, len(pattern))
		for y := 0; y < len(pattern); y++ {
			transposed[x][y] = pattern[len(pattern)-y-1][x]
		}
	}

	input := make([]string, len(transposed))
	for i, row := range transposed {
		input[i] = string(row)
	}

	return input
}

func fold(pattern []string, smudges int) int {
	for i := 0; i < len(pattern)-1; i++ {
		// fold from inside outwards
		diff := 0
		for a, b := i, i+1; a >= 0 && b < len(pattern); a, b = a-1, b+1 {

			// check each character
			for x := 0; x < len(pattern[a]); x++ {
				if pattern[a][x] != pattern[b][x] {
					diff++
				}
			}
		}

		if diff == smudges {
			return i + 1
		}
	}

	return 0
}

func solve(input []string, smudges int) int {

	patterns := parsePatterns(input)

	s := 0
	for _, pattern := range patterns {
		s += fold(pattern, smudges) * 100
		s += fold(rotate(pattern), smudges)
	}

	return s
}

func main() {

	input := parseInput()

	fmt.Println("Part 1:", solve(input, 0))
	fmt.Println("Part 2:", solve(input, 1))
}
