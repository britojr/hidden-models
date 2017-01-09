// Package contingecytable implements Contingecy Table
package contingecytable

// Sparse contingey table
type Sparse struct {
	strideMap map[int]int
	countMap  map[int]int
}

// LoadFromData creates a new contingey table from data
func LoadFromData(data [][]int, cardinality []int) *Sparse {
	sp := Sparse{make(map[int]int), make(map[int]int)}
	stride := 1
	for k, v := range cardinality {
		sp.strideMap[k] = stride
		stride *= v
	}
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
	mt := Sparse{make(map[int]int), make(map[int]int)}
	return &mt
}

// Reduce creates new contingecy summing out the variable with stride of one
func (sp *Sparse) Reduce() *Sparse {
	mt := Sparse{make(map[int]int), make(map[int]int)}
	return &mt
}
