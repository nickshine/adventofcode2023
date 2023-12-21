package main

import (
	"fmt"
	"os"
	"slices"
	"strings"
)

func parseInput() []string {
	raw, _ := os.ReadFile(os.Args[1])
	input := strings.Split(strings.Trim(string(raw), "\n"), "\n")
	return input
}

const (
	FLIPFLOP int = iota
	CONJUNCTION
	BROADCASTER
)

type pulse int

func (p pulse) String() string {
	if p == LOW {
		return "low"
	} else if p == HIGH {
		return "high"
	}

	return ""
}

const (
	_ pulse = iota
	LOW
	HIGH
)

type module struct {
	name   string
	typ    int
	on     bool
	dests  []string         // list of module names
	inputs map[string]pulse // map of connected input names to pulse types
}

type state struct {
	from, to module
	pulse    pulse
}

func parseModules(in []string) map[string]module {

	modules := make(map[string]module, len(in))

	for _, line := range in {

		parts := strings.Split(line, " -> ")
		name := parts[0]
		dests := strings.Split(parts[1], ", ")

		var typ int

		switch {
		case name[0] == '%':
			typ = FLIPFLOP
			name = name[1:]
		case name[0] == '&':
			typ = CONJUNCTION
			name = name[1:]
		case name == "broadcaster":
			typ = BROADCASTER
		}
		m := module{name: name, typ: typ, on: false, dests: dests, inputs: make(map[string]pulse)}

		modules[name] = m

	}

	// initialize connected inputs
	for k, m := range modules {
		for _, dest := range m.dests {
			dm := modules[dest]
			if dm.inputs == nil {
				dm.inputs = make(map[string]pulse)
			}
			dm.inputs[k] = LOW
			modules[dest] = dm
		}
	}

	return modules
}

func solve(in []string, part1 bool) int {

	modules := parseModules(in)
	modules["output"] = module{name: "output"}

	var cycles []int

	high, low := 0, 0

OUTER:
	for i := 0; ; i++ {
		if i == 1000 && part1 {
			break
		}
		queue := []state{{from: module{}, to: modules["broadcaster"], pulse: LOW}}

		for len(queue) > 0 {
			cur := queue[0]
			if cur.pulse == LOW {
				low++
			} else {
				high++
			}
			queue = queue[1:]

			// manually inspected input to see all inputs to rx inv
			// cl -> lx
			// rp -> lx   -> rx
			// lb -> lx
			// nj -> lx

			// if any of these modules hit, a cycle has occurred
			if slices.Contains([]string{"cl", "rp", "lb", "nj"}, cur.to.name) && cur.pulse == LOW {
				cycles = append(cycles, i+1)
			}

			if len(cycles) == 4 {
				break OUTER
			}

			switch cur.to.typ {
			case BROADCASTER:
				for _, d := range cur.to.dests {
					m := modules[d]
					if m.typ == FLIPFLOP {
						m.on = !m.on
						modules[m.name] = m
					}
					queue = append(queue, state{from: cur.to, to: m, pulse: LOW})
				}
			case FLIPFLOP:
				switch cur.pulse {
				case LOW:
					for _, d := range cur.to.dests {
						m := modules[d]

						if cur.to.on {
							queue = append(queue, state{from: cur.to, to: m, pulse: HIGH})
						} else {
							if m.typ == FLIPFLOP {
								m.on = !m.on
								modules[m.name] = m
							}
							queue = append(queue, state{from: cur.to, to: m, pulse: LOW})
						}
					}

				case HIGH:
					continue
				}
			case CONJUNCTION:
				cur.to.inputs[cur.from.name] = cur.pulse

				allHigh := true
				for _, pulse := range cur.to.inputs {
					if pulse == LOW {
						allHigh = false
						break
					}
				}

				for _, d := range cur.to.dests {
					m := modules[d]
					if allHigh {
						if m.typ == FLIPFLOP {
							m.on = !m.on
							modules[m.name] = m
						}
						queue = append(queue, state{from: cur.to, to: m, pulse: LOW})
					} else {
						queue = append(queue, state{from: cur.to, to: m, pulse: HIGH})
					}
				}
			}
		}
	}

	if part1 {
		return high * low
	}
	return cycles[0] * cycles[1] * cycles[2] * cycles[3]
}

func main() {
	input := parseInput()

	fmt.Println("Part 1:", solve(input, true))
	fmt.Println("Part 2:", solve(input, false))
}
