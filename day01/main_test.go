package main

import (
	"bufio"
	"bytes"
	"testing"

	_ "embed"

	"github.com/matryer/is"
)

//go:embed example1.txt
var example1 []byte

//go:embed example2.txt
var example2 []byte

//go:embed input.txt
var input []byte

func TestExample1(t *testing.T) {
	is := is.New(t)
	scan := bufio.NewScanner(bytes.NewReader(example1))

	result, err := run(scan)
	is.NoErr(err)

	t.Log(result)
	is.Equal(result.sum, 142)
}

func TestExample2(t *testing.T) {
	is := is.New(t)
	scan := bufio.NewScanner(bytes.NewReader(example2))

	result, err := run(scan)
	is.NoErr(err)

	t.Log(result)
	is.Equal(result.sum2, 281)
}
func TestInput(t *testing.T) {
	is := is.New(t)
	scan := bufio.NewScanner(bytes.NewReader(input))

	result, err := run(scan)
	is.NoErr(err)

	t.Log(result)
	is.Equal(result.sum, 54573)
	is.Equal(result.sum2, 54591)
}
