// Package contingecytable implements Contingecy Table
package contingecytable

// Sparse contingey table
type Sparse struct {
	strideMap   map[int]int
	countMap    map[int]int
	varOrdering []int
	cardinality map[int]int
}

// LoadFromData creates a new contingey table from data
func LoadFromData(data [][]int, cardinality []int) *Sparse {
	sp := Sparse{}
	stride := 1
	sp.strideMap = make(map[int]int)
	sp.cardinality = make(map[int]int)
	sp.varOrdering = make([]int, len(cardinality))
	for k, v := range cardinality {
		sp.strideMap[k] = stride
		stride *= v
		sp.cardinality[k] = v
		sp.varOrdering[k] = k
	}
	sp.countMap = make(map[int]int)
	for i := 0; i < len(data); i++ {
		pos := 0
		for j := 0; j < len(data[i]); j++ {
			pos += data[i][j] * sp.strideMap[j]
		}
		sp.countMap[pos]++
	}
	return &sp
}

// Marginalize creates new contingecy containing only the given variables
func (sp *Sparse) Marginalize(vars ...int) *Sparse {
	mt := Sparse{}
	return &mt
}

// Reduce creates new contingecy summing out the variable with stride of one
func (sp *Sparse) Reduce() *Sparse {
	mt := Sparse{}
	x := sp.varOrdering[0]
	numEntry := 1
	for k, v := range sp.cardinality {
		if k != x {
			numEntry *= v
		}
	}
	stride := 1
	mt.strideMap = make(map[int]int)
	for _, v := range sp.varOrdering {
		if v != x {
			mt.strideMap[v] = stride
			stride *= sp.cardinality[v]
		}
	}
	mt.countMap = make(map[int]int)
	c := sp.cardinality[x]
	for i := 0; i < numEntry; i++ {
		aux := 0
		for j := 0; j < c; j++ {
			aux += sp.countMap[(i*c)+j]
		}
		if aux != 0 {
			mt.countMap[i] = aux
		}
	}
	return &mt
}
