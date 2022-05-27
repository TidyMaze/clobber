package main

import (
	"math/rand"
	"testing"
	"time"
)

func Test(t *testing.T) {
	// random seed to current datetime
	rand.Seed(time.Now().UnixNano())

	// empty
	// whites
	// blacks
	findBestActionFromInitialGrid(1000000)
}

var InitialWhiteGrid [64]bool = [64]bool{
	false, true, false, true, false, true, false, true,
	true, false, true, false, true, false, true, false,
	false, true, false, true, false, true, false, true,
	true, false, true, false, true, false, true, false,
	false, true, false, true, false, true, false, true,
	true, false, true, false, true, false, true, false,
	false, true, false, true, false, true, false, true,
	true, false, true, false, true, false, true, false,
}

var InitialBlackGrid [64]bool = [64]bool{
	true, false, true, false, true, false, true, false,
	false, true, false, true, false, true, false, true,
	true, false, true, false, true, false, true, false,
	false, true, false, true, false, true, false, true,
	true, false, true, false, true, false, true, false,
	false, true, false, true, false, true, false, true,
	true, false, true, false, true, false, true, false,
	false, true, false, true, false, true, false, true,
}

const MaxIterationsBench = 5000

var startGrid = Grid{
	[64]bool{},
	InitialWhiteGrid,
	InitialBlackGrid,
}

var state = State{
	grid:   startGrid,
	turn:   1,
	player: WhitePlayer,
}

func findBestActionFromInitialGrid(maxIterations int) {
	startTime := time.Now().UnixMilli()
	_, _ = runMCTSSearch(state, startTime, MaxTimeMsLocal, maxIterations)
	// debug("best", bestAction, "value", bestValue, "playouts", playouts)
}

// func BenchmarkFindBestActionFromInitialGrid(b *testing.B) {
// 	for i := 0; i < b.N; i++ {
// 		findBestActionFromInitialGrid(MaxIterationsBench)
// 	}
// }
