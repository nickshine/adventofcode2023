package main

import (
	"fmt"
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

	total := 0
	for _, race := range races {
		wins := 0
		for i := 1; i < race.time; i++ {
			// hold i, mv len(race.time) - i
			if i*(race.time-i) > race.distance {
				wins++
			}
		}
		if total == 0 {
			total = wins
		} else {
			total *= wins
		}

	}

	return total
}

func part2(input []string) int {
	race := parse2(input)
	fmt.Printf("Race: %#v\n", race)

	total := 0
	wins := 0
	for i := 1; i < race.time; i++ {
		if i*(race.time-i) > race.distance {
			wins++
		}
	}
	if total == 0 {
		total = wins
	} else {
		total *= wins
	}

	return total
}

func main() {

	input := parseInput()

	fmt.Println("Part 1:", part1(input))
	fmt.Println("Part 2:", part2(input))
}
