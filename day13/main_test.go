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
	is.Equal(result.valuePT1, 405)
	is.Equal(result.valuePT2, 400)
}

func TestSolution(t *testing.T) {
	is := is.New(t)
	scan := bufio.NewScanner(bytes.NewReader(input))

	result, err := run(scan)
	is.NoErr(err)

	t.Log(result)
	is.True(result.valuePT1 > 1704) // attempt 1
	is.Equal(result.valuePT1, 30705)
	is.Equal(result.valuePT2, 44615)
}
