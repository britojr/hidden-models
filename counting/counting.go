// Package counting implements Counting Tables
package counting

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

// Marginalize creates new contingecy containing only the given variables
func (s *SparseTable) Marginalize(vars ...int) *SparseTable {
	panic("Not implemented")
}

// Reduce creates new contingecy summing out the variable with stride of one
func (s *SparseTable) Reduce() (r *SparseTable) {
	//r = NewSparse()
	x := s.varOrdering[0]
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
