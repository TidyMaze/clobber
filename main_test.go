package main

import (
	"testing"
	"time"
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

	best := runMonteCarloSearch(state, time.Now().UnixMilli(), MAX_TIME_MS_LOCAL)
	debug("best", displayCoord(best.From)+displayCoord(best.To))
}
