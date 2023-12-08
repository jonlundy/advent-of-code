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

//go:embed example3.txt
var example3 []byte

//go:embed input.txt
var input []byte

func TestExample1(t *testing.T) {
	is := is.New(t)
	scan := bufio.NewScanner(bytes.NewReader(example1))

	result, err := run(scan)
	is.NoErr(err)

	t.Log(result.stepsPT1)
	is.Equal(result.stepsPT1, 2)
}

func TestExample2(t *testing.T) {
	is := is.New(t)
	scan := bufio.NewScanner(bytes.NewReader(example2))

	result, err := run(scan)
	is.NoErr(err)

	t.Log(result.stepsPT1)
	is.Equal(result.stepsPT1, 6)
}

func TestExample3(t *testing.T) {
	is := is.New(t)
	scan := bufio.NewScanner(bytes.NewReader(example3))

	result, err := run(scan)
	is.NoErr(err)

	t.Log(result.stepsPT2)
	is.Equal(result.stepsPT2, uint64(6))
}

func TestInput(t *testing.T) {
	is := is.New(t)
	scan := bufio.NewScanner(bytes.NewReader(input))

	result, err := run(scan)
	is.NoErr(err)

	t.Log("part1 solution", result.stepsPT1)
	is.Equal(result.stepsPT1, 14429)
	t.Log("part2 solution", result.stepsPT2)
	is.Equal(result.stepsPT2, uint64(10921547990923))
}

// first: 14429
// second: 10921547990923