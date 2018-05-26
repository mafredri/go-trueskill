package gaussian

import (
	"math"
	"testing"

	"github.com/mafredri/go-trueskill/mathextra"
)

const epsilon = 1e-13 // Precision for floating point comparison

func sq(x float64) float64 { return x * x }

func TestGaussianNewFromMeanAndStdDev(t *testing.T) {
	mu := 25.0
	sigma := 25.0 / 3.0
	g := NewFromMeanAndStdDev(mu, sigma)

	// Values taken from Ralf Herbrich's F# TrueSkill implementation
	wantPrecisionMean := 0.36
	wantPrecision := 0.0144
	wantMean := 25.0
	wantStdDev := 8.33333333333333
	wantVariance := 69.4444444444445

	if !mathextra.Float64AlmostEq(wantPrecisionMean, g.PrecisionMean, epsilon) {
		t.Errorf("PrecisionMean == %.13f, want %.13f", g.PrecisionMean, wantPrecisionMean)
	}

	if !mathextra.Float64AlmostEq(wantPrecision, g.Precision, epsilon) {
		t.Errorf("Precision == %.13f, want %.13f", g.Precision, wantPrecision)
	}

	if !mathextra.Float64AlmostEq(wantMean, g.Mean(), epsilon) {
		t.Errorf("Mean() == %.13f, want %.13f", g.Mean(), wantMean)
	}

	if !mathextra.Float64AlmostEq(wantStdDev, g.StdDev(), epsilon) {
		t.Errorf("StdDev() == %.13f, want %.13f", g.StdDev(), wantStdDev)
	}

	if !mathextra.Float64AlmostEq(wantVariance, g.Variance(), epsilon) {
		t.Errorf("Variance() == %.13f, want %.13f", g.Variance(), wantVariance)
	}
}

func TestGaussianMult(t *testing.T) {
	// Values taken from Ralf Herbrich's F# TrueSkill implementation
	a := NewFromPrecision(0.2879769618430530, 0.0115190784737221)
	b := NewFromPrecision(0.2319194248106950, 0.0055296783672657)

	prod := a.Mul(b)

	wantPrecisionMean := 0.5198963866537470
	wantPrecision := 0.0170487568409878

	if !mathextra.Float64AlmostEq(wantPrecisionMean, prod.PrecisionMean, epsilon) {
		t.Errorf("PrecisionMean == %.13f, want %.13f", prod.PrecisionMean, wantPrecisionMean)
	}

	if !mathextra.Float64AlmostEq(wantPrecision, prod.Precision, epsilon) {
		t.Errorf("Precision == %.13f, want %.13f", prod.Precision, wantPrecision)
	}
}

func TestGaussianMultNormalAndShifted(t *testing.T) {
	// From moserware/Skills:
	// > Verified against the formula at http://www.tina-vision.net/tina-knoppix/tina-memo/2003-003.pdf
	norm := NewFromMeanAndStdDev(0, 1)
	shifted := NewFromMeanAndStdDev(2, 3)

	prod := norm.Mul(shifted)

	wantMean := 0.2
	wantStdDev := 3.0 / math.Sqrt(10)

	if !mathextra.Float64AlmostEq(wantMean, prod.Mean(), epsilon) {
		t.Errorf("Mean() == %.13f, want %.13f", prod.Mean(), wantMean)
	}

	if !mathextra.Float64AlmostEq(wantStdDev, prod.StdDev(), epsilon) {
		t.Errorf("StdDev() == %.13f, want %.13f", prod.StdDev(), wantStdDev)
	}
}

func TestGaussianMult4567(t *testing.T) {
	// From moserware/Skills:
	// > Verified against the formula at http://www.tina-vision.net/tina-knoppix/tina-memo/2003-003.pdf
	a := NewFromMeanAndStdDev(4, 5)
	b := NewFromMeanAndStdDev(6, 7)

	prod := a.Mul(b)

	wantMean := (4*sq(7) + 6*sq(5)) / (sq(5) + sq(7))
	wantStdDev := math.Sqrt((sq(5) * sq(7)) / (sq(5) + sq(7)))

	if !mathextra.Float64AlmostEq(wantMean, prod.Mean(), epsilon) {
		t.Errorf("Mean() == %.13f, want %.13f", prod.Mean(), wantMean)
	}

	if !mathextra.Float64AlmostEq(wantStdDev, prod.StdDev(), epsilon) {
		t.Errorf("StdDev() == %.13f, want %.13f", prod.StdDev(), wantStdDev)
	}
}

func TestGaussianDiv(t *testing.T) {
	// Values taken from Ralf Herbrich's F# TrueSkill implementation
	a := NewFromPrecision(0.5198963866537470, 0.0170487568409878)
	b := NewFromPrecision(0.2319194248106950, 0.0055296783672657)

	prodDiv := a.Div(b)

	wantPrecisionMean := 0.2879769618430530
	wantPrecision := 0.0115190784737221

	if !mathextra.Float64AlmostEq(wantPrecisionMean, prodDiv.PrecisionMean, epsilon) {
		t.Errorf("Mean() == %.13f, want %.13f", prodDiv.PrecisionMean, wantPrecisionMean)
	}

	if !mathextra.Float64AlmostEq(wantPrecision, prodDiv.Precision, epsilon) {
		t.Errorf("StdDev() == %.13f, want %.13f", prodDiv.Precision, wantPrecision)
	}
}

