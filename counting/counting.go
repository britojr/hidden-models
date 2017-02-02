// Package counting implements Counting Tables
package counting

// Table interface for counting tables
type Table interface {
	Reduce()
}

// SparseTable is a sparse implementation of counting table
type SparseTable struct {
	strideMap   map[int]int
	countMap    map[int]int
	cardinality map[int]int
	varOrdering []int
}

// NewSparse creates new sparse table
func NewSparse() *SparseTable {
	return new(SparseTable)
}

// LoadFromData initializes a new sparse table from data provided as an int matrix
func (s *SparseTable) LoadFromData(data [][]int, cardinality []int) {
	s.strideMap = make(map[int]int)
	s.countMap = make(map[int]int)
	s.cardinality = make(map[int]int)
	s.varOrdering = make([]int, len(cardinality))
	stride := 1
	for k, v := range cardinality {
		s.strideMap[k] = stride
		s.cardinality[k] = v
		s.varOrdering[k] = k
		stride *= v
	}
	for i := 0; i < len(data); i++ {
		pos := 0
		for j := 0; j < len(data[i]); j++ {
			pos += data[i][j] * s.strideMap[j]
		}
		s.countMap[pos]++
	}
}

// Reduce creates new table summing out the variable with the lesser stride
func (s *SparseTable) Reduce() (r *SparseTable) {
	r = NewSparse()
	r.varOrdering = append([]int(nil), s.varOrdering[1:]...)
	stride := 1
	r.strideMap = make(map[int]int)
	r.cardinality = make(map[int]int)
	for _, v := range r.varOrdering {
		r.strideMap[v] = stride
		stride *= s.cardinality[v]
		r.cardinality[v] = s.cardinality[v]
	}
	r.countMap = make(map[int]int)
	c := s.cardinality[s.varOrdering[0]]
	for i := 0; i < stride; i++ {
		aux := 0
		for j := 0; j < c; j++ {
			aux += s.countMap[(i*c)+j]
		}
		if aux != 0 {
			r.countMap[i] = aux
		}
	}
	return
}

// Eliminate creates new contingecy table summing out the given variable
func (s *SparseTable) Eliminate(x int) (r *SparseTable) {
	r = NewSparse()
	r.varOrdering = make([]int, len(s.varOrdering)-1)
	r.cardinality = make(map[int]int)
	r.strideMap = make(map[int]int)
	j := 0
	stride := 1
	for _, v := range s.varOrdering {
		if v != x {
			r.varOrdering[j] = v
			r.cardinality[v] = s.cardinality[v]
			r.strideMap[v] = stride
			stride *= s.cardinality[v]
			j++
		}
	}
	r.countMap = make(map[int]int)
	c := s.cardinality[x]
	for k := 0; k < stride; k += (s.strideMap[x] * c) {
		for i := 0; i < s.strideMap[x]; i++ {
			base := i + k
			aux := 0
			for j := 0; j < c; j++ {
				aux += s.countMap[base+j*s.strideMap[x]]
			}
			if aux != 0 {
				r.countMap[i] = aux
			}
		}
	}
	return
}

// Marginalize creates new contingecy containing only the given variables
func (s *SparseTable) Marginalize(vars ...int) *SparseTable {
	panic("Not implemented")
}

// sumOut sum out a given variable
func (s *SparseTable) sumOut(x int) (r *SparseTable) {
	r = NewSparse()
	tableSize := 1
	for k, v := range s.cardinality {
		if k != x {
			tableSize *= v
		}
	}
	stride := 1
	r.strideMap = make(map[int]int)
	for _, v := range s.varOrdering {
		if v != x {
			r.strideMap[v] = stride
			stride *= s.cardinality[v]
		}
	}
	r.countMap = make(map[int]int)
	c := s.cardinality[x]
	for i := 0; i < tableSize; i++ {
		aux := 0
		for j := 0; j < c; j++ {
			aux += s.countMap[(i*c)+j]
		}
		if aux != 0 {
			r.countMap[i] = aux
		}
	}
	return
}
