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

type edge struct {
	from, to string
}

func (e edge) String() string {
	return fmt.Sprintf("%s --> %s", e.from, e.to)
}

type graph struct {
	nodes map[string]map[edge]struct{}
}

func (g *graph) addEdge(from, to string) {
	if _, ok := g.nodes[from]; !ok {
		g.nodes[from] = make(map[edge]struct{})
	}
	g.nodes[from][edge{from, to}] = struct{}{}

	// undirected
	if _, ok := g.nodes[to]; !ok {
		g.nodes[to] = make(map[edge]struct{})
	}
	g.nodes[to][edge{to, from}] = struct{}{}
}

func (g *graph) removeEdge(e edge) {
	from, to := e.from, e.to
	fromEdges := g.nodes[from]
	toEdges := g.nodes[to]

	delete(fromEdges, e)
	// undirected
	delete(toEdges, edge{e.to, e.from})
}

func (g *graph) listEdges() []edge {
	var out []edge
	for _, edges := range g.nodes {
		for e := range edges {
			out = append(out, e)
		}
	}

	return out
}

func newGraph(input []string) *graph {
	g := &graph{
		nodes: make(map[string]map[edge]struct{}),
	}

	for _, line := range input {
		parts := strings.Split(line, ": ")
		from := parts[0]
		edges := strings.Split(parts[1], " ")
		for _, to := range edges {
			g.addEdge(from, to)
		}
	}

	return g
}

func findPath(g *graph, u, end string) []string {
	seen := make(map[string]struct{})

	var dfs func(u string) []string

	dfs = func(u string) []string {
		seen[u] = struct{}{}

		if u == end {
			return []string{u}
		}

		for e := range g.nodes[u] {
			next := e.to
			if _, ok := seen[next]; ok {
				continue
			}

			path := dfs(next)
			if len(path) > 0 {
				path = append(path, u)
				return path
			}
		}
		return nil
	}

	return dfs(u)
}

func findClusters(g *graph) [][]string {

	seen := make(map[string]struct{})
	var clusters [][]string

	var dfs func(u string)

	dfs = func(u string) {
		seen[u] = struct{}{}
		clusters[len(clusters)-1] = append(clusters[len(clusters)-1], u)

		for e := range g.nodes[u] {
			if _, ok := seen[e.to]; !ok {
				dfs(e.to)
			}
		}
	}

	for node := range g.nodes {
		if _, ok := seen[node]; !ok {
			clusters = append(clusters, []string{})
			dfs(node)
		}
	}

	return clusters
}

func part1(in []string) int {
	g := newGraph(in)
	gg := &graph{nodes: make(map[string]map[edge]struct{})}

	for _, e := range g.listEdges() {
		ng := newGraph(in)
		ng.removeEdge(e)

		safe := true
		for i := 0; i < 3; i++ {
			path := findPath(ng, e.from, e.to)
			if len(path) == 0 {
				safe = false // there are not 3 other paths
				break
			}
			// remove path so another can be found (3 total)
			for j := 0; j < len(path)-1; j++ {
				ng.removeEdge(edge{path[j], path[j+1]})
			}
		}

		if safe {
			gg.addEdge(e.from, e.to)
		}
	}

	clusters := findClusters(gg)
	fmt.Printf("Clusters: %+v\n", clusters)
	return len(clusters[0]) * len(clusters[1])
}

func part2(in []string) int {
	return 0
}

func main() {
	input := parseInput()

	fmt.Println("Part 1:", part1(input))
	fmt.Println("Part 2:", part2(input))
}
