package factor

import (
	"testing"

	"github.com/mafredri/go-mathextra"
)

// Test values taken from Ralf Herbrich's F# TrueSkill implementation

func TestVGreaterThan(t *testing.T) {
	want := 0.4181660649773850
	tVar := 0.7495591915280050
	eps := 0.0631282276750071
	r := VGreaterThan(tVar, eps)

	if !mathextra.Float64AlmostEq(want, r, 1e-6) {
		t.Errorf("VGreaterThan(%f, %f) == %.6f, want %.6f", tVar, eps, r, want)
	}
}

func TestWGreaterThan(t *testing.T) {
	want := 0.4619049929317120
	tVar := 0.7495591915280050
	eps := 0.0631282276750071
	r := WGreaterThan(tVar, eps)

	if !mathextra.Float64AlmostEq(want, r, 1e-6) {
		t.Errorf("WGreaterThan(%f, %f) == %.6f, want %.6f", tVar, eps, r, want)
	}
}

func TestVWithin(t *testing.T) {
	want := -0.7485644072749330
	tVar := 0.7495591915280050
	eps := 0.0631282276750071
	r := VWithin(tVar, eps)

	if !mathextra.Float64AlmostEq(want, r, 1e-6) {
		t.Errorf("VWithin(%f, %f) == %.6f, want %.6f", tVar, eps, r, want)
	}
}

func TestWWithin(t *testing.T) {
	want := 0.9986734210033660
	tVar := 0.7495591915280050
	eps := 0.0631282276750071
	r := WWithin(tVar, eps)

	if !mathextra.Float64AlmostEq(want, r, 1e-6) {
		t.Errorf("WWithin(%f, %f) == %.6f, want %.6f", tVar, eps, r, want)
	}
}
