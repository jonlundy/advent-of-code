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
	is.Equal(result.valuePT1, 46)
	is.Equal(result.valuePT2, 51)
}

func TestSolution(t *testing.T) {
	is := is.New(t)
	scan := bufio.NewScanner(bytes.NewReader(input))

	result, err := run(scan)
	is.NoErr(err)

	t.Log(result)

	is.Equal(result.valuePT1, 8098)
	is.Equal(result.valuePT2, 8335)
}

func TestStack(t *testing.T) {
	is := is.New(t)

	s := stack[int]{}
	s.Push(5)
	is.Equal(s.Pop(), 5)
}