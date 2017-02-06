package bitcounter

import "github.com/willf/bitset"

// BitCounter ..
type BitCounter struct {
	vars   *bitset.BitSet // wich variables are representend
	cardin *[]int         // cardinality of all variables
	vals   []*valToLine   // all the possible values for a variable
}

type valToLine map[int]*bitset.BitSet

// NewBitCounter creates new BitCounter
func NewBitCounter() *BitCounter {
	b := new(BitCounter)
	return b
}

// LoadFromData initializes the BitCounter from a given dataset
func (b *BitCounter) LoadFromData(dataset [][]int, cardinality []int) {
	lin, col := len(dataset), len(dataset[0])
	b.vars = bitset.New(uint(col)).Complement()
	b.vals = make([]*valToLine, col)
	aux := append([]int(nil), cardinality...)
	b.cardin = &aux
	for i, c := range cardinality {
		b.vals[i] = new(valToLine)
		*b.vals[i] = make(map[int]*bitset.BitSet)
		for j := 0; j < c; j++ {
			(*b.vals[i])[j] = bitset.New(uint(lin))
		}
	}
	for i := 0; i < lin; i++ {
		for j := 0; j < col; j++ {
			(*b.vals[j])[dataset[i][j]].Set(uint(i))
		}
	}
}

// GetOccurrences ..
func (b *BitCounter) GetOccurrences(varset *bitset.BitSet) []int {
	// TODO: add caching
	v := make([]int, 0)
	assig := make([]int, varset.Count())
	for assig != nil {
		v = append(v, b.countAssignment(assig, varset))
		b.nextAssignment(&assig, varset)
	}
	return v
}

func (b *BitCounter) nextAssignment(assig *[]int, varset *bitset.BitSet) {
	i := 0
	(*assig)[i]++
	j, _ := varset.NextSet(0)
	for (*assig)[i] == (*b.cardin)[j] {
		(*assig)[i] = 0
		i++
		if i >= len(*assig) {
			*assig = nil
			return
		}
		j, _ = varset.NextSet(j + 1)
		(*assig)[i]++
	}
}

func (b *BitCounter) countAssignment(assig []int, varset *bitset.BitSet) int {
	j, _ := varset.NextSet(0)
	aux := (*b.vals[j])[assig[0]].Clone()
	for i := 1; i < len(assig); i++ {
		j, _ = varset.NextSet(j + 1)
		aux.InPlaceIntersection((*b.vals[j])[assig[i]])
	}
	return int(aux.Count())
}

// TODO: revise/remove bellow

// Marginalize ..
func (b *BitCounter) Marginalize(vars ...int) (r *BitCounter) {
	r = b.Clone()
	auxvars := bitset.New(r.vars.Len())
	for _, v := range vars {
		auxvars.Set(uint(v))
	}
	r.vars.InPlaceIntersection(auxvars)
	return
}

// SumOut ..
func (b *BitCounter) SumOut(x int) (r *BitCounter) {
	r = b.Clone()
	r.vars.Clear(uint(x))
	return
}

// Clone ..
func (b *BitCounter) Clone() (r *BitCounter) {
	r = new(BitCounter)
	r.cardin = b.cardin
	r.vals = b.vals
	r.vars = b.vars.Clone()
	return
}

// ValueIterator ..
func (b *BitCounter) ValueIterator() (f func() *int) {
	val := make([]int, b.vars.Count())
	f = func() *int {
		if val == nil {
			return nil
		}
		v := b.getCount(val)
		b.nextValuation(&val)
		return &v
	}
	return
}

// GetCountSlice ..
func (b *BitCounter) GetCountSlice() []int {
	valoration := make([]int, b.vars.Count())
	values := []int(nil)
	for valoration != nil {
		values = append(values, b.getCount(valoration))
		b.nextValuation(&valoration)
	}
	return values
}

// ValueIteratorNonZero ..
func (b *BitCounter) ValueIteratorNonZero() (f func() *int) {
	val := make([]int, b.vars.Count())
	f = func() *int {
		var v int
		for val != nil && v == 0 {
			v = b.getCount(val)
			b.nextValuation(&val)
		}
		if v != 0 {
			return &v
		}
		return nil
	}
	return
}

func (b *BitCounter) nextValuation(val *[]int) {
	i := 0
	(*val)[i]++
	j, _ := b.vars.NextSet(0)
	for (*val)[i] == (*b.cardin)[j] {
		(*val)[i] = 0
		i++
		if i == len(*val) {
			*val = nil
			return
		}
		j, _ = b.vars.NextSet(j + 1)
		(*val)[i]++
	}
}

func (b *BitCounter) getCount(val []int) int {
	j, _ := b.vars.NextSet(0)
	aux := (*b.vals[j])[val[0]].Clone()
	for i := 1; i < len(val); i++ {
		j, _ = b.vars.NextSet(j + 1)
		aux.InPlaceIntersection((*b.vals[j])[val[i]])
	}
	return int(aux.Count())
}

func (b *BitCounter) getCardinality(x int) int {
	return (*b.cardin)[x]
}
