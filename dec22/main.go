package main

import (
	"fmt"
	"os"
	"slices"
	"strconv"
	"strings"
)

type grid [][][]int

type brick struct {
	id    int
	start [3]int
	end   [3]int
}

func parseInput() (grid, []brick) {
	raw, _ := os.ReadFile(os.Args[1])
	lines := strings.Split(strings.Trim(string(raw), "\n"), "\n")

	max := 0

	var bricks []brick
	for i, l := range lines {
		parts := strings.Split(l, "~")
		startRaw := strings.Split(parts[0], ",")
		endRaw := strings.Split(parts[1], ",")
		sx, _ := strconv.Atoi(startRaw[0])
		sy, _ := strconv.Atoi(startRaw[1])
		sz, _ := strconv.Atoi(startRaw[2])
		ex, _ := strconv.Atoi(endRaw[0])
		ey, _ := strconv.Atoi(endRaw[1])
		ez, _ := strconv.Atoi(endRaw[2])

		start := [3]int{sx, sy, sz}
		end := [3]int{ex, ey, ez}

		bricks = append(bricks, brick{i, start, end})

		if ex > max {
			max = sx
		}
		if ey > max {
			max = ey
		}
		if ez > max {
			max = ez
		}
	}

	grid := make(grid, max+1)
	for x := range grid {
		grid[x] = make([][]int, max+1)
		for y := range grid[x] {
			grid[x][y] = make([]int, max+1)
			for z := range grid[x][y] {
				grid[x][y][z] = -1
			}
		}
	}
	grid.fillBricks(bricks)
	return grid, bricks
}

func (g grid) displayXZ() {
	fmt.Println("  x  ")

	n := len(g)

	for z := range g {
		for x := range g {
			// for each y layer, find the first brick if any
			v := -1
			for y := range g {
				if g[x][y][n-1-z] >= 0 {
					v = g[x][y][n-1-z]
					break
				}
			}
			if v == -1 {
				fmt.Printf(".")
			} else {
				fmt.Printf("%c", 'A'+v)
			}
		}
		fmt.Println(" ", n-1-z)
	}
}

func (g grid) displayYZ() {
	fmt.Println("  y  ")

	n := len(g)

	for z := range g {
		for y := range g {
			// for each x layer, find the first brick if any
			v := -1
			for x := range g {
				if g[x][y][n-1-z] >= 0 {
					v = g[x][y][n-1-z]
					break
				}
			}
			if v == -1 {
				fmt.Printf(".")
			} else {
				fmt.Printf("%c", 'A'+v)
			}
		}
		fmt.Println(" ", n-1-z)
	}
}

func (g grid) fillBrick(b brick, id int) {
	sx, ex := b.start[0], b.end[0]
	sy, ey := b.start[1], b.end[1]
	sz, ez := b.start[2], b.end[2]

	dx := ex - sx
	dy := ey - sy
	dz := ez - sz

	g[sx][sy][sz] = id

	for d := 1; d <= dx; d++ {
		g[sx+d][sy][sz] = id
	}
	for d := 1; d <= dy; d++ {
		g[sx][sy+d][sz] = id
	}
	for d := 1; d <= dz; d++ {
		g[sx][sy][sz+d] = id
	}

}

func (g grid) fillBricks(bricks []brick) {
	for i, b := range bricks {
		g.fillBrick(b, i)
	}
}

func (g grid) removeBrick(b brick) {
	sx, ex := b.start[0], b.end[0]
	sy, ey := b.start[1], b.end[1]
	sz, ez := b.start[2], b.end[2]

	dx := ex - sx
	dy := ey - sy
	dz := ez - sz

	g[sx][sy][sz] = -1

	for d := 1; d <= dx; d++ {
		g[sx+d][sy][sz] = -1
	}
	for d := 1; d <= dy; d++ {
		g[sx][sy+d][sz] = -1
	}
	for d := 1; d <= dz; d++ {
		g[sx][sy][sz+d] = -1
	}
}

