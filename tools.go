package aoc

import (
	"bufio"
	"fmt"
	"os"
	"path/filepath"
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

type integer interface {
	int | int8 | int16 | int32 | int64 | uint | uint8 | uint16 | uint32 | uint64
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
	arr := make([]int, len(fields))
	for i, s := range fields {
		if v, err := strconv.Atoi(s); err == nil {
			arr[i] = v
		}
	}
	return arr
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

func Repeat(s string, i int) []string {
	lis := make([]string, i)
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
func SumFunc[T any,U integer](fn func(T) U, input ...T) U {
	return Sum(SliceMap(fn, input...)...)
}
func SumIFunc[T any,U integer](fn func(int, T) U, input ...T) U {
	return Sum(SliceIMap(fn, input...)...)
}

func Power2(n int) int {
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
