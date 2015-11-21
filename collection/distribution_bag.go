package collection

import "github.com/mafredri/go-gaussian"

// DistributionBag .
type DistributionBag struct {
	prior gaussian.Gaussian
	bag   *[]gaussian.Gaussian
}

// NewDistributionBag .
func NewDistributionBag(prior gaussian.Gaussian) DistributionBag {
	return DistributionBag{prior, &[]gaussian.Gaussian{}}
}

// NextIndex .
func (db *DistributionBag) NextIndex() int {
	*db.bag = append(*db.bag, db.prior)
	return db.Len() - 1
}

// Reset .
func (db *DistributionBag) Reset() {
	db.bag = &[]gaussian.Gaussian{}
}

// Len .
func (db DistributionBag) Len() int {
	return len(*db.bag)
}

// Get .
func (db DistributionBag) Get(i int) gaussian.Gaussian {
	bag := *db.bag
	return bag[i]
}

// Put .
func (db DistributionBag) Put(i int, g gaussian.Gaussian) {
	bag := *db.bag
	bag[i] = g
}

// PutPriorAt .
func (db DistributionBag) PutPriorAt(i int) {
	db.Put(i, db.prior)
}
