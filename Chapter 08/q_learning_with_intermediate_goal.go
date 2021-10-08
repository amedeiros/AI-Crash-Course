package main

import (
	"fmt"
	"math/rand"

	"gonum.org/v1/gonum/floats"
	"gonum.org/v1/gonum/mat"
)

// Setting the parameters gamma and alpha for Q-Learning
const (
	gamma = 0.75
	alpha = 0.9
)

var (
	R               = rand.New(rand.NewSource(99))
	locationToState = map[string]int{
		"A": 0,
		"B": 1,
		"C": 2,
		"D": 3,
		"E": 4,
		"F": 5,
		"G": 6,
		"H": 7,
		"I": 8,
		"J": 9,
		"K": 10,
		"L": 11,
	}
	stateToLocation = map[int]string{
		0:  "A",
		1:  "B",
		2:  "C",
		3:  "D",
		4:  "E",
		5:  "F",
		6:  "G",
		7:  "H",
		8:  "I",
		9:  "J",
		10: "K",
		11: "L",
	}
	actions = []int{0, 1, 2, 3, 4, 5, 6, 7, 8, 9, 10, 11}
	X       = newDenseFromMatrix([][]float64{
		{0, 1, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0},
		{1, 0, 1, 0, 0, 1, 0, 0, 0, 0, 0, 0},
		{0, 1, 0, 0, 0, 0, 1, 0, 0, 0, 0, 0},
		{0, 0, 0, 0, 0, 0, 0, 1, 0, 0, 0, 0},
		{0, 0, 0, 0, 0, 0, 0, 0, 1, 0, 0, 0},
		{0, 1, 0, 0, 0, 0, 0, 0, 0, 1, 0, 0},
		{0, 0, 1, 0, 0, 0, 1, 1, 0, 0, 0, 0},
		{0, 0, 0, 1, 0, 0, 1, 0, 0, 0, 0, 1},
		{0, 0, 0, 0, 1, 0, 0, 0, 0, 1, 0, 0},
		{0, 0, 0, 0, 0, 1, 0, 0, 1, 0, 1, 0},
		{0, 0, 0, 0, 0, 0, 0, 0, 0, 1, 0, 1},
		{0, 0, 0, 0, 0, 0, 0, 1, 0, 0, 1, 0},
	})

	trainedResults = map[string]*mat.Dense{}
)

func main() {
	for endingLocation := range locationToState {
		train(endingLocation)
	}

	fmt.Println(route("E", "K"))
}

func route(startingLocation, endingLocation string) []string {
	r := []string{startingLocation}
	nextLocation := startingLocation
	Q := trainedResults[endingLocation]
	for nextLocation != endingLocation {
		startingState := locationToState[startingLocation]
		nextState := floats.MaxIdx(mat.Row(nil, startingState, Q))
		nextLocation = stateToLocation[nextState]
		r = append(r, nextLocation)
		startingLocation = nextLocation
	}

	return r
}

func train(endingLocation string) {
	NewX := mat.DenseCopyOf(X)
	endingState := locationToState[endingLocation]
	NewX.Set(endingState, endingState, 1000)
	Q := mat.NewDense(12, 12, nil)
	for i := 0; i < 1000; i++ {
		currentState := R.Intn(12)
		var playableActions []int
		for j := 0; j < 12; j++ {
			if NewX.At(currentState, j) > 0 {
				playableActions = append(playableActions, j)
			}
		}
		nextState := playableActions[R.Intn(len(playableActions))]
		TD := NewX.At(currentState, nextState) + gamma*Q.At(nextState, floats.MaxIdx(mat.Row(nil, nextState, Q))) - Q.At(currentState, nextState)
		Q.Set(currentState, nextState, Q.At(currentState, nextState)+alpha*TD)
	}

	trainedResults[endingLocation] = Q
}

// Helper function to convert a matrice into a Dense matrix
func newDenseFromMatrix(m [][]float64) *mat.Dense {
	var data []float64
	for _, mm := range m {
		data = append(data, mm...)
	}

	return mat.NewDense(len(m), len(m[0]), data)
}
