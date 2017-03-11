// Package utils provides general use functions
package utils

import (
	"fmt"
	"log"
	"math"
	"strconv"

	"github.com/britojr/kbn/assignment"
	"github.com/willf/bitset"
)

// Counter returns a counting value for an assignment
type Counter interface {
	Count(assig *assignment.Assignment) (count int, ok bool)
	CountAssignments(varlist []int) []int
	Cardinality() []int
	NumTuples() int
}

const epslon = 1e-10

// FuzzyEqual compares float numbers with a tolerance
func FuzzyEqual(a, b float64, delta ...float64) bool {
	eps := epslon
	if len(delta) > 0 {
		eps = delta[0]
	}
	if math.Abs(a-b) < eps {
		return true
	}
	return false
}

// ErrCheck validates error and prints a log message
func ErrCheck(err error, message string) {
	if err != nil {
		log.Printf("%v: err(%v)\n", message, err)
		panic(err)
	}
}

// SliceAtoi creates an int slice from a string slice
func SliceAtoi(ss []string) []int {
	arr := make([]int, len(ss))
	var err error
	for k, v := range ss {
		arr[k], err = strconv.Atoi(v)
		ErrCheck(err, fmt.Sprintf("Can't convert %v to int", v))
	}
	return arr
}

// SliceItoU64 creates an uint64 array from an int array
func SliceItoU64(is []int) []uint64 {
	arr := make([]uint64, len(is))
	for i, v := range is {
		arr[i] = uint64(v)
	}
	return arr
}

// SliceItoF64 creates an float64 array from an int array
func SliceItoF64(is []int) []float64 {
	arr := make([]float64, len(is))
	for i, v := range is {
		arr[i] = float64(v)
	}
	return arr
}

// SliceSplit creates two slices a with elements < n and b with elements >=n
func SliceSplit(s []int, n int) ([]int, []int) {
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

// SetSlice sets all the bits given on the slice
func SetSlice(b *bitset.BitSet, varlist []int) {
	for _, u := range varlist {
		b.Set(uint(u))
	}
}

// ClearSlice clears all the bits given on the slice
func ClearSlice(b *bitset.BitSet, varlist []int) {
	for _, u := range varlist {
		b.Clear(uint(u))
	}
}

// SliceSumFloat64 returns the sum of all slice values
func SliceSumFloat64(values []float64) (sum float64) {
	for _, v := range values {
		sum += v
	}
	return
}

// NormalizeSlice normalizes the slice so all values sum to one
func NormalizeSlice(values []float64) {
	var sum float64
	for _, v := range values {
		sum += v
	}
	if sum == 0 {
		panic("trying to normalize zero factor")
	}
	for i, v := range values {
		values[i] = v / sum
	}
}

// NormalizeIntSlice create a new normalized slice from an int slice
func NormalizeIntSlice(values []int) []float64 {
	var sum float64
	for _, v := range values {
		sum += float64(v)
	}
	norm := make([]float64, len(values))
	if sum == 0 {
		return norm
	}
	for i, v := range values {
		norm[i] = float64(v) / sum
	}
	return norm
}

// SliceFromBitSet returns the corresponding int slice from a set
func SliceFromBitSet(b *bitset.BitSet) []int {
	s := make([]int, 0, b.Count())
	for i, ok := b.NextSet(0); ok; i, ok = b.NextSet(i + 1) {
		s = append(s, int(i))
	}
	return s
}

// ListIntersection creates an intersection of a list of bitsets
func ListIntersection(setlist []*bitset.BitSet) *bitset.BitSet {
	if len(setlist) == 0 {
		return NewBitSet()
	}
	r := setlist[0].Clone()
	for i := 1; i < len(setlist); i++ {
		r.InPlaceIntersection(setlist[i])
	}
	return r
}

// NewBitSet creates new bitset pointer with optional hint size
func NewBitSet(size ...uint) *bitset.BitSet {
	var x uint
	if len(size) != 0 {
		x = size[0]
	}
	return bitset.New(x)
}

// SliceUnion returns the union of slices
func SliceUnion(l1, l2 []int, size ...uint) []int {
	b := NewBitSet(size...)
	SetSlice(b, l1)
	SetSlice(b, l2)
	return SliceFromBitSet(b)
}

// SliceDifference returns the difference of slices (l1 and not l2)
func SliceDifference(l1, l2 []int, size ...uint) []int {
	b := NewBitSet(size...)
	SetSlice(b, l1)
	ClearSlice(b, l2)
	return SliceFromBitSet(b)
}
