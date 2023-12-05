package main

import (
	"bufio"
	"bytes"
	"testing"

	_ "embed"

	"github.com/matryer/is"
)

//go:embed example.txt
var example []byte

func TestExample(t *testing.T) {
	is := is.New(t)
	scan := bufio.NewScanner(bytes.NewReader(example))

	points, cards := run(scan)
	is.Equal(points, 13)
	is.Equal(cards, 30)
}

func TestSolution(t *testing.T) {
	is := is.New(t)
	scan := bufio.NewScanner(bytes.NewReader(input))

	points, cards := run(scan)
	is.Equal(points, 23235)
	is.Equal(cards, 5920640)
}
