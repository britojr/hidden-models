package factor

import (
	"github.com/britojr/kbn/assignment"
	"github.com/britojr/kbn/utils"
)

// Factor ..
type Factor struct {
	cardin  []int
	varlist []int
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

// Get ..
func (f *Factor) Get(assig assignment.Assignment) float64 {
	x := 0
	for i := range assig {
		x += assig.Value(i) * f.stride[assig.Var(i)]
	}
	return f.values[x]
}

// Product ..
func (f *Factor) Product(g *Factor) *Factor {
	h := new(Factor)
	h.cardin = f.cardin
	h.varlist = utils.UnionSlice(f.varlist, g.varlist)
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

// SumOut ..
func (f *Factor) SumOut(x int) *Factor {
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
	for k := 0; k < len(f.values); k += sp {
		for i := 0; i < s; i++ {

		}
	}
	return h
}
