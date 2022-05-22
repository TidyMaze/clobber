package main

import (
	"fmt"
	"math"
	"math/rand"
	"os"
	"strconv"
	"strings"
	"time"
)

const DEBUG = false

const MAX_TIME_MS_CG = 135
const MAX_TIME_MS_LOCAL = 10 * 1000

var node_count = 0

var playouts = 0

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
	From, To int8
}

func (a Action) String() string {
	return fmt.Sprintf("%s%s", displayCoord(a.From), displayCoord(a.To))
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
	id       int
	state    *State
	action   *Action
	visits   int
	wins     int
	parent   *MCTSNode
	children []*MCTSNode
}

func uctMCTS(node *MCTSNode) float64 {
	if node.visits == 0 {
		return math.Inf(1)
	} else {
		return float64(node.wins)/float64(node.visits) + 1.41*math.Sqrt(math.Log(float64(node.parent.visits))/float64(node.visits))
	}
}

func showNode(node *MCTSNode) string {
	//grid := fmt.Sprintf("%v", node.state.grid)

	uct := float64(-1)
	if node.parent != nil {
		uct = uctMCTS(node)
	}

	parentNodeId := "nil"
	if node.parent != nil {
		parentNodeId = strconv.Itoa(node.parent.id)
	}

	action := "nil"
	if node.action != nil {
		action = fmt.Sprintf("%v", *node.action)
	}

	return fmt.Sprintf("node %d wins %d visits %d uct %f parent %s action %s", node.id, node.wins, node.visits, uct, parentNodeId, action)
}

func selectionMCTS(node *MCTSNode) *MCTSNode {
	if DEBUG {
		debug(fmt.Sprintf("selectionMCTS from   %s", showNode(node)))
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

	//if DEBUG {
	//	debug(fmt.Sprintf("selectionMCTS child %s value %f", showNode(bestChild), bestValue))
	//}

	return selectionMCTS(bestChild)
}

func expandMCTS(node *MCTSNode) {
	actions := getValidActions(node.state)

	max := len(*actions)

	for i := 0; i < max; i++ {
		childState := applyAction(*node.state, &(*actions)[i])
		node_count++

		if DEBUG {
			//debug(fmt.Sprintf("expandMCTS child %s", showNode(child)))
		}

		newNode := newNode(node, childState, &(*actions)[i])
		addToChildren(node, newNode)
	}
}

func addToChildren(node *MCTSNode, newNode *MCTSNode) {
	node.children = append(node.children, newNode)
}

func newNode(node *MCTSNode, childState *State, action *Action) *MCTSNode {
	return &MCTSNode{node_count, childState, action, 0, 0, node, make([]*MCTSNode, 0)}
}

func simulateMCTS(node *MCTSNode) (*MCTSNode, Player) {
	playouts++
	if len(node.children) == 0 {
		return node, getOpponent(Player(node.state.player))
	}

	child := node.children[rand.Intn(len(node.children))]

	if DEBUG {
		debug(fmt.Sprintf("simulateMCTS picked child %s", showNode(child)))
		showTree(node, 0)
	}

	return child, playUntilEnd(*child.state)
}

func backPropagateMCTS(node *MCTSNode, winner Player) {

	if winner == getOpponent(node.state.player) {
		node.wins++
	} else {
		node.wins--
	}
	node.visits++

	if DEBUG {
		debug(fmt.Sprintf("backPropagateMCTS   %s with winner %d  ", showNode(node), winner))
	}

	if node.parent != nil {
		backPropagateMCTS(node.parent, winner)
	}
}

func showTree(node *MCTSNode, padding int) {
	if node.visits > 0 {
		debug(strings.Repeat(" ", padding) + showNode(node))
	}
	for _, child := range node.children {
		showTree(child, padding+2)
	}
}

func mcts(node *MCTSNode, startTime int64, maxTimeMs int64) *MCTSNode {
	playouts = 0
	if DEBUG {
		debug("initial node", showNode(node))
		showTree(node, 0)
	}

	for i := 0; (time.Now().UnixMilli() - startTime) < maxTimeMs; i++ {
		selectedNode := selectionMCTS(node)
		expandMCTS(selectedNode)
		child, winner := simulateMCTS(selectedNode)
		backPropagateMCTS(child, winner)

		if DEBUG {
			showTree(node, 0)
			debug("==== end of iteration ====", i)
			debug()
		}
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
			turn:   0,
			player: myPlayer, // which player has to play after this turn
		}

		// actionsCount: number of legal actions
		var actionsCount int
		fmt.Scan(&actionsCount)

		validActions := getValidActions(&state)
		//debug("validActions", validActions)

		if len(*validActions) != actionsCount {
			panic("invalid number of actions: " + strconv.Itoa(len(*validActions)) + " != " + strconv.Itoa(actionsCount))
		}

		bestAction, bestValue := runMCTSSearch(state, startTime, MAX_TIME_MS_CG)
		debug("bestAction", bestAction, "bestValue", bestValue, "after", playouts, "playouts")

		fmt.Println(fmt.Sprintf("%s %.2f", bestAction, bestValue))
		turn++
	}
}

