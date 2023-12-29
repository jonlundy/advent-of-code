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

//go:embed input.txt
var input []byte

func TestExample(t *testing.T) {
	is := is.New(t)
	scan := bufio.NewScanner(bytes.NewReader(example))

	result, err := run(scan)
	is.NoErr(err)

	t.Log(result)
	is.Equal(result.valuePT1, 136)
	is.Equal(result.valuePT2, 64)
}

func TestSolution(t *testing.T) {
	is := is.New(t)
	scan := bufio.NewScanner(bytes.NewReader(input))

	result, err := run(scan)
	is.NoErr(err)

	t.Log(result)
	is.True(result.valuePT2 < 87286) // first submission
	is.True(result.valuePT2 < 87292) // second submission
	is.True(result.valuePT2 < 87287) // third submission

	is.Equal(result.valuePT1, 110407)
	is.Equal(result.valuePT2, 87273)
}
