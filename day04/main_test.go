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

	r, err := run(scan)
	is.NoErr(err)
	is.Equal(r.points, 13)
	is.Equal(r.cards, 30)
}

func TestSolution(t *testing.T) {
	is := is.New(t)
	scan := bufio.NewScanner(bytes.NewReader(input))

	r, err := run(scan)
	is.NoErr(err)
	is.Equal(r.points, 23235)
	is.Equal(r.cards, 5920640)
}
