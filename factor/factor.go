package factor

// Factor .
type Factor struct {
	UpdateMessage    func(i int) float64
	LogNormalization func() float64
	NumMessages      int
	ResetMarginals   func()
	SendMessage      func(i int) float64
}

// List .
type List struct {
	list []Factor
}

// NewList .
func NewList() List {
	return List{[]Factor{}}
}

// Add .
func (fl *List) Add(f Factor) Factor {
	fl.list = append(fl.list, f)
	return f
}

// LogNormalization .
func (fl List) LogNormalization() float64 {
	for _, f := range fl.list {
		f.ResetMarginals()
	}

	var sumLogZ float64
	for _, f := range fl.list {
		for j := 0; j < f.NumMessages; j++ {
			sumLogZ += f.SendMessage(j)
			// log.Printf("sumLogZ %f", sumLogZ)
		}
	}

	var sumLogS float64
	for _, f := range fl.list {
		sumLogS += f.LogNormalization()
		// log.Printf("sumLogS %f", sumLogS)
	}

	return sumLogZ + sumLogS
}
