package factor

import (
	"math"

	"github.com/mafredri/go-gaussian"
	"github.com/mafredri/go-trueskill/collection"
)

// GaussianFactors is used to perform all skill related gaussian operations and turning them into factors capable of updating the factor graph.
type GaussianFactors struct {
	msgBag collection.DistributionBag
}

// NewGaussianFactors initializes a gaussian factor with a distribution bag and returns it.
func NewGaussianFactors() GaussianFactors {
	prior := gaussian.NewFromPrecision(0, 0)
	return GaussianFactors{
		msgBag: collection.NewDistributionBag(prior),
	}
}

func sendMessageHelper(msgIdx, varIdx int, msgBag, varBag collection.DistributionBag) float64 {
	mar := varBag.Get(varIdx)
	msg := msgBag.Get(msgIdx)
	varBag.Put(varIdx, mar.Mul(msg))

	// logZ
	return gaussian.LogProdNorm(mar, msg)
}

// GaussianPrior calculates the prior for the factor graph.
func (gf GaussianFactors) GaussianPrior(mu, sigmaSquared float64, varIdx int,
	varBag collection.DistributionBag) Factor {
	msgIdx := gf.msgBag.NextIndex()
	newMsg := gaussian.NewFromMeanAndVariance(mu, sigmaSquared)

	updateMessage := func(i int) float64 {
		if i != 0 {
			panic("Index out of range")
		}

		oldMarginal := varBag.Get(varIdx)
		oldMsg := gf.msgBag.Get(msgIdx)
		newMarginal := gaussian.NewFromPrecision(oldMarginal.PrecisionMean+newMsg.PrecisionMean-oldMsg.PrecisionMean,
			oldMarginal.Precision+newMsg.Precision-oldMsg.Precision)
		varBag.Put(varIdx, newMarginal)
		gf.msgBag.Put(msgIdx, newMsg)

		delta := oldMarginal.Sub(newMarginal)
		return delta
	}
	sendMessage := func(i int) float64 {
		if i != 0 {
			panic("Index out of range")
		}

		return sendMessageHelper(msgIdx, varIdx, gf.msgBag, varBag)
	}

	return Factor{
		UpdateMessage:    updateMessage,
		LogNormalization: func() float64 { return 0 },
		NumMessages:      1,
		ResetMarginals:   func() { varBag.PutPriorAt(varIdx) },
		SendMessage:      sendMessage,
	}
}

// GaussianLikeliehood calculates the likeliehood for the factor graph.
func (gf GaussianFactors) GaussianLikeliehood(betaSquared float64, varIdx1, varIdx2 int, varBag1, varBag2 collection.DistributionBag) Factor {
	msgIdx1 := gf.msgBag.NextIndex()
	msgIdx2 := gf.msgBag.NextIndex()

	prec := 1.0 / betaSquared

	updateHelper := func(m1, m2, v1, v2 int, bag1, bag2 collection.DistributionBag) float64 {
		msg1 := gf.msgBag.Get(m1)
		msg2 := gf.msgBag.Get(m2)
		mar1 := bag1.Get(v1)
		mar2 := bag2.Get(v2)

		a := prec / (prec + mar2.Precision - msg2.Precision)
		newMsg := gaussian.NewFromPrecision(a*(mar2.PrecisionMean-msg2.PrecisionMean),
			a*(mar2.Precision-msg2.Precision))
		oldMarginalWithoutMsg := mar1.Div(msg1)
		newMarginal := oldMarginalWithoutMsg.Mul(newMsg)

		gf.msgBag.Put(m1, newMsg)
		bag1.Put(v1, newMarginal)

		delta := newMarginal.Sub(mar1)
		return delta
	}
	updateMessage := func(i int) float64 {
		switch i {
		case 0:
			return updateHelper(msgIdx1, msgIdx2, varIdx1, varIdx2, varBag1, varBag2)
		case 1:
			return updateHelper(msgIdx2, msgIdx1, varIdx2, varIdx1, varBag2, varBag1)
		default:
			panic("Index out of range")
		}
	}
	logNormalization := func() float64 {
		logNorm := gaussian.LogRatioNorm(varBag1.Get(varIdx1), gf.msgBag.Get(msgIdx1))
		return logNorm
	}
	resetMarginals := func() {
		varBag1.PutPriorAt(varIdx1)
		varBag2.PutPriorAt(varIdx2)
	}
	sendMessage := func(i int) float64 {
		switch i {
		case 0:
			return sendMessageHelper(msgIdx1, varIdx1, gf.msgBag, varBag1)
		case 1:
			return sendMessageHelper(msgIdx2, varIdx2, gf.msgBag, varBag2)
		default:
			panic("Index out of range")
		}
	}

	return Factor{
		UpdateMessage:    updateMessage,
		LogNormalization: logNormalization,
		NumMessages:      2,
		ResetMarginals:   resetMarginals,
		SendMessage:      sendMessage,
	}
}

