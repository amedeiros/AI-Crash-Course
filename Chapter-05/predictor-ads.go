package main

import (
	"fmt"
	"log"
	"math/rand"
	"os"
	"strconv"

	"encoding/csv"

	"gonum.org/v1/gonum/floats"
	"gonum.org/v1/gonum/mathext"
)

var (
	R = rand.New(rand.NewSource(99))
)


func main() {
	csvFile, err := os.Open("ads.csv")
	if err != nil {
		fmt.Println(err)
	}
	defer csvFile.Close()

	csvLines, err := csv.NewReader(csvFile).ReadAll()
	if err != nil {
		log.Fatal(err)
	}

	// Total ad strategies.
	d := len(csvLines[0])
	// Total customers
	N := len(csvLines)

	// Rewards
	nPosRewards := make([]float64, d)
	nNegRewards := make([]float64, d)

	// Taking our ad through beta distibution and updating its losses and wins
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

		v, _ := strconv.ParseInt(csvLines[i][selected], 10, 64)
		if v == 1 {
			nPosRewards[selected] += 1
		} else {
			nNegRewards[selected] += 1
		}
	}

	// // Find the most selected
	nSelected := make([]float64, len(nPosRewards))
	nSelected = floats.AddTo(nSelected, nPosRewards, nNegRewards)

	for i := 0; i < len(nSelected); i++ {
		log.Printf("Ad number %d was selected %d times", i+1, int64(nSelected[i]))
	}

	log.Printf("Conclusion: Best ad is ad number %d", floats.MaxIdx(nSelected)+1)
}
