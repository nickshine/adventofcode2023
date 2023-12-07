package main

import (
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

type hand struct {
	v   string
	typ int
	bid int
}

const (
	HIGHCARD int = iota
	ONEPAIR
	TWOPAIR
	THREEOFKIND
	FULLHOUSE
	FOUROFKIND
	FIVEOFKIND
)

var cardValues = map[rune]int{'2': 2, '3': 3, '4': 4, '5': 5, '6': 6, '7': 7, '8': 8, '9': 9, 'T': 10, 'J': 11, 'Q': 12, 'K': 13, 'A': 14}
var cardValuesJoker = map[rune]int{'J': 1, '2': 2, '3': 3, '4': 4, '5': 5, '6': 6, '7': 7, '8': 8, '9': 9, 'T': 10, 'Q': 11, 'K': 12, 'A': 13}

func parseType(in string, joker bool) int {
	h := []rune(in)
	counts := make(map[rune]int)

	for _, r := range h {
		counts[r]++
	}

	j := counts['J']

	switch len(counts) {
	case 1:
		return FIVEOFKIND // 22222
	case 2:
		if joker && j > 0 {
			return FIVEOFKIND // 222JJ -> 22222
		}
		for _, c := range counts {
			if c == 4 { // 22223
				return FOUROFKIND
			}
		}
		return FULLHOUSE
	case 3: // TWOPAIR or THREEOFKIND // 22334, 22234,
		for _, c := range counts {
			if c == 3 {
				if joker && j > 0 { // 2233J -> 22333
					return FOUROFKIND
				}
				return THREEOFKIND
			}
		}
		if joker && j > 1 { // JJ233 -> 33233
			return FOUROFKIND
		} else if joker && j > 0 { // J2333 -> 22333
			return FULLHOUSE
		}
		return TWOPAIR
	case 4:
		if joker && j > 0 { // 234JJ -> 23444
			return THREEOFKIND
		}
		return ONEPAIR // 22345
	case 5:
		if joker && j > 0 { // 2345J -> 23455
			return ONEPAIR
		}
		return HIGHCARD // 23456
	}
	return 0
}

func parseHands(in []string, joker bool) []hand {

	var hands []hand

	for _, line := range in {
		parts := strings.Fields(line)
		v := parts[0]
		bid, _ := strconv.Atoi(parts[1])
		typ := parseType(parts[0], joker)

		hands = append(hands, hand{
			v:   v,
			typ: typ,
			bid: bid,
		})
	}

	return hands
}

func sortHands(h []hand, joker bool) {
	cards := cardValues
	if joker {
		cards = cardValuesJoker
	}
	slices.SortFunc(h, func(a, b hand) int {
		d := a.typ - b.typ
		switch {
		case d == 0:
			for i := 0; i < len(a.v); i++ {
				d = cards[rune(a.v[i])] - cards[rune(b.v[i])]
				switch {
				case d == 0:
					continue
				case d > 0:
					return 1
				case d < 0:
					return -1
				}
			}

			return 0
		case d > 0:
			return 1
		default:
			return -1
		}
	})
}

func part1(input []string) int {

	hands := parseHands(input, false)
	sortHands(hands, false)
	// fmt.Printf("sorted hands: %#v\n", hands)

	total := 0
	for i, h := range hands {
		rank := i + 1
		total += (h.bid * rank)
	}

	return total
}

func part2(input []string) int {
	hands := parseHands(input, true)
	sortHands(hands, true)
	// fmt.Printf("sorted hands: %#v\n", hands)

	total := 0
	for i, h := range hands {
		rank := i + 1
		total += (h.bid * rank)
	}

	return total
}

func main() {

	input := parseInput()

	fmt.Println("Part 1:", part1(input))
	fmt.Println("Part 2:", part2(input))
}

// 251287184 too high
// 250005367 too low
// 249797187
// 249884506
// 249884049
