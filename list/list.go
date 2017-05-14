package list

import "github.com/willf/bitset"

// Split creates two slices a with elements < n and b with elements >=n
func Split(s []int, n int) ([]int, []int) {
	a, b := make([]int, 0, len(s)), make([]int, 0, len(s))
	for _, v := range s {
		if v < n {
			a = append(a, v)
		} else {
			b = append(b, v)
		}
	}
	return a, b
}

// OrderedDiff returns three lists: intersection, b-a and a-b for two ordered slices a,b
func OrderedDiff(a, b []int) (inter, in, out []int) {
	n, m := len(a), len(b)
	i, j := 0, 0
	for i < n && j < m {
		switch {
		case a[i] < b[j]:
			out = append(out, a[i])
			i++
		case a[i] > b[j]:
			in = append(in, b[j])
			j++
		default:
			inter = append(inter, a[i])
			i, j = i+1, j+1
		}
	}
	for i < n {
		out = append(out, a[i])
		i++
	}
	for j < m {
		in = append(in, b[j])
		j++
	}
	return
}

// Union returns the union of two slices
func Union(l1, l2 []int, size ...uint) []int {
	b := NewBitSet(size...)
	SetBits(b, l1)
	SetBits(b, l2)
	return FromBitSet(b)
}

// Difference returns the difference of slices (l1 and not l2)
func Difference(l1, l2 []int, size ...uint) []int {
	b := NewBitSet(size...)
	SetBits(b, l1)
	ClearBits(b, l2)
	return FromBitSet(b)
}

// bits

// NewBitSet creates new bitset pointer with optional hint size
func NewBitSet(size ...uint) *bitset.BitSet {
	var x uint
	if len(size) != 0 {
		x = size[0]
	}
	return bitset.New(x)
}

// FromBitSet returns the corresponding int slice from a bitset
func FromBitSet(b *bitset.BitSet) []int {
	s := make([]int, 0, b.Count())
	for i, ok := b.NextSet(0); ok; i, ok = b.NextSet(i + 1) {
		s = append(s, int(i))
	}
	return s
}

// SetBits sets bits in all positions given on the slice
func SetBits(b *bitset.BitSet, varlist []int) {
	for _, u := range varlist {
		b.Set(uint(u))
	}
}

// ClearBits clears bits in all positions given on the slice
func ClearBits(b *bitset.BitSet, varlist []int) {
	for _, u := range varlist {
		b.Clear(uint(u))
	}
}

// IntersectionBits creates an intersection of a list of bitsets
func IntersectionBits(setlist []*bitset.BitSet) *bitset.BitSet {
	if len(setlist) == 0 {
		return NewBitSet()
	}
	r := setlist[0].Clone()
	for i := 1; i < len(setlist); i++ {
		r.InPlaceIntersection(setlist[i])
	}
	return r
}