// GaussianWeightedSum calculates the weighted sum for the facor graph.
func (gf GaussianFactors) GaussianWeightedSum(a1, a2 float64, varIdx0, varIdx1, varIdx2 int,
	varBag0, varBag1, varBag2 collection.DistributionBag) Factor {

	msgIdx0 := gf.msgBag.NextIndex()
	msgIdx1 := gf.msgBag.NextIndex()
	msgIdx2 := gf.msgBag.NextIndex()
	weights0 := []float64{a1, a2}
	weights0Squared := []float64{weights0[0] * weights0[0], weights0[1] * weights0[1]}
	weights1 := []float64{-a2 / a1, 1.0 / a1}
	weights1Squared := []float64{weights1[0] * weights1[0], weights1[1] * weights1[1]}
	weights2 := []float64{-a1 / a2, 1.0 / a2}
	weights2Squared := []float64{weights2[0] * weights2[0], weights2[1] * weights2[1]}

	updateHelper := func(w, wS []float64, m1, m2, m3, v1, v2, v3 int, bag1, bag2, bag3 collection.DistributionBag) float64 {
		d0 := bag2.Get(v2).Div(gf.msgBag.Get(m2))
		d1 := bag3.Get(v3).Div(gf.msgBag.Get(m3))
		msg1 := gf.msgBag.Get(m1)
		mar1 := bag1.Get(v1)
		denom := wS[0]*d1.Precision + wS[1]*d0.Precision
		newPrecisionMean := (w[0]*d1.Precision*d0.PrecisionMean + w[1]*d0.Precision*d1.PrecisionMean) / denom
		newPrecision := d0.Precision * d1.Precision / denom
		newMsg := gaussian.NewFromPrecision(newPrecisionMean, newPrecision)
		oldMarginalWithoutMsg := mar1.Div(msg1)
		newMarginal := oldMarginalWithoutMsg.Mul(newMsg)

		gf.msgBag.Put(m1, newMsg)
		bag1.Put(v1, newMarginal)

		return newMarginal.Sub(mar1)
	}
	updateMessage := func(i int) float64 {
		switch i {
		case 0:
			return updateHelper(weights0, weights0Squared, msgIdx0, msgIdx1, msgIdx2, varIdx0, varIdx1, varIdx2,
				varBag0, varBag1, varBag2)
		case 1:
			return updateHelper(weights1, weights1Squared, msgIdx1, msgIdx2, msgIdx0, varIdx1, varIdx2, varIdx0,
				varBag1, varBag2, varBag0)
		case 2:
			return updateHelper(weights2, weights2Squared, msgIdx2, msgIdx1, msgIdx0, varIdx2, varIdx1, varIdx0,
				varBag2, varBag1, varBag0)
		default:
			panic("Index out of range.")
		}
	}
	logNormalization := func() float64 {
		ratio1 := gaussian.LogRatioNorm(varBag1.Get(varIdx1), gf.msgBag.Get(msgIdx1))
		ratio2 := gaussian.LogRatioNorm(varBag2.Get(varIdx2), gf.msgBag.Get(msgIdx2))
		return ratio1 + ratio2
	}
	resetMarginals := func() {
		varBag0.PutPriorAt(varIdx0)
		varBag1.PutPriorAt(varIdx1)
		varBag2.PutPriorAt(varIdx2)
	}
	sendMessage := func(i int) float64 {
		switch i {
		case 0:
			return sendMessageHelper(msgIdx0, varIdx0, gf.msgBag, varBag0)
		case 1:
			return sendMessageHelper(msgIdx1, varIdx1, gf.msgBag, varBag1)
		case 2:
			return sendMessageHelper(msgIdx2, varIdx2, gf.msgBag, varBag2)
		default:
			panic("Index out of range")
		}
	}

	return Factor{
		UpdateMessage:    updateMessage,
		LogNormalization: logNormalization,
		NumMessages:      3,
		ResetMarginals:   resetMarginals,
		SendMessage:      sendMessage,
	}
}

