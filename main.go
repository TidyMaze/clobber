package main

import (
	"fmt"
	"math/rand"
	"os"
	"strconv"
)

const NB_GAMES_PER_ROOT_ACTION = 10

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

type MonteCarloResult struct {
	games int
	wins  int
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
		//debug("validActions", validActions)

		if len(validActions) != actionsCount {
			panic("invalid number of actions: " + strconv.Itoa(len(validActions)) + " != " + strconv.Itoa(actionsCount))
		}

		debug("Starting Monte Carlo")
		bestAction := runMonteCarloSearch(grid, myPlayer)
		debug("bestAction", bestAction)

		fmt.Println(bestAction.From.x, bestAction.From.y, bestAction.To.x, bestAction.To.y)
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
					destCoord := Coord{int8(j) + d.x, int8(i) + d.y}

					if isInMap(destCoord) && isValidMove(grid, Coord{int8(j), int8(i)}, destCoord) {
						actions = append(actions, Action{
							From: Coord{int8(j), int8(i)},
							To:   destCoord,
						})
					}
				}
			}
		}
	}
	return actions
}

func applyAction(grid Grid, action Action) Grid {
	grid[action.To.y][action.To.x] = grid[action.From.y][action.From.x]
	grid[action.From.y][action.From.x] = Empty
	return grid
}

func isValidMove(grid Grid, from Coord, to Coord) bool {
	return grid[from.x][from.y] != Empty &&
		grid[to.x][to.y] != Empty &&
		grid[to.x][to.y] != grid[from.x][from.y]
}

func getOpponent(p Player) Player {
	switch p {
	case WhitePlayer:
		return BlackPlayer
	case BlackPlayer:
		return WhitePlayer
	}
	panic("invalid player value " + string(p))
}

func runMonteCarloSearch(grid Grid, player Player) Action {
	rootActions := getValidActions(grid, player)
	//debug("rootActions", rootActions)

	rootResults := make(map[Action]MonteCarloResult)

	// run 1000 full games per root action, and store the winning rate
	// the loser is the player that is unable to play
	for iAction, rootAction := range rootActions {
		debug("Sampling root action", rootAction, "(", iAction, "/", len(rootActions), ")")
		var wins int
		var games int
		for i := 0; i < NB_GAMES_PER_ROOT_ACTION; i++ {
			currentGrid := grid
			currentPlayer := player

			depth := 0
			for depth = 0; ; depth++ {
				validActions := getValidActions(currentGrid, currentPlayer)
				if len(validActions) == 0 {
					break
				}

				action := validActions[rand.Intn(len(validActions))]
				currentGrid = applyAction(currentGrid, action)
				currentPlayer = getOpponent(currentPlayer)
			}

			if currentPlayer != player {
				wins++
			}
			games++

			debug("Game", i, "/", NB_GAMES_PER_ROOT_ACTION, "finished at depth", depth)
		}
		rootResults[rootAction] = MonteCarloResult{
			wins:  wins,
			games: games,
		}

		debug("rootAction", rootAction, "wins", wins, "games", games)
	}

	// find the action with the highest win rate
	var bestAction Action
	var bestRate float64
	for _, rootAction := range rootActions {
		rate := float64(rootResults[rootAction].wins) / float64(rootResults[rootAction].games)
		if rate > bestRate {
			bestRate = rate
			bestAction = rootAction

			debug("bestAction", bestAction, "bestRate", bestRate)
		}
	}
	return bestAction
}
