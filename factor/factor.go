package factor

import "github.com/britojr/kbn/assignment"

// Factor ..
type Factor struct {
	cardin  []int
	varlist []int
	stride  map[int]int
	values  float64
}

// New ..
func New(varlist []int, cardin []int, values []float64) *Factor {
	f := new(Factor)
	return f
}

// Get ..
func (f *Factor) Get(assig assignment.Assignment) float64 {
	return 0.0
}

// Product ..
func (f *Factor) Product(g *Factor) *Factor {
	h := new(Factor)
	return h
}

// SumOut ..
func (f *Factor) SumOut(x int) *Factor {
	h := new(Factor)
	return h
}
