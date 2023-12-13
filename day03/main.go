package main

import (
	"bufio"
	_ "embed"
	"fmt"
	"strconv"

	aoc "go.sour.is/advent-of-code-2023"
)

func main() { aoc.MustResult(aoc.Runner(run)) }

type partNumber struct {
	number    int
	row       int
	col       int
	end       int
	hasSymbol bool
}

type symbol struct {
	symbol        rune
	row           int
	col           int
	adjacentParts []*partNumber
}
type symbolTab map[int]map[int]*symbol

func (tab symbolTab) hasSymbol(row, col int, p partNumber) bool {
	if cols, ok := tab[row]; ok {
		s, ok := cols[col]
		if ok {
			s.adjacentParts = append(s.adjacentParts, &p)
			cols[col] = s
		}
		return ok
	}
	return false
}
func (tab symbolTab) scanSymbol(p partNumber) bool {
	rowStart, rowEnd := max(p.row-1, 0), p.row+1
	colStart, colEnd := max(p.col-1, 0), p.end+1

	for i := rowStart; i <= rowEnd; i++ {
		for j := colStart; j <= colEnd; j++ {
			ok := tab.hasSymbol(i, j, p)
			if ok {
				return true
			}
		}
	}
	return false
}

// 553079
// 84363105


type result struct {
	valuePT1 int
	valuePT2 int
}

func run(scan *bufio.Scanner) (*result, error) {
	parts := []partNumber{}
	symbols := make(symbolTab)
	symbolList := []*symbol{}
	row := 0
	for scan.Scan() {
		text := scan.Text()
		row += 1

		slice := make([]rune, 0, 3)
		var col int
		for i, a := range text {
			col = i
			if a >= '0' && a <= '9' {
				slice = append(slice, a)
				continue
			}
			if v, err := strconv.Atoi(string(slice)); err == nil {
				parts = append(parts, partNumber{number: v, row: row, col: col - len(slice), end: col - 1})
				slice = slice[:0]
			}

			if a != '.' {
				cols, ok := symbols[row]
				if !ok {
					cols = make(map[int]*symbol)
				}
				s := &symbol{row: row, col: col, symbol: a}
				cols[col] = s
				symbols[row] = cols
				symbolList = append(symbolList, s)
			}
		}
		if v, err := strconv.Atoi(string(slice)); err == nil {
			parts = append(parts, partNumber{number: v, row: row, col: col - len(slice), end: col - 1})
			slice = slice[:0]
			_ = slice
		}
	}

	sum := aoc.SumIFunc(
		func(i int, p partNumber) int {
			ok := symbols.scanSymbol(p)
			parts[i].hasSymbol = ok
			if ok {
				return p.number
			}
			return 0
		}, parts...,)

	sumGears := aoc.SumFunc(
		func(s *symbol) int {
			if s.symbol == '*' && len(s.adjacentParts) == 2 {
				return s.adjacentParts[0].number * s.adjacentParts[1].number
			}
			return 0
		}, symbolList...)

	// fmt.Println(parts)
	// fmt.Println(symbols)
	// fmt.Println(symbolList)
	fmt.Println("part1:", sum)
	fmt.Println("part2:", sumGears)
	return &result{sum, sumGears}, nil
}
