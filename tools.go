package aoc

import (
	"bufio"
	"cmp"
	"fmt"
	"os"
	"path/filepath"
	"sort"
	"strconv"
	"strings"
)

func Runner[R any, F func(*bufio.Scanner) (R, error)](run F) (R, error) {
	if len(os.Args) != 2 {
		Log("Usage:", filepath.Base(os.Args[0]), "FILE")
		os.Exit(22)
	}

	input, err := os.Open(os.Args[1])
	if err != nil {
		Log(err)
		os.Exit(1)
	}

	scan := bufio.NewScanner(input)
	return run(scan)
}

func MustResult[T any](result T, err error) {
	if err != nil {
		fmt.Println("ERR", err)
		os.Exit(1)
	}

	Log("result", result)
}

func Log(v ...any) { fmt.Fprintln(os.Stderr, v...) }
func Logf(format string, v ...any) {
	if !strings.HasSuffix(format, "\n") {
		format += "\n"
	}
	fmt.Fprintf(os.Stderr, format, v...)
}

func Reverse[T any](arr []T) []T {
	for i := 0; i < len(arr)/2; i++ {
		arr[i], arr[len(arr)-i-1] = arr[len(arr)-i-1], arr[i]
	}
	return arr
}

type uinteger interface{
	uint | uint8 | uint16 | uint32 | uint64
}
type sinteger interface{
	int | int8 | int16 | int32 | int64
}
type integer interface {
	sinteger | uinteger
}

// type float interface {
// 	complex64 | complex128 | float32 | float64
// }
// type number interface{ integer | float }

// greatest common divisor (GCD) via Euclidean algorithm
func GCD[T integer](a, b T) T {
	for b != 0 {
		t := b
		b = a % b
		a = t
	}
	return a
}

// find Least Common Multiple (LCM) via GCD
func LCM[T integer](integers ...T) T {
	if len(integers) == 0 {
		return 0
	}
	if len(integers) == 1 {
		return integers[0]
	}

	a, b := integers[0], integers[1]
	result := a * b / GCD(a, b)

	for _, c := range integers[2:] {
		result = LCM(result, c)
	}

	return result
}

func ReadStringToInts(fields []string) []int {
	return SliceMap(Atoi, fields...)
}

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

func SliceMap[T, U any](fn func(T) U, in ...T) []U {
	lis := make([]U, len(in))
	for i := range lis {
		lis[i] = fn(in[i])
	}
	return lis
}
func SliceIMap[T, U any](fn func(int, T) U, in ...T) []U {
	lis := make([]U, len(in))
	for i := range lis {
		lis[i] = fn(i, in[i])
	}
	return lis
}

func Atoi(s string) int {
	i, _ := strconv.Atoi(s)
	return i
}

func Repeat[T any](s T, i int) []T {
	lis := make([]T, i)
	for i := range lis {
		lis[i] = s
	}
	return lis
}

func Sum[T integer](arr ...T) T {
	var acc T
	for _, a := range arr {
		acc += a
	}
	return acc
}
func SumFunc[T any, U integer](fn func(T) U, input ...T) U {
	return Sum(SliceMap(fn, input...)...)
}
func SumIFunc[T any, U integer](fn func(int, T) U, input ...T) U {
	return Sum(SliceIMap(fn, input...)...)
}

func Power2(n int) int {
	if n == 0 {
		return 1
	}
	p := 2
	for ; n > 1; n-- {
		p *= 2
	}
	return p
}

func ABS(i int) int {
	if i < 0 {
		return -i
	}
	return i
}

func Transpose[T any](matrix [][]T) [][]T {
	rows, cols := len(matrix), len(matrix[0])

	m := make([][]T, cols)
	for i := range m {
		m[i] = make([]T, rows)
	}

	for i := 0; i < cols; i++ {
		for j := 0; j < rows; j++ {
			m[i][j] = matrix[j][i]
		}
	}
	return m
}

func Reduce[T, U any](fn func(int, T, U) U, u U, list ...T) U {
	for i, t := range list {
		u = fn(i, t, u)
	}
	return u
}

func Max[T cmp.Ordered](a T, v ...T) T {
	for _, b := range v {
		if b > a {
			a = b
		}
	}
	return a
}
func Min[T cmp.Ordered](a T, v ...T) T {
	for _, b := range v {
		if b < a {
			a = b
		}
	}
	return a
}

type PQElem[T any, I integer] struct {
	Value    T
	Priority I
}
type PQList[T any, I integer] []PQElem[T, I]

func (pq PQList[T, I]) Len() int {
	return len(pq)
}
func (pq PQList[T, I]) Less(i int, j int) bool {
	return pq[i].Priority < pq[j].Priority
}
func (pq PQList[T, I]) Swap(i int, j int) {
	pq[i], pq[j] = pq[j], pq[i]
}

var _ sort.Interface = (*PQList[rune, int])(nil)

type PriorityQueue[T any, I integer] struct {
	elem PQList[T, I]
}

func (pq *PriorityQueue[T, I]) Enqueue(elem T, priority I) {
	pq.elem = append(pq.elem, PQElem[T, I]{elem, priority})
	sort.Sort(pq.elem)
}
func (pq *PriorityQueue[T, I]) IsEmpty() bool {
	return len(pq.elem) == 0
}
func (pq *PriorityQueue[T, I]) Dequeue() (T, bool) {
	var elem T
	if pq.IsEmpty() {
		return elem, false
	}

	elem, pq.elem = pq.elem[0].Value, pq.elem[1:]
	return elem, true
}

type Vertex[V comparable, I integer] struct {
	to    V
	score I
}
type graph[V comparable, I uinteger] struct {
	adj map[V][]Vertex[V, I]
}
func Graph[V comparable, I uinteger](size int) *graph[V,I] {
	return &graph[V,I]{
		adj:  make(map[V][]Vertex[V,I], size),
	}
}
func (g *graph[V,I]) AddEdge(u, v V, w I) {
	g.adj[u] = append(g.adj[u], Vertex[V, I]{to: v, score: w})
	g.adj[v] = append(g.adj[v], Vertex[V, I]{to: u, score: w})
}
func (g *graph[V,I]) Dijkstra(src V) {
	pq := PriorityQueue[V,I]{}
	dist := make(map[V]I, len(g.adj))
	visited := make(map[V]bool, len(g.adj))
	var INF I
	INF = ^INF>>1

	pq.Enqueue(src, 0)
	dist[src] = 0

	for !pq.IsEmpty() {
		u, _ := pq.Dequeue()

		if _, ok := visited[u]; ok {
			continue
		}
		visited[u] = true

		for _, v := range g.adj[u] {
			_, ok := visited[v.to]
			var du, dv I
			if d, inf := dist[u]; !inf {
				du=INF
			} else {
				du = d
			}
			if d, inf := dist[v.to]; !inf {
				dv=INF
			} else {
				dv = d
			} 

			if !ok && du + v.score < dv {
				dist[v.to] = du + v.score
				pq.Enqueue(v.to, du + v.score)
			}
		}
	}
	for v, w := range dist {
		fmt.Printf("%v, %v\n", v, w)
	}
}