package main

import (
	"fmt"
	"math"
	"math/rand"
	"os"
	"strconv"
	"time"
)

const DEBUG = true

const MAX_TIME_MS_CG = 150
const MAX_TIME_MS_LOCAL = 10 * 1000

type Grid = [64]Cell

type State struct {
	grid   Grid
	turn   int
	player Player
}

func (s State) Clone() State {
	gridCopy := s.grid

	return State{gridCopy, s.turn, s.player}
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

type MonteCarloResult struct {
	games int
	wins  int
}

type MCTSNode struct {
	state    State
	action   *Action
	visits   int
	wins     int
	parent   *MCTSNode
	children []*MCTSNode
}

func uctMCTS(node *MCTSNode) float64 {
	if node.parent == nil {
		return 0
	} else {
		return float64(node.wins)/float64(node.visits) + 1.41*math.Sqrt(math.Log(float64(node.parent.visits))/float64(node.visits))
	}
}

func showNode(node *MCTSNode) string {
	grid := fmt.Sprintf("%v", node.state.grid)
	uct := uctMCTS(node)
	childrenCount := len(node.children)
	return fmt.Sprintf("grid %s wins %d visits %d uct %f children size %d", grid, node.wins, node.visits, uct, childrenCount)
}

func selectionMCTS(node *MCTSNode) *MCTSNode {
	if DEBUG {
		debug(fmt.Sprintf("selectionMCTS from %s", showNode(node)))
	}

	if len(node.children) == 0 {
		return node
	}

	var bestChild *MCTSNode
	var bestValue float64

	for _, child := range node.children {
		value := uctMCTS(child)

		if bestChild == nil || value > bestValue {
			bestChild = child
			bestValue = value
		}
	}

	if DEBUG {
		debug(fmt.Sprintf("selectionMCTS child %s value %f", showNode(bestChild), bestValue))
	}

	return selectionMCTS(bestChild)
}

func expandMCTS(node *MCTSNode) {
	var children []*MCTSNode

	actions := getValidActions(&node.state)
	for _, action := range actions {
		childState := applyAction(node.state, &action)
		child := &MCTSNode{childState, &action, 0, 0, node, []*MCTSNode{}}

		if DEBUG {
			//debug(fmt.Sprintf("expandMCTS child %s", showNode(child)))
		}

		children = append(children, child)
	}

	node.children = children
}

func simulateMCTS(node *MCTSNode, myPlayer Player) Player {
	// pick a random child node and simulate a game
	if len(node.children) == 0 {
		return Player(node.state.player)
	}

	child := node.children[rand.Intn(len(node.children))]
	return playUntilEnd(child.state)
}

func backPropagateMCTS(node *MCTSNode, winner Player) {
	if winner == node.state.player {
		node.wins++
	}
	node.visits++

	if node.parent != nil {
		backPropagateMCTS(node.parent, winner)
	}
}

func searchMCTS(node *MCTSNode, myPlayer Player, iterations int) *MCTSNode {
	for i := 0; i < iterations; i++ {
		selectedNode := selectionMCTS(node)
		expandMCTS(selectedNode)
		winner := simulateMCTS(selectedNode, myPlayer)
		backPropagateMCTS(selectedNode, winner)
	}

	var bestChild *MCTSNode
	var bestValue int
	for _, child := range node.children {
		if bestChild == nil || child.visits > bestValue {
			bestChild = child
			bestValue = child.visits
		}
	}

	return bestChild
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
				grid[i*8+j] = charToCell(line[j])
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
			player: myPlayer,
		}

		// actionsCount: number of legal actions
		var actionsCount int
		fmt.Scan(&actionsCount)

		validActions := getValidActions(&state)
		//debug("validActions", validActions)

		if len(validActions) != actionsCount {
			panic("invalid number of actions: " + strconv.Itoa(len(validActions)) + " != " + strconv.Itoa(actionsCount))
		}

		_ = startTime

		//debug("Starting Monte Carlo")
		rootNode := MCTSNode{state, nil, 0, 0, nil, []*MCTSNode{}}
		bestNode := searchMCTS(&rootNode, myPlayer, 1000)
		bestAction := bestNode.action
		debug("bestAction", bestAction, "uct", uctMCTS(bestNode))

		fmt.Println(displayCoord(bestAction.From) + displayCoord(bestAction.To))
		turn++
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

func getValidActions(state *State) []Action {
	currentPlayerCell := getCellOfPlayer(state.player)
	actions := make([]Action, 0, 64)
	for i := 0; i < 8; i++ {
		for j := 0; j < 8; j++ {
			if state.grid[i*8+j] == currentPlayerCell {
				for _, d := range directions {
					dX := int8(j) + d.x
					dY := int8(i) + d.y

					if inMap(dX, dY) && isValidMove(&state.grid, int8(j), int8(i), dX, dY) {
						if len(actions) == 128 {
							panic("too many actions")
						}

						actions = append(actions, Action{
							From: Coord{int8(j), int8(i)},
							To:   Coord{dX, dY},
						})
					}
				}
			}
		}
	}
	return actions
}

func inMap(dX int8, dY int8) bool {
	return dX >= 0 && dX < 8 && dY >= 0 && dY < 8
}

func applyAction(state State, action *Action) State {
	state.grid[action.To.y*8+action.To.x] = state.grid[action.From.y*8+action.From.x]
	state.grid[action.From.y*8+action.From.x] = Empty
	state.turn = state.turn + 1
	state.player = getOpponent(state.player)
	return state
}

func isValidMove(grid *Grid, fX int8, fY int8, tX int8, tY int8) bool {
	fromCell := grid[fY*8+fX]
	toCell := grid[tY*8+tX]
	return fromCell != Empty && toCell != Empty && toCell != fromCell
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

func runMonteCarloSearch(state State, startTime int64, maxTimeMs int64) Action {
	rootActions := getValidActions(&state)
	rootResults := make(map[Action]MonteCarloResult)
	actionRobin := 0

	for (time.Now().UnixMilli() - startTime) < int64(maxTimeMs) {
		rootAction := rootActions[actionRobin%len(rootActions)]
		currentState := applyAction(state, &rootAction)
		winner := playUntilEnd(currentState)

		winScore := 0
		if winner == state.player {
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

func playUntilEnd(currentState State) Player {
	for depth := 0; ; depth++ {
		if depth > 8*8 {
			panic("depth too high")
		}

		validActions := getValidActions(&currentState)
		if len(validActions) == 0 {
			return getOpponent(currentState.player)
		}

		randAction := randomAction(validActions)
		currentState = applyAction(currentState, &randAction)
	}
}

func randomAction(validActions []Action) Action {
	return validActions[rand.Intn(len(validActions))]
}

func displayCoord(c Coord) string {
	// x maps to board columns from a to h
	// y maps to board rows from 1 to 8 but reversed top to bottom
	column := string(byte('a' + c.x))
	row := strconv.Itoa(8 - int(c.y))
	return column + row
}
