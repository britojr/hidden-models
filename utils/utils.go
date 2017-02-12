// Package utils provides general use functions
package utils

import (
	"fmt"
	"log"
	"strconv"

	"github.com/willf/bitset"
)

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

// NewBitSetFromSlice creates new bitset with informed size and elements
func NewBitSetFromSlice(size int, vars []int) *bitset.BitSet {
	varset := bitset.New(uint(size))
	for _, u := range vars {
		varset.Set(uint(u))
	}
	return varset
}

// ErrCheck validates error and prints a log message
func ErrCheck(err error, message string) {
	if err != nil {
		log.Printf("%v: err(%v)\n", message, err)
		panic(err)
	}
}
