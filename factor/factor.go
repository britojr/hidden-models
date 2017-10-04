package factor

import (
	"errors"
	"math"
	"math/rand"
	"time"

	"github.com/britojr/kbn/assignment"
	"github.com/britojr/kbn/list"
	"github.com/britojr/utl/stats"
)

// Factor ..
type Factor struct {
	cardin  []int
	varlist []int // list of variables starting from the lowest stride
	stride  map[int]int
	values  []float64
}

// NewFactorValues creates a new factor with specified values
func NewFactorValues(varlist []int, cardin []int, values []float64) *Factor {
	f := new(Factor)
	f.cardin = cardin
	f.varlist = varlist
	f.values = values
	f.stride = makeStride(varlist, cardin)
	return f
}

// NewFactor creates a factor with zero values
func NewFactor(varlist []int, cardin []int) *Factor {
	f := new(Factor)
	f.cardin = cardin
	f.varlist = varlist
	f.stride = makeStride(varlist, cardin)
	size := 1
	if len(f.varlist) > 0 {
		size = f.cardin[f.varlist[len(f.varlist)-1]] * f.stride[f.varlist[len(f.varlist)-1]]
	}
	f.values = make([]float64, size)
	return f
}

func makeStride(varlist, cardin []int) map[int]int {
	stride := make(map[int]int)
	if len(varlist) > 0 {
		stride[varlist[0]] = 1
		for i := 1; i < len(varlist); i++ {
			stride[varlist[i]] = cardin[varlist[i-1]] * stride[varlist[i-1]]
		}
	}
	return stride
}

// Variables returns the list of variables
func (f *Factor) Variables() []int {
	return f.varlist
}

// Cardinality returns the cardinalities
func (f *Factor) Cardinality() []int {
	return f.cardin
}

// Values returns the slice of values
func (f *Factor) Values() []float64 {
	return f.values
}

// Get returnrs the value corresponding to the given assignment
func (f *Factor) Get(assig *assignment.Assignment) float64 {
	return f.values[assig.Index(f.stride)]
}

// GetEvidValue returnrs the value corresponding to the given evidence
func (f *Factor) GetEvidValue(evid []int) float64 {
	x := 0
	for _, v := range f.varlist {
		x += evid[v] * f.stride[v]
	}
	return f.values[x]
}

// SetValues updates the slice of values
func (f *Factor) SetValues(values []float64) {
	f.values = values
}

// SetUniform changes factor value to uniformly distributed
func (f *Factor) SetUniform() *Factor {
	for i := range f.values {
		f.values[i] = 1.0 / float64(len(f.values))
	}
	return f
}

// SetRandom sets the factor with normalized random values
func (f *Factor) SetRandom() *Factor {
	rand.Seed(time.Now().UTC().UnixNano())
	for i := range f.values {
		f.values[i] = rand.Float64()
	}
	stats.Normalize(f.values)
	return f
}

// SetDirichlet sets the factor with normalized Dirichlet distribution
func (f *Factor) SetDirichlet(alpha float64) *Factor {
	stats.Dirichlet1(alpha, f.values)
	return f
}

// ClearCopy creates a copy factor with zero values
func (f *Factor) ClearCopy() *Factor {
	h := new(Factor)
	h.cardin = f.cardin
	h.varlist = f.varlist
	h.stride = f.stride
	h.values = make([]float64, len(f.values))
	return h
}

// Clone returns a copy of the current factor
func (f *Factor) Clone() *Factor {
	h := new(Factor)
	h.cardin = f.cardin
	h.varlist = f.varlist
	h.stride = f.stride
	h.values = append([]float64(nil), f.values...)
	return h
}

// Set a value to the current assignment
func (f *Factor) Set(assig *assignment.Assignment, v float64) {
	f.values[assig.Index(f.stride)] = v
}

// Add add a value to the current assignment
func (f *Factor) Add(assig *assignment.Assignment, v float64) {
	f.values[assig.Index(f.stride)] += v
}

