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
	op       func(n int) bool
}

func (r rule) String() string {
	return fmt.Sprintf("%c: dest:%s", r.category, r.dest)
}

func parseRules(rawRules []string) []rule {
	var rules []rule
	for _, r := range rawRules {
		ruleParts := strings.Split(r, ":")
		if len(ruleParts) == 1 {
			fmt.Printf("  standalone rule %s\n", ruleParts[0])
			rules = append(rules, rule{dest: ruleParts[0]})
			continue
		}

		matches := conditionRE.FindStringSubmatch(ruleParts[0])
		dest := ruleParts[1]

		left := matches[1]     // x, m, a, s
		operator := matches[2] // < or >
		right := matches[3]
		value, _ := strconv.Atoi(right)
		fmt.Printf("  left: %s, op: %s, right: %s, dest: %s\n", left, operator, right, dest)

		var fn func(n int) bool

		switch operator {
		case "<":
			fn = func(n int) bool {
				return n < value
			}
		case ">":
			fn = func(n int) bool {
				return n > value
			}
		default:
			panic("invalid input")
		}

		newRule := rule{
			category: rune(left[0]),
			op:       fn,
			dest:     dest,
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
		fmt.Printf("name: %s\n", name)
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
		if rule.op != nil {
			// if condition satisfied, go to dest
			if rule.op(p.rating(rule.category)) {
				return process(flows, rule.dest, p)
			}
			// otherwise, go to next rule
			continue
		}

		// standalone rule (no condition)
		return process(flows, rule.dest, p)
	}

	return false
}

func part1(workflows []string, ratings []string) int {

	flows := parseWorkflows(workflows)
	parts := parseRatings(ratings)

	// for flow, rules := range flows {
	// 	fmt.Printf("workflow: %s\n", flow)
	// 	for i, rule := range rules {
	// 		fmt.Printf("  rule %d - %s\n", i, rule)
	// 	}
	// }

	sum := 0
	for _, p := range parts {
		if process(flows, "in", p) {
			sum += p.sum()
		}
	}

	return sum
}

func part2(workflows []string, ratings []string) int {

	return 0
}

func main() {
	workflows, ratings := parseInput()

	fmt.Println("Part 1:", part1(workflows, ratings))
	fmt.Println("Part 2:", part2(workflows, ratings))
}
