package main

import (
	"fmt"
	"os"
)

type Grid = [8][8]Cell

type Cell uint8

const (
	Empty Cell = iota
	White
	Black
)

type Player uint8

const (
	WhitePlayer Player = iota
	BlackPlayer
)

type Coord struct {
	x, y int8
}

type Action struct {
	From, To Coord
}

var directions = [4]Coord{
	{1, 0},
	{0, 1},
	{-1, 0},
	{0, -1},
}

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

func parsePlayer(c byte) Player {
	switch c {
	case 'w':
		return WhitePlayer
	case 'b':
		return BlackPlayer
	}
	panic("invalid player value " + string(c))
}

func main() {
	// boardSize: height and width of the board
	var boardSize int
	fmt.Scan(&boardSize)

	// color: current color of your pieces ("w" or "b")
	var color string
	fmt.Scan(&color)

	myPlayer := parsePlayer(color[0])

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

		validActions := getValidActions(grid, myPlayer)

		if len(validActions) != actionsCount {
			panic("invalid number of actions: " + string(len(validActions)) + " != " + string(actionsCount))
		}

		// fmt.Fprintln(os.Stderr, "Debug messages...")
		fmt.Println("random") // e.g. e2e3 (move piece at e2 to e3)
	}

}

func debug(v ...interface{}) {
	fmt.Fprintln(os.Stderr, v...)
}

func getCellOfPlayer(p Player) Cell {
	switch p {
	case WhitePlayer:
		return White
	case BlackPlayer:
		return Black
	}
	panic("invalid player value " + string(p))
}

func isInMap(coord Coord) bool {
	return coord.x >= 0 && coord.x < 8 && coord.y >= 0 && coord.y < 8
}

func getValidActions(grid Grid, player Player) []Action {
	var actions []Action
	for i := 0; i < 8; i++ {
		for j := 0; j < 8; j++ {
			if grid[i][j] == getCellOfPlayer(player) {
				for _, d := range directions {
					destCoord := Coord{int8(i) + d.x, int8(j) + d.y}

					if isInMap(destCoord) && isValidMove(grid, Coord{int8(i), int8(j)}, destCoord) {
						actions = append(actions, Action{
							From: Coord{int8(i), int8(j)},
							To:   Coord{int8(i) + d.x, int8(j) + d.y},
						})
					}
				}
			}
		}
	}
	return actions
}

func isValidMove(grid Grid, from Coord, to Coord) bool {
	return grid[from.x][from.y] != Empty &&
		grid[to.x][to.y] != Empty &&
		grid[to.x][to.y] != grid[from.x][from.y]
}
