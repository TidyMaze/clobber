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

	rootNode := MCTSNode{node_count, state, nil, 0, 0, nil, []*MCTSNode{}}
	node_count++

	best := searchMCTS(&rootNode, state.player, 5).action

	debug("best", displayCoord(best.From)+displayCoord(best.To))
}
