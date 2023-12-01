package main

import (
	"fmt"
	"os"
	"strconv"
	"strings"
	"unicode"
)

func parseInput() []string {
	raw, _ := os.ReadFile(os.Args[1])
	return strings.Split(strings.Trim(string(raw), "\n"), "\n")
}

func part1(input []string) int {
	var sum int

	for _, line := range input {
		var f, l string
		for _, c := range line {
			if unicode.IsDigit(c) {
				f = string(c)
				break
			}
		}
		for i := len(line) - 1; i >= 0; i-- {
			if unicode.IsDigit(rune(line[i])) {
				l = string(line[i])
				break
			}

		}

		n, _ := strconv.Atoi(f + l)
		sum += n
	}

	return sum
}

func part2(input []string) int {
	var sum int

	nmap := map[string]string{
		"one":   "1",
		"two":   "2",
		"three": "3",
		"four":  "4",
		"five":  "5",
		"six":   "6",
		"seven": "7",
		"eight": "8",
		"nine":  "9",
	}

	for _, line := range input {

		var f, l string
	FIRST:
		for i, c := range line {
			// check if digit
			if unicode.IsDigit(c) {
				f = string(c)
				break
			}

			// check if letters
			for k, v := range nmap {
				if strings.HasPrefix(line[i:], k) {
					f = v
					break FIRST
				}
			}
		}

	LAST:
		for i := len(line) - 1; i >= 0; i-- {
			if unicode.IsDigit(rune(line[i])) {
				l = string(line[i])
				break
			}

			for k, v := range nmap {
				if strings.HasSuffix(line[:i+1], k) {
					l = v
					break LAST
				}
			}
		}

		n, _ := strconv.Atoi(f + l)
		sum += n

	}

	return sum
}

func main() {

	input := parseInput()

	fmt.Println("Part 1:", part1(input))
	fmt.Println("Part 2:", part2(input))
}
