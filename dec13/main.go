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

func foldHorizontal(pattern []string) int {
OUTER:
	for i := 0; i < len(pattern)-1; i++ {
		// fmt.Printf("Checking row %d:\n", i)
		// fold from inside outwards
		for a, b := i, i+1; a >= 0 && b < len(pattern); a, b = a-1, b+1 {

			if pattern[a] != pattern[b] {
				// fmt.Printf("row %d DOES NOT mirror row %d\n", a, b) // move on to next fold point
				// fmt.Printf("  %s\n  %s\n", pattern[a], pattern[b])
				continue OUTER
			} else {
				fmt.Printf("row %d matches row %d\n", a, b)
			}
		}

		// fmt.Printf("All rows have matched for index %d\n", i)
		return i + 1
	}

	return 0
}

func foldVertical(pattern []string) int {

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

	return foldHorizontal(input)
}

func part1(input []string) int {

	patterns := parsePatterns(input)

	sh := 0
	for _, pattern := range patterns {
		fmt.Printf("Looking for folds on pattern:\n%s\n", pattern)
		hv := foldHorizontal(pattern)
		if hv > 0 {
			sh += (hv * 100)
		}

		vv := foldVertical(pattern)
		sh += vv
	}

	return sh
}

func part2(input []string) int {

	return 0
}

func main() {

	input := parseInput()

	fmt.Println("Part 1:", part1(input))
	fmt.Println("Part 2:", part2(input))
}
