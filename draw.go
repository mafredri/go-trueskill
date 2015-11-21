package trueskill

import (
	"math"

	"github.com/mafredri/go-gaussian"
)

func drawProbability(beta, drawMargin float64) float64 {
	return 2*gaussian.NormCdf(drawMargin/(math.Sqrt(1+1)*beta)) - 1
}

func drawMargin(beta, drawProb float64) float64 {
	// Number of players on each team
	n1 := 1.0
	n2 := 1.0
	return -math.Sqrt((n1+n2)*beta*beta) * gaussian.NormPpf((1.0-drawProb)/2.0)
}
