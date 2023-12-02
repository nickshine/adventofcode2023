package main

import (
	"fmt"
	"os"
	"regexp"
	"strconv"
	"strings"
)

func parseInput() []string {
	raw, _ := os.ReadFile(os.Args[1])
	return strings.Split(strings.Trim(string(raw), "\n"), "\n")
}

type cubes struct {
	red   int
	green int
	blue  int
}

type switchFunc func(in cubes, count int, color string) cubes

var inputRE = regexp.MustCompile(`^Game (\d+): (.*)$`)

func parseSet(in string, f switchFunc) cubes {
	var out cubes
	set := strings.Split(in, ",")
	for _, sv := range set {
		sv = strings.TrimSpace(sv)
		svs := strings.Split(sv, " ")
		if len(svs) != 2 {
			panic("bad input")
		}
		amount, color := svs[0], svs[1]
		count, _ := strconv.Atoi(amount)

		out = f(out, count, color)
	}

	return out
}

func isPossible(setstring string, bag cubes) bool {
	sets := strings.Split(setstring, ";")
	for _, s := range sets {
		set := parseSet(s, func(c cubes, count int, color string) cubes {
			switch color {
			case "red":
				c.red = count
			case "blue":
				c.blue = count
			case "green":
				c.green = count

			default:
				panic("invalid input")
			}

			return c
		})

		if bag.red < set.red || bag.green < set.green || bag.blue < set.blue {
			return false
		}
	}

	return true
}

func part1(input []string) int {
	bag := cubes{red: 12, green: 13, blue: 14}
	var sum int

	for _, line := range input {
		res := inputRE.FindStringSubmatch(line)
		id, _ := strconv.Atoi(res[1])
		if isPossible(res[2], bag) {
			sum += id
		}
	}

	return sum
}

func part2(input []string) int {
	var sum int

	for _, line := range input {
		res := inputRE.FindStringSubmatch(line)
		setstring := res[2]
		sets := strings.Split(setstring, ";")

		// for each cubeset, find the max of each color - that is the minimum bag
		maxR, maxG, maxB := 0, 0, 0
		for _, s := range sets {
			parseSet(s, func(c cubes, count int, color string) cubes {
				switch color {
				case "red":
					c.red = count
					if count > maxR {
						maxR = count
					}
				case "blue":
					c.blue = count
					if count > maxB {
						maxB = count
					}
				case "green":
					c.green = count
					if count > maxG {
						maxG = count
					}

				default:
					panic("invalid input")
				}
				return c
			})
		}
		sum += maxR * maxG * maxB
	}

	return sum
}

func main() {

	input := parseInput()

	fmt.Println("Part 1:", part1(input))
	fmt.Println("Part 2:", part2(input))
}
