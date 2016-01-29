package trueskill

import "math"

func calculate2PlayerMatchQuality(ts Config, p1 Player, p2 Player) float64 {
	betaSquared := ts.Beta * ts.Beta
	p1SigmaSquared := p1.Sigma() * p1.Sigma()
	p2SigmaSquared := p2.Sigma() * p2.Sigma()
	pMean := p1.Mu() - p2.Mu()

	sqrt := math.Sqrt(2 * betaSquared / (2*betaSquared + p1SigmaSquared + p2SigmaSquared))
	exp := math.Exp((-1 * (pMean * pMean)) / (2 * (2*betaSquared + p1SigmaSquared + p2SigmaSquared)))

	return sqrt * exp
}
