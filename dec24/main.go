package main

import (
	"fmt"
	"os"
	"regexp"
	"strconv"
	"strings"
)

var re = regexp.MustCompile(`^(\d+),\s+(\d+),\s+(\d+)\s+@\s+(-?\d+),\s+(-?\d+),\s+(-?\d+)`)

type hailstone struct {
	px, py, pz float64
	vx, vy, vz float64
}

func (h hailstone) String() string {
	return fmt.Sprintf("(%.0f,%.0f,%.0f), (%.0f,%.0f,%.0f)", h.px, h.py, h.pz, h.vx, h.vy, h.vz)
}

func parseInput() []hailstone {
	raw, _ := os.ReadFile(os.Args[1])
	lines := strings.Split(strings.Trim(string(raw), "\n"), "\n")

	var hailstones []hailstone

	for _, line := range lines {
		match := re.FindStringSubmatch(line)

		out := make([]float64, 6)

		for i := 0; i < len(match)-1; i++ {
			n, _ := strconv.Atoi(match[i+1])
			out[i] = float64(n)

		}
		px, py, pz, vx, vy, vz := out[0], out[1], out[2], out[3], out[4], out[5]
		hailstones = append(hailstones, hailstone{px, py, pz, vx, vy, vz})
	}

	return hailstones
}

func slope(vx, vy float64) float64 {
	if vx != 0 {
		return vy / vx
	}

	return 0
}

func intersectTime(a, b hailstone) float64 {
	// x(t) = x + vx * t
	// y(t) = y + vy * t
	// a.px + a.vx*t = b.px + b.vx*t
	// a.vx*t - b.vx*t = b.px - a.px
	// t(a.vx - b.vx) = b.px - a.px
	// t = (b.px - a.px) / (a.vx - b.vx)

	// if a.vx == b.vx {
	// 	return 0
	// }

	t := (b.px - a.px) / (a.vx - b.vx)
	return t
}

func part1(h []hailstone) int {

	const (
		MIN = 200000000000000
		MAX = 400000000000000
	)

	sum := 0
	for i := 0; i < len(h); i++ {
		for j := i + 1; j < len(h); j++ {
			// fmt.Printf("Hailstone A: %s\n", h[i])
			// fmt.Printf("Hailstone B: %s\n", h[j])

			a, b := h[i], h[j]

			// Line: y = mx + b
			// m = slope
			// b = y-intercept

			// slope = dy/dx
			ma := slope(a.vx, a.vy)
			mb := slope(b.vx, b.vy)
			// y-intercept
			ba := a.py - ma*a.px
			bb := b.py - mb*b.px

			// Line A: y = mx+b
			// Line B: y = mx+b
			// ma*x + ba = mb*x+bb
			// ma*x - mb*x = bb - ba
			// x*(ma - mb) = bb - ba
			x := (bb - ba) / (ma - mb)
			y := ma*x + ba
			// fmt.Printf("Intersection X: %.3f\n", x)
			// fmt.Printf("Intersection Y: %.3f\n", y)

			if x < MIN || x > MAX || y < MIN || y > MAX {
				continue
			}

			if x < a.px && a.vx > 0 || x > a.px && a.vx < 0 ||
				y < a.py && a.vy > 0 || y > a.py && a.vy < 0 || // A in the past
				x < b.px && b.vx > 0 || x > b.px && b.vx < 0 ||
				y < b.py && b.vy > 0 || y > b.py && b.vy < 0 { // B in the past
				continue
			}

			sum++
		}

		fmt.Println()
	}

	return sum
}

func part2(h []hailstone) int {

	// For each candidate rock velocity (-500 to 500), considering just the first two hailstones and only the x/y coordinates.
	// Assuming the rock intercepts both these hailstones, use the same approach as part 1 to calculate the interception point & time,
	// and from that calculate the start position of the rock. We then test this origin & velocity to see if it fits all x/y/z coordinates
	// for all hailstones.

	const RANGE = 500

	h1 := h[0]
	h2 := h[1]
	// ma := slope(h1.vx, h1.vy)
	// mb := slope(h2.vx, h2.vy)
	// // y-intercept
	// ba := h1.py - ma*h1.px
	// bb := h2.py - mb*h2.px

	for i := -RANGE; i <= RANGE; i++ {
		for j := -RANGE; j <= RANGE; j++ {
			for k := -RANGE; k <= RANGE; k++ {

				if i == 0 || j == 0 || k == 0 {
					continue
				}

				vx, vy, vz := float64(i), float64(j), float64(k)
				// simultaneous linear equation (from part 1):
				// H1:  x = A + a*t   y = B + b*t
				// H2:  x = C + c*u   y = D + d*u
				//
				//  t = [ d ( C - A ) - c ( D - B ) ] / ( a * d - b * c )
				//
				// Solve for origin of rock intercepting both hailstones in x,y:
				//     x = A + a*t - vx*t   y = B + b*t - vy*t
				//     x = C + c*u - vx*u   y = D + d*u - vy*u

				// x(t) = x + vx * t
				// y(t) = y + vy * t

				A, a := h1.px, h1.vx-vx
				B, b := h1.py, h1.vy-vy
				C, c := h2.px, h2.vx-vx
				D, d := h2.py, h2.vy-vy

				if c == 0 || (a*d)-(b*c) == 0 {
					continue
				}

				// rock intercepts H1 at time t
				t := (d*(C-A) - c*(D-B)) / ((a * d) - (b * c))

				// Calculate starting position of rock from intercept point
				x := h1.px + h1.vx*t - vx*t
				y := h1.py + h1.vy*t - vy*t
				z := h1.pz + h1.vz*t - vz*t

				// check all hailstones

				valid := true
				for _, hs := range h {

					var u float64
					if hs.vx != vx {
						u = (x - hs.px) / (hs.vx - vx)
					} else if hs.vy != vy {
						u = (y - hs.py) / (hs.vy - vy)
					} else if hs.vz != vz {
						u = (z - hs.pz) / (hs.vz - vz)
					} else {
						panic("invalid input")
					}

					if (x+u*vx != hs.px+u*hs.vx) || (y+u*vy != hs.py+u*hs.vy) || (z+u*vz != hs.pz+u*hs.vz) {
						valid = false
						break
					}
				}

				if valid {
					fmt.Printf("%0.0f,%0.0f,%0.0f, %.0f,%.0f,%.0f\n", x, y, z, vx, vy, vz)
					return int(x + y + z)
				}
			}
		}

	}

	return 0
}

func main() {
	input := parseInput()

	fmt.Println("Part 1:", part1(input))
	fmt.Println("Part 2:", part2(input))
}
