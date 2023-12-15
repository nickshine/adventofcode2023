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

type platform [][]rune

func parsePlatform(in []string) platform {

	platform := make([][]rune, len(in))

	for i, line := range in {
		platform[i] = make([]rune, len(line))
		for j, col := range line {
			platform[i][j] = col
		}
	}

	return platform
}

func (p platform) String() string {
	var sb strings.Builder
	for _, row := range p {
		for _, c := range row {
			sb.WriteRune(c)
		}
		sb.WriteRune('\n')
	}

	return sb.String()
}

func (p platform) calcLoad() int {
	sum := 0
	for y, row := range p {
		for _, col := range row {
			if col == 'O' {
				sum += len(p) - y
			}
		}
	}

	return sum
}

func (p platform) tilt() {
	for x := 0; x < len(p[0]); x++ {
		fillIdx := 0
		for y := 0; y < len(p) && fillIdx < len(p); {
			switch p[fillIdx][x] {
			case 'O', '#':
				fillIdx++
				y = fillIdx
				continue
			default:
			}

			switch p[y][x] {
			case '#':
				fillIdx = y + 1
				y++
			case '.':
				if fillIdx >= y {
					fillIdx = y
				}
				y++
			case 'O':
				if fillIdx >= 0 && fillIdx < y {
					p[fillIdx][x] = 'O' // slide the rock
					p[y][x] = '.'
					fillIdx++
					continue
				}
				y++
			default:
				panic("invalid input")
			}
		}
	}
}

func rotate(p platform) platform {
	transposed := make(platform, len(p[0]))

	for y := 0; y < len(p[0]); y++ {
		transposed[y] = make([]rune, len(p))
		for x := 0; x < len(p); x++ {
			transposed[y][x] = p[len(p)-x-1][y]
		}
	}

	return transposed
}

func part1(input []string) int {

	platform := parsePlatform(input)
	platform.tilt()
	fmt.Println(platform)

	return platform.calcLoad()
}

func part2(input []string) int {
	p := parsePlatform(input)
	const N = 1000000000
	hm := make(map[string][]int) // map of platform strings to indexes of occurence

	for i := 0; i < N; i++ {
		s := p.String()
		hm[s] = append(hm[s], i)
		if h, ok := hm[s]; ok {
			if len(h) >= 2 && (N-i)%(i-h[0]) == 0 {
				break
			}
		}

		// do a cycle
		for j := 0; j < 4; j++ {
			p.tilt()
			p = rotate(p)
		}
	}

	return p.calcLoad()
}

func main() {

	input := parseInput()

	fmt.Println("Part 1:", part1(input))
	fmt.Println("Part 2:", part2(input))
}
