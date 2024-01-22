package aoc_test

import (
	"sort"
	"testing"

	"github.com/matryer/is"
	aoc "go.sour.is/advent-of-code"
)

func TestSet(t *testing.T) {
	is := is.New(t)

	s := aoc.Set(1, 2, 3)
	is.True(!s.Has(0))
	is.True(s.Has(1))
	is.True(s.Has(2))
	is.True(s.Has(3))
	is.True(!s.Has(4))

	s.Add(4)
	is.True(s.Has(4))

	items := s.Items()
	sort.Ints(items)
	is.Equal(items, []int{1, 2, 3, 4})
}
