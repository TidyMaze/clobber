package main

import (
	"testing"
)

func Test(t *testing.T) {
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
		winner: 0,
		player: WhitePlayer,
	}

	rootNode := MCTSNode{state, nil, 0, 0, nil, []*MCTSNode{}}
	best := searchMCTS(&rootNode, state.player, 1000).action

	debug("best", displayCoord(best.From)+displayCoord(best.To))
}
