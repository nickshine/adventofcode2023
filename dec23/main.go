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

type point struct {
	x, y int
}

func (p point) String() string {
	return fmt.Sprintf("(%d,%d)", p.x, p.y)
}

var (
	UP    = [2]int{0, -1}
	DOWN  = [2]int{0, 1}
	LEFT  = [2]int{-1, 0}
	RIGHT = [2]int{1, 0}

	directions = map[rune][][2]int{
		'^': {UP},
		'v': {DOWN},
		'<': {LEFT},
		'>': {RIGHT},
		'.': {UP, DOWN, LEFT, RIGHT},
	}

	directions2 = map[rune][][2]int{
		'^': {UP, DOWN, LEFT, RIGHT},
		'v': {UP, DOWN, LEFT, RIGHT},
		'<': {UP, DOWN, LEFT, RIGHT},
		'>': {UP, DOWN, LEFT, RIGHT},
		'.': {UP, DOWN, LEFT, RIGHT},
	}
)

func draw(grid []string, seen map[point]struct{}) {
	fmt.Println()
	for y, row := range grid {
		for x, c := range row {
			if _, ok := seen[point{x, y}]; ok {
				fmt.Printf("O")
			} else {
				fmt.Printf("%c", c)
			}
		}
		fmt.Println()
	}
}

func part1(grid []string, start, end point) int {
	var dfs func(p point, step int, seen map[point]struct{}) int

	dfs = func(p point, step int, seen map[point]struct{}) int {
		if p == end {
			return step
		}
		seen[p] = struct{}{}
		dirs := directions[rune(grid[p.y][p.x])]

		maxDepth := -1
		for _, d := range dirs {
			x, y := p.x+d[0], p.y+d[1]

			if y < 0 || y >= len(grid) || x < 0 || x >= len(grid[0]) {
				continue
			}

			if grid[y][x] == '#' {
				continue
			}

			next := point{x, y}
			if _, ok := seen[next]; ok {
				continue
			}
			maxDepth = max(maxDepth, dfs(next, step+1, seen))
		}

		delete(seen, p)
		return maxDepth
	}

	return dfs(start, 0, make(map[point]struct{}))
}

type edge struct {
	from, to point
	weight   int
}

func (e edge) String() string {
	return fmt.Sprintf("%s --%d--> %s", e.from, e.weight, e.to)
}

type graph struct {
	edges map[point]map[edge]struct{}
}

func (g *graph) addEdge(from, to point, weight int) {
	if _, ok := g.edges[from]; !ok {
		g.edges[from] = make(map[edge]struct{})
	}
	g.edges[from][edge{from, to, weight}] = struct{}{}

	// undirected
	if _, ok := g.edges[to]; !ok {
		g.edges[to] = make(map[edge]struct{})
	}

	g.edges[to][edge{to, from, weight}] = struct{}{}
}

func (g *graph) removePoint(p point) {
	edges := g.edges[p]
	var points []point
	weight := 0
	for e := range edges {
		points = append(points, e.to)
		weight += e.weight
		delete(edges, e)
	}
	delete(g.edges, p)

	// undirected
	for _, pp := range points {
		edges := g.edges[pp]

		for e := range edges {
			if e.to == p {
				delete(edges, e)
			}
		}
	}

	// add the consolidated edge
	g.addEdge(points[0], points[1], weight)
}

func (g *graph) compress() {
	var points []point

	for p, edges := range g.edges {
		// remove only points that are not an intersection
		if len(edges) == 2 {
			points = append(points, p)
		}
	}

	for _, p := range points {
		g.removePoint(p)
	}
}

func newGraph(grid []string) *graph {
	g := &graph{
		edges: make(map[point]map[edge]struct{}),
	}

	for y, row := range grid {
		for x, c := range row {
			if c == '#' {
				continue
			}
			a := point{x, y}
			for _, d := range directions2[rune(grid[y][x])] {
				nx, ny := x+d[0], y+d[1]

				if ny < 0 || ny >= len(grid) || nx < 0 || nx >= len(grid[0]) {
					continue
				}

				if grid[ny][nx] == '#' {
					continue
				}

				b := point{nx, ny}
				g.addEdge(a, b, 1)
			}
		}
	}

	return g
}

func part2(grid []string, start, end point) int {
	g := newGraph(grid)
	g.compress()

	// // display intersections
	// for p, edges := range g.edges {
	// 	if len(edges) > 2 {
	// 		for edge := range edges {
	// 			fmt.Printf("Intesection at point %s, has edge %s\n", p, edge)
	// 		}
	// 	} else {
	// 		fmt.Printf("point %s, DOES NOT have 2 edges: %+v\n", p, edges)
	// 	}
	// }

	var dfs func(p point, step int, seen map[point]struct{}) int

	dfs = func(p point, step int, seen map[point]struct{}) int {
		if p == end {
			return step
		}

		seen[p] = struct{}{}
		maxDepth := -1
		for e := range g.edges[p] {

			next := e.to
			if _, ok := seen[next]; ok {
				continue
			}

			maxDepth = max(maxDepth, dfs(next, step+e.weight, seen))
		}

		delete(seen, p)
		return maxDepth
	}

	total := dfs(start, 0, make(map[point]struct{}))
	return total
}

func main() {
	grid := parseInput()
	start, end := point{1, 0}, point{len(grid[0]) - 2, len(grid) - 1}

	fmt.Println("Part 1:", part1(grid, start, end))
	fmt.Println("Part 2:", part2(grid, start, end))
}