func (g grid) canFall(b brick) bool {
	sx, sy, sz := b.start[0], b.start[1], b.start[2]
	ex, ey := b.end[0], b.end[1]

	dx := ex - sx
	dy := ey - sy

	if sz == 1 { // on ground
		return false
	}

	if g[sx][sy][sz-1] > -1 { // another brick below
		return false
	}

	for d := 1; d <= dx; d++ { // need to check z-1 on all x
		if g[sx+d][sy][sz-1] > -1 {
			return false
		}
	}
	for d := 1; d <= dy; d++ { // need to check z-1 on all y
		if g[sx][sy+d][sz-1] > -1 {
			return false
		}
	}

	return true

}

// fall returns a set of the bricks that fell
func (g grid) fall(bricks []brick) {

	for {
		done := true
		slices.SortFunc[[]brick](bricks, func(a, b brick) int {
			return a.start[2] - b.start[2]
		})

		for i, b := range bricks {
			g.removeBrick(b)
			sx, sy, sz := b.start[0], b.start[1], b.start[2]
			ex, ey := b.end[0], b.end[1]

			if sz > 1 {
				var stop bool
				for j := 0; j < i; j++ { // for each brick below i
					stop = false
					bb := bricks[j]
					bsx, bsy := bb.start[0], bb.start[1]
					bex, bey, bez := bb.end[0], bb.end[1], bb.end[2]
					if bez == sz-1 { // if brick is just below
						if sx <= bex && ex >= bsx && sy <= bey && ey >= bsy {
							stop = true
							break
						}
					}
				}
				if !stop {
					done = false
					bricks[i].start[2]--
					bricks[i].end[2]--
				}
			}
		}
		if done {
			break
		}
	}
	g.fillBricks(bricks)
}

func fall(bricks []brick) {

	for {
		done := true
		slices.SortFunc[[]brick](bricks, func(a, b brick) int {
			return a.start[2] - b.start[2]
		})

		for i, b := range bricks {
			sx, sy, sz := b.start[0], b.start[1], b.start[2]
			ex, ey := b.end[0], b.end[1]

			if sz > 1 {
				var stop bool
				for j := 0; j < i; j++ { // for each brick below i
					stop = false
					bb := bricks[j]
					bsx, bsy := bb.start[0], bb.start[1]
					bex, bey, bez := bb.end[0], bb.end[1], bb.end[2]
					if bez == sz-1 { // if brick is just below
						if sx <= bex && ex >= bsx && sy <= bey && ey >= bsy {
							stop = true
							break
						}
					}
				}
				if !stop {
					done = false
					bricks[i].start[2]--
					bricks[i].end[2]--
				}
			}
		}
		if done {
			break
		}
	}
}

func (g grid) disintegrate(idx int, bricks []brick) int {
	brick := bricks[idx]
	g.removeBrick(brick) // remove b temporarily

	result := 0

	for i, b := range bricks {
		if i == idx { // skip the removed brick
			continue
		}

		if g.canFall(b) {
			result++
		}
	}

	g.fillBrick(brick, idx) // put brick back
	return result
}

func part1() int {
	grid, bricks := parseInput()
	// grid.displayXZ()
	// fmt.Println()
	// grid.displayYZ()

	grid.fall(bricks)

	sum := 0
	for i := range bricks {
		n := grid.disintegrate(i, bricks)
		if n == 0 {
			sum++
		}
	}

	return sum
}

func part2() int {
	_, bricks := parseInput()
	fall(bricks)

	total := 0
	for i := range bricks {
		copiedBricks := make([]brick, len(bricks))
		copy(copiedBricks, bricks)
		copiedBricks = append(copiedBricks[:i], copiedBricks[i+1:]...)
		fallBricks := make([]brick, len(copiedBricks))
		copy(fallBricks, copiedBricks)
		fall(fallBricks)

		for j := range copiedBricks {
			for k := range fallBricks {
				if copiedBricks[j].id == fallBricks[k].id {
					if copiedBricks[j].start != fallBricks[k].start || copiedBricks[j].end != fallBricks[k].end {
						total++
					}
					break
				}
			}
		}
	}

	return total
}

func copyGrid(g grid) grid {

	n := len(g)
	c := make(grid, n)

	for x := range c {
		c[x] = make([][]int, n)
		for y := range c[x] {
			c[x][y] = make([]int, n)
			for z := range c[x][y] {
				c[x][y][z] = g[x][y][z]
			}
		}
	}

	return c
}

func main() {
	fmt.Println("Part 1:", part1())
	fmt.Println("Part 2:", part2())
}

// 2287 too low
