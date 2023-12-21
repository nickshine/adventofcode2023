package main

import (
	"fmt"
	"os"
	"strings"
)

func parseInput() (grid []string, start point) {
	raw, _ := os.ReadFile(os.Args[1])
	grid = strings.Split(strings.Trim(string(raw), "\n"), "\n")

	for y := range grid {
		for x := range grid[y] {
			if grid[y][x] == 'S' {
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

func dfs(grid []string, p point, step, maxStep int, seen map[state]struct{}) int {

	if step == maxStep {
		return 1
	}

	x, y := p.x, p.y
	choices := []point{
		{x: x, y: y - 1}, // up
		{x: x, y: y + 1}, // down
		{x: x - 1, y: y}, // left
		{x: x + 1, y: y}, // right
	}

	total := 0
	n := len(grid)
	for _, c := range choices {
		// if c.y < 0 || c.y >= len(grid) || c.x < 0 || c.x >= len(grid[0]) {
		// 	continue
		// }

		// go modulo can have negatives, so need extra % - https://github.com/golang/go/issues/448#issuecomment-66049769
		if grid[(c.y%n+n)%n][(c.x%n+n)%n] == '#' {
			continue
		}

		point := point{c.x, c.y}
		state := state{point, step}
		if _, ok := seen[state]; ok {
			continue
		}

		seen[state] = struct{}{}
		total += dfs(grid, point, step+1, maxStep, seen)
	}

	return total
}

func bfs(grid []string, start point, maxStep int) int {
	seen := map[state]struct{}{{start, 0}: {}}
	evens := map[point]struct{}{start: {}}
	queue := []state{{start, 0}}

	for len(queue) > 0 {
		cur := queue[0]
		queue = queue[1:]

		if cur.step > maxStep {
			continue
		}

		x, y := cur.p.x, cur.p.y

		choices := []point{
			{x: x, y: y - 1}, // up
			{x: x, y: y + 1}, // down
			{x: x - 1, y: y}, // left
			{x: x + 1, y: y}, // right
		}

		n := len(grid)
		for _, p := range choices {
			// out of bounds
			// if p.y < 0 || p.y >= len(grid) || p.x < 0 || p.x >= len(grid[0]) {
			// 	continue
			// }

			// if grid[p.y][p.x] == '#' {
			// 	continue
			// }

			// go modulo can have negatives, so need extra % - https://github.com/golang/go/issues/448#issuecomment-66049769
			if grid[(p.y%n+n)%n][(p.x%n+n)%n] == '#' {
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

func part1(grid []string, start point) int {
	a := dfs(grid, start, 0, 64, make(map[state]struct{}))
	b := bfs(grid, start, 64)

	fmt.Println("dfs:", a)
	fmt.Println("bfs:", b)

	return b
}

func part2(grid []string, start point) int {
	const maxSteps = 26501365
	// s1 := bfs(grid, start, start.y)
	// s2 := bfs(grid, start, start.y+len(grid))
	// s3 := bfs(grid, start, start.y+(len(grid)*2))
	// fmt.Println(s1, s2, s3)

	// Starting point is on an empty row in center of grid (row 65)
	s1 := dfs(grid, start, 0, start.y, make(map[state]struct{}))
	s2 := dfs(grid, start, 0, start.y+len(grid), make(map[state]struct{}))
	s3 := dfs(grid, start, 0, start.y+(len(grid)*2), make(map[state]struct{}))

	// quadratic - x is the number of repeated grids in x direction
	// ax^2 + bx + c
	x := maxSteps / len(grid)

	// ((diff between s2 and s3) - (diff between s1 and s2)) / 2
	// (s3-s2)-(s2-s1)
	a := (s3 - 2*s2 + s1) / 2
	b := s2 - s1 - a
	c := s1

	return (a*x*x + b*x + c)
}

func main() {
	grid, start := parseInput()

	fmt.Println("Part 1:", part1(grid, start))
	fmt.Println("Part 2:", part2(grid, start))
}
