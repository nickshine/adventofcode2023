package main

import (
	"fmt"
	"os"
	"strings"
)

func parseInput() (grid [][]rune, start point) {
	raw, _ := os.ReadFile(os.Args[1])
	rows := strings.Split(strings.Trim(string(raw), "\n"), "\n")

	grid = make([][]rune, len(rows))

	for y, row := range rows {
		grid[y] = make([]rune, len(row))
		for x, c := range row {
			grid[y][x] = c
			if c == 'S' {
				start = point{x, y}
			}
		}
	}

	return grid, start
}

type point struct {
	x, y int
}

func (p point) String() string {
	return fmt.Sprintf("(%d,%d)", p.x, p.y)
}

type state struct {
	p    point
	step int
}

func part1(grid [][]rune, start point) int {

	// can be any plot that is even number steps away, up to 64 away
	const MAXSTEP = 64

	seen := map[state]struct{}{{start, 0}: {}}
	evens := map[point]struct{}{start: {}}
	queue := []state{{start, 0}}

	for len(queue) > 0 {
		cur := queue[0]
		queue = queue[1:]

		if cur.step > MAXSTEP {
			continue
		}

		x, y := cur.p.x, cur.p.y

		// proceed an even number of steps
		choices := []point{
			{x: x, y: y - 1}, // up
			{x: x, y: y + 1}, // down
			{x: x - 1, y: y}, // left
			{x: x + 1, y: y}, // right
		}

		for _, p := range choices {
			// out of bounds
			if p.y < 0 || p.y >= len(grid) || p.x < 0 || p.x >= len(grid[0]) {
				continue
			}

			if grid[p.y][p.x] == '#' {
				continue
			}

			next := state{p, cur.step + 1}

			if _, ok := seen[next]; ok {
				continue
			}

			seen[next] = struct{}{}
			if next.step%2 == 0 {
				evens[next.p] = struct{}{}
			}

			queue = append(queue, next)
		}
	}

	return len(evens)
}

func part2(in [][]rune, start point) int {
	return 0
}

func main() {
	grid, start := parseInput()

	fmt.Println("Part 1:", part1(grid, start))
	fmt.Println("Part 2:", part2(grid, start))
}