func gaussianGreaterThanOrWithinUpdateMessage(epsilon float64, msgIdx, varIdx int,
	msgBag, varBag collection.DistributionBag, vFunc, wFunc func(t, epsilon float64) float64) float64 {
	oldMarginal := varBag.Get(varIdx)
	oldMsg := msgBag.Get(msgIdx)
	msgFromVar := oldMarginal.Div(oldMsg)
	c := msgFromVar.Precision
	d := msgFromVar.PrecisionMean
	sqrtC := math.Sqrt(c)
	dOnSqrtC := d / sqrtC
	epsTimesSqrtC := epsilon * sqrtC
	denom := 1.0 - wFunc(dOnSqrtC, epsTimesSqrtC)
	newPrecision := c / denom
	newPrecisionMean := (d + sqrtC*vFunc(dOnSqrtC, epsTimesSqrtC)) / denom
	newMarginal := gaussian.NewFromPrecision(newPrecisionMean, newPrecision)
	newMsg := oldMsg.Mul(newMarginal).Div(oldMarginal)

	msgBag.Put(msgIdx, newMsg)
	varBag.Put(varIdx, newMarginal)

	return newMarginal.Sub(oldMarginal)
}

// GaussianGreaterThan calculates the greater than margin for the factor graph.
func (gf GaussianFactors) GaussianGreaterThan(epsilon float64, varIdx int, varBag collection.DistributionBag) Factor {
	msgIdx := gf.msgBag.NextIndex()

	updateMessage := func(i int) float64 {
		if i != 0 {
			panic("Index out of range.")
		}

		return gaussianGreaterThanOrWithinUpdateMessage(epsilon, msgIdx, varIdx, gf.msgBag, varBag,
			VGreaterThan, WGreaterThan)
	}
	logNormalization := func() float64 {
		marginal := varBag.Get(varIdx)
		msg := gf.msgBag.Get(msgIdx)
		msgFromVar := marginal.Div(msg)
		logProdNorm := gaussian.LogProdNorm(msgFromVar, msg)
		return -logProdNorm + math.Log(gaussian.NormCdf((msgFromVar.Mean()-epsilon)/msgFromVar.StdDev()))
	}
	sendMessage := func(i int) float64 {
		if i != 0 {
			panic("Index out of range")
		}

		return sendMessageHelper(msgIdx, varIdx, gf.msgBag, varBag)
	}

	return Factor{
		UpdateMessage:    updateMessage,
		LogNormalization: logNormalization,
		NumMessages:      1,
		ResetMarginals:   func() { varBag.PutPriorAt(varIdx) },
		SendMessage:      sendMessage,
	}
}

// GaussianWithin calculates the within margin for the factor graph.
func (gf GaussianFactors) GaussianWithin(epsilon float64, varIdx int, varBag collection.DistributionBag) Factor {
	msgIdx := gf.msgBag.NextIndex()

	updateMessage := func(i int) float64 {
		if i != 0 {
			panic("Index out of range.")
		}

		return gaussianGreaterThanOrWithinUpdateMessage(epsilon, msgIdx, varIdx, gf.msgBag, varBag, VWithin, WWithin)
	}
	logNormalization := func() float64 {
		marginal := varBag.Get(varIdx)
		msg := gf.msgBag.Get(msgIdx)
		msgFromVar := marginal.Div(msg)
		logProdNorm := gaussian.LogProdNorm(msgFromVar, msg)
		mean := msgFromVar.Mean()
		stdDev := msgFromVar.StdDev()
		z := gaussian.NormCdf((epsilon-mean)/stdDev) - gaussian.NormCdf((-epsilon-mean)/stdDev)
		return -logProdNorm + math.Log(z)
	}
	sendMessage := func(i int) float64 {
		if i != 0 {
			panic("Index out of range")
		}

		return sendMessageHelper(msgIdx, varIdx, gf.msgBag, varBag)
	}

	return Factor{
		UpdateMessage:    updateMessage,
		LogNormalization: logNormalization,
		NumMessages:      1,
		ResetMarginals:   func() { varBag.PutPriorAt(varIdx) },
		SendMessage:      sendMessage,
	}
}
