package main

import (
	"bufio"
	"bytes"
	_ "embed"
	"fmt"
	"strings"
)

//go:embed input.txt
var input []byte

func main() {
	buf := bytes.NewReader(input)
	scan := bufio.NewScanner(buf)

	sum := 0
	for scan.Scan() {
		var first, last int

		orig := scan.Text()
		_ = orig

		text := scan.Text()

		slice := make([]rune, 5)
		for i := range text {
			copy(slice, []rune(text[i:]))

			switch {
			case slice[0] >= '0' && slice[0] <= '9':
				first = int(slice[0] - '0')
			case strings.HasPrefix(string(slice), "one"):
				first = 1
			case strings.HasPrefix(string(slice), "two"):
				first = 2
			case strings.HasPrefix(string(slice), "three"):
				first = 3
			case strings.HasPrefix(string(slice), "four"):
				first = 4
			case strings.HasPrefix(string(slice), "five"):
				first = 5
			case strings.HasPrefix(string(slice), "six"):
				first = 6
			case strings.HasPrefix(string(slice), "seven"):
				first = 7
			case strings.HasPrefix(string(slice), "eight"):
				first = 8
			case strings.HasPrefix(string(slice), "nine"):
				first = 9

			}
			if first != 0 {
				break
			}
		}

		text = string(reverse([]rune(text)))

		for i := range text {
			copy(slice, []rune(text[i:]))
			slice = reverse(slice)

			switch {
			case slice[4] >= '0' && slice[4] <= '9':
				last = int(slice[4] - '0')
			case strings.HasSuffix(string(slice), "one"):
				last = 1
			case strings.HasSuffix(string(slice), "two"):
				last = 2
			case strings.HasSuffix(string(slice), "three"):
				last = 3
			case strings.HasSuffix(string(slice), "four"):
				last = 4
			case strings.HasSuffix(string(slice), "five"):
				last = 5
			case strings.HasSuffix(string(slice), "six"):
				last = 6
			case strings.HasSuffix(string(slice), "seven"):
				last = 7
			case strings.HasSuffix(string(slice), "eight"):
				last = 8
			case strings.HasSuffix(string(slice), "nine"):
				last = 9

			}
			if last != 0 {
				break
			}
		}

		sum += first*10 + last
	}
	if err := scan.Err(); err != nil {
		panic(err)
	}

	fmt.Println(sum)
}

func reverse[T any](arr []T) []T{
	for i := 0; i < len(arr)/2; i++ {
		arr[i], arr[len(arr)-i-1] = arr[len(arr)-i-1], arr[i]
	}
	return arr
}


// type sorter[T rune | int] []T

// func (s sorter[T]) Less(i, j int) bool { return s[i] < s[j] }
// func (s sorter[T]) Len() int      { return len(s) }
// func (s sorter[T]) Swap(i, j int) { s[i], s[j] = s[j], s[i] }

/*
0
1
2
3
4
5
6
7
8
9
one
two
three
four
five
six
seven
eight
nine
ten


*/
