package aoc

import (
	"bufio"
	"fmt"
	"os"
	"path/filepath"
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


func ReadStringToInts(fields []string) []int {
	return SliceMap(Atoi, fields...)
}
