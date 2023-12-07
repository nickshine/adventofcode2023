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

var (
	cardValues      = map[rune]int{'2': 2, '3': 3, '4': 4, '5': 5, '6': 6, '7': 7, '8': 8, '9': 9, 'T': 10, 'J': 11, 'Q': 12, 'K': 13, 'A': 14}
	cardValuesJoker = map[rune]int{'J': 1, '2': 2, '3': 3, '4': 4, '5': 5, '6': 6, '7': 7, '8': 8, '9': 9, 'T': 10, 'Q': 11, 'K': 12, 'A': 13}
)

func parseType(h string, joker bool) int {
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
	case 3:
		for _, c := range counts {
			if c == 3 { // 23444, 22444
				if joker && j > 0 {
					return FOUROFKIND // 2J444 -> 24444, 22JJ3 -> 22223, 23JJJ -> 23333
				}
				return THREEOFKIND // 22234
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

func optimalCard(counts map[rune]int) rune {
	var maxCount int
	var card rune
	for r, count := range counts {
		if count >= maxCount {
			maxCount = count
			card = r
		}
	}

	return card
}

func parseTypeOptimised(h string, joker bool) int {
	counts := make(map[rune]int)

	for _, r := range h {
		counts[r]++
	}

	if joker {
		j := counts['J']
		if j > 0 {
			delete(counts, 'J')
			counts[optimalCard(counts)] += j
		}

	}

	// counts
	// len 1
	// AAAAA [A: 5] -> 1 (length of counts) -> FIVEOFKIND
	// len 2
	// 2222A -> [2: 4, A: 1] -> 2 -> FOUROFKIND
	// 222AA -> [2: 3, A: 2] -> 2 -> FULLHOUSE
	// len 3
	// 22334 -> [2: 2, 3: 2, 4: 1] -> TWOPAIR
	// 23444 -> [2: 1, 3: 1, 4: 3] -> THREEOFKIND
	// len 4
	// 23455 -> [2: 1, 3: 1, 4: 1, 5: 2] -> ONEPAIR
	// len 5
	// 23456 -> HIGHCARD
	strength := map[int]int{1: FIVEOFKIND, 2: FULLHOUSE, 3: TWOPAIR, 4: ONEPAIR, 5: HIGHCARD}

	for _, c := range counts {
		if c == 4 {
			strength[2] = FOUROFKIND
		}

		if c == 3 {
			strength[3] = THREEOFKIND
		}
	}

	return strength[len(counts)]

}

func parseHands(in []string, joker bool) []hand {

	var hands []hand

	for _, line := range in {
		parts := strings.Fields(line)
		v := parts[0]
		bid, _ := strconv.Atoi(parts[1])
		// typ := parseType(parts[0], joker)
		typ := parseTypeOptimised(parts[0], joker)

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
		total += h.bid * rank
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
		total += h.bid * rank
	}

	return total
}

func main() {

	input := parseInput()

	fmt.Println("Part 1:", part1(input))
	fmt.Println("Part 2:", part2(input))
}
