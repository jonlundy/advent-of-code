package aoc_test

import (
	"fmt"
	"sort"
	"testing"

	"github.com/matryer/is"
	aoc "go.sour.is/advent-of-code"
)

func TestReverse(t *testing.T) {
	is := is.New(t)

	is.Equal(aoc.Reverse([]int{1, 2, 3, 4}), []int{4, 3, 2, 1})
}

func TestReadStringToInts(t *testing.T) {
	is := is.New(t)

	is.Equal(aoc.ReadStringToInts([]string{"1", "2", "3"}), []int{1, 2, 3})
}

func TestRepeat(t *testing.T) {
	is := is.New(t)

	is.Equal(aoc.Repeat(5, 3), []int{5, 5, 5})
}


func TestTranspose(t *testing.T) {
	is := is.New(t)

	is.Equal(
		aoc.Transpose(
			[][]int{
				{1, 1},
				{0, 0},
				{1, 1},
			},
		),
		[][]int{
			{1, 0, 1},
			{1, 0, 1},
		},
	)
}