package main

import (
	"fmt"
	"os"
	"regexp"
	"strconv"
	"strings"
)

func parseInput() []string {
	raw, _ := os.ReadFile(os.Args[1])
	return strings.Split(strings.Trim(string(raw), "\n"), "\n")
}

type card struct {
	winning map[int]struct{}
	have    []int
	id      int
}

var re = regexp.MustCompile(`^Card +(\d+): +(.+) +\| +(.+)$`)

func parseCard(in string) card {
	c := card{
		winning: make(map[int]struct{}),
		have:    nil,
	}

	parts := re.FindStringSubmatch(in)

	c.id, _ = strconv.Atoi(parts[1])
	winning := strings.Fields(parts[2])
	have := strings.Fields(parts[3])

	for _, n := range winning {
		rn, _ := strconv.Atoi(n)
		c.winning[rn] = struct{}{}
	}

	for _, n := range have {
		rn, _ := strconv.Atoi(n)
		c.have = append(c.have, rn)
	}

	return c
}

func part1(input []string) int {
	sum := 0

	for _, line := range input {
		points, multiple := 0, 0
		c := parseCard(line)

		for _, n := range c.have {
			if _, ok := c.winning[n]; ok {
				points = 1 << multiple
				multiple++
			}
		}
		sum += points
	}

	return sum
}

func part2(input []string) int {
	sum := 0

	// map of card ids to instances
	instances := make(map[int]int)

	for _, line := range input {
		matches := 0
		c := parseCard(line)
		instances[c.id]++

		for _, n := range c.have {
			if _, ok := c.winning[n]; ok {
				matches++
			}
		}

		for i := 0; i < matches; i++ {
			instances[c.id+i+1] += instances[c.id]
		}
	}

	for _, v := range instances {
		sum += v
	}

	return sum
}

func main() {

	input := parseInput()

	fmt.Println("Part 1:", part1(input))
	fmt.Println("Part 2:", part2(input))
}
