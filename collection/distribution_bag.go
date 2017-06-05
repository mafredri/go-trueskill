package collection

import "github.com/mafredri/go-gaussian"

// DistributionBag is a storage for gaussian distributions.
type DistributionBag struct {
	prior gaussian.Gaussian
	bag   []gaussian.Gaussian
}

// NewDistributionBag returns a new distribution bag.
func NewDistributionBag(prior gaussian.Gaussian) *DistributionBag {
	return &DistributionBag{prior: prior}
}

// NextIndex initializes the next free slot in the bag with the default gaussian
// and returns the index.
func (db *DistributionBag) NextIndex() int {
	db.bag = append(db.bag, db.prior)
	return db.Len() - 1
}

// Reset empties the distribution bag by replacing it with an empty one.
func (db *DistributionBag) Reset() {
	db.bag = nil
}

// Len returns the length (size) of the bag.
func (db *DistributionBag) Len() int {
	return len(db.bag)
}

// Get returns the gaussian from the bag at the given position.
func (db *DistributionBag) Get(i int) gaussian.Gaussian {
	return db.bag[i]
}

// Put a gaussian into the given position of the bag.
func (db *DistributionBag) Put(i int, g gaussian.Gaussian) {
	db.bag[i] = g
}

// PutPriorAt puts the default gaussian into the given position of the bag.
func (db *DistributionBag) PutPriorAt(i int) {
	db.Put(i, db.prior)
}
