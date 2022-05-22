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

	startTime := time.Now().UnixMilli()
	bestAction, bestValue := runMCTSSearch(state, startTime, MAX_TIME_MS_LOCAL)
	debug("best", bestAction, "value", bestValue, "playouts", playouts)
}
