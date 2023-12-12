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

type record struct {
	v     string
	sizes []int
}

func parseRecords(in []string, factor int) []record {
	var records []record

	for _, line := range in {
		parts := strings.Split(line, " ")
		rawRecord := parts[0]
		rawSizes := strings.Split(parts[1], ",")
		var sizes []int
		for _, s := range rawSizes {
			n, _ := strconv.Atoi(s)
			sizes = append(sizes, n)
		}

		var unfoldedRecs []string
		var unfoldedSizes []int
		for i := 0; i < factor; i++ {
			unfoldedRecs = append(unfoldedRecs, rawRecord)
			unfoldedSizes = append(unfoldedSizes, sizes...)
		}
		rec := strings.Join(unfoldedRecs, "?")
		records = append(records, record{rec, unfoldedSizes})
	}

	return records
}

var re = regexp.MustCompile(`#+`)

func calcSizes(record string) []int {
	var sizes []int
	parts := re.FindAllString(record, -1)

	for _, p := range parts {
		sizes = append(sizes, len(p))
	}
	return sizes
}

func isEqual(a, b []int) bool {
	if len(a) != len(b) {
		return false
	}

	for i, v := range a {
		if v != b[i] {
			return false
		}
	}
	return true
}

func combinations(record string, unknowns, sizes []int) int {
	// base case
	if len(unknowns) == 0 {
		comboSizes := calcSizes(record)
		if isEqual(comboSizes, sizes) {
			return 1
		}
		return 0
	}

	s := 0
	l := len(unknowns) - 1
	i := unknowns[l]
	unknowns = unknowns[:l]

	b := []byte(record)
	b[i] = '.'
	s += combinations(string(b), unknowns, sizes)
	b[i] = '#'
	s += combinations(string(b), unknowns, sizes)
	b[i] = '?'
	unknowns = append(unknowns, i)
	return s
}

func findCombos(r record) int {
	var unknowns []int
	for i, c := range r.v {
		if c == '?' {
			unknowns = append(unknowns, i)
		}
	}

	return combinations(r.v, unknowns, r.sizes)

}

func part1(input []string) int {

	records := parseRecords(input, 1)
	for _, r := range records {
		fmt.Printf("Records: %s, sizes: %#v\n", r.v, r.sizes)
	}

	total := 0
	for _, rec := range records {
		total += findCombos(rec)
	}

	return total

}

func part2(input []string) int {
	records := parseRecords(input, 5)
	for _, r := range records {
		fmt.Printf("Records: %s, sizes: %#v\n", r.v, r.sizes)
	}

	//TODO - dynamic programming

	return 0
}

func main() {

	input := parseInput()

	fmt.Println("Part 1:", part1(input))
	fmt.Println("Part 2:", part2(input))
}
