package bitset

import "github.com/willf/bitset"

// BitSet is a set of bits
type BitSet struct {
	*bitset.BitSet
}

// New creates new bitset pointer with optional hint size
func New(size ...uint) *BitSet {
	var x uint
	if len(size) != 0 {
		x = size[0]
	}
	b := &BitSet{bitset.New(x)}
	return b
}

/* add these functions
// UnionSlice Returns an int slice with the union of both slices given
func UnionSlice(a []int, b []int, size int) []int {
	c := make([]int, 0, len(a)+len(b))
	varset := bitset.New(uint(size))
	SetFromSlice(varset, a)
	SetFromSlice(varset, b)
	v, ok := varset.NextSet(0)
	for ok {
		c = append(c, int(v))
		v, ok = varset.NextSet(v + 1)
	}
	return c
}

// SetSubtract Returns a Slice with the result of subtraction A - B
func SetSubtract(a, b *bitset.BitSet) []int {
	return SliceFromSet(a.Difference(b))
}

// SliceFromSet ..
func SliceFromSet(varset *bitset.BitSet) []int {
	c := make([]int, 0, varset.Count())
	v, ok := varset.NextSet(0)
	for ok {
		c = append(c, int(v))
		v, ok = varset.NextSet(v + 1)
	}
	return c
}

// SetFromSlice ..
func SetFromSlice(varset *bitset.BitSet, vars []int) {
	for _, u := range vars {
		varset.Set(uint(u))
	}
}*/
