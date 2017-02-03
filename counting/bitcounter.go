// Package counting implements Counting Tables
package counting

import "github.com/willf/bitset"

// BitCounter ..
type BitCounter struct {
	vars map[int]*valToLine
}

type valToLine map[int]bitset.BitSet

// NewBitCounter creates new BitCounter
func NewBitCounter() *BitCounter {
	panic("Not implemented")
}

// LoadFromData initializes the BitCounter from a given dataset
func (c *BitCounter) LoadFromData(dataset [][]int, cardinality []int) {
	panic("Not implemented")
}

// Marginalize ..
func (c *BitCounter) Marginalize(vars ...int) (r *BitCounter) {
	panic("Not implemented")
}

// SumOut ..
func (c *BitCounter) SumOut(vars ...int) (r *BitCounter) {
	panic("Not implemented")
}

// ValueIterator ..
func (c *BitCounter) ValueIterator() func() *int {
	panic("Not implemented")
}
