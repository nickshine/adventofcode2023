package main

import (
	"fmt"
	"math"
	"os"
	"strconv"
	"strings"
)

func parseInput() []string {
	raw, _ := os.ReadFile(os.Args[1])
	return strings.Split(strings.Trim(string(raw), "\n"), "\n")
}

type race struct {
	time     int
	distance int
}

func parse(in []string) []race {

	var races []race

	times := strings.Fields(in[0])[1:]
	dists := strings.Fields(in[1])[1:]

	for i := 0; i < len(times); i++ {

		time, _ := strconv.Atoi(times[i])
		dist, _ := strconv.Atoi(dists[i])

		races = append(races, race{time, dist})
	}

	return races
}

func parse2(in []string) race {

	time := strings.Join(strings.Fields(in[0])[1:], "")
	dist := strings.Join(strings.Fields(in[1])[1:], "")
	t, _ := strconv.Atoi(time)
	d, _ := strconv.Atoi(dist)

	return race{t, d}
}

func part1(input []string) int {

	races := parse(input)

	total := 1
	for _, race := range races {
		wins := 0
		for i := 1; i < race.time; i++ {
			// hold i, mv len(race.time) - i
			if i*(race.time-i) > race.distance {
				wins++
			}
		}
		total *= wins

	}

	return total
}

func part2(input []string) int {
	race := parse2(input)

	wins, total := 0, 1
	for i := 1; i < race.time; i++ {
		if i*(race.time-i) > race.distance {
			wins++
		}
	}
	total *= wins

	return total
}

func part2Optimised(input []string) int {
	race := parse2(input)

	// hold x, mv time-x = d

	// x * (t-x) > d
	// (x * (t-x) - d = 0)
	// x = -(sqrt(t^2-4*d)-t)/2, x = (sqrt(t^2-4*d)+t)/2
	var a, b, c, discriminant float64

	a = 1
	b = float64(race.time)
	c = float64(race.distance)

	discriminant = (b * b) - (4 * a * c)

	x0 := -(math.Sqrt(discriminant) - b) / (2 * a)
	x1 := (math.Sqrt(discriminant) + b) / (2 * a)
	// l0 := int(x0)
	// l1 := l0 + 1
	// r0 := int(x1)
	// r1 := r0 + 1

	// fmt.Printf("x0: %f\n", x0)
	// fmt.Printf("x1: %f\n", x1)

	// f := func(x int) int {
	// 	return x*(race.time-x) - race.distance
	// }

	// fmt.Printf("solve for %d: %d\n", l0, f(l0))
	// fmt.Printf("solve for %d: %d\n", l1, f(l1))
	// fmt.Printf("solve for %d: %d\n", r0, f(r0))
	// fmt.Printf("solve for %d: %d\n", r1, f(r1))
	return int(x1) - int(x0)
}

func main() {

	input := parseInput()

	fmt.Println("Part 1:", part1(input))
	fmt.Println("Part 2:", part2(input))
	fmt.Println("Part 2:", part2Optimised(input))
}
