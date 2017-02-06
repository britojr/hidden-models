package contingency

// Table is a static implementation of contingency table
type Table struct {
	variables   []int
	occurrences []int
	cardinality *[]int
}

// NewTable ..
func NewTable(variables []int, occurrences []int, cardin *[]int) *Table {
	t := new(Table)
	t.variables = variables
	t.occurrences = occurrences
	t.cardinality = cardin
	return t
}

// SetVariables ..
func (t *Table) SetVariables(vars ...int) {
	t.variables = append([]int(nil), vars...)
}

// SetCardinality ..
func (t *Table) SetCardinality(cardin *[]int) {
	t.cardinality = cardin
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
