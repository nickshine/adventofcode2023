package main

import (
	"fmt"
	"os"
	"regexp"
	"strconv"
	"strings"
)

func parseInput() ([]string, []string) {
	raw, _ := os.ReadFile(os.Args[1])
	parts := strings.Split(strings.Trim(string(raw), "\n"), "\n\n")

	workflows := strings.Split(strings.Trim(parts[0], "\n"), "\n")
	ratings := strings.Split(strings.Trim(parts[1], "\n"), "\n")

	return workflows, ratings

}

var (
	workflowRE  = regexp.MustCompile(`^(\w+){(.*)}$`)
	ratingsRE   = regexp.MustCompile(`^{x=(\d+),m=(\d+),a=(\d+),s=(\d+)}$`)
	conditionRE = regexp.MustCompile(`^([xmas])([<>])(\d+)$`)
)

type part struct {
	x, m, a, s int
}

func (p part) sum() int {
	return p.x + p.m + p.a + p.s
}

func (p part) rating(category rune) int {

	switch category {
	case 'x':
		return p.x
	case 'm':
		return p.m
	case 'a':
		return p.a
	case 's':
		return p.s
	default:
		panic("invalid input")
	}
}

type rule struct {
	category rune // x, m, a, s
	dest     string
	operator rune
	value    int
}

func (r rule) String() string {
	return fmt.Sprintf("%c: dest:%s", r.category, r.dest)
}

func parseRules(rawRules []string) []rule {
	var rules []rule
	for _, r := range rawRules {
		ruleParts := strings.Split(r, ":")
		if len(ruleParts) == 1 {
			rules = append(rules, rule{dest: ruleParts[0]})
			continue
		}

		matches := conditionRE.FindStringSubmatch(ruleParts[0])
		dest := ruleParts[1]

		left := matches[1]     // x, m, a, s
		operator := matches[2] // < or >
		right := matches[3]
		value, _ := strconv.Atoi(right)

		newRule := rule{
			category: rune(left[0]),
			dest:     dest,
			value:    value,
			operator: rune(operator[0]),
		}

		rules = append(rules, newRule)
	}

	return rules
}

func parseWorkflows(workflows []string) map[string][]rule {

	flows := make(map[string][]rule, len(workflows))

	for _, line := range workflows {
		parts := workflowRE.FindStringSubmatch(line)
		name := parts[1]
		rawRules := strings.Split(parts[2], ",")
		rules := parseRules(rawRules)
		flows[name] = rules
	}

	return flows
}

func parseRatings(ratings []string) []part {
	var parts []part

	for _, line := range ratings {
		match := ratingsRE.FindStringSubmatch(line)

		x, _ := strconv.Atoi(match[1])
		m, _ := strconv.Atoi(match[2])
		a, _ := strconv.Atoi(match[3])
		s, _ := strconv.Atoi(match[4])

		parts = append(parts, part{x, m, a, s})
	}

	return parts
}

func process(flows map[string][]rule, flowName string, p part) bool {
	// base case
	if flowName == "A" {
		return true
	} else if flowName == "R" {
		return false
	}

	rules := flows[flowName]

	for _, rule := range rules {
		switch rule.operator {
		case '<':
			if p.rating(rule.category) < rule.value {
				return process(flows, rule.dest, p)
			}
			continue
		case '>':
			if p.rating(rule.category) > rule.value {
				return process(flows, rule.dest, p)
			}
			continue
		default: // standalone rule
			return process(flows, rule.dest, p)
		}
	}

	return false
}

func part1(workflows []string, ratings []string) int {

	flows := parseWorkflows(workflows)
	parts := parseRatings(ratings)

	sum := 0
	for _, p := range parts {
		if process(flows, "in", p) {
			sum += p.sum()
		}
	}

	return sum
}

func copyMap(in map[rune][2]int) map[rune][2]int {
	cp := make(map[rune][2]int)
	for k, v := range in {
		cp[k] = v
	}

	return cp
}

func processCombos(flows map[string][]rule, flowName string, ranges map[rune][2]int) int {
	if flowName == "R" {
		return 0
	} else if flowName == "A" { // return the product of ranges
		product := 1
		for _, minMax := range ranges {
			product *= (minMax[1] - minMax[0] + 1)
		}
		return product
	}

	var total int
	rules := flows[flowName]
	for _, rule := range rules {

		minMax := ranges[rule.category]
		newRanges := copyMap(ranges)
		newMinMax := minMax
		switch rule.operator {
		case '<':
			newMinMax[1] = rule.value - 1 // set max to < value (successful range)
			minMax[0] = rule.value        // set min to value (non-successful range)

		case '>':
			newMinMax[0] = rule.value + 1 // set min to > value (successful range)
			minMax[1] = rule.value        //set max to value (non-successful range)

		default: // standalone rule
			total += processCombos(flows, rule.dest, ranges)
			continue
		}

		ranges[rule.category] = minMax
		newRanges[rule.category] = newMinMax
		total += processCombos(flows, rule.dest, newRanges) // (successful)
	}

	return total
}

func part2(workflows []string, ratings []string) int {

	flows := parseWorkflows(workflows)
	categories := map[rune][2]int{
		'x': {1, 4000},
		'm': {1, 4000},
		'a': {1, 4000},
		's': {1, 4000},
	}

	return processCombos(flows, "in", categories)
}

func main() {
	workflows, ratings := parseInput()

	fmt.Println("Part 1:", part1(workflows, ratings))
	fmt.Println("Part 2:", part2(workflows, ratings))
}
