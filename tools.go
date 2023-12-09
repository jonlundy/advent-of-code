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
