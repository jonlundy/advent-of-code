package main

import (
	"bufio"
	_ "embed"
	"fmt"
	"strings"

	aoc "go.sour.is/advent-of-code"
)

func main() { aoc.MustResult(aoc.Runner(run)) }

type result struct {
	sum  int
	sum2 int
}

func (r result) String() string {
	return fmt.Sprintln("result pt1:", r.sum, "\nresult pt2:", r.sum2)
}

var numbers = []string{"zero", "one", "two", "three", "four", "five", "six", "seven", "eight", "nine"}

func run(scan *bufio.Scanner) (*result, error) {
	result := &result{}

	for scan.Scan() {
		var first, last int
		var first2, last2 int

		text := scan.Text()

		slice := make([]rune, 5)
		for i := range text {
			copy(slice, []rune(text[i:]))

			switch {
			case slice[0] >= '0' && slice[0] <= '9':
				if first == 0 {
					first = int(slice[0] - '0')
				}
				if first2 == 0 {
					first2 = int(slice[0] - '0')
				}
			default:
				if first2 != 0 {
					continue
				}

				for i, s := range numbers {
					if strings.HasPrefix(string(slice), s) {
						first2 = i
						break
					}
				}
			}
			if first != 0 && first2 != 0 {
				break
			}
		}

		text = string(aoc.Reverse([]rune(text)))

		for i := range text {
			copy(slice, []rune(text[i:]))
			slice = aoc.Reverse(slice)

			switch {
			case slice[4] >= '0' && slice[4] <= '9':
				if last == 0 {
					last = int(slice[4] - '0')
				}
				if last2 == 0 {
					last2 = int(slice[4] - '0')
				}
			default:
				if last2 != 0 {
					continue
				}
				for i, s := range numbers {
					if strings.HasSuffix(string(slice), s) {
						last2 = i
						break
					}
				}

			}
			if last != 0 && last2 != 0 {
				break
			}
		}

		result.sum += first*10 + last
		result.sum2 += first2*10 + last2
	}
	if err := scan.Err(); err != nil {
		return nil, err
	}

	return result, nil
}
