package collection

import (
	"testing"

	"github.com/mafredri/go-gaussian"
)

func TestDistributionBag(t *testing.T) {
	wantIdx := 0
	wantItem := gaussian.NewFromPrecision(0, 0)

	prior := gaussian.NewFromPrecision(0, 0)
	bag := NewDistributionBag(prior)
	idx := bag.NextIndex()
	item := bag.Get(idx)

	if idx != wantIdx {
		t.Errorf("NextIndex() == %#v, want %#v", idx, wantIdx)
	}
	if !item.Equals(wantItem) {
		t.Errorf("Get(%d) == %#v, want %#v", idx, item, wantItem)
	}
}

func TestDistributionBagPutAndGet(t *testing.T) {
	wantItem := gaussian.NewFromMeanAndStdDev(25.0, 25.0/3.0)

	prior := gaussian.NewFromPrecision(0, 0)
	bag := NewDistributionBag(prior)
	idx := bag.NextIndex()
	bag.Put(idx, wantItem)

	item := bag.Get(idx)

	if !item.Equals(wantItem) {
		t.Errorf("Get(%d) == %#v, want %#v", idx, item, wantItem)
	}
}
