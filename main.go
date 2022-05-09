package main

import (
	"fmt"
	"math/rand"
	"os"
	"strconv"
	"time"
)

const NB_GAMES_PER_ROOT_ACTION_TOTAL = 2000
const IS_CG = true

type Grid = [8][8]Cell

type State struct {
	grid         Grid
	turn         int
	winner       int
	player       Player
	validActions []Action
}

func (s State) Clone() State {
	gridCopy := s.grid

	validActionsCopy := make([]Action, len(s.validActions))
	copy(validActionsCopy, s.validActions)

	return State{gridCopy, s.turn, s.winner, s.player, validActionsCopy}
}

type Cell uint8

const (
	Empty Cell = iota
	White
	Black
)

type Player uint8

const (
	WhitePlayer Player = iota + 1
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

var startGrid = Grid{
	{White, Black, White, Black, White, Black, White, Black},
	{Black, White, Black, White, Black, White, Black, White},
	{White, Black, White, Black, White, Black, White, Black},
	{Black, White, Black, White, Black, White, Black, White},
	{White, Black, White, Black, White, Black, White, Black},
	{Black, White, Black, White, Black, White, Black, White},
	{White, Black, White, Black, White, Black, White, Black},
	{Black, White, Black, White, Black, White, Black, White},
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
	// random seed to current datetime
	rand.Seed(time.Now().UnixNano())

	if IS_CG {

		// boardSize: height and width of the board
		var boardSize int
		fmt.Scan(&boardSize)

		// color: current color of your pieces ("w" or "b")
		var color string
		fmt.Scan(&color)

		myPlayer := parsePlayer(color[0])

		turn := 0

		for {
			turn++

			grid := Grid{}

			startTime := int64(0)
			for i := 0; i < boardSize; i++ {
				// line: horizontal row
				var line string
				fmt.Scan(&line)
				startTime = time.Now().UnixMilli()

				for j := 0; j < boardSize; j++ {
					grid[i][j] = charToCell(line[j])
				}
			}

			//debug("grid", grid)

			// lastAction: last action made by the opponent ("null" if it's the first turn)
			var lastAction string
			fmt.Scan(&lastAction)

			if lastAction == "null" {
				turn = 1
			}

			state := State{
				grid:   grid,
				turn:   turn,
				winner: 0,
				player: myPlayer,
			}

			// actionsCount: number of legal actions
			var actionsCount int
			fmt.Scan(&actionsCount)

			state.validActions = getValidActions(state)
			//debug("validActions", validActions)

			if len(state.validActions) != actionsCount {
				panic("invalid number of actions: " + strconv.Itoa(len(state.validActions)) + " != " + strconv.Itoa(actionsCount))
			}

			//debug("Starting Monte Carlo")
			bestAction := runMonteCarloSearch(state, startTime)
			debug("bestAction", bestAction)

			fmt.Println(displayCoord(bestAction.From) + displayCoord(bestAction.To))
			turn++
		}
	} else {
		state := State{
			grid:   startGrid,
			turn:   1,
			winner: 0,
			player: WhitePlayer,
		}

		best := runMonteCarloSearch(state, 0)
		debug("best", best)
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

func getValidActions(state State) []Action {
	var actions []Action
	for i := 0; i < 8; i++ {
		for j := 0; j < 8; j++ {
			if state.grid[i][j] == getCellOfPlayer(state.player) {
				for _, d := range directions {
					fromCoord := Coord{int8(j), int8(i)}
					destCoord := Coord{int8(j) + d.x, int8(i) + d.y}

					if isInMap(destCoord) && isValidMove(state.grid, fromCoord, destCoord) {
						actions = append(actions, Action{
							From: fromCoord,
							To:   destCoord,
						})
					}
				}
			}
		}
	}
	return actions
}

func applyAction(state State, action Action) State {
	newState := state.Clone()
	newState.grid[action.To.y][action.To.x] = newState.grid[action.From.y][action.From.x]
	newState.grid[action.From.y][action.From.x] = Empty
	newState.turn = state.turn + 1
	newState.player = getOpponent(state.player)

	validActions := getValidActions(newState)
	newState.validActions = validActions

	if len(validActions) == 0 {
		newState.winner = int(getOpponent(newState.player))
	}

	return newState
}

func isValidMove(grid Grid, from Coord, to Coord) bool {
	return grid[from.y][from.x] != Empty &&
		grid[to.y][to.x] != Empty &&
		grid[to.y][to.x] != grid[from.y][from.x]
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

func getRemainingPieces(grid Grid, p Player) int {
	var count int
	for i := 0; i < 8; i++ {
		for j := 0; j < 8; j++ {
			if grid[i][j] == getCellOfPlayer(p) {
				count++
			}
		}
	}
	return count
}

func runMonteCarloSearch(state State, startTime int64) Action {
	rootActions := getValidActions(state)
	//nbGamesPerRootAction := NB_GAMES_PER_ROOT_ACTION_TOTAL / len(rootActions)
	rootResults := make(map[Action]MonteCarloResult)

	actionRobin := 0

	for (time.Now().UnixMilli() - startTime) < 150 {
		rootAction := rootActions[actionRobin%len(rootActions)]
		currentState := applyAction(state, rootAction)

		depth := 0
		for depth = 0; ; depth++ {
			if depth > 8*8 {
				panic("depth too high")
			}
			if currentState.winner != 0 {
				break
			}
			currentState = applyAction(currentState, randomAction(currentState))
		}

		isWinning := currentState.winner == int(state.player)

		winScore := 0
		if isWinning {
			winScore = 1
		}

		existing, ok := rootResults[rootAction]
		if !ok {
			rootResults[rootAction] = MonteCarloResult{
				wins:  winScore,
				games: 1,
			}
		} else {
			existing.wins += winScore
			existing.games++
			rootResults[rootAction] = existing
		}

		actionRobin++
	}

	// find the action with the highest win rate
	var bestAction Action
	var bestRate = float64(-1)
	for _, rootAction := range rootActions {
		rate := float64(rootResults[rootAction].wins) / float64(rootResults[rootAction].games)
		if rate > bestRate {
			bestRate = rate
			bestAction = rootAction

			debug("bestAction", bestAction, "bestRate", fmt.Sprintf("%.2f", bestRate), "(", rootResults[bestAction].wins, "/", rootResults[bestAction].games, ")")
		}
	}
	return bestAction
}

func randomAction(currentState State) Action {
	return currentState.validActions[rand.Intn(len(currentState.validActions))]
}

func displayCoord(c Coord) string {
	// x maps to board columns from a to h
	// y maps to board rows from 1 to 8 but reversed top to bottom
	column := string('a' + c.x)
	row := strconv.Itoa(8 - int(c.y))
	return column + row
}
