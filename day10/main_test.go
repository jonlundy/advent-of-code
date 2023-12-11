package main

import (
	"bufio"
	"bytes"
	_ "embed"
	"testing"

	"github.com/matryer/is"
)

//go:embed example1.txt
var example1 []byte

//go:embed example2.txt
var example2 []byte

//go:embed example3.txt
var example3 []byte

//go:embed example4.txt
var example4 []byte

//go:embed example5.txt
var example5 []byte

//go:embed input.txt
var input []byte

func TestExample1(t *testing.T) {
	is := is.New(t)
	scan := bufio.NewScanner(bytes.NewReader(example1))

	result, err := run(scan)
	is.NoErr(err)

	// t.Log(result.valuePT1)
	is.Equal(result.valuePT1, 4)
}

func TestExample2(t *testing.T) {
	is := is.New(t)
	scan := bufio.NewScanner(bytes.NewReader(example2))

	result, err := run(scan)
	is.NoErr(err)

	// t.Log(result.valuePT1)
	is.Equal(result.valuePT1, 8)
}

func TestExample3(t *testing.T) {
	is := is.New(t)
	scan := bufio.NewScanner(bytes.NewReader(example3))

	result, err := run(scan)
	is.NoErr(err)

	// t.Log(result.valuePT1)
	is.Equal(result.valuePT1, 23)

	// t.Log(result.valuePT2)
	is.Equal(result.valuePT2, 4)
}

func TestExample4(t *testing.T) {
	is := is.New(t)
	scan := bufio.NewScanner(bytes.NewReader(example4))

	result, err := run(scan)
	is.NoErr(err)

	// t.Log(result.valuePT1)
	is.Equal(result.valuePT1, 70)

	// t.Log(result.valuePT2)
	is.Equal(result.valuePT2, 8)
}

func TestExample5(t *testing.T) {
	is := is.New(t)
	scan := bufio.NewScanner(bytes.NewReader(example5))

	result, err := run(scan)
	is.NoErr(err)

	// t.Log(result.valuePT1)
	is.Equal(result.valuePT1, 80)

	// t.Log(result.valuePT2)
	is.Equal(result.valuePT2, 10)
}

func TestInput(t *testing.T) {
	is := is.New(t)
	scan := bufio.NewScanner(bytes.NewReader(input))

	result, err := run(scan)
	is.NoErr(err)

	// t.Log(result.valuePT1)
	is.True(result.valuePT1 != 51)
	is.Equal(result.valuePT1, 6649)

	t.Log(result.valuePT2)
	is.True(result.valuePT2 != 0)
	is.Equal(result.valuePT2, 601)
}

// first: 51 false
