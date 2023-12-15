package main

import (
	"fmt"
	"os"
	"strconv"
	"strings"
)

func parseInput() []string {
	raw, _ := os.ReadFile(os.Args[1])
	return strings.Split(strings.Trim(string(raw), "\n"), ",")
}

type lens struct {
	label       string
	focalLength int
}

func (l lens) String() string {
	return fmt.Sprintf("[%s %d]", l.label, l.focalLength)
}

func hash(step string) int {

	cv := 0
	for i := range step {
		cv += int(step[i])
		cv *= 17
		cv %= 256
	}

	return cv
}

func part1(steps []string) int {

	s := 0
	for _, step := range steps {
		s += hash(step)
	}

	return s
}

func remove(label string, lenses []lens) []lens {
	for i, v := range lenses {
		if v.label == label {
			return append(lenses[0:i], lenses[i+1:]...)
		}
	}

	return lenses
}

func add(v lens, lenses []lens) []lens {

	for i, l := range lenses {
		l := l
		if l.label == v.label {
			lenses[i] = v
			return lenses
		}
	}

	return append(lenses, v)
}

func display(h map[int][]lens) {
	for k, lenses := range h {
		fmt.Printf("Box %d: ", k)
		for _, l := range lenses {
			fmt.Printf("%s ", l)
		}
		fmt.Println()
	}
}

func part2(steps []string) int {

	boxMap := make(map[int][]lens, 256)

	for _, step := range steps {
		var label string
		var focalLength int
		parts := strings.Split(step, "=")
		if len(parts) == 2 {
			label = parts[0]
			focalLength, _ = strconv.Atoi(parts[1])
			h := hash(label)
			box := boxMap[h]
			boxMap[h] = add(lens{label, focalLength}, box)
		} else {
			label = strings.TrimRight(step, "-")
			h := hash(label)
			boxMap[h] = remove(label, boxMap[h])
		}
	}

	s := 0
	// display(boxMap)
	for box, lenses := range boxMap {
		for i, l := range lenses {
			s += (1 + box) * (i + 1) * l.focalLength
		}
	}

	return s
}

func main() {

	input := parseInput()

	fmt.Println("Part 1:", part1(input))
	fmt.Println("Part 2:", part2(input))
}
