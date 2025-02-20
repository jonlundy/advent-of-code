package aoc_test

import (
	"fmt"
	"sort"
	"testing"

	"github.com/matryer/is"
	aoc "go.sour.is/advent-of-code"
)


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
	less := func(b, a *elem) bool {
		return (*a)[0] < (*b)[0]
	}

	pq := aoc.PriorityQueue(less)

	pq.Insert(&elem{1, 4})
	pq.Insert(&elem{3, 2})
	pq.Insert(&elem{2, 3})
	pq.Insert(&elem{4, 1})

	v := pq.ExtractMin()
	is.True(v != nil)
	is.Equal(v, &elem{4, 1})

	v = pq.ExtractMin()
	is.True(v != nil)
	is.Equal(v, &elem{3, 2})

	v = pq.ExtractMin()
	is.True(v != nil)
	is.Equal(v, &elem{2, 3})

	v = pq.ExtractMin()
	is.True(v != nil)
	is.Equal(v, &elem{1, 4})

	v = pq.ExtractMin()
	is.True(v == nil)
}


func ExamplePriorityQueue() {
	type memo struct {
		pt    int
		score int
	}
	less := func(a, b *memo) bool { return a.score < b.score }

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
	pq.Insert(&memo{0, 0})

	for !pq.IsEmpty() {
		m := pq.ExtractMin()

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
				pq.Insert(&memo{v, du + w})
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

	g := aoc.Graph(aoc.WithAdjacencyList[int, int](adjacencyList))
	is.Equal(g.Neighbors(1), []int{2, 4})
	is.Equal(map[int][]int(g.AdjacencyList()), adjacencyList)
}

func ExampleFibHeap() {
	type memo struct {
		pt    int
		score int
	}
	less := func(a, b *memo) bool { return (*a).score < (*b).score }

	adj := map[int][][2]int{
		0: {{1, 2}, {2, 6}},
		1: {{3, 5}},
		2: {{3, 8}},
		3: {{4, 10}, {5, 15}},
		4: {{6, 2}},
		5: {{6, 6}},
	}

	pq := aoc.FibHeap(less)
	visited := aoc.Set([]int{}...)
	dist := aoc.DefaultMap[int](int(^uint(0) >> 1))

	dist.Set(0, 0)
	pq.Insert(&memo{0, 0})

	for !pq.IsEmpty() {
		m := pq.ExtractMin()

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
				pq.Insert(&memo{v, du + w})
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

func TestFibHeap(t *testing.T) {
	is := is.New(t)

	type elem [2]int
	less := func(a, b *elem) bool {
		return (*a)[0] < (*b)[0]
	}

	pq := aoc.FibHeap(less)

	pq.Insert(&elem{1, 4})
	pq.Insert(&elem{3, 2})
	pq.Insert(&elem{2, 3})
	pq.Insert(&elem{4, 1})

	v := pq.ExtractMin()
	is.True(v != nil)
	is.Equal(v, &elem{1, 4})

	pq.Insert(&elem{5, 8})
	pq.Insert(&elem{6, 7})
	pq.Insert(&elem{7, 6})
	pq.Insert(&elem{8, 5})

	v = pq.ExtractMin()
	is.True(v != nil)
	is.Equal(v, &elem{2, 3})

	v = pq.ExtractMin()
	is.True(v != nil)
	is.Equal(v, &elem{3, 2})

	v = pq.ExtractMin()
	is.True(v != nil)
	is.Equal(v, &elem{4, 1})

	v = pq.ExtractMin()
	is.True(v != nil)
	is.Equal(v, &elem{5, 8})

	m := aoc.FibHeap(less)
	m.Insert(&elem{1, 99})
	m.Insert(&elem{12, 9})
	m.Insert(&elem{11, 10})
	m.Insert(&elem{10, 11})
	m.Insert(&elem{9, 12})

	pq.Merge(m)

	v = pq.Find(func(t *elem) bool {
		return (*t)[0] == 6
	})
	is.Equal(v, &elem{6, 7})

	v = pq.Find(func(t *elem) bool {
		return (*t)[0] == 12
	})
	is.Equal(v, &elem{12, 9})

	v = pq.ExtractMin()
	is.True(v != nil)
	is.Equal(v, &elem{1, 99})

	pq.DecreaseKey(
		func(t *elem) bool { return t[0] == 12 },
		func(t *elem) { t[0] = 3 },
	)

	v = pq.ExtractMin()
	is.True(v != nil)
	is.Equal(v, &elem{3, 9})

	var keys []int
	for !pq.IsEmpty() {
		v := pq.ExtractMin()
		fmt.Println(v)
		keys = append(keys, v[0])
	}
	is.Equal(keys, []int{6, 7, 8, 9, 10, 11})
}
