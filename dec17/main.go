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
	loc point
	d   direction
	mv  int
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

type grid [][]int

func (g grid) inbounds(x, y int) bool {
	if len(g) == 0 {
		return false
	}

	return x >= 0 && x < len(g[0]) && y >= 0 && y < len(g)
}

func parseGrid(input []string) [][]int {
	var grid [][]int

	for y, row := range input {
		grid = append(grid, make([]int, len(row)))
		for x, c := range row {
			n, _ := strconv.Atoi(string(c))
			grid[y][x] = n
		}
	}

	return grid
}

// choices filters the available next path steps based on current direction and move count.
func choices(g grid, cur *node, states map[state]int, minMv, maxMv int) []*node {
	var choices []*node

	x, y := cur.state.loc.x, cur.state.loc.y

	for _, d := range directions {
		if !g.inbounds(x+d.x, y+d.y) {
			continue
		}

		// if we can no longer travel in direction d, or d is not a valid turn, continue
		if cur.state.mv == maxMv && !slices.Contains(turns[cur.state.d], d) || opposites[cur.state.d] == d {
			continue
		}

		// must travel minMv steps in same direction before turning
		if d != cur.state.d && cur.state.mv < minMv {
			continue
		}

		nx, ny := x+d.x, y+d.y

		nextMv := 1 // assumes a turn - start over moves when turning
		if d == cur.state.d {
			nextMv = cur.state.mv%maxMv + 1
		}

		nextState := state{loc: point{nx, ny}, d: d, mv: nextMv}
		nextLoss := g[ny][nx]
		// if next loc is visited with the exact same direction and mv count, and it has a lower loss, abandon this path
		if loss, ok := states[nextState]; ok && loss <= cur.loss+nextLoss {
			continue
		}

		states[nextState] = cur.loss + nextLoss
		choices = append(choices, &node{state: nextState, loss: cur.loss + nextLoss})
	}

	return choices

}

func dijkstra(g grid, start, end point, minMv, maxMv int) int {
	start1 := state{loc: start, d: RIGHT, mv: 0}
	start2 := state{loc: start, d: DOWN, mv: 0}

	queue := queue{
		&node{loss: 0, state: start1},
		&node{loss: 0, state: start2},
	}

	// map of states to minimum loss
	states := map[state]int{start1: 0, start2: 0}
	heap.Init(&queue)

	for len(queue) > 0 {
		cur := heap.Pop(&queue).(*node)
		// abandon path if there is already a lower loss path
		if states[cur.state] < cur.loss {
			continue
		}

		// if at end and able to stop
		if cur.state.loc == end && cur.state.mv >= minMv {
			return cur.loss
		}

		for _, c := range choices(g, cur, states, minMv, maxMv) {
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
