package gaussian

import (
	"math"

	"github.com/mafredri/go-trueskill/mathextra"
)

// NormCdf returns the cumulative gaussian distribution (cdf) at the point of interest.
func NormCdf(t float64) float64 {
	return mathextra.Erfc(-t/math.Sqrt2) / 2.0
}

// NormPdf returns the probability density function (pdf) at the point of interest.
func NormPdf(t float64) float64 {
	return mathextra.InvSqrt2Pi * math.Exp(-(t * t / 2.0))
}

// NormPpf returns the percent point function (ppf, the inverse of cdf) at the point of interest.
func NormPpf(p float64) float64 {
	return -math.Sqrt2 * mathextra.InvErfc(2.0*p)
}
