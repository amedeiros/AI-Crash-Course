package main

import (
	"errors"
	"log"
	"math"
	"math/rand"
)

var (
	CONVERSION_RATES = [5]float64{0.15, 0.04, 0.13, 0.11, 0.05}
	R                = rand.New(rand.NewSource(99))
)

const (
	N = int64(10_000)
	d = int64(len(CONVERSION_RATES))
)

func main() {
	// Create the dataset
	X := zeroMatrix(N, d)
	for i := int64(0); i < N; i++ {
		for j := int64(0); j < d; j++ {
			if R.NormFloat64() < CONVERSION_RATES[j] {
				X[i][j] = 1.0
			}
		}
	}

	// Rewards
	nPosRewards := make([]float64, d)
	nNegRewards := make([]float64, d)

	// Taking our best slot machine through beta distibution and updating its losses and wins
	for i := int64(0); i < N; i++ {
		selected := int64(0)
		maxRandom := 0.0
		for j := int64(0); j < d; j++ {
			randomBeta := Beta(nPosRewards[j]+1, nNegRewards[j]+1)
			if randomBeta > maxRandom {
				maxRandom = randomBeta
				selected = j
			}
		}
		if X[i][selected] == 1 {
			nPosRewards[selected] += 1
		} else {
			nNegRewards[selected] += 1
		}
	}

	nSelected, err := addFloatLists(nPosRewards, nNegRewards)
	if err != nil {
		log.Fatal(err)
	}

	for i := 0; i < len(nSelected); i++ {
		log.Printf("Machine number %d was selected %d times", i+1, int64(nSelected[i]))
	}

	_, idx := max(nSelected)
	log.Printf("Conclusion: Best machine is machine number %d", idx+1)
}

// Create a zero based matrix similar to numpy `np.zeros((N, d))`
func zeroMatrix(rows, columns int64) [][]float64 {
	matrix := make([][]float64, rows)
	for i := int64(0); i < rows; i++ {
		matrix[i] = make([]float64, columns)
	}

	return matrix
}

func addFloatLists(a, b []float64) ([]float64, error) {
	if len(a) != len(b) {
		return nil, errors.New("both lists must be equal")
	}

	r := make([]float64, len(a))
	for i := 0; i < len(a); i++ {
		r[i] = a[i] + b[i]
	}

	return r, nil
}

func max(a []float64) (float64, int64) {
	m := a[0]
	index := 0

	for i := 0; i < len(a); i++ {
		if a[i] > m {
			m = a[i]
			index = i
		}
	}

	return m, int64(index)
}

// STOLE THIS BETA AND LBETA FROM https://github.com/gonum/gonum/blob/master/mathext/beta.go
func Beta(a, b float64) float64 {
	return math.Exp(Lbeta(a, b))
}

func Lbeta(a, b float64) float64 {
	switch {
	case math.IsInf(a, +1) || math.IsInf(b, +1):
		return math.NaN()
	case a == 0 && b == 0:
		return math.NaN()
	case a < 0 || b < 0:
		return math.NaN()
	case math.IsNaN(a) || math.IsNaN(b):
		return math.NaN()
	case a == 0 || b == 0:
		return math.Inf(+1)
	}

	la, _ := math.Lgamma(a)
	lb, _ := math.Lgamma(b)
	lab, _ := math.Lgamma(a + b)
	return la + lb - lab
}
