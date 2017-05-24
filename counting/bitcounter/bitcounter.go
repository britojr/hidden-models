package bitcounter

import (
	"fmt"

	"github.com/britojr/kbn/assignment"
	"github.com/britojr/kbn/list"
	"github.com/willf/bitset"
)

// BitCounter manages the counting occurrences of sets of variables in a dataset
type BitCounter struct {
	cardin []int            // cardinality of each variable
	values []*valToLine     // all assignable values for each variable
	cache  map[string][]int // cached occurrence counting slices for different varlists
	lines  int              // number of lines in the dataset
}

type valToLine map[int]*bitset.BitSet

// NewBitCounter creates new BitCounter
func NewBitCounter() *BitCounter {
	return new(BitCounter)
}

// LoadFromData initializes the BitCounter from a given dataset and cardinality array
func (b *BitCounter) LoadFromData(dataset [][]int, cardinality []int) {
	lin, col := len(dataset), len(dataset[0])
	b.values = make([]*valToLine, col)
	b.cardin = append([]int(nil), cardinality...)
	for i, c := range cardinality {
		b.values[i] = new(valToLine)
		*b.values[i] = make(map[int]*bitset.BitSet)
		for j := 0; j < c; j++ {
			(*b.values[i])[j] = bitset.New(uint(lin))
		}
	}
	for i := 0; i < lin; i++ {
		for j := 0; j < col; j++ {
			(*b.values[j])[dataset[i][j]].Set(uint(i))
		}
	}
	b.lines = lin
	// initialize empty cache
	b.cache = make(map[string][]int)

}

// CountAssignments returns array with the counting of each possible assignment
// of the given set of variables
func (b *BitCounter) CountAssignments(varlist []int) (v []int) {
	if len(varlist) <= 0 {
		return
	}
	strvarlist := fmt.Sprint(varlist)
	v, ok := b.cache[strvarlist]
	if !ok {
		assig := assignment.New(varlist, b.cardin)
		for assig.Next() {
			if count, ok := b.Count(assig); ok {
				v = append(v, count)
			} else {
				return
			}
		}
		b.cache[strvarlist] = v
	}
	return
}

// Count returns number of lines
func (b *BitCounter) Count(assig *assignment.Assignment) (int, bool) {
	setlist := make([]*bitset.BitSet, 0, len(assig.Variables()))
	for i := range assig.Variables() {
		if assig.Var(i) < len(b.cardin) {
			setlist = append(setlist, (*b.values[assig.Var(i)])[assig.Value(i)])
		}
	}
	if len(setlist) > 0 {
		return int(list.IntersectionBits(setlist).Count()), true
	}
	// TODO: what to send when the clique is all of hidden variables?
	return -1, false
	//return b.lines, false
	//return 0, false
}

// Cardinality returns cardinality slice
func (b *BitCounter) Cardinality() []int {
	return b.cardin
}

// NumTuples returns number of lines
func (b *BitCounter) NumTuples() int {
	return b.lines
}
