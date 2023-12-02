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

func part1(input []string) int {
	var inputRE = regexp.MustCompile(`^Game (\d+): (.*)$`)
	bag := cubes{red: 12, green: 13, blue: 14}
	var sum int

OUTER:
	for _, line := range input {
		res := inputRE.FindStringSubmatch(line)
		id, _ := strconv.Atoi(res[1])
		setstring := res[2]
		cubeSetStrings := strings.Split(setstring, ";")
		for _, s := range cubeSetStrings {
			set := strings.Split(s, ",")
			draw := cubes{}
			for _, sv := range set {
				sv = strings.TrimSpace(sv)
				svs := strings.Split(sv, " ")
				if len(svs) != 2 {
					panic("bad input")
				}
				amount, color := svs[0], svs[1]
				count, _ := strconv.Atoi(amount)

				switch color {
				case "red":
					draw.red = count
				case "blue":
					draw.blue = count
				case "green":
					draw.green = count

				default:
					panic("invalid input")
				}
			}

			if bag.red < draw.red || bag.green < draw.green || bag.blue < draw.blue {
				// game not possible
				continue OUTER
			}
		}
		sum += id
	}

	return sum
}

func part2(input []string) int {
	var inputRE = regexp.MustCompile(`^Game (\d+): (.*)$`)
	var sum int

	for _, line := range input {
		res := inputRE.FindStringSubmatch(line)
		setstring := res[2]
		cubeSetStrings := strings.Split(setstring, ";")

		// for each cubeset, find the max of each color - that is the minimum bag
		maxR, maxG, maxB := 0, 0, 0
		for _, s := range cubeSetStrings {
			set := strings.Split(s, ",")
			draw := cubes{}
			for _, sv := range set {
				sv = strings.TrimSpace(sv)
				svs := strings.Split(sv, " ")
				if len(svs) != 2 {
					panic("bad input")
				}
				amount, color := svs[0], svs[1]
				count, _ := strconv.Atoi(amount)

				switch color {
				case "red":
					draw.red = count
					if count > maxR {
						maxR = count
					}
				case "blue":
					draw.blue = count
					if count > maxB {
						maxB = count
					}
				case "green":
					draw.green = count
					if count > maxG {
						maxG = count
					}

				default:
					panic("invalid input")
				}
			}
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
