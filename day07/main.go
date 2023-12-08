package main

import (
	"bufio"
	"fmt"
	"os"
	"sort"
)

func main() {
	if len(os.Args) != 2 {
		fmt.Fprintln(os.Stderr, "Usage: day07 FILE")
	}

	input, err := os.Open(os.Args[1])
	if err != nil {
		fmt.Fprintln(os.Stderr, err)
	}

	scan := bufio.NewScanner(input)

	score1, score2 := run(scan)

	fmt.Println("score 1", score1)
	fmt.Println("score 2", score2)
}

func run(scan *bufio.Scanner) (uint64, uint64) {
	var game Game

	for scan.Scan() {
		var cards string
		var bid int
		_, err := fmt.Sscanf(scan.Text(), "%s %d", &cards, &bid)
		if err != nil {
			panic(err)
		}

		fmt.Println("cards", cards, "bid", bid)
		game.plays = append(game.plays, Play{bid, []rune(cards), &game})
	}

	game.cardOrder = getOrder(cardTypes1)
	product1 := calcProduct(game)

	game.cardOrder = getOrder(cardTypes2)
	game.wildCard = 'J'
	product2 := calcProduct(game)

	return product1, product2
}

var cardTypes1 = []rune{'A', 'K', 'Q', 'J', 'T', '9', '8', '7', '6', '5', '4', '3', '2'}
var cardTypes2 = []rune{'A', 'K', 'Q', 'T', '9', '8', '7', '6', '5', '4', '3', '2', 'J'}

func calcProduct(game Game) uint64 {
	sort.Sort(game.plays)

	var product uint64
	for i, play := range game.plays {
		rank := i + 1
		fmt.Printf("play %d %s %d %s %x\n", rank, string(play.hand), play.bid, play.HandType(), play.HandStrength())
		product += uint64(play.bid * rank)
	}
	return product
}

func getOrder(cardTypes []rune) map[rune]int {
	cardOrder := make(map[rune]int, len(cardTypes))
	for i, r := range cardTypes {
		cardOrder[r] = len(cardTypes) - i
	}
	return cardOrder
}

type Game struct {
	plays     Plays
	cardOrder map[rune]int
	wildCard  rune
}

type Play struct {
	bid  int
	hand Hand
	game *Game
}

type Hand []rune

func (h Play) HandType() string {
	hc, _ := h.HighCard()
	switch {
	case h.IsFiveOfKind():
		return "5K-" + string(hc)
	case h.IsFourOfKind():
		return "4K-" + string(hc)
	case h.IsFullHouse():
		return "FH-" + string(hc)
	case h.IsThreeOfKind():
		return "3K-" + string(hc)
	case h.IsTwoPair():
		return "2P-" + string(hc)
	case h.IsOnePair():
		return "1P-" + string(hc)
	case h.IsHighCard():
		return "HC-" + string(hc)
	}
	return "Uno"
}

func (h Play) HandStrength() int {
	_, v := h.HighCard()
	switch {
	case h.IsFiveOfKind():
		return 0x700000 | v
	case h.IsFourOfKind():
		return 0x600000 | v
	case h.IsFullHouse():
		return 0x500000 | v
	case h.IsThreeOfKind():
		return 0x400000 | v
	case h.IsTwoPair():
		return 0x300000 | v
	case h.IsOnePair():
		return 0x200000 | v
	case h.IsHighCard():
		return 0x100000 | v
	}
	return 0
}

func (h Play) IsFiveOfKind() bool {
	_, _, _, _, has5 := h.game.hasSame(h.hand)
	return has5
}
func (h Play) IsFourOfKind() bool {
	_, _, _, has4, _ := h.game.hasSame(h.hand)
	return has4
}
func (h Play) IsFullHouse() bool {
	_, has2, has3, _, _ := h.game.hasSame(h.hand)
	return has3 && has2
}
func (h Play) IsThreeOfKind() bool {
	has1, _, has3, _, _ := h.game.hasSame(h.hand)
	return has3 && has1
}
func (h Play) IsTwoPair() bool {
	_, has2, has3, _, _ := h.game.hasSame(h.hand)
	return !has3 && has2 && h.game.pairs(h.hand) == 2
}
func (h Play) IsOnePair() bool {
	_, has2, has3, _, _ := h.game.hasSame(h.hand)
	return !has3 && has2 && h.game.pairs(h.hand) == 1
}
func (h Play) IsHighCard() bool {
	has1, has2, has3, has4, _ := h.game.hasSame(h.hand)
	return has1 && !has2 && !has3 && !has4
}
func (h Play) HighCard() (rune, int) {
	var i int
	pairs := make(Pairs, 5)
	cnt := h.game.Counts(h.hand)
	for r, c := range cnt {
		pairs[i].c = c
		pairs[i].r = r
		pairs[i].o = h.game.cardOrder[r]

		i++
	}

	sort.Sort(sort.Reverse(pairs))

	value := 0
	for i, r := range h.hand {
		if r == 0 {
			continue
		}

		value |= h.game.cardOrder[r] << (4 * (4 - i))
	}

	return pairs[0].r, value
}

type Pairs []struct {
	r rune
	c int
	o int
}

func (p Pairs) Len() int { return len(p) }
func (p Pairs) Less(i, j int) bool {
	if p[i].c == p[j].c {
		return p[i].o < p[j].o
	}
	return p[i].c < p[j].c
}
func (p Pairs) Swap(i, j int) { p[i], p[j] = p[j], p[i] }

type Plays []Play

func (p Plays) Len() int           { return len(p) }
func (p Plays) Less(i, j int) bool { return p[i].HandStrength() < p[j].HandStrength() }
func (p Plays) Swap(i, j int)      { p[i], p[j] = p[j], p[i] }

func (g *Game) Counts(cards []rune) map[rune]int {
	m := make(map[rune]int, len(g.cardOrder))
	for _, c := range cards {
		m[c]++
	}
	if g.wildCard != 0 && m[g.wildCard] > 0 {
		var maxK rune
		var maxV int
		for k, v := range m {
			if k != g.wildCard && v > maxV {
				maxK, maxV = k, v
			}
		}
		if maxK != 0 {
			m[maxK] += m[g.wildCard]
			delete(m, g.wildCard)
		}
	}
	return m
}
func (g *Game) hasSame(cards []rune) (has1, has2, has3, has4, has5 bool) {
	cnt := g.Counts(cards)
	for _, c := range cnt {
		switch c {
		case 1:
			has1 = true

		case 2:
			has2 = true

		case 3:
			has3 = true

		case 4:
			has4 = true

		case 5:
			has5 = true
		}
	}
	return
}
func (g *Game) pairs(cards []rune) int {
	pairs := 0
	for _, n := range g.Counts(cards) {
		if n == 2 {
			pairs++
		}
	}
	return pairs
}
