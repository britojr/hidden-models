package bitset

import "math/big"

// BitSet is a set of bits
type BitSet struct {
	s *big.Int
}

// New creates new bitset pointer
func New() *BitSet {
	b := &BitSet{new(big.Int)}
	return b
}

// Set sets the ith bit to true
func (b *BitSet) Set(i int) {
	// (*big.Int)(b).SetBit((*big.Int)(b), i, 1)
	b.s.SetBit(b.s, i, 1)
}

// Clear clears the ith bit
func (b *BitSet) Clear(i int) {
	// (*big.Int)(b).SetBit((*big.Int)(b), i, 0)
	b.s.SetBit(b.s, i, 0)
}

// Get returns the value (0 or 1) o the ith bit
func (b *BitSet) Get(i int) int {
	// return int((*big.Int)(b).Bit(i))
	return int(b.s.Bit(i))
}

// Test returns true if ith bit is set and false otherwise
func (b *BitSet) Test(i int) bool {
	// return (*big.Int)(b).Bit(i) == 1
	return b.s.Bit(i) == 1
}

// Count returns the number of bits that are set
func (b *BitSet) Count() int {
	return 0
}
