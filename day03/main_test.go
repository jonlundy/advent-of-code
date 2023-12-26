package main

import (
	"bufio"
	"bytes"
	_ "embed"
	"testing"

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
	is.Equal(result.valuePT1, int(4361))

	is.Equal(result.valuePT2, 467835)
}

func TestInput(t *testing.T) {
	is := is.New(t)
	scan := bufio.NewScanner(bytes.NewReader(input))

	result, err := run(scan)
	is.NoErr(err)

	t.Log(result.valuePT1)
	is.Equal(result.valuePT1, int(553079))

	t.Log(result.valuePT2)
	is.Equal(result.valuePT2, 84363105)
}
