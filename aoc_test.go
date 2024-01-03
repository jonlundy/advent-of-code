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

func TestPriorityQueue(t *testing.T) {
	is := is.New(t)

	type elem [2]int
	less := func(a, b elem) bool {
		return a[0] < b[0]
	}

	pq := aoc.PriorityQueue(less)

	pq.Enqueue(elem{1, 4})
	pq.Enqueue(elem{3, 2})
	pq.Enqueue(elem{2, 3})
	pq.Enqueue(elem{4, 1})

	v, ok := pq.Dequeue()
	is.True(ok)
	is.Equal(v, elem{4, 1})

	v, ok = pq.Dequeue()
	is.True(ok)
	is.Equal(v, elem{3, 2})

	v, ok = pq.Dequeue()
	is.True(ok)
	is.Equal(v, elem{2, 3})

	v, ok = pq.Dequeue()
	is.True(ok)
	is.Equal(v, elem{1, 4})

	v, ok = pq.Dequeue()
	is.True(!ok)
	is.Equal(v, elem{})
}

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

func ExamplePriorityQueue() {
	type memo struct {
		pt    int
		score int
	}
	less := func(a, b memo) bool { return b.score < a.score }

	adj := map[int][][2]int{
		0: {{1, 2}, {2, 6}},
		1: {{3, 5}},
		2: {{3, 8}},
		3: {{4, 10}, {5, 15}},
		4: {{6, 2}},
		5: {{6, 6}},
	}

	pq := aoc.PriorityQueue(less)
	visited := aoc.Set([]int{}...)
	dist := aoc.DefaultMap[int](int(^uint(0) >> 1))

	dist.Set(0, 0)
	pq.Enqueue(memo{0, 0})

	for !pq.IsEmpty() {
		m, _ := pq.Dequeue()

		u := m.pt
		if visited.Has(u) {
			continue
		}
		visited.Add(u)

		du, _ := dist.Get(u)

		for _, edge := range adj[u] {
			v, w := edge[0], edge[1]
			dv, _ := dist.Get(v)

			if !visited.Has(v) && du+w < dv {
				dist.Set(v, du+w)
				pq.Enqueue(memo{v, du + w})
			}
		}
	}

	items := dist.Items()
	sort.Slice(items, func(i, j int) bool { return items[i].K < items[j].K })
	for _, v := range items {
		fmt.Printf("point %d is %d steps away.\n", v.K, v.V)
	}

	// Output:
	// point 0 is 0 steps away.
	// point 1 is 2 steps away.
	// point 2 is 6 steps away.
	// point 3 is 7 steps away.
	// point 4 is 17 steps away.
	// point 5 is 22 steps away.
	// point 6 is 19 steps away.
}

func TestStack(t *testing.T) {
	is := is.New(t)

	s := aoc.Stack(1,2,3,4)
	is.True(!s.IsEmpty())
	is.Equal(s.Pop(), 4)
	is.Equal(s.Pop(), 3)
	is.Equal(s.Pop(), 2)
	is.Equal(s.Pop(), 1)
	is.True(s.IsEmpty())
	s.Push(4,3,2,1)
	is.True(!s.IsEmpty())
	is.Equal(s.Pop(), 1)
	is.Equal(s.Pop(), 2)
	is.Equal(s.Pop(), 3)
	is.Equal(s.Pop(), 4)
	is.True(s.IsEmpty())
}

func TestGraph(t *testing.T) {
	is := is.New(t)

	var adjacencyList = map[int][]int{
		2: {3, 5, 1},
		1: {2, 4},
		3: {6, 2},
		4: {1, 5, 7},
		5: {2, 6, 8, 4},
		6: {3, 0, 9, 5},
		7: {4, 8},
		8: {5, 9, 7},
		9: {6, 0, 8},
	}

	g := aoc.Graph(aoc.WithAdjacencyList[int,int](adjacencyList))
	is.Equal(g.Neighbors(1), []int{2,4})
	is.Equal(map[int][]int(g.AdjacencyList()), adjacencyList)
}