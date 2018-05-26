package mathextra

import "math"

// Erfc returns the complementary error function of x (from math.Erfc).
func Erfc(x float64) float64 {
	return math.Erfc(x)
}

func invErfcYAboveOrBelow(q float64) (float64, float64) {
	q1 := 0.005504751339936943
	q1 = q1*q + 0.2279687217114118
	q1 = q1*q + 1.697592457770869
	q1 = q1*q + 1.802933168781950
	q1 = q1*q + -3.093354679843504
	q1 = q1*q - 2.077595676404383
	q2 := 0.007784695709041462
	q2 = q2*q + 0.3224671290700398
	q2 = q2*q + 2.445134137142996
	q2 = q2*q + 3.754408661907416
	q2 = q2*q + 1.0

	return q1, q2
}

// InvErfc returns the inverse complementary error function of y.
func InvErfc(y float64) float64 {
	switch {
	case y < 0 || y > 2 || math.IsNaN(y):
		return math.NaN()
	case y == 0:
		return math.Inf(1)
	case y == 2:
		return math.Inf(-1)
	}

	var x float64
	if y >= 0.0485 && y <= 1.9515 {
		q := y - 1.0
		r := q * q
		r1 := 0.01370600482778535
		r1 = r1*r - 0.3051415712357203
		r1 = r1*r + 1.524304069216834
		r1 = r1*r - 3.057303267970988
		r1 = r1*r + 2.710410832036097
		r1 = r1*r - 0.8862269264526915
		r2 := -0.05319931523264068
		r2 = r2*r + 0.6311946752267222
		r2 = r2*r - 2.432796560310728
		r2 = r2*r + 4.175081992982483
		r2 = r2*r - 3.320170388221430
		r2 = r2*r + 1.0
		x = r1 * q / r2
	} else if y < 0.0485 {
		q := math.Sqrt(-2.0 * math.Log(y/2.0))
		q1, q2 := invErfcYAboveOrBelow(q)
		x = q1 / q2
	} else if y > 1.9515 {
		q := math.Sqrt(-2.0 * math.Log(1.0-y/2.0))
		q1, q2 := invErfcYAboveOrBelow(q)
		x = -q1 / q2
	} else {
		x = 0.0
	}

	u := (Erfc(x) - y) / (-2.0 / math.Sqrt(math.Pi) * math.Exp(-x*x))

	return x - u/(1.0+x*u)
}
