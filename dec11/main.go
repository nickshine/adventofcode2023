package main

import (
	"fmt"
	"math"
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

type edge struct {
	from, to point
	weight   int
}

func (e edge) ID() string {
	return fmt.Sprintf("%s-%s", e.from, e.to)
}

func (e edge) String() string {
	return fmt.Sprintf("%s--%s\t(%d)", e.from, e.to, e.weight)
}

type graph struct {
	nodes map[point]struct{}
	edges map[string]*edge
}

func (g graph) String() string {
	var sb strings.Builder

	sb.WriteString("Nodes:\n")
	for node := range g.nodes {
		sb.WriteString(fmt.Sprintf("%s ", node))
	}
	sb.WriteString("\nEdges:\n")
	for _, edg := range g.edges {
		sb.WriteString(fmt.Sprintf("\t%s\n", edg))
	}
	sb.WriteString("\n")
	return sb.String()
}

func (g *graph) getEdge(a, b point) (edge, bool) {

	edg := edge{from: a, to: b}
	if _, ok := g.edges[edg.ID()]; ok {
		return edg, true
	}

	edg = edge{from: b, to: a}
	if _, ok := g.edges[edg.ID()]; ok {
		return edg, true
	}

	return edge{}, false
}

func expandRows(rows []string) []string {
	var expanded []string
	for _, row := range rows {
		empty := true
		for _, c := range row {
			if c != '.' {
				empty = false
				break
			}
		}
		expanded = append(expanded, row)
		if empty {
			expanded = append(expanded, row)
		}
	}

	return expanded
}

func expandColumns(rows []string) []string {
	// get empty column indices
	var emptyColumns []int
	for x := 0; x < len(rows[0]); x++ {
		empty := true
		for y := 0; y < len(rows); y++ {
			if rune(rows[y][x]) != '.' {
				empty = false
				break
			}
		}
		if empty {
			emptyColumns = append(emptyColumns, x)
		}
	}

	// expand columns
	var expanded []string
	for _, row := range rows {
		var sb strings.Builder
		l := 0
		for _, c := range emptyColumns {
			sb.WriteString(row[l:c])
			sb.WriteRune(rune(row[c]))
			l = c
		}
		sb.WriteString(row[l:])

		expanded = append(expanded, sb.String())
	}

	return expanded

}

func expandSpace(rows []string) []string {

	expandedRows := expandRows(rows)
	return expandColumns(expandedRows)
}

func parseGraph(rows []string) *graph {

	g := &graph{
		nodes: make(map[point]struct{}),
		edges: make(map[string]*edge),
	}

	// fmt.Println("Original graph:")
	for _, r := range rows {
		fmt.Println(r)
	}

	// fmt.Println("Expanded graph:")
	expandedRows := expandSpace(rows)
	for _, r := range expandedRows {
		fmt.Println(r)
	}

	// add nodes to graph
	for y, row := range expandedRows {
		for x, r := range row {
			switch r {
			case '.':
				continue
			case '#':
				g.nodes[point{x, y}] = struct{}{}
			}
		}
	}

	// add edges (pairs) to graph
	for a := range g.nodes {
		for b := range g.nodes {
			if a == b {
				continue
			}
			if _, exists := g.getEdge(a, b); !exists {
				e := &edge{from: a, to: b}
				g.edges[e.ID()] = e
			}
		}
	}

	return g
}

func (g *graph) calcShortestPaths() {
	for _, edge := range g.edges {

		a, b := edge.from, edge.to
		dx := math.Abs(float64(a.x - b.x))
		dy := math.Abs(float64(a.y - b.y))

		edge.weight = int(dx + dy)
	}
}

func part1(input []string) int {

	g := parseGraph(input)
	g.calcShortestPaths()

	sum := 0
	for _, edge := range g.edges {
		sum += edge.weight
	}

	return sum
}

func part2(input []string, expansion int) int {

	var galaxies []point
	var rowOccupied = make(map[int]struct{})
	var colOccupied = make(map[int]struct{})

	for y, row := range input {
		for x, c := range row {
			if c == '#' {
				galaxies = append(galaxies, point{x, y})
				rowOccupied[y] = struct{}{}
				colOccupied[x] = struct{}{}
			}
		}
	}

	// distance from row 0 to row i
	rd := []int{0}
	// distance from col 0 to col i
	cd := []int{0}

	for i := range input {
		if _, ok := rowOccupied[i]; ok {
			rd = append(rd, rd[len(rd)-1]+1)
		} else {
			rd = append(rd, rd[len(rd)-1]+expansion)
		}
	}

	for i := range input[0] {
		if _, ok := colOccupied[i]; ok {
			cd = append(cd, cd[len(cd)-1]+1)
		} else {
			cd = append(cd, cd[len(cd)-1]+expansion)
		}
	}

	var sum float64

	for i, ip := range galaxies {
		for _, jp := range galaxies[:i] {
			sum += math.Abs(float64(rd[ip.y] - rd[jp.y]))
			sum += math.Abs(float64(cd[ip.x] - cd[jp.x]))
		}
	}

	return int(sum)
}

func main() {

	input := parseInput()

	fmt.Println("Part 1:", part1(input))
	fmt.Println("Part 2:", part2(input, 1000000))
}
