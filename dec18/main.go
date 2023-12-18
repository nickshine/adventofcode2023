package main

import (
	"fmt"
	"os"
	"strconv"
	"strings"
)

func parseInput() []string {
	raw, _ := os.ReadFile(os.Args[1])
	return strings.Split(strings.Trim(string(raw), "\n"), "\n")
}

type point struct {
	x, y int
}

type direction point

var (
	UP    = direction{0, -1}
	DOWN  = direction{0, 1}
	LEFT  = direction{-1, 0}
	RIGHT = direction{1, 0}
)

func parseDirection(d string) direction {
	if len(d) != 1 {
		panic("invalid input")
	}

	r := d[0]

	var out direction
	switch r {
	case 'R', '0':
		out = RIGHT
	case 'D', '1':
		out = DOWN
	case 'L', '2':
		out = LEFT
	case 'U', '3':
		out = UP
	default:
		panic("invalid input")
	}

	return out
}

func parsePoints(in []string) []point {

	var points []point
	x, y := 0, 0

	for _, line := range in {
		parts := strings.Split(line, " ")
		dir := parseDirection(parts[0])
		count, _ := strconv.Atoi(parts[1])

		for i := 0; i < count; i++ {
			x, y = x+dir.x, y+dir.y
			points = append(points, point{x, y})
		}
	}

	return points
}

func parseHexPoints(in []string) []point {

	var points []point
	x, y := 0, 0

	for _, line := range in {
		parts := strings.Split(line, " ")
		hex := strings.TrimPrefix(parts[2], "(#")
		hex = strings.TrimSuffix(hex, ")")
		count, _ := strconv.ParseInt(hex[:5], 16, 64)
		dir := parseDirection(hex[5:])

		for i := 0; i < int(count); i++ {
			x, y = x+dir.x, y+dir.y
			points = append(points, point{x, y})
		}
	}

	return points
}

func solve(points []point) int {
	// https://en.wikipedia.org/wiki/Shoelace_formula
	sum := 0
	for i := 0; i < len(points)-1; i++ {

		sum += (points[i].y + points[i+1].y) * (points[i].x - points[i+1].x)
	}
	sum += (points[len(points)-1].y + points[0].y) * (points[len(points)-1].x - points[0].x)
	area := sum / 2

	// https://en.wikipedia.org/wiki/Pick's_theorem
	inside := area - len(points)/2 + 1
	outside := len(points)

	return inside + outside

}

func part1(in []string) int {
	points := parsePoints(in)
	return solve(points)
}

func part2(in []string) int {
	points := parseHexPoints(in)
	return solve(points)
}

func main() {
	input := parseInput()

	fmt.Println("Part 1:", part1(input))
	fmt.Println("Part 2:", part2(input))
}
