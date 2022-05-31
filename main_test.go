package main

import (
	"testing"
	"time"
)

// func Test(t *testing.T) {
// 	// random seed to current datetime
// 	rand.Seed(time.Now().UnixNano())

// 	// empty
// 	// whites
// 	// blacks
// 	findBestActionFromInitialGrid(1000000)
// }

const InitialWhiteGrid uint64 = 0b1010101001010101101010100101010110101010010101011010101001010101
const InitialBlackGrid uint64 = 0b0101010110101010010101011010101001010101101010100101010110101010

const MaxIterationsBench = 5000

var startGrid = Grid{
	0,
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
	rootNode := &MCTSNode{uint32(nodeCount), state, Action{-1, -1}, 0, 0, nil, []*MCTSNode{}}
	_, _ = runMCTSSearch(rootNode, startTime, MaxTimeMsLocal, maxIterations)
	// debug("best", bestAction, "value", bestValue, "playouts", playouts)
}

func BenchmarkFindBestActionFromInitialGrid(b *testing.B) {
	initNeighborsCache()
	initMaskCache()
	for i := 0; i < b.N; i++ {
		findBestActionFromInitialGrid(MaxIterationsBench)
	}
}
