package main

import (
	"fmt"
	"os"
)

type Grid = [8][8]Cell

type Cell int8

const (
	Empty Cell = iota
	White
	Black
)

func charToCell(c byte) Cell {
	switch c {
	case '.':
		return Empty
	case 'w':
		return White
	case 'b':
		return Black
	}
	panic("invalid cell value " + string(c))
}

func main() {
	// boardSize: height and width of the board
	var boardSize int
	fmt.Scan(&boardSize)

	// color: current color of your pieces ("w" or "b")
	var color string
	fmt.Scan(&color)

	for {
		grid := Grid{}

		for i := 0; i < boardSize; i++ {
			// line: horizontal row
			var line string
			fmt.Scan(&line)

			for j := 0; j < boardSize; j++ {
				grid[i][j] = charToCell(line[j])
			}
		}

		debug("grid", grid)

		// lastAction: last action made by the opponent ("null" if it's the first turn)
		var lastAction string
		fmt.Scan(&lastAction)

		// actionsCount: number of legal actions
		var actionsCount int
		fmt.Scan(&actionsCount)

		// fmt.Fprintln(os.Stderr, "Debug messages...")
		fmt.Println("random") // e.g. e2e3 (move piece at e2 to e3)
	}

}

func debug(v ...interface{}) {
	fmt.Fprintln(os.Stderr, v...)
}
