package aoc

import "cmp"

type uinteger interface {
	~uint | ~uint8 | ~uint16 | ~uint32 | ~uint64
}
type sinteger interface {
	~int | ~int8 | ~int16 | ~int32 | ~int64
}
type integer interface {
	sinteger | uinteger
}

type float interface {
	complex64 | complex128 | float32 | float64
}
type number interface{ integer | float }

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

func Sum[T number](arr ...T) T {
	var acc T
	for _, a := range arr {
		acc += a
	}
	return acc
}
func SumFunc[T any, U number](fn func(T) U, input ...T) U {
	return Sum(SliceMap(fn, input...)...)
}
func SumIFunc[T any, U number](fn func(int, T) U, input ...T) U {
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

func ABS[I integer](i I) I {
	if i < 0 {
		return -i
	}
	return i
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
