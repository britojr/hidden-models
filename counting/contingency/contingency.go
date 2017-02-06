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

// SetValues ..
func (t *Table) SetValues(values []int) {
	t.occurrences = values
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
