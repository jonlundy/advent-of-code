package main

import (
	"bufio"
	"bytes"
	_ "embed"
	"fmt"
	"log/slog"
	"os"
	"strconv"
	"strings"
)

//go:embed input.txt
var input []byte

type card struct {
	card    int
	winner  []int
	scratch []int
	copies  int
}

func main() {
	var level slog.Level
	if err := level.UnmarshalText([]byte(os.Getenv("DEBUG_LEVEL"))); err == nil && level != 0 {
		slog.SetDefault(slog.New(slog.NewTextHandler(os.Stderr, &slog.HandlerOptions{Level: level})))
	}

	buf := bytes.NewReader(input)
	scan := bufio.NewScanner(buf)

	points, cards := run(scan)

	fmt.Println("points:", points)
	fmt.Println("cards:", cards)
}

func run(scan *bufio.Scanner) (int, int) {
	cards := []*card{}

	for scan.Scan() {
		pfx, text, ok := strings.Cut(scan.Text(), ":")

		if !ok || !strings.HasPrefix(pfx, "Card ") {
			continue
		}

		num, _ := strconv.Atoi(strings.TrimSpace(strings.SplitN(pfx, " ", 2)[1]))
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
			num, _ = strconv.Atoi(string(buf))
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

	return sumPoints, sumCards
}
