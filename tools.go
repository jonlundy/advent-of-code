package aoc

import (
	"bufio"
	"fmt"
	"os"
)

func Runner[R any, F func(*bufio.Scanner) (R, error)](run F) (R, error) {
	if len(os.Args) != 2 {
		fmt.Fprintln(os.Stderr, "Usage:", os.Args[0] , "FILE")
	}

	input, err := os.Open(os.Args[1])
	if err != nil {
		fmt.Fprintln(os.Stderr, err)
	}

	scan := bufio.NewScanner(input)
	return run(scan)
}

func Reverse[T any](arr []T) []T {
	for i := 0; i < len(arr)/2; i++ {
		arr[i], arr[len(arr)-i-1] = arr[len(arr)-i-1], arr[i]
	}
	return arr
}


// greatest common divisor (GCD) via Euclidean algorithm
func GCD(a, b uint64) uint64 {
	for b != 0 {
		t := b
		b = a % b
		a = t
	}
	return a
}

// find Least Common Multiple (LCM) via GCD
func LCM(integers ...uint64) uint64 {
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
