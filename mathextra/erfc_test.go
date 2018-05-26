package mathextra

import "testing"

const epsilon = 1e-15

func TestErfc(t *testing.T) {
	want := 0.47950012218695346231725334610803547126354842424203629994119427 // http://www.wolframalpha.com/input/?i=Erfc%5B0.5%5D
	x := 0.5
	r := Erfc(x)
	if !Float64AlmostEq(want, r, epsilon) {
		t.Errorf("Erfc(%f) == %.15f, want %.15f", x, r, want)
	}
}

func TestErfcAlt(t *testing.T) {
	want := 0.89999999998335327818240403573915314758450635451947634620529197 // http://www.wolframalpha.com/input/?i=Erfc%5B0.0888559905091274%5D
	x := 0.0888559905091274
	r := Erfc(x)
	if !Float64AlmostEq(want, r, epsilon) {
		t.Errorf("Erfc(%f) == %.15f, want %.15f", x, r, want)
	}
}

func TestInvErfc(t *testing.T) {
	want := 0.47693627620446987338141835364313055980896974905947064470388270 // http://www.wolframalpha.com/input/?i=InvErfc%5B0.5%5D
	x := 0.5
	r := InvErfc(x)
	if !Float64AlmostEq(want, r, epsilon) {
		t.Errorf("InvErfc(%f) == %.15f, want %.15f", x, r, want)
	}
}

func TestInvErfcAlt(t *testing.T) {
	want := 0.08885599049425768701573725056779177757205224433319690375562870 // http://www.wolframalpha.com/input/?i=InvErfc%5B0.9%5D
	x := 0.9
	r := InvErfc(x)
	if !Float64AlmostEq(want, r, epsilon) {
		t.Errorf("InvErfc(%f) == %.15f, want %.15f", x, r, want)
	}
}

func TestInvErfcBelow(t *testing.T) {
	//if y >= 0.0485 && y <= 1.9515 {
	want := 1.39571526948237878870944781195435062874302991448558450474285005 // http://www.wolframalpha.com/input/?i=InvErfc%5B0.0484%5D
	x := 0.0484
	r := InvErfc(x)
	if !Float64AlmostEq(want, r, epsilon) {
		t.Errorf("InvErfc(%f) == %.15f, want %.15f", x, r, want)
	}
}

func TestInvErfcAbove(t *testing.T) {
	want := -1.39571526948237878870944781195435062874302991448558450474285005 // http://www.wolframalpha.com/input/?i=InvErfc%5B1.9516%5D
	x := 1.9516
	r := InvErfc(x)
	if !Float64AlmostEq(want, r, epsilon) {
		t.Errorf("InvErfc(%f) == %.15f, want %.15f", x, r, want)
	}
}