func TestGaussianDivNormal(t *testing.T) {
	// From moserware/Skills:
	// > Since the multiplication was worked out by hand, we use the same numbers but work backwards
	prod := NewFromMeanAndStdDev(0.2, 3.0/math.Sqrt(10))
	norm := NewFromMeanAndStdDev(0, 1)

	prodDiv := prod.Div(norm)

	wantMean := 2.0
	wantStdDev := 3.0

	if !mathextra.Float64AlmostEq(wantMean, prodDiv.Mean(), epsilon) {
		t.Errorf("Mean() == %.13f, want %.13f", prodDiv.Mean(), wantMean)
	}

	if !mathextra.Float64AlmostEq(wantStdDev, prodDiv.StdDev(), epsilon) {
		t.Errorf("StdDev() == %.13f, want %.13f", prodDiv.StdDev(), wantStdDev)
	}
}

func TestGaussianDiv45(t *testing.T) {
	// From moserware/Skills:
	// > Since the multiplication was worked out by hand, we use the same numbers but work backwards
	prodMu := (4*sq(7) + 6*sq(5)) / (sq(5) + sq(7))
	prodSigma := math.Sqrt(((sq(5) * sq(7)) / (sq(5) + sq(7))))
	prod := NewFromMeanAndStdDev(prodMu, prodSigma)
	m4s5 := NewFromMeanAndStdDev(4, 5)

	prodDiv := prod.Div(m4s5)

	wantMean := 6.0
	wantStdDev := 7.0

	if !mathextra.Float64AlmostEq(wantMean, prodDiv.Mean(), epsilon) {
		t.Errorf("Mean() == %.13f, want %.13f", prodDiv.Mean(), wantMean)
	}

	if !mathextra.Float64AlmostEq(wantStdDev, prodDiv.StdDev(), epsilon) {
		t.Errorf("StdDev() == %.13f, want %.13f", prodDiv.StdDev(), wantStdDev)
	}
}

func TestGaussianAbsDiff(t *testing.T) {
	// Values taken from Ralf Herbrich's F# TrueSkill implementation
	a := NewFromPrecision(0.5198963866537470, 0.0170487568409878)
	b := NewFromPrecision(0.2879769618430530, 0.0115190784737221)
	diff := AbsDiff(a, b)

	want := 0.2319194248106950

	if !mathextra.Float64AlmostEq(want, diff, epsilon) {
		t.Errorf("AbsDiff(norm, norm) == %.13f, want %.13f", diff, want)
	}

}

func TestGaussianAbsDiffNormal(t *testing.T) {
	// From moserware/Skills:
	// > Verified with Ralf Herbrich's F# implementation
	norm := NewFromMeanAndStdDev(0, 1)
	diff := AbsDiff(norm, norm)

	want := 0.0

	if !mathextra.Float64AlmostEq(want, diff, epsilon) {
		t.Errorf("AbsDiff(norm, norm) == %.13f, want %.13f", diff, want)
	}
}

func TestGaussianAbsDiff1234(t *testing.T) {
	// From moserware/Skills:
	// > Verified with Ralf Herbrich's F# implementation
	a := NewFromMeanAndStdDev(1, 2)
	b := NewFromMeanAndStdDev(3, 4)
	diff := AbsDiff(a, b)

	want := 0.4330127018922193

	if !mathextra.Float64AlmostEq(want, diff, epsilon) {
		t.Errorf("AbsDiff(norm, norm) == %.13f, want %.13f", diff, want)
	}
}

func TestLogProdNorm(t *testing.T) {
	// Values taken from Ralf Herbrich's F# TrueSkill implementation
	a := NewFromPrecision(0.2879769618430530, 0.0115190784737221)
	b := NewFromPrecision(0.0445644935525882, 0.0055296783672657)
	logZ := LogProdNorm(a, b)

	want := -4.2499118707392800

	if !mathextra.Float64AlmostEq(logZ, want, epsilon) {
		t.Errorf("LogProdNorm(a, a) == %.13f, want %.13f", logZ, want)
	}
}

func TestLogProdNormNormal(t *testing.T) {
	// From moserware/Skills:
	// > Verified with Ralf Herbrich's F# implementation
	norm := NewFromMeanAndStdDev(0, 1)
	logZ := LogProdNorm(norm, norm)

	want := -1.2655121234846454

	if !mathextra.Float64AlmostEq(logZ, want, epsilon) {
		t.Errorf("LogProdNorm(a, a) == %.13f, want %.13f", logZ, want)
	}
}

func TestLogProdNorm1234(t *testing.T) {
	// From moserware/Skills:
	// > Verified with Ralf Herbrich's F# implementation
	a := NewFromMeanAndStdDev(1, 2)
	b := NewFromMeanAndStdDev(3, 4)
	logZ := LogProdNorm(a, b)

	want := -2.5168046699816684

	if !mathextra.Float64AlmostEq(logZ, want, epsilon) {
		t.Errorf("LogProdNorm(a, a) == %.13f, want %.13f", logZ, want)
	}
}

func TestLogRatioNorm(t *testing.T) {
	// Values taken from Ralf Herbrich's F# TrueSkill implementation
	a := NewFromPrecision(0.5198963866537470, 0.0170487568409878)
	b := NewFromPrecision(0.2879769618430530, 0.0115190784737221)
	logZ := LogRatioNorm(a, b)

	want := 4.2499118707392800

	if !mathextra.Float64AlmostEq(logZ, want, epsilon) {
		t.Errorf("LogProdNorm(a, a) == %.13f, want %.13f", logZ, want)
	}
}

func TestLogRatioNorm1234(t *testing.T) {
	// From moserware/Skills:
	// > Verified with Ralf Herbrich's F# implementation
	a := NewFromMeanAndStdDev(1, 2)
	b := NewFromMeanAndStdDev(3, 4)
	logZ := LogRatioNorm(a, b)

	want := 2.6157405972171204

	if !mathextra.Float64AlmostEq(logZ, want, epsilon) {
		t.Errorf("LogProdNorm(a, a) == %.13f, want %.13f", logZ, want)
	}
}
