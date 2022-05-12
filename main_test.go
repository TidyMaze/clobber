package main

import (
	"math/rand"
	"testing"
	"time"
)

func Test(t *testing.T) {
	// random seed to current datetime
	rand.Seed(time.Now().UnixNano())

	var startGrid = Grid{
		White, Black, White, Black, White, Black, White, Black,
		Black, White, Black, White, Black, White, Black, White,
		White, Black, White, Black, White, Black, White, Black,
		Black, White, Black, White, Black, White, Black, White,
		White, Black, White, Black, White, Black, White, Black,
		Black, White, Black, White, Black, White, Black, White,
		White, Black, White, Black, White, Black, White, Black,
		Black, White, Black, White, Black, White, Black, White,
	}

	state := State{
		grid:   startGrid,
		turn:   1,
		player: WhitePlayer,
	}

	rootNode := MCTSNode{node_count, &state, nil, 0, 0, nil, []MCTSNode{}}
	node_count++

	startTime := time.Now().UnixMilli()

	bestNode := searchMCTS(&rootNode, startTime, MAX_TIME_MS_LOCAL)
	best := bestNode.action
	debug("best", displayCoord(best.From)+displayCoord(best.To), showNode(bestNode), "after", playouts, "playouts")
}
