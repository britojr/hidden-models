package contingency

import "github.com/willf/bitset"

// Table is a static implementation of contingency table
type Table struct {
	variables   *bitset.BitSet // wich variables are representend
	cardinality []int          // cardinality of all variables
	stride      []int          // stride of all variables
	occurrences []int          // number of occurrences of each attribution
	size        int            // number of possible assignments
}

// NewTable ..
func NewTable(variables *bitset.BitSet, cardinality []int) *Table {
	t := new(Table)
	t.variables = variables
	t.cardinality = cardinality

	var i, j uint
	i, _ = variables.NextSet(0)
	t.stride[i] = 1
	j, ok := variables.NextSet(i + 1)
	for ok {
		t.stride[j] = t.stride[i] * t.cardinality[i]
		i = j
		j, ok = variables.NextSet(i + 1)
	}
	t.size = t.stride[i] * t.cardinality[i]
	return t
}

// SetOccurrences ..
func (t *Table) SetOccurrences(occurrences []int) {
	t.occurrences = occurrences
}

// Size ..
func (t *Table) Size() int {
	return len(t.occurrences)
}

// Get ..
func (t *Table) Get(i int) int {
	return t.occurrences[i]
}

// GetOccurrences ..
func (t *Table) GetOccurrences() []int {
	return t.occurrences
}

// SumOut ..
func (t *Table) SumOut(x int) (r *Table) {
	assignment = make([]int, t.variables.Count())
	j := 0
	for i := 0; l < t.Size(); i++ {

	}

	// TODO: fix and test this sumout
	stride := 1
	index := 0
	for i, v := range t.variables {
		if v == x {
			index = i
			break
		}
		stride *= (*t.cardinality)[v]
	}
	values := []int(nil)
	base := 0
	v := 0
	for j := 0; j < (*t.cardinality)[x]; j++ {
		v += t.occurrences[base+(stride*j)]
	}
	values = append(values, v)

	auxvar := append([]int(nil), t.variables[:index]...)
	auxvar = append(auxvar, t.variables[index+1:]...)
	r = NewTable(auxvar, values, t.cardinality)
	return
}
