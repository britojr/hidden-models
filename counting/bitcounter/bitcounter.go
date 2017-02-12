package bitcounter

import (
	"github.com/britojr/kbn/assignment"
	"github.com/britojr/kbn/utils"
	"github.com/willf/bitset"
)

// BitCounter manages the counting occurrences of sets of variables in a dataset
type BitCounter struct {
	varset *bitset.BitSet   // wich variables are representend
	cardin []int            // cardinality of each variable
	values []*valToLine     // all assignable values for each variable
	cache  map[string][]int // cached occurence counting slices for different varsets
}

type valToLine map[int]*bitset.BitSet

// NewBitCounter creates new BitCounter
func NewBitCounter() *BitCounter {
	return new(BitCounter)
}

// LoadFromData initializes the BitCounter from a given dataset and cardinality array
func (b *BitCounter) LoadFromData(dataset [][]int, cardinality []int) {
	lin, col := len(dataset), len(dataset[0])
	// varset containing all variables
	b.varset = bitset.New(uint(col)).Complement()
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
	// initialize empty cache
	b.cache = make(map[string][]int)
}

// GetOccurrences returns array with the counting of each possible assignment
// of the given set of variables
func (b *BitCounter) GetOccurrences(varlist []int) (v []int) {
	if len(varlist) <= 0 {
		return
	}
	varset := utils.NewBitSetFromSlice(len(b.cardin), varlist)
	varsetstring := varset.String()
	v, ok := b.cache[varsetstring]
	if !ok {
		assig := assignment.New(varlist, b.cardin)
		for assig != nil {
			v = append(v, b.countAssignment(assig))
			assig.Next()
		}
		b.cache[varsetstring] = v
	}
	return
}

func (b *BitCounter) countAssignment(assig assignment.Assignment) int {
	aux := (*b.values[assig.Var(0)])[assig.Value(0)].Clone()
	for i := 1; i < len(assig); i++ {
		aux.InPlaceIntersection((*b.values[assig.Var(i)])[assig.Value(i)])
	}
	return int(aux.Count())
}
