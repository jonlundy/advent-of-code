package aoc_test

import (
	"testing"

	"github.com/matryer/is"
	aoc "go.sour.is/advent-of-code"
)

func TestReverse(t *testing.T) {
	is := is.New(t)

	is.Equal(aoc.Reverse([]int{1, 2, 3, 4}), []int{4, 3, 2, 1})
}

func TestLCM(t *testing.T) {
	is := is.New(t)

	is.Equal(aoc.LCM([]int{}...), 0)
	is.Equal(aoc.LCM(5), 5)
	is.Equal(aoc.LCM(5, 3), 15)
	is.Equal(aoc.LCM(5, 3, 2), 30)
}

func TestReadStringToInts(t *testing.T) {
	is := is.New(t)

	is.Equal(aoc.ReadStringToInts([]string{"1", "2", "3"}), []int{1, 2, 3})
}

func TestRepeat(t *testing.T) {
	is := is.New(t)

	is.Equal(aoc.Repeat(5, 3), []int{5, 5, 5})
}

func TestPower2(t *testing.T) {
	is := is.New(t)

	is.Equal(aoc.Power2(0), 1)
	is.Equal(aoc.Power2(1), 2)
	is.Equal(aoc.Power2(2), 4)
}

func TestABS(t *testing.T) {
	is := is.New(t)

	is.Equal(aoc.ABS(1), 1)
	is.Equal(aoc.ABS(0), 0)
	is.Equal(aoc.ABS(-1), 1)
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

func TestList(t *testing.T) {
	is := is.New(t)

	lis := aoc.NewList[int](nil)
	lis.Add(5, 0)

	a, _ := lis.Head().Value()

	is.Equal(a, 5)
}

func TestGraph(t *testing.T) {
	g := aoc.Graph[int, uint](7)
	g.AddEdge(0, 1, 2)
	g.AddEdge(0, 2, 6)
	g.AddEdge(1, 3, 5)
	g.AddEdge(2, 3, 8)
	g.AddEdge(3, 4, 10)
	g.AddEdge(3, 5, 15)
	g.AddEdge(4, 6, 2)
	g.AddEdge(5, 6, 6)
	g.Dijkstra(0)
}
