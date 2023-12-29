package aoc

import "fmt"

type Node[T any] struct {
	value T
	pos   int
	left  *Node[T]
}

func (n *Node[T]) add(a *Node[T]) *Node[T] {
	if a == nil {
		return n
	}

	if n == nil {
		return a
	}

	n.left = a
	return a
}

func (n *Node[T]) Value() (value T, ok bool) {
	if n == nil {
		return
	}
	return n.value, true
}

func (n *Node[T]) Position() int {
	if n == nil {
		return -1
	}
	return n.pos
}
func (n *Node[T]) SetPosition(i int) {
	if n == nil {
		return
	}
	n.pos = i
}
func (n *Node[T]) Next() *Node[T] {
	if n == nil {
		return nil
	}
	return n.left
}

func (n *Node[T]) String() string {
	if n == nil {
		return "EOL"
	}
	return fmt.Sprintf("node %v", n.value)
}

type List[T any] struct {
	head *Node[T]
	n    *Node[T]
	p    map[int]*Node[T]
}

func NewList[T any](a *Node[T]) *List[T] {
	lis := &List[T]{
		head: a,
		n:    a,
		p:    make(map[int]*Node[T]),
	}
	lis.add(a)

	return lis
}
func (l *List[T]) Add(value T, pos int) {
	a := &Node[T]{value: value, pos: pos}
	l.add(a)
}
func (l *List[T]) add(a *Node[T]) {
	if l.head == nil {
		l.head = a
	}
	if a == nil {
		return
	}

	l.n = l.n.add(a)
	l.p[a.pos] = a
}
func (l *List[T]) Get(pos int) *Node[T] {
	return l.p[pos]
}
func (l *List[T]) GetN(pos ...int) []*Node[T] {
	lis := make([]*Node[T], len(pos))
	for i, p := range pos {
		lis[i] = l.p[p]
	}
	return lis
}
func (l *List[T]) Head() *Node[T] {
	return l.head
}
