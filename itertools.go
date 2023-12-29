package aoc

import (
	"strconv"
)

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

func Reduce[T, U any](fn func(int, T, U) U, u U, list ...T) U {
	for i, t := range list {
		u = fn(i, t, u)
	}
	return u
}

func Reverse[T any](arr []T) []T {
	for i := 0; i < len(arr)/2; i++ {
		arr[i], arr[len(arr)-i-1] = arr[len(arr)-i-1], arr[i]
	}
	return arr
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

// Pairwise iterates over a list pairing i, i+1
func Pairwise[T any](arr []T) [][2]T {
	var pairs [][2]T
	for i := range arr[:len(arr)-1] {
		pairs = append(pairs, [2]T{arr[i], arr[i+1]})
	}
	return pairs
}
