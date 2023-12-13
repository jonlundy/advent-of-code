package main

import (
	"bufio"
	_ "embed"
	"log/slog"
	"strconv"
	"strings"

	aoc "go.sour.is/advent-of-code-2023"
)

//go:embed input.txt
var input []byte

type card struct {
	card    int
	winner  []int
	scratch []int
	copies  int
}

func main() { aoc.MustResult(aoc.Runner(run)) }

type result struct {
	points int
	cards  int
}

func run(scan *bufio.Scanner) (result, error) {
	cards := []*card{}

	for scan.Scan() {
		pfx, text, ok := strings.Cut(scan.Text(), ":")

		if !ok || !strings.HasPrefix(pfx, "Card ") {
			continue
		}

		num := aoc.Atoi(strings.TrimSpace(strings.SplitN(pfx, " ", 2)[1]))
		cards = append(cards, &card{card: num})
		buf := make([]rune, 0, 4)
		winner := true

		for _, a := range text {
			if a >= '0' && a <= '9' {
				buf = append(buf, a)
				continue
			}
			if a == ' ' && len(buf) > 0 {
				num, _ = strconv.Atoi(string(buf))
				buf = buf[:0]
				if winner {
					cards[len(cards)-1].winner = append(cards[len(cards)-1].winner, num)
				} else {
					cards[len(cards)-1].scratch = append(cards[len(cards)-1].scratch, num)
				}
			}
			if a == '|' {
				winner = false
			}
		}

		if len(buf) > 0 {
			num = aoc.Atoi(string(buf))
			buf = buf[:0]
			_ = buf // ignore
			cards[len(cards)-1].scratch = append(cards[len(cards)-1].scratch, num)
		}
	}

	var sumPoints int
	var sumCards int

	for _, card := range cards {
		m := []int{}

		for _, w := range card.winner {
			for _, s := range card.scratch {
				if w == s {
					m = append(m, w)
				}
			}
		}

		if len(m) > 0 {
			sumPoints += 1 << (len(m) - 1)
			for i, c := range cards[card.card:min(card.card+len(m), len(cards))] {
				c.copies += 1 + 1*card.copies
				slog.Debug("cards", "card", card.card, "wins", i+1, "with", card.copies, "copies, wins", len(m), "giving", 1+1*card.copies, "for card", c.card, "total", c.copies)
			}
		}

		sumCards += 1 + 1*card.copies
		slog.Debug("cards", "card", card.card, "wins", len(m), "as", 1, "+", 1*card.copies, "for", 1+1*card.copies, "copies. total", sumCards)
		slog.Debug("points", "card", card.card, "match", m, "score", sumPoints)
	}

	return result{sumPoints, sumCards}, nil
}
