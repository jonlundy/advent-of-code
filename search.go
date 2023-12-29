package aoc

import (
	"sort"
)

type priorityQueue[T any, U []T] struct {
	elems U
	less  func(a, b T) bool
}

func PriorityQueue[T any, U []T](less func(a, b T) bool) *priorityQueue[T, U] {
	return &priorityQueue[T, U]{less: less}
}
func (pq *priorityQueue[T, U]) Enqueue(elem T) {
	pq.elems = append(pq.elems, elem)
	sort.Slice(pq.elems, func(i, j int) bool { return pq.less(pq.elems[j], pq.elems[i]) })
}
func (pq *priorityQueue[T, I]) IsEmpty() bool {
	return len(pq.elems) == 0
}
func (pq *priorityQueue[T, I]) Dequeue() (T, bool) {
	var elem T
	if pq.IsEmpty() {
		return elem, false
	}

	pq.elems, elem = pq.elems[:len(pq.elems)-1], pq.elems[len(pq.elems)-1]
	return elem, true
}

type DS[T comparable] struct {
	*priorityQueue[T, []T]
	*set[T]
}
