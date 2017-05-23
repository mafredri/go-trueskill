package trueskill

import (
	"math"

	"github.com/mafredri/go-gaussian"
)

func drawProbability(beta, drawMargin, totalPlayers float64) float64 {
	return 2*gaussian.NormCdf(drawMargin/(math.Sqrt(totalPlayers)*beta)) - 1
}

func drawMargin(beta, drawProb, totalPlayers float64) float64 {
	// totalPlayers represents the total number of players for both team A
	// and team B. Considering a match between 3 teams, A, B and C, we would
	// call drawMargin two times. One time with len(A) + len(B), the other
	// time with len(B) + len(C).
	return -math.Sqrt((totalPlayers)*beta*beta) * gaussian.NormPpf((1-drawProb)/2)
}
