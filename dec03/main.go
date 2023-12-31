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

type point struct {
	x, y int
}

func parseMap(input []string) (numbers, symbols map[point]string, bounds point) {

	// a mapping of coordinates to numbers/symbols
	numbers = make(map[point]string)
	symbols = make(map[point]string)

	for y, row := range input {

		// scan for numbers
		x, r := 0, 0
		for r < len(row) {

			// find end of number
			if unicode.IsDigit(rune(row[r])) {
				r++
				continue
			}

			// end of a number
			if x != r {
				numbers[point{x, y}] = row[x:r] // capture the number
			}

			// must be "." or symbol
			if string(row[r]) != "." {
				symbols[point{r, y}] = string(row[r])
			}

			x = r + 1 // move starting point
			r = x
		}
	}

	return numbers, symbols, point{len(input[0]), len(input)}
}

func adjacents(p, bounds point, numLength int) []point {
	var a []point

	if p.x-1 >= 0 {
		a = append(a, point{p.x - 1, p.y}) // left
	}

	if p.x+numLength < bounds.x {
		a = append(a, point{p.x + numLength, p.y}) // right
	}
	// add up/down xrange
	for x := p.x - 1; x <= p.x+numLength; x++ {
		if p.y-1 >= 0 {
			a = append(a, point{x, p.y - 1}) // up
		}
		if p.y+1 < bounds.y {
			a = append(a, point{x, p.y + 1}) // down
		}
	}

	return a
}

func part1(input []string) int {

	numbers, symbols, bounds := parseMap(input)

	var sum int

	for p, v := range numbers {

		num, _ := strconv.Atoi(v)

		for _, a := range adjacents(p, bounds, len(v)) {
			// if adjacent to a symbol
			if _, ok := symbols[a]; ok {
				sum += num
				break
			}
		}
	}

	return sum
}

func part2(input []string) int {
	numbers, symbols, bounds := parseMap(input)

	// map of gears to adjacent numbers
	gears := make(map[point][]int)

	for p, v := range numbers {

		num, _ := strconv.Atoi(v)

		for _, a := range adjacents(p, bounds, len(v)) {
			// if adjacent to a symbol
			if s, ok := symbols[a]; ok && s == "*" {
				gears[a] = append(gears[a], num)
				break
			}
		}
	}

	var sum int

	for _, nums := range gears {

		// discard invalid gears
		if len(nums) != 2 {
			continue
		}

		sum += nums[0] * nums[1]
	}

	return sum
}

func main() {

	input := parseInput()

	fmt.Println("Part 1:", part1(input))
	fmt.Println("Part 2:", part2(input))
}
