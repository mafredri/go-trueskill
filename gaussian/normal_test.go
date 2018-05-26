package gaussian

import (
	"testing"

	"github.com/mafredri/go-trueskill/mathextra"
)

func TestNormCdf(t *testing.T) {
	want := 0.4775931529267244 // http://www.wolframalpha.com/input/?i=CDF%5BNormalDistribution%5B0%2C1%5D%2C-0.0561951988878394%5D
	x := -0.0561951988878394
	r := NormCdf(x)
	if !mathextra.Float64AlmostEq(want, r, 1e-16) {
		t.Errorf("NormCdf(%f) == %.16f, want %.16f", x, r, want)
	}
}

func TestNormCdfAlt(t *testing.T) {
	want := 0.691462 // http://www.wolframalpha.com/input/?i=CDF%5BNormalDistribution%5B0%2C1%5D%2C+0.5%5D
	x := 0.5
	r := NormCdf(x)
	if !mathextra.Float64AlmostEq(want, r, 1e-6) {
		t.Errorf("NormCdf(%f) == %.6f, want %.6f", x, r, want)
	}
}

func TestNormPdf(t *testing.T) {
	want := 0.13118422734394626 // http://www.wolframalpha.com/input/?i=PDF%5BNormalDistribution%5B0%2C+1%5D%2C+1.4914517054932300%5D
	x := 1.4914517054932300
	r := NormPdf(x)

	if !mathextra.Float64AlmostEq(want, r, 1e-16) {
		t.Errorf("NormPdf(%f) == %.16f, want %.16f", x, r, want)
	}
}

func TestNormPdfAlt(t *testing.T) {
	want := 0.352065 // http://www.wolframalpha.com/input/?i=PDF%5BNormalDistribution%5B0%2C+1%5D%2C+0.5%5D
	x := 0.5
	r := NormPdf(x)

	if !mathextra.Float64AlmostEq(want, r, 1e-6) {
		t.Errorf("NormPdf(%f) == %.6f, want %.6f", x, r, want)
	}
}

func TestNormPpf(t *testing.T) {
	want := -1.28155 // http://www.wolframalpha.com/input/?i=InverseCDF%5BNormalDistribution%5B0%2C+1%5D%2C+0.1%5D
	x := 0.1
	r := NormPpf(x)

	if !mathextra.Float64AlmostEq(want, r, 1e-5) {
		t.Errorf("NormPpf(%f) == %.5f, want %.5f", x, r, want)
	}
}

func TestNormPpfAlt(t *testing.T) {
	want := -0.125661 // http://www.wolframalpha.com/input/?i=InverseCDF%5BNormalDistribution%5B0%2C+1%5D%2C+0.45%5D
	x := 0.45
	r := NormPpf(x)

	if !mathextra.Float64AlmostEq(want, r, 1e-6) {
		t.Errorf("NormPpf(%f) == %.6f, want %.6f", x, r, want)
	}
}
