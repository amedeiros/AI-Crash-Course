package main

import (
	"log"
	"math/rand"

	"gonum.org/v1/gonum/floats"
	"gonum.org/v1/gonum/mat"
	"gonum.org/v1/gonum/mathext"
)

var (
	CONVERSION_RATES = [5]float64{0.15, 0.04, 0.13, 0.11, 0.05}
	R                = rand.New(rand.NewSource(99))
)

const (
	N = 10_000
	d = len(CONVERSION_RATES)
)

func main() {
	// Create the dataset
	X := mat.NewDense(N, d, nil)
	for i := 0; i < N; i++ {
		for j := 0; j < d; j++ {
			if R.NormFloat64() < CONVERSION_RATES[j] {
				X.Set(i, j, 1.0)
			}
		}
	}

	// Rewards
	nPosRewards := make([]float64, d)
	nNegRewards := make([]float64, d)

	// Taking our best slot machine through beta distibution and updating its losses and wins
	for i := 0; i < N; i++ {
		selected := 0
		maxRandom := 0.0
		for j := 0; j < d; j++ {
			randomBeta := mathext.Beta(nPosRewards[j]+1, nNegRewards[j]+1)
			if randomBeta > maxRandom {
				maxRandom = randomBeta
				selected = j
			}
		}

		if X.At(i, selected) == 1 {
			nPosRewards[selected] += 1
		} else {
			nNegRewards[selected] += 1
		}
	}

	// Find the most selected
	nSelected := make([]float64, len(nPosRewards))
	nSelected = floats.AddTo(nSelected, nPosRewards, nNegRewards)

	for i := 0; i < len(nSelected); i++ {
		log.Printf("Machine number %d was selected %d times", i+1, int64(nSelected[i]))
	}

	log.Printf("Conclusion: Best machine is machine number %d", floats.MaxIdx(nSelected)+1)
}
