package factor

import (
	"math/rand"
	"time"

	"github.com/britojr/kbn/assignment"
	"github.com/britojr/kbn/utils"
)

// Factor ..
type Factor struct {
	cardin  []int
	varlist []int // list of variables starting from the lowest stride
	stride  map[int]int
	values  []float64
}

// New ..
func New(varlist []int, cardin []int, values []float64) *Factor {
	f := new(Factor)
	f.cardin = cardin
	f.varlist = varlist
	f.values = values
	f.stride = make(map[int]int)
	f.stride[varlist[0]] = 1
	for i := 1; i < len(varlist); i++ {
		f.stride[varlist[i]] = cardin[varlist[i-1]] * f.stride[varlist[i-1]]
	}
	return f
}

// NewFactor creates a factor with zero values
func NewFactor(varlist []int, cardin []int) *Factor {
	f := new(Factor)
	f.cardin = cardin
	f.varlist = varlist
	f.stride = make(map[int]int)
	f.stride[varlist[0]] = 1
	for i := 1; i < len(varlist); i++ {
		f.stride[varlist[i]] = cardin[varlist[i-1]] * f.stride[varlist[i-1]]
	}
	size := f.cardin[f.varlist[len(f.varlist)-1]] * f.stride[f.varlist[len(f.varlist)-1]]
	f.values = make([]float64, size)
	return f
}

// Clone returns a copy of the current factor
func (f *Factor) Clone() *Factor {
	g := NewFactor(f.varlist, f.cardin)
	g.values = append([]float64(nil), f.values...)
	return f
}

// SetValues ..
func (f *Factor) SetValues(values []float64) {
	f.values = values
}

// SetUniform changes factor value to uniformly distributed
func (f *Factor) SetUniform() {
	for i := range f.values {
		f.values[i] = 1.0 / float64(len(f.values))
	}
}

// SetRandom sets the factor with random values
func (f *Factor) SetRandom() {
	rand.Seed(time.Now().UTC().UnixNano())
	for i := range f.values {
		f.values[i] = rand.Float64()
	}
	f.Normalize()
}

// Variables ..
func (f *Factor) Variables() []int {
	return f.varlist
}

// Cardinality ..
func (f *Factor) Cardinality() []int {
	return f.cardin
}

// Get ..
func (f *Factor) Get(assig assignment.Assignment) float64 {
	x := 0
	for i := range assig {
		x += assig.Value(i) * f.stride[assig.Var(i)]
	}
	return f.values[x]
}

// Set a value to the current assignment
func (f *Factor) Set(assig assignment.Assignment, v float64) {
	x := 0
	for i := range assig {
		x += assig.Value(i) * f.stride[assig.Var(i)]
	}
	f.values[x] = v
}

// Add add a value to the current assignment
func (f *Factor) Add(assig assignment.Assignment, v float64) {
	x := 0
	for i := range assig {
		x += assig.Value(i) * f.stride[assig.Var(i)]
	}
	f.values[x] += v
}

// Product ..
func (f *Factor) Product(g *Factor) *Factor {
	h := new(Factor)
	h.cardin = f.cardin
	h.varlist = utils.SliceUnion(f.varlist, g.varlist, uint(len(f.cardin)))
	h.stride = make(map[int]int)
	h.stride[h.varlist[0]] = 1
	for i := 1; i < len(h.varlist); i++ {
		h.stride[h.varlist[i]] = h.cardin[h.varlist[i-1]] * h.stride[h.varlist[i-1]]
	}
	size := h.cardin[h.varlist[len(h.varlist)-1]] * h.stride[h.varlist[len(h.varlist)-1]]
	h.values = make([]float64, size)
	assig := assignment.New(h.varlist, h.cardin)
	for i := 0; i < size; i++ {
		h.values[i] = f.Get(assig) * g.Get(assig)
		assig.Next()
	}
	return h
}

// SumOutOne ..
func (f *Factor) SumOutOne(x int) *Factor {
	h := new(Factor)
	h.cardin = f.cardin
	h.varlist = make([]int, 0, len(f.varlist)-1)
	for _, v := range f.varlist {
		if v != x {
			h.varlist = append(h.varlist, v)
		}
	}
	h.stride = make(map[int]int)
	h.stride[h.varlist[0]] = 1
	for i := 1; i < len(h.varlist); i++ {
		h.stride[h.varlist[i]] = h.cardin[h.varlist[i-1]] * h.stride[h.varlist[i-1]]
	}
	size := h.cardin[h.varlist[len(h.varlist)-1]] * h.stride[h.varlist[len(h.varlist)-1]]
	h.values = make([]float64, size)
	c := f.cardin[x]
	s := f.stride[x]
	sp := c * s
	index := 0
	for k := 0; k < len(f.values); k += sp {
		for i := 0; i < s; i++ {
			for j := 0; j < c; j++ {
				h.values[index] += f.values[k+i+(j*s)]
			}
			index++
		}
	}
	return h
}

// SumOut ..
func (f *Factor) SumOut(vars []int) *Factor {
	// TODO: create better implementation for SumOut and Marginalize
	q := f
	for _, x := range vars {
		q = q.SumOutOne(x)
	}
	return q
}

// Marginalize ..
func (f *Factor) Marginalize(vars []int) *Factor {
	// TODO: create better implementation for SumOut and Marginalize
	diff := utils.SliceDifference(f.varlist, vars, uint(len(f.cardin)))
	q := f
	for _, x := range diff {
		q = q.SumOutOne(x)
	}
	return q
}

// Restrict mutes out every value that is not consistent with a given evidence tuple
func (f *Factor) Restrict(evid []int) *Factor {
	h := new(Factor)
	h.cardin = f.cardin
	h.varlist = f.varlist
	h.stride = f.stride
	h.values = make([]float64, len(f.values))
	assig := assignment.New(h.varlist, h.cardin)
	for i := range h.values {
		if assig.Consistent(evid) {
			h.values[i] = f.values[i]
		}
		assig.Next()
	}
	return h
}

// Normalize normalizes the factor so all values sum to 1
func (f *Factor) Normalize() {
	var tot float64
	for i := range f.values {
		tot += f.values[i]
	}
	for i := range f.values {
		f.values[i] /= tot
	}
}
