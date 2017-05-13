// Package utils provides general use functions
package utils

import (
	"fmt"
	"log"
	"math"
	"math/rand"
	"sort"
	"strconv"
	"time"

	"github.com/dtromb/gogsl/randist"
	"github.com/dtromb/gogsl/rng"
	"github.com/willf/bitset"
)

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

// Max returns the max number
func Max(xs []float64) float64 {
	v := xs[0]
	for _, x := range xs {
		if x > v {
			v = x
		}
	}
	return v
}

// Min returns the min number
func Min(xs []float64) float64 {
	v := xs[0]
	for _, x := range xs {
		if x < v {
			v = x
		}
	}
	return v
}

// Mode calculates the mode of a float64 slice
func Mode(xs []float64) (v float64) {
	d := make(map[float64]int)
	c := 0
	for _, x := range xs {
		d[x]++
		if d[x] > c {
			c = d[x]
			v = x
		}
	}
	return
}

// Median calculates the media of a float64 slice
func Median(xs []float64) (v float64) {
	aux := append([]float64(nil), xs...)
	sort.Float64s(aux)
	i := len(aux) / 2
	if len(aux)%2 != 0 {
		v = aux[i]
	} else {
		v = (aux[i] + aux[i-1]) / 2
	}
	return
}

// Mean calculates the Mean of a float64 slice
func Mean(xs []float64) (v float64) {
	// return stats.Mean(xs, 1, len(xs))
	for _, x := range xs {
		v += x
	}
	v /= float64(len(xs))
	return
}

// Variance calculates the variance of a float64 slice
func Variance(xs []float64) (v float64) {
	m := Mean(xs)
	for _, x := range xs {
		v += (m - x) * (m - x)
	}
	v /= float64(len(xs))
	return
}

// Stdev calculates the standard deviation of a float64 slice
func Stdev(xs []float64) float64 {
	// return stats.Sd(xs, 1, len(xs))
	return math.Sqrt(Variance(xs))
}

// Dirichlet sets values as a Dirichlet distribution
func Dirichlet(alpha, values []float64) {
	rand.Seed(time.Now().UTC().UnixNano())
	rng.EnvSetup()
	r := rng.RngAlloc(rng.DefaultRngType())
	rng.Set(r, rand.Int())
	randist.Dirichlet(r, len(alpha), alpha, values)
}

// Uniform sets values uniformly
func Uniform(values []float64) {
	for i := range values {
		values[i] = 1.0 / float64(len(values))
	}
}

// Random sets random values
func Random(values []float64) {
	rand.Seed(time.Now().UTC().UnixNano())
	for i := range values {
		values[i] = rand.Float64()
	}
	NormalizeSlice(values)
}

// Atoi converst string to int
func Atoi(s string) int {
	i, err := strconv.Atoi(s)
	ErrCheck(err, fmt.Sprintf("Can't convert %v to int", s))
	return i
}

// AtoF64 converst string to float64
func AtoF64(s string) float64 {
	i, err := strconv.ParseFloat(s, 64)
	ErrCheck(err, fmt.Sprintf("Can't convert %v to float64", s))
	return i
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

// SliceAtoF64 creates a float64 slice from a string slice
func SliceAtoF64(ss []string) []float64 {
	arr := make([]float64, len(ss))
	var err error
	for k, v := range ss {
		arr[k], err = strconv.ParseFloat(v, 64)
		ErrCheck(err, fmt.Sprintf("Can't convert %v to float64", v))
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

// OrderedSliceDiff returns three lists: intersection, b-a and a-b for ordered slices a,b
func OrderedSliceDiff(a, b []int) (inter, in, out []int) {
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
