package assignment

// Assignment is used to create and iterate through variable assignments
type Assignment []struct {
	variable int
	value    int
	card     int
}

// New assignment creates a new assignment with the specified order of the variables
func New(varlist []int, cardinality []int) Assignment {
	if len(varlist) == 0 {
		return nil
	}
	a := make(Assignment, len(varlist))
	for i, v := range varlist {
		a[i].variable = v
		if v < len(cardinality) {
			a[i].card = cardinality[v]
		}
	}
	return a
}

// Next creates the next valuation for the assignment
// or returns nil after the last possible value
func (a *Assignment) Next() {
	i := 0
	(*a)[i].value++
	for (*a)[i].value >= (*a)[i].card {
		(*a)[i].value = 0
		i++
		if i >= len(*a) {
			*a = nil
			return
		}
		(*a)[i].value++
	}
}

// Var returns the id of variable at ith position
func (a Assignment) Var(i int) int {
	return a[i].variable
}

// Value returns the value assigned to the variable at ith position
func (a Assignment) Value(i int) int {
	return a[i].value
}

// Consistent returns true if the current assignment is consistent with a given valoration
func (a Assignment) Consistent(values []int) bool {
	for _, v := range a {
		if v.variable >= len(values) || values[v.variable] == -1 {
			continue
		}
		if values[v.variable] != v.value {
			return false
		}
	}
	return true
}
