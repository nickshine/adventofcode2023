package main

import (
	"fmt"
	"os"
	"regexp"
	"strings"
)

func parseInput() []string {
	raw, _ := os.ReadFile(os.Args[1])
	return strings.Split(strings.Trim(string(raw), "\n"), "\n")
}

var re = regexp.MustCompile(`^(\w+) = \((\w+), (\w+)\)$`)

// parseMap returns the map instructions and nodes map.
func parseMap(input []string) (string, map[string][2]string) {

	nodes := make(map[string][2]string)

	instructions := strings.Trim(input[0], "\n")
	for i := 2; i < len(input); i++ {
		parts := re.FindStringSubmatch(input[i])
		nodes[parts[1]] = [2]string{parts[2], parts[3]}
	}

	return instructions, nodes
}

func part1(input []string) int {

	instructions, nodes := parseMap(input)
	step := 0
	cur := "AAA"

	for {
		var next string
		instruction := instructions[step%len(instructions)]
		if instruction == 'L' {
			next = nodes[cur][0]
		} else {
			next = nodes[cur][1]
		}

		if next == "ZZZ" {
			break
		}

		cur = next
		step++
	}

	return step
}

// https://siongui.github.io/2017/06/03/go-find-lcm-by-gcd/
func GCD(a, b int) int {
	for b != 0 {
		t := b
		b = a % b
		a = t
	}
	return a
}

// find Least Common Multiple (LCM) via GCD
func LCM(a, b int, integers ...int) int {
	result := a * b / GCD(a, b)

	for i := 0; i < len(integers); i++ {
		result = LCM(result, integers[i])
	}

	return result
}

func part2(input []string) int {
	instructions, nodes := parseMap(input)

	count := 0

	var cur []string
	for k := range nodes {
		if k[2] == 'A' {
			cur = append(cur, k)
		}
	}

	cycles := map[int]int{}

	for {
		for i, v := range cur {
			if v[2] == 'Z' {
				if _, ok := cycles[i]; !ok {
					cycles[i] = count
				}
			}
		}

		if len(cycles) >= len(cur) {
			break
		}

		var next []string

		if instructions[count%len(instructions)] == 'L' {
			for _, v := range cur {
				next = append(next, nodes[v][0])
			}
		} else {
			for _, v := range cur {
				next = append(next, nodes[v][1])
			}
		}

		count++

		cur = next
	}

	fmt.Println("Cycles:", cycles)
	return LCM(cycles[0], cycles[1], cycles[2], cycles[3], cycles[4], cycles[5])
}

func main() {

	input := parseInput()

	fmt.Println("Part 1:", part1(input))
	fmt.Println("Part 2:", part2(input))
}