// Product multiply two factors and return a new factor as the result
func (f *Factor) Product(g *Factor) *Factor {
	h := new(Factor)
	h.cardin = f.cardin
	h.varlist = list.Union(f.varlist, g.varlist, uint(len(f.cardin)))
	h.stride = makeStride(h.varlist, h.cardin)
	size := 1
	if len(h.varlist) > 0 {
		size = h.cardin[h.varlist[len(h.varlist)-1]] * h.stride[h.varlist[len(h.varlist)-1]]
	}
	h.values = make([]float64, size)
	// TODO: create version without assig an bench to see how better it is
	assig := assignment.New(h.varlist, h.cardin)
	for i := range h.values {
		assig.Next()
		h.values[i] = f.Get(assig) * g.Get(assig)
	}
	return h
}

// Division divide factor f by factor g and return a new factor as the result
func (f *Factor) Division(g *Factor) *Factor {
	h := new(Factor)
	h.cardin = f.cardin
	h.varlist = f.varlist
	h.stride = f.stride
	h.values = append([]float64(nil), f.values...)
	_, in, _ := list.OrderedDiff(f.varlist, g.varlist)
	g = g.SumOut(in)
	// TODO: create version without assig an bench to see how better it is
	assig := assignment.New(h.varlist, h.cardin)
	for i := range h.values {
		assig.Next()
		v := g.Get(assig)
		if v != 0 {
			h.values[i] /= v
		} else {
			h.values[i] = 0
		}
	}
	return h
}

// SumOutOne returns a factor with the given variable summed out
func (f *Factor) SumOutOne(x int) *Factor {
	h := new(Factor)
	h.cardin = f.cardin
	h.varlist = make([]int, 0, len(f.varlist)-1)
	var c, s, sp int
	for _, v := range f.varlist {
		if v == x {
			c = f.cardin[x]
			s = f.stride[x]
			sp = c * s
			continue
		}
		h.varlist = append(h.varlist, v)
	}
	h.stride = makeStride(h.varlist, h.cardin)
	size := 1
	if len(h.varlist) > 0 {
		size = h.cardin[h.varlist[len(h.varlist)-1]] * h.stride[h.varlist[len(h.varlist)-1]]
	}
	h.values = make([]float64, size)
	if sp > 0 {
		index := 0
		for k := 0; k < len(f.values); k += sp {
			for i := 0; i < s; i++ {
				for j := 0; j < c; j++ {
					h.values[index] += f.values[k+i+(j*s)]
				}
				index++
			}
		}
	} else {
		copy(h.values, f.values)
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

// Reduce mutes out every value that is not consistent with a given evidence tuple
func (f *Factor) Reduce(evid []int) *Factor {
	h := new(Factor)
	h.cardin = f.cardin
	h.varlist = f.varlist
	h.stride = f.stride
	h.values = make([]float64, len(f.values))
	// TODO: think of better way to to this part
	k := 0
	free := []int(nil)
	for _, v := range h.varlist {
		if v >= len(evid) || evid[v] == -1 {
			free = append(free, v)
		} else {
			k += h.stride[v] * evid[v]
		}
	}
	if len(free) > 0 {
		assig := assignment.New(free, h.cardin)
		for assig.Next() {
			i := k + assig.Index(h.stride)
			h.values[i] = f.values[i]
		}
	} else {
		h.values[k] = f.values[k]
	}
	return h
}

// Normalize normalizes the factor so all values sum to 1
func (f *Factor) Normalize() *Factor {
	if len(f.values) > 0 {
		stats.Normalize(f.values)
	}
	return f
}

// MaxDifference calculates the max difference between two lists of factors
func MaxDifference(f, g []*Factor) (diff float64, num, val int, err error) {
	for i := range f {
		if f[i] == nil && g[i] == nil {
			continue
		}
		if !(f[i] != nil && g[i] != nil) {
			err = errors.New("incompatible list of factors")
			return
		}
		q := f[i].Values()
		for j, v := range g[i].Values() {
			if d := math.Abs(q[j] - v); d > diff {
				diff = d
				num, val = i, j
			}
		}
	}
	return
}
