package main

import (
	"fmt"
	"os"
	"strings"
)

func parseInput() []string {
	raw, _ := os.ReadFile(os.Args[1])
	return strings.Split(strings.Trim(string(raw), "\n"), "\n")
}

func parseGrid(input []string) [][]tile {

	grid := make([][]tile, len(input))
	for y, row := range input {
		grid[y] = make([]tile, len(row))
		for x, c := range row {
			grid[y][x] = tile{v: c, beams: make(map[beam]struct{})}
		}
	}

	return grid
}

type point struct {
	x, y int
}

func (p point) String() string {
	return fmt.Sprintf("(%d,%d)", p.x, p.y)
}

type direction point

var (
	RIGHT = direction{1, 0}
	DOWN  = direction{0, 1}
	LEFT  = direction{-1, 0}
	UP    = direction{0, -1}
)

func (d direction) String() string {
	switch d {
	case RIGHT:
		return ">"
	case DOWN:
		return "v"
	case LEFT:
		return "<"
	case UP:
		return "^"
	}
	return ""
}

type beam direction

type tile struct {
	v     rune
	beams map[beam]struct{}
}

func (t tile) String() string {
	if t.v == '.' {
		switch len(t.beams) {
		case 0:
			return string(t.v)
		case 1:
			for b := range t.beams {
				return fmt.Sprintf("%s", direction(b))
			}
		default:
			return fmt.Sprintf("%d", len(t.beams))
		}
	}

	return string(t.v)
}

func energize(grid [][]tile, p point, b beam) {
	if p.y < 0 || p.y >= len(grid) || p.x < 0 || p.x >= len(grid[p.y]) { // out of bounds
		return
	}

	x, y := p.x, p.y
	dx, dy := b.x, b.y
	t := grid[y][x]

	// base case - if a beam with same direction has been seen, stop recursion
	if _, ok := t.beams[b]; ok {
		return
	}

	t.beams[b] = struct{}{}

	switch t.v {

	case '.':
		energize(grid, point{x + dx, y + dy}, b)
	case '|':
		switch direction(b) {
		case RIGHT, LEFT:
			energize(grid, point{x + UP.x, y + UP.y}, beam(UP))
			energize(grid, point{x + DOWN.x, y + DOWN.y}, beam(DOWN))
		case UP, DOWN:
			energize(grid, point{x + dx, y + dy}, b)
		}
	case '-':
		switch direction(b) {
		case RIGHT, LEFT:
			energize(grid, point{x + dx, y + dy}, b)
		case UP, DOWN:
			energize(grid, point{x + LEFT.x, y + LEFT.y}, beam(LEFT))
			energize(grid, point{x + RIGHT.x, y + RIGHT.y}, beam(RIGHT))
		}
	case '/':
		switch direction(b) {
		case RIGHT:
			energize(grid, point{x + UP.x, y + UP.y}, beam(UP))
		case DOWN:
			energize(grid, point{x + LEFT.x, y + LEFT.y}, beam(LEFT))
		case LEFT:
			energize(grid, point{x + DOWN.x, y + DOWN.y}, beam(DOWN))
		case UP:
			energize(grid, point{x + RIGHT.x, y + RIGHT.y}, beam(RIGHT))
		}
	case '\\':
		switch direction(b) {
		case RIGHT:
			energize(grid, point{x + DOWN.x, y + DOWN.y}, beam(DOWN))
		case DOWN:
			energize(grid, point{x + RIGHT.x, y + RIGHT.y}, beam(RIGHT))
		case LEFT:
			energize(grid, point{x + UP.x, y + UP.y}, beam(UP))
		case UP:
			energize(grid, point{x + LEFT.x, y + LEFT.y}, beam(LEFT))
		}
	default:
		panic("invalid input")

	}
}

func display(grid [][]tile) {
	for _, row := range grid {
		for _, t := range row {
			fmt.Printf("%s", t)
		}
		fmt.Println()
	}
}

func copyTiles(in [][]tile) [][]tile {
	var out [][]tile
	for y, row := range in {
		out = append(out, make([]tile, len(row)))
		for x, c := range row {
			out[y][x] = tile{v: c.v, beams: make(map[beam]struct{})}
		}
	}
	return out
}

func energized(grid [][]tile) int {
	visited := 0
	for _, row := range grid {
		for _, t := range row {
			if len(t.beams) > 0 {
				visited++
			}
		}
	}
	return visited
}

func getStarts(grid [][]tile) map[point]direction {

	edges := make(map[point]direction)

	for y, row := range grid {
		for x := range row {
			if y == 0 { // top edge
				edges[point{x, y}] = DOWN
			} else if y == len(grid)-1 { // bottom edge
				edges[point{x, y}] = UP
			} else if x == 0 { // left edge
				edges[point{x, y}] = RIGHT
			} else if x == len(row)-1 {
				edges[point{x, y}] = LEFT // right edge
			}
		}
	}

	return edges
}

func part1(input []string) int {
	grid := parseGrid(input)

	energize(grid, point{0, 0}, beam(RIGHT))
	return energized(grid)
}

func part2(input []string) int {
	grid := parseGrid(input)
	edges := getStarts(grid)

	max := 0
	for p, d := range edges {
		g := copyTiles(grid)
		energize(g, p, beam(d))
		count := energized(g)
		if count > max {
			max = count
		}
	}
	return max
}

func main() {

	input := parseInput()

	fmt.Println("Part 1:", part1(input))
	fmt.Println("Part 2:", part2(input))
}
