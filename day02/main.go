package main

import (
	"bufio"
	"bytes"
	_ "embed"
	"fmt"
	"strconv"
	"strings"
)

//go:embed input.txt
var input []byte

type gameResult struct {
	red, green, blue int
}

func main() {
	buf := bytes.NewReader(input)
	scan := bufio.NewScanner(buf)

	// only 12 red cubes, 13 green cubes, and 14 blue cubes
	maxCounts := gameResult{
		red:   12,
		green: 13,
		blue:  14,
	}

	games := [][]gameResult{}
	games = append(games, []gameResult{})

	for scan.Scan() {
		pfx, text, ok := strings.Cut(scan.Text(), ":")
		if !ok || !strings.HasPrefix(pfx, "Game") {
			continue
		}

		games = append(games, []gameResult{})

		for _, round := range strings.Split(text, ";") {
			game := gameResult{}

			for _, result := range strings.Split(round, ",") {
				ns, color, _ := strings.Cut(strings.TrimSpace(result), " ")
				n, err := strconv.Atoi(ns)
				if err != nil {
					panic(err)
				}

				switch color {
				case "red":
					game.red = n
				case "green":
					game.green = n
				case "blue":
					game.blue = n
				}

			}
			games[len(games)-1] = append(games[len(games)-1], game)
		}
	}

	fmt.Println(games)
	fmt.Println(len(games))

	sum := 0
	powerSum := 0
	for i, game := range games {
		if i == 0 {
			continue
		}

		mins := gameResult{}
		ok := true

		for _, round := range game {
			mins.red = max(mins.red, round.red)
			mins.green = max(mins.green, round.green)
			mins.blue = max(mins.blue, round.blue)

			if maxCounts.red < round.red {
				fmt.Println("game", i, round, "too many red", round.red)
				ok = false
			} else if maxCounts.blue < round.blue {
				fmt.Println("game", i, round, "too many blue", round.blue)
				ok = false
			} else if maxCounts.green < round.green {
				fmt.Println("game", i, round, "too many green", round.green)
				ok = false
			}
			fmt.Println("game", i, round, ok)
		}
		if ok {
			sum += i
			fmt.Println("game", i, "passes", sum)
		}
		power := mins.red*mins.blue*mins.green
		fmt.Println("game", i, "mins", mins, power)
		powerSum += power
	}
	fmt.Println("sum", sum, "power", powerSum)
}
