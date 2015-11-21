package factor

import (
	"math"

	"github.com/mafredri/go-gaussian"
)

// VGreaterThan returns the additive correction for a single-sided truncated gaussian with unit variance
func VGreaterThan(t, epsilon float64) float64 {
	denom := gaussian.NormCdf(t - epsilon)
	if denom < 2.222758749e-162 {
		return -t + epsilon
	}

	return gaussian.NormPdf(t-epsilon) / denom
}

// WGreaterThan returns the multiplicative correction for a single-sided truncated gaussian with unit variance
func WGreaterThan(t, epsilon float64) float64 {
	var denom = gaussian.NormCdf(t - epsilon)
	if denom < 2.222758749e-162 {
		if t < 0.0 {
			return 1.0
		}
		return 0.0
	}

	vt := VGreaterThan(t, epsilon)
	return vt * (vt + t - epsilon)
}

// VWithin returns the additive correction for a double-sided truncated gaussian with unit variance
func VWithin(t, epsilon float64) float64 {
	v := math.Abs(t)
	denom := gaussian.NormCdf(epsilon-v) - gaussian.NormCdf(-epsilon-v)
	if denom < 2.222758749e-162 {
		if t < 0.0 {
			return -t - epsilon
		}
		return -t + epsilon
	}

	num := gaussian.NormPdf(-epsilon-v) - gaussian.NormPdf(epsilon-v)
	if t < 0.0 {
		return -num / denom
	}
	return num / denom
}

// WWithin returns the multiplicative correction for a double-sided truncated gaussian with unit variance
func WWithin(t, epsilon float64) float64 {
	v := math.Abs(t)
	denom := gaussian.NormCdf(epsilon-v) - gaussian.NormCdf(-epsilon-v)
	if denom < 2.222758749e-162 {
		return 1.0
	}

	vt := VWithin(t, epsilon)
	return vt*vt + ((epsilon-v)*gaussian.NormPdf(epsilon-v)-(-epsilon-v)*gaussian.NormPdf(-epsilon-v))/denom
}
