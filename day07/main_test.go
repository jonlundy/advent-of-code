package main

import (
	"bufio"
	"bytes"
	"fmt"
	"testing"

	_ "embed"

	"github.com/matryer/is"
)

// AKQJT98765432
// dcba987654321

//go:embed example.txt
var example []byte

//go:embed input.txt
var input []byte

func TestHands(t *testing.T) {
	is := is.New(t)

	var game Game
	game.cardOrder = getOrder(cardTypes1)

	h := Play{hand: []rune("AAA23"), game: &game}
	// h.generateCounts()
	is.Equal(h.HandType(), "3K-A")

	h = Play{hand: []rune("JJJJJ"), game:&game}
	h.generateCounts()

	is.Equal(h.HandType(), "5K-J")
	is.Equal(fmt.Sprintf("%x", h.HandStrength()), "7aaaaa")

	h = Play{hand: []rune("KKKKJ"), game: &game}
	is.Equal(h.HandType(), "4K-K")
	is.Equal(fmt.Sprintf("%x", h.HandStrength()), "6cccca")

	h = Play{hand: []rune("QQQJA"), game: &game}
	is.Equal(h.HandType(), "3K-Q")
	is.Equal(fmt.Sprintf("%x", h.HandStrength()), "4bbbad")
}

func TestPower(t *testing.T) {
	for i := 1; i <= 13; i++ {
		for j := 100; j < 800; j += 100 {
			t.Log(i, j, i+j)
		}
	}
}

func TestExample(t *testing.T) {
	is := is.New(t)
	scan := bufio.NewScanner(bytes.NewReader(example))

	score1, score2 := run(scan)
	is.Equal(score1, uint64(6440))
	is.Equal(score2, uint64(5905))
}

func TestSolution(t *testing.T) {
	is := is.New(t)
	scan := bufio.NewScanner(bytes.NewReader(input))

	score1, score2 := run(scan)
	t.Log("score1", score1)
	is.Equal(score1, uint64(248559379))

	t.Log("score2", score2)
	is.Equal(score2, uint64(249631254))
}
