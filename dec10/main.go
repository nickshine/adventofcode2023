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

type tiles [][]point

func (t tiles) String() string {
	var sb strings.Builder

	for _, y := range t {
		for _, x := range y {
			sb.WriteRune(x.r)
		}
		sb.WriteString("\n")
	}
	return sb.String()
}

type point struct {
	x, y int
	r    rune
}

func (p point) String() string {
	return fmt.Sprintf("%c(%d,%d)", p.r, p.x, p.y)
}

type state struct {
	p    point
	step int
}

// parse returns the staring point and tiles grid.
func parse(in []string) (point, tiles) {

	var t [][]point
	var start point

	for y, line := range in {
		var row []point
		for x, r := range line {
			p := point{x, y, r}
			if r == 'S' {
				start = p
			}
			row = append(row, p)
		}

		t = append(t, row)
	}

	return start, t
}

// isConnectedLeft returns true if b is left of a, and connected.
func (a point) isConnectedLeft(b point) bool {
	if a.y == b.y && a.x == b.x+1 {
		if (a.r == '-' || a.r == 'J' || a.r == '7' || a.r == 'S') && (b.r == '-' || b.r == 'L' || b.r == 'F' || b.r == 'S') {
			return true
		}
	}
	return false
}

// isConnectedLeft returns true if b is right of a, and connected.
func (a point) isConnectedRight(b point) bool {
	if a.y == b.y && a.x == b.x-1 {
		if (a.r == '-' || a.r == 'L' || a.r == 'F' || a.r == 'S') && (b.r == '-' || b.r == 'J' || b.r == '7' || b.r == 'S') {
			return true
		}
	}

	return false
}

// isConnectedUp returns true if b is above of a, and connected.
func (a point) isConnectedUp(b point) bool {
	if a.x == b.x && a.y == b.y+1 {
		if (a.r == '|' || a.r == 'L' || a.r == 'J' || a.r == 'S') && (b.r == '|' || b.r == 'F' || b.r == '7' || b.r == 'S') {
			return true
		}
	}
	return false
}

// isConnectedDown returns true if b is below a, and connected.
func (a point) isConnectedDown(b point) bool {
	if a.x == b.x && a.y == b.y-1 {
		if (a.r == '|' || a.r == 'F' || a.r == '7' || a.r == 'S') && (b.r == '|' || b.r == 'L' || b.r == 'J' || b.r == 'S') {
			return true
		}
	}

	return false
}

func (a point) isConnected(b point) bool {
	return a.isConnectedLeft(b) || a.isConnectedRight(b) || a.isConnectedUp(b) || a.isConnectedDown(b)
}

func (s state) String() string {
	return fmt.Sprintf("(%s, step %d)", s.p, s.step)
}

// bfs returns the loop points and max distance from start for the loop
func bfs(t tiles, start point) (map[point]struct{}, int) {
	queue := []state{{start, 0}}
	seen := map[point]struct{}{start: {}}
	maxStep := 0
	for len(queue) > 0 {
		cur := queue[0]
		queue = queue[1:]
		// fmt.Printf("Current; %s, queue: %v\n", cur, queue)

		if cur.step > maxStep {
			maxStep = cur.step
		}

		x, y := cur.p.x, cur.p.y
		neighbors := []point{
			{x: x, y: y - 1}, // up
			{x: x, y: y + 1}, // down
			{x: x - 1, y: y}, // left
			{x: x + 1, y: y}, // right
		}

		for _, p := range neighbors {
			if p.y < 0 || p.y >= len(t) || p.x < 0 || p.x >= len(t[0]) {
				// fmt.Printf("invalid neighbor %s\n", p)
				continue
			}

			p = t[p.y][p.x]
			next := state{p, cur.step + 1}
			// fmt.Printf("next: %s\n", next)

			if _, ok := seen[next.p]; ok {
				// fmt.Printf("SEEN: %s, continuing\n", next.p)
				continue
			}

			if cur.p.isConnected(p) {
				// fmt.Printf("Adding neighbor %s\n", next)
				seen[p] = struct{}{}
				queue = append(queue, next)
			}
		}
	}
	return seen, maxStep
}

func part1(input []string) int {

	start, tiles := parse(input)
	fmt.Println(tiles)
	_, maxDistance := bfs(tiles, start)
	return maxDistance
}

// https://en.wikipedia.org/wiki/Point_in_polygon
// if the ray cast intersects the polygon an even number of times, it is outside
// if the ray cast intersects the polygon an odd number of times, it is inside
// if it enters the polygon boundary, it does not "intersect" until it leaves boundary
func part2(input []string) int {

	start, tiles := parse(input)
	fmt.Println(tiles)
	polygon, _ := bfs(tiles, start)

	total := 0
	for i, row := range tiles {
		fmt.Printf("___________scanning row %d\n", i)
		intersections := 0
		for _, p := range row {
			if _, ok := polygon[p]; ok {
				switch p.r {
				case '|', 'L', 'J', 'S': // use the bottoms of the vertical edges to only
					fmt.Printf("Intersection found for point %s\n", p)
					intersections++
				}
			} else {
				if intersections%2 == 1 {
					total++
				}
			}
		}
	}

	return total
}

func main() {

	input := parseInput()

	fmt.Println("Part 1:", part1(input))
	fmt.Println("Part 2:", part2(input))
}
