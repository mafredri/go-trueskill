// Package gaussian provides
package gaussian

import (
	"math"

	"github.com/mafredri/go-trueskill/mathextra"
)

// Gaussian represents a gaussian based on a precision and a precision adjusted mean.
type Gaussian struct {
	PrecisionMean float64 // PrecisionMean (pi, π = μ/σ^2) is the precision adjusted mean.
	Precision     float64 // Precision (tau, τ = 1/σ2) is the inverse of the variance.
}

// NewFromMeanAndStdDev create a new gaussian from the mean and standard deviation.
func NewFromMeanAndStdDev(mean, stdDev float64) Gaussian {
	variance := stdDev * stdDev
	return NewFromMeanAndVariance(mean, variance)
}

// NewFromMeanAndVariance create a new gaussian from the mean and variance.
func NewFromMeanAndVariance(mean, variance float64) Gaussian {
	return NewFromPrecision(mean/variance, 1.0/variance)
}

// NewFromPrecision create a new gaussian from the precision adjusted mean and the precision.
func NewFromPrecision(precisionMean, precision float64) Gaussian {
	return Gaussian{precisionMean, precision}
}

// Mean returns the mean (mu, μ) of a gaussian by dividing the precision adjusted mean with the precision.
func (a Gaussian) Mean() float64 {
	return a.PrecisionMean / a.Precision
}

// Variance is the variance (σ^2) of a gaussian derived from dividing one with the precision.
func (a Gaussian) Variance() float64 {
	return 1.0 / a.Precision
}

// StdDev is the standard deviation (sigma, σ) of a gaussian derived from the square root of the variance.
func (a Gaussian) StdDev() float64 {
	return math.Sqrt(a.Variance())
}

// Mul multiplies two gaussians by adding their precision adjusted means and precisions.
func (a Gaussian) Mul(b Gaussian) Gaussian {
	return NewFromPrecision(a.PrecisionMean+b.PrecisionMean, a.Precision+b.Precision)
}

// Div divides two gaussians by subtracting the precision adjusted means and precisions.
func (a Gaussian) Div(b Gaussian) Gaussian {
	return NewFromPrecision(a.PrecisionMean-b.PrecisionMean, a.Precision-b.Precision)
}

// Sub subtracts two gaussians by returning their absolute difference.
func (a Gaussian) Sub(b Gaussian) float64 {
	return AbsDiff(a, b)
}

// Equals compares the precision adjusted mean and precision of two gaussians. Returns true if they are equal.
func (a Gaussian) Equals(b Gaussian) bool {
	return a.PrecisionMean == b.PrecisionMean && a.Precision == b.Precision
}

// AbsDiff returns the largest difference between the subtracted precision adjusted means and the precisions of two gaussians.
func AbsDiff(a, b Gaussian) float64 {
	pm := math.Abs(a.PrecisionMean - b.PrecisionMean)
	p := math.Sqrt(math.Abs(a.Precision - b.Precision))
	return math.Max(pm, p)
}

// LogProdNorm returns the log product normalization of two gaussians.
func LogProdNorm(a, b Gaussian) float64 {
	if a.Precision == 0.0 || b.Precision == 0.0 {
		return 0.0
	}

	varSum := a.Variance() + b.Variance()
	meanDiff := a.Mean() - b.Mean()

	return -mathextra.LogSqrt2Pi - math.Log(varSum)/2.0 - meanDiff*meanDiff/(2.0*varSum)
}

// LogRatioNorm returns the log ratio normalization of two gaussians.
func LogRatioNorm(a, b Gaussian) float64 {
	if a.Precision == 0.0 || b.Precision == 0.0 {
		return 0.0
	}

	bVar := b.Variance()
	varDiff := bVar - a.Variance()
	if varDiff == 0.0 {
		return 0.0
	}

	meanDiff := a.Mean() - b.Mean()
	return math.Log(bVar) + mathextra.LogSqrt2Pi - math.Log(varDiff)/2.0 + meanDiff*meanDiff/(2.0*varDiff)
}
