package main

import (
	"container/heap"
	"fmt"
	"os"
	"slices"
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

func (p point) String() string {
	return fmt.Sprintf("(%d,%d)", p.x, p.y)
}

type node struct {
	loss  int
	state state
}

type state struct {
	p  point
	d  direction
	mv int
}

type direction point

var (
	UP    = direction{0, -1}
	DOWN  = direction{0, 1}
	LEFT  = direction{-1, 0}
	RIGHT = direction{1, 0}
)

var directions = [4]direction{UP, DOWN, LEFT, RIGHT}

var turns = map[direction][]direction{
	UP:    {LEFT, RIGHT},
	DOWN:  {LEFT, RIGHT},
	LEFT:  {UP, DOWN},
	RIGHT: {UP, DOWN},
}

var opposites = map[direction]direction{
	LEFT:  RIGHT,
	RIGHT: LEFT,
	UP:    DOWN,
	DOWN:  UP,
}

// priority queue - https://pkg.go.dev/container/heap
type queue []*node

func (q queue) Len() int           { return len(q) }
func (q queue) Less(i, j int) bool { return q[i].loss < q[j].loss }
func (q queue) Swap(i, j int)      { q[i], q[j] = q[j], q[i] }

func (q *queue) Push(x any) {
	*q = append(*q, x.(*node))
}

func (q *queue) Pop() any {
	old := *q
	n := len(old)
	item := old[n-1]
	old[n-1] = nil // avoid memory leak
	*q = old[0 : n-1]
	return item
}

func parseGrid(input []string) [][]node {
	var grid [][]node

	for y, row := range input {
		grid = append(grid, make([]node, len(row)))
		for x, c := range row {
			n, _ := strconv.Atoi(string(c))
			grid[y][x] = node{loss: n}
		}
	}

	return grid
}

func choices(g [][]node, cur *node, lm map[state]int, minMv, maxMv int) []*node {
	var choices []*node

	x, y := cur.state.p.x, cur.state.p.y

	for _, d := range directions {
		if x+d.x < 0 || x+d.x >= len(g[0]) || y+d.y < 0 || y+d.y >= len(g) { // out of bounds
			continue
		}

		// if we can no longer travel in the same direction and d is not a valid turn, continue
		if cur.state.mv == maxMv && !slices.Contains(turns[cur.state.d], d) || opposites[cur.state.d] == d {
			continue
		}

		nx, ny := x+d.x, y+d.y

		var nextMv int

		if d != cur.state.d {
			if cur.state.mv < minMv {
				continue // must travel in same direction if < minMv
			}
			nextMv = 1 // start over moves when turning
		} else {
			nextMv = cur.state.mv%maxMv + 1
		}

		nextState := state{p: point{nx, ny}, d: d, mv: nextMv}
		nextLoss := g[ny][nx].loss
		// if we've already been to the next point with the exact same direction and mv count, and it has a lower loss, abandon this path
		if ns, ok := lm[nextState]; ok && ns <= cur.loss+nextLoss {
			continue
		}

		lm[nextState] = cur.loss + nextLoss
		choices = append(choices, &node{state: nextState, loss: cur.loss + nextLoss})
	}

	return choices

}

func dijkstra(g [][]node, start, end point, minMv, maxMv int) int {
	start1 := state{p: start, d: RIGHT, mv: 0}
	start2 := state{p: start, d: DOWN, mv: 0}

	queue := queue{
		&node{loss: 0, state: start1},
		&node{loss: 0, state: start2},
	}

	// map of states to minimum loss
	lossMap := map[state]int{start1: 0, start2: 0}
	heap.Init(&queue)

	for len(queue) > 0 {
		cur := heap.Pop(&queue).(*node)
		// abandon path if there is already a lower loss path
		if lossMap[cur.state] < cur.loss {
			continue
		}

		// if at end and able to stop
		if cur.state.p == end && cur.state.mv >= minMv {
			return cur.loss
		}

		for _, c := range choices(g, cur, lossMap, minMv, maxMv) {
			heap.Push(&queue, c)
		}
	}

	return -1
}

func part1(input []string) int {
	grid := parseGrid(input)

	start := point{0, 0}
	end := point{len(grid) - 1, len(grid[0]) - 1}
	weight := dijkstra(grid, start, end, 0, 3)
	return weight
}

func part2(input []string) int {
	grid := parseGrid(input)
	start := point{0, 0}
	end := point{len(grid) - 1, len(grid[0]) - 1}
	weight := dijkstra(grid, start, end, 4, 10)
	return weight
}

func main() {

	input := parseInput()

	fmt.Println("Part 1:", part1(input))
	fmt.Println("Part 2:", part2(input))
}
