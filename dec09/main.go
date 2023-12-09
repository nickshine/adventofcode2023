package main

import (
	"fmt"
	"os"
	"strconv"
	"strings"
)

func parseInput() []string {
	raw, _ := os.ReadFile(os.Args[1])
	return strings.Split(strings.Trim(string(raw), "\n"), "\n")
}

func parse(in []string) [][]int {

	var histories [][]int

	for _, line := range in {
		var history []int
		vals := strings.Fields(line)
		for _, v := range vals {
			n, _ := strconv.Atoi(v)
			history = append(history, n)
		}

		histories = append(histories, history)
	}

	return histories
}

// extrapolate returns the prev and next values for the sequence
func extrapolate(sequence []int) (int, int) {
	// base case
	done := true
	for _, v := range sequence {
		if v != 0 {
			done = false
		}
	}
	if done {
		return 0, 0
	}

	var seq []int
	for i := 0; i < len(sequence)-1; i++ {
		diff := sequence[i+1] - sequence[i]
		seq = append(seq, diff)
	}

	prev, next := extrapolate(seq)
	return sequence[0] - prev, sequence[len(sequence)-1] + next
}

func part1(input []string) int {
	histories := parse(input)
	total := 0

	for _, history := range histories {
		_, next := extrapolate(history)
		total += next
	}

	return total
}

func part2(input []string) int {
	histories := parse(input)
	total := 0

	for _, history := range histories {
		prev, _ := extrapolate(history)
		total += prev
	}

	return total
}

func main() {

	input := parseInput()

	fmt.Println("Part 1:", part1(input))
	fmt.Println("Part 2:", part2(input))
}
