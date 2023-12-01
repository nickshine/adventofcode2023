package main

import (
	"fmt"
	"os"
	"strconv"
	"strings"
	"unicode"
)

func parseInput() []string {
	raw, err := os.ReadFile(os.Args[1])
	if err != nil {
		panic(err)
	}

	return strings.Split(strings.Trim(string(raw), "\n"), "\n")
}

func part1(lines []string) int {
	var digits []string

	for _, line := range lines {
		var f, l string
		for i := 0; i < len(line); i++ {
			if unicode.IsDigit(rune(line[i])) {
				f = string(line[i])
				break
			}
		}
		for i := len(line) - 1; i >= 0; i-- {
			if unicode.IsDigit(rune(line[i])) {
				l = string(line[i])
				break
			}

		}

		digits = append(digits, f+l)
	}

	sum := 0

	for _, v := range digits {
		i, _ := strconv.Atoi(v)
		sum += i
	}

	return sum
}

func part2(lines []string) int {

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

	var digits []string

	for _, line := range lines {

		var f, l string
	FIRST:
		for i := 0; i < len(line); i++ {
			// check if digit
			if unicode.IsDigit(rune(line[i])) {
				f = string(line[i])
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

	SECOND:
		for i := len(line) - 1; i >= 0; i-- {
			if unicode.IsDigit(rune(line[i])) {
				l = string(line[i])
				break
			}

			for k, v := range nmap {
				if strings.HasSuffix(line[:i+1], k) {
					l = v
					break SECOND
				}
			}
		}

		digits = append(digits, f+l)
	}

	sum := 0

	for _, v := range digits {
		i, _ := strconv.Atoi(v)
		sum += i
	}

	return sum

}

func main() {

	input := parseInput()

	fmt.Println("Part 1:", part1(input))
	fmt.Println("Part 2:", part2(input))
}
