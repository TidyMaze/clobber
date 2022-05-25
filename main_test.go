package main

import (
	"math/rand"
	"testing"
	"time"
)

func Test(t *testing.T) {
	// random seed to current datetime
	rand.Seed(time.Now().UnixNano())

	const InitialWhiteGrid uint64 = 0b1010101001010101101010100101010110101010010101011010101001010101
	const InitialBlackGrid uint64 = 0b0101010110101010010101011010101001010101101010100101010110101010

	var startGrid = Grid{
		// empty
		0,
		// whites
		InitialWhiteGrid,
		// blacks
		InitialBlackGrid,
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
