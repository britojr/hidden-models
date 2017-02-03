// Package counting implements Counting Tables
package counting

import "github.com/willf/bitset"

// BitCounter ..
type BitCounter struct {
	vars map[int]*valToLine
}

type valToLine map[int]*bitset.BitSet

// NewBitCounter creates new BitCounter
func NewBitCounter() *BitCounter {
	return new(BitCounter)
}

// LoadFromData initializes the BitCounter from a given dataset
func (b *BitCounter) LoadFromData(dataset [][]int, cardinality []int) {
	b.vars = make(map[int]*valToLine)
	for i, c := range cardinality {
		b.vars[i] = new(valToLine)
		*b.vars[i] = make(map[int]*bitset.BitSet)
		for j := 0; j < c; j++ {
			(*b.vars[i])[j] = bitset.New(uint(len(dataset)))
		}
	}
	for i := 0; i < len(dataset); i++ {
		for j := 0; j < len(dataset[0]); j++ {
			(*b.vars[j])[dataset[i][j]].Set(uint(i))
		}
	}
}

// Marginalize ..
func (b *BitCounter) Marginalize(vars ...int) (r *BitCounter) {
	panic("Not implemented")
}

// SumOut ..
func (b *BitCounter) SumOut(vars ...int) (r *BitCounter) {
	panic("Not implemented")
}

// ValueIterator ..
func (b *BitCounter) ValueIterator() func() *int {
	panic("Not implemented")
}
