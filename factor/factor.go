package factor

// Factor is a factor capable of updating the factor graph.
type Factor struct {
	UpdateMessage    func(i int) float64
	LogNormalization func() float64
	NumMessages      int
	ResetMarginals   func()
	SendMessage      func(i int) float64
}

// List is a list of all factors, used to get the log normalization for the factor graph.
type List struct {
	list []Factor
}

// NewList returns a new list of factors.
func NewList() List {
	return List{[]Factor{}}
}

// Add a factor to the graph.
func (fl *List) Add(f Factor) Factor {
	fl.list = append(fl.list, f)
	return f
}

// LogNormalization returns the log normalization of all factors in the factor graph.
func (fl List) LogNormalization() float64 {
	for _, f := range fl.list {
		f.ResetMarginals()
	}

	var sumLogZ float64
	for _, f := range fl.list {
		for j := 0; j < f.NumMessages; j++ {
			sumLogZ += f.SendMessage(j)
		}
	}

	var sumLogS float64
	for _, f := range fl.list {
		sumLogS += f.LogNormalization()
	}

	return sumLogZ + sumLogS
}
