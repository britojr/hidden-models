package bitcounter

import "github.com/willf/bitset"

// BitCounter ..
type BitCounter struct {
	vars  map[int]*valToLine
	order []int
}

type valToLine map[int]*bitset.BitSet

// NewBitCounter creates new BitCounter
func NewBitCounter() *BitCounter {
	b := new(BitCounter)
	b.vars = make(map[int]*valToLine)
	b.order = make([]int, 0)
	return b
}

// LoadFromData initializes the BitCounter from a given dataset
func (b *BitCounter) LoadFromData(dataset [][]int, cardinality []int) {
	for i, c := range cardinality {
		b.order = append(b.order, i)
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
	r = NewBitCounter()
	for _, v := range vars {
		r.vars[v] = b.vars[v]
		r.order = append(r.order, v)
	}
	return
}

// SumOut ..
func (b *BitCounter) SumOut(x int) (r *BitCounter) {
	r = NewBitCounter()
	for _, v := range b.order {
		if v == x {
			continue
		}
		r.vars[v] = b.vars[v]
		r.order = append(r.order, v)
	}
	return
}

// ValueIterator ..
func (b *BitCounter) ValueIterator() (f func() *int) {
	val := make([]int, len(b.order))
	f = func() *int {
		if val == nil {
			return nil
		}
		v := b.getCount(val)
		val = b.nextValuation(val)
		return &v
	}
	return
}

func (b *BitCounter) nextValuation(val []int) []int {
	i := 0
	val[i]++
	for val[i] == b.getCardinality(b.order[i]) {
		val[i] = 0
		i++
		if i == len(val) {
			return nil
		}
		val[i]++
	}
	return val
}

func (b *BitCounter) getCardinality(x int) int {
	return len(*b.vars[x])
}

func (b *BitCounter) getCount(val []int) int {
	aux := (*b.vars[b.order[0]])[val[0]].Clone()
	for i, v := range val {
		aux.InPlaceIntersection((*b.vars[b.order[i]])[v])
	}
	return int(aux.Count())
}