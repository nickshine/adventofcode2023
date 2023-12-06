package main

import (
	"fmt"
	"math"
	"os"
	"strconv"
	"strings"
	"sync"
)

func parseInput() []string {
	raw, _ := os.ReadFile(os.Args[1])
	return strings.Split(strings.Trim(string(raw), "\n"), "\n\n")
}

// conversion contains a source, destination, and range.
type conversion struct {
	src, dest, ran int
}

func (c conversion) String() string {
	return fmt.Sprintf("src: %d, dest: %d, range: %d", c.src, c.dest, c.ran)
}

type categoryMap struct {
	conversions []conversion
}

func (c *categoryMap) toDest(src int) int {

	for _, cv := range c.conversions {
		min, max := cv.src, cv.src+cv.ran
		if src >= min && src < max {

			dest := src + (cv.dest - cv.src)
			return dest
		}
	}

	return src
}

func parse(in []string) ([]int, []categoryMap) {
	var seeds []int
	for _, v := range strings.Fields(strings.TrimPrefix(in[0], "seeds: ")) {
		sv, _ := strconv.Atoi(v)
		seeds = append(seeds, sv)
	}

	var cmaps []categoryMap
	for _, cmap := range in[1:] {

		var cs categoryMap

		parts := strings.Split(cmap, "\n")

		for _, v := range parts[1:] {
			fields := strings.Fields(v)
			dest, _ := strconv.Atoi(fields[0])
			src, _ := strconv.Atoi(fields[1])
			ran, _ := strconv.Atoi(fields[2])
			cs.conversions = append(cs.conversions, conversion{src, dest, ran})
		}

		cmaps = append(cmaps, cs)
	}

	return seeds, cmaps
}

func part1(input []string) int {
	seeds, categoryMaps := parse(input)

	min := math.MaxInt32
	var src, dest int

	for _, seed := range seeds {
		src = seed
		for _, c := range categoryMaps {
			dest = c.toDest(src)
			src = dest
		}

		loc := dest

		if loc < min {
			min = loc
		}
	}

	return min
}

// go run main.go input.txt  152.89s user 0.47s system 99% cpu 2:33.38 total
func part2Slow(input []string) int {
	seedRanges, categoryMaps := parse(input)

	min := math.MaxInt32
	var src, dest int

	for i := 0; i < len(seedRanges); i += 2 {
		start, length := seedRanges[i], seedRanges[i+1]
		fmt.Printf("seedRange %d/%d, start: %d, length: %d\n", i, len(seedRanges), start, length)

		for seed := start; seed < start+length; seed++ {
			src = seed
			// fmt.Println("seed:", seed)
			for _, c := range categoryMaps {
				dest = c.toDest(src)
				src = dest
			}

			loc := dest

			if loc < min {
				min = loc
			}
		}
	}

	return min
}

func part2Concurrency(input []string) int {
	seedRanges, categoryMaps := parse(input)
	concurrency := 10

	seedChan := make(chan int)
	locations := make(chan int)

	worker := func(id int) {
		var src, dest int
		for seed := range seedChan {
			// fmt.Printf("worker %d: received seed: %d\n", id, seed)
			src = seed

			for _, c := range categoryMaps {
				dest = c.toDest(src)
				src = dest
			}

			locations <- dest
		}
	}

	var wg sync.WaitGroup
	// create goroutines
	for i := 0; i < concurrency; i++ {
		wg.Add(1)
		i := i
		go func() {
			defer wg.Done()
			worker(i)
		}()
	}

	// send work to the workers in another goroutine
	go func() {
		for i := 0; i < len(seedRanges); i += 2 {
			start, length := seedRanges[i], seedRanges[i+1]
			fmt.Printf("seedRange %d/%d, start: %d, length: %d\n", i, len(seedRanges), start, length)

			for seed := start; seed < start+length; seed++ {
				seedChan <- seed
			}
		}
		close(seedChan)
	}()

	go func() {
		wg.Wait()
		close(locations)
	}()

	min := math.MaxInt32

	for loc := range locations {
		if loc < min {
			min = loc
		}
	}

	return min
}

func main() {

	input := parseInput()

	fmt.Println("Part 1:", part1(input))
	fmt.Println("Part 2:", part2Slow(input))
	fmt.Println("Part 2:", part2Concurrency(input))
}
