// Package utils provides general use functions
package utils

import (
	"fmt"
	"log"
	"math"
	"strconv"

	"github.com/willf/bitset"
)

const epslon = 1e-10

// FuzzyEqual compares float numbers with a tolerance
func FuzzyEqual(a, b float64) bool {
	if math.Abs(a-b) < epslon {
		return true
	}
	return false
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
}

// ErrCheck validates error and prints a log message
func ErrCheck(err error, message string) {
	if err != nil {
		log.Printf("%v: err(%v)\n", message, err)
		panic(err)
	}
}