func runMCTSSearch(state State, startTime int64, maxTime int64) (*Action, float64) {
	node_count = 0
	rootNode := MCTSNode{node_count, &state, nil, 0, 0, nil, []*MCTSNode{}}
	node_count++
	bestNode := mcts(&rootNode, startTime, maxTime)
	return bestNode.action, float64(bestNode.visits)
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

func getValidActions(state *State) *[]Action {
	currentPlayerCell := getCellOfPlayer(state.player)
	opponentCell := getCellOfPlayer(getOpponent(state.player))

	actions := make([]Action, 0, 128)
	a := Action{}

	for i := 0; i < 8; i++ {
		for j := 0; j < 8; j++ {
			a.From = int8(i*8 + j)

			if state.grid[a.From] == currentPlayerCell {
				for id := 0; id < len(directions); id++ {
					dX := int8(j) + directions[id].x
					if dX < 0 || dX >= 8 {
						continue
					}

					dY := int8(i) + directions[id].y
					if dY < 0 || dY >= 8 {
						continue
					}

					a.To = dY*8 + dX

					if state.grid[a.To] == opponentCell {
						actions = append(actions, a)
					}
				}
			}
		}
	}
	return &actions
}

func inMap(dX int8, dY int8) bool {
	return dX >= 0 && dX < 8 && dY >= 0 && dY < 8
}

func applyAction(state State, action *Action) *State {
	state.grid[action.To] = state.grid[action.From]
	state.grid[action.From] = Empty
	state.turn = state.turn + 1
	state.player = getOpponent(state.player)
	return &state
}

func applyActionMut(state *State, action *Action) {
	state.grid[action.To] = state.grid[action.From]
	state.grid[action.From] = Empty
	state.turn = state.turn + 1
	state.player = getOpponent(state.player)
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

func stateEval(state *State, myPlayer Player, nextActions *[]Action) float64 {
	actionsCount := len(*nextActions)

	eval := math.Inf(-1)

	if actionsCount == 0 && myPlayer == state.player {
		eval = -1000000.0 + float64(state.turn)
	} else if actionsCount == 0 && myPlayer != state.player {
		eval = 1000000.0 - float64(state.turn)
	} else if myPlayer == state.player {
		eval = float64(actionsCount)*1000 + float64(state.turn)
	} else if myPlayer != state.player {
		eval = float64(actionsCount)*-1000 - float64(state.turn)
	}

	return eval
}

func runMinimaxSearch(state *State, maxDepth int) (Action, float64) {
	rootActions := getValidActions(state)

	var bestAction Action = Action{From: -1, To: -1}
	bestValue := math.Inf(-1)

	if DEBUG {
		debug("Taking max", maxDepth)
	}

	for _, action := range *rootActions {
		stateCopy := *state
		applyActionMut(&stateCopy, &action)

		value := -negamax(&stateCopy, maxDepth-1, state.player, math.Inf(-1), math.Inf(1), -1)
		if value > bestValue {
			bestValue = value
			bestAction = action
			//if DEBUG {
			debug("New best GLOBAL value", bestValue, "for action", action)
			//}
		}
	}

	return bestAction, bestValue
}

func negamax(state *State, maxDepth int, myPlayer Player, alpha float64, beta float64, color int) float64 {
	nextActions := getValidActions(state)

	if maxDepth == 0 {
		eval := stateEval(state, myPlayer, nextActions)
		if DEBUG {
			//debug("Reaching max depth", maxDepth, "eval", eval)
		}
		return float64(color) * eval
	}

	if len(*nextActions) == 0 {
		eval := stateEval(state, myPlayer, nextActions)
		if DEBUG {
			debug("Reaching leaf node", maxDepth, "eval", eval)
		}
		return float64(color) * eval
	}

	whoPlayed := getOpponent(state.player)

	if DEBUG {
		debug("Next player", state.player, "who played", whoPlayed, "Depth", maxDepth)
	}
	value := math.Inf(-1)
	for iNextAction := 0; iNextAction < len(*nextActions); iNextAction++ {
		nextAction := &(*nextActions)[iNextAction]
		nextState := applyAction(*state, nextAction)
		value = math.Max(value, -negamax(nextState, maxDepth-1, myPlayer, -beta, -alpha, -color))

		alpha = math.Max(alpha, value)

		if alpha >= beta {
			if DEBUG {
				debug("Cutoff", value, "for action", nextAction)
			}
			break
		}
	}
	return value
}

func runMonteCarloSearch(state State, startTime int64, maxTimeMs int64) Action {
	rootActions := getValidActions(&state)
	rootResults := make(map[Action]MonteCarloResult)
	actionRobin := 0

	for (time.Now().UnixMilli() - startTime) < int64(maxTimeMs) {
		rootAction := (*rootActions)[actionRobin%len(*rootActions)]
		currentState := applyAction(state, &rootAction)
		winner := playUntilEnd(*currentState)

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
	for _, rootAction := range *rootActions {
		rate := float64(rootResults[rootAction].wins) / float64(rootResults[rootAction].games)
		if rate > bestRate {
			bestRate = rate
			bestAction = rootAction

			debug("bestAction", bestAction, "bestRate", fmt.Sprintf("%.2f", bestRate), "(", rootResults[bestAction].wins, "/", rootResults[bestAction].games, ")")
		}
	}
	return bestAction
}

func playUntilEnd(s State) Player {
	currentState := &s
	for depth := 0; ; depth++ {
		if depth > 8*8 {
			panic("depth too high")
		}

		validActions := getValidActions(currentState)
		if len(*validActions) == 0 {
			return getOpponent(currentState.player)
		}

		randAction := randomAction(*validActions)
		applyActionMut(currentState, &randAction)
	}
}

func randomAction(validActions []Action) Action {
	return validActions[rand.Intn(len(validActions))]
}

func displayCoord(c int8) string {
	// x maps to board columns from a to h
	// y maps to board rows from 1 to 8 but reversed top to bottom
	column := string(byte('a' + c%8))
	row := strconv.Itoa(8 - int(c/8))
	return column + row
}
