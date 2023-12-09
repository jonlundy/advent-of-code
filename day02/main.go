package main

import (
	"bufio"
	_ "embed"
	"fmt"
	"os"
	"strconv"
	"strings"

	aoc "go.sour.is/advent-of-code-2023"
)

func main() {
	result, err := aoc.Runner(run)
	if err != nil {
		aoc.Log("ERR", err)
		os.Exit(1)
	}

	aoc.Log(result)
}

type result struct {
	sum      int
	powerSum int
}

func (r result) String() string {
	return fmt.Sprintln("result pt1:", r.sum, "\nresult pt2:", r.powerSum)
}

type gameResult struct {
	red, green, blue int
}

func run(scan *bufio.Scanner) (*result, error) {
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

	aoc.Log(games)
	aoc.Log(len(games))

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
				aoc.Log("game", i, round, "too many red", round.red)
				ok = false
			} else if maxCounts.blue < round.blue {
				aoc.Log("game", i, round, "too many blue", round.blue)
				ok = false
			} else if maxCounts.green < round.green {
				aoc.Log("game", i, round, "too many green", round.green)
				ok = false
			}
			aoc.Log("game", i, round, ok)
		}
		if ok {
			sum += i
			aoc.Log("game", i, "passes", sum)
		}
		power := mins.red * mins.blue * mins.green
		aoc.Log("game", i, "mins", mins, power)
		powerSum += power
	}
	return &result{sum, powerSum}, nil
}
