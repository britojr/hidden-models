package assignment

// Assignment is used to create and iterate through variable assignments
type Assignment struct {
	varlist []int // list of variables
	values  []int // value[i] = value of the ith variable
	card    []int // card[i] = cardinality of ith variable
}

// New assignment creates a new assignment with the specified order of the variables
func New(varlist []int, cardinality []int) *Assignment {
	a := &Assignment{
		varlist: varlist,
		values:  make([]int, len(varlist)),
		card:    make([]int, len(varlist)),
	}
	for i, v := range varlist {
		if v < len(cardinality) {
			a.card[i] = cardinality[v]
		}
	}
	// starts before the first value
	if len(a.values) > 0 {
		a.values[0] = -1
	}
	return a
}

// Next creates the next valuation for the assignment
// or returns false after the last possible value
func (a *Assignment) Next() bool {
	if len(a.values) == 0 {
		return false
	}
	i := 0
	a.values[i]++
	for a.values[i] >= a.card[i] {
		a.values[i] = 0
		i++
		if i >= len(a.values) {
			return false
		}
		a.values[i]++
	}
	return true
}

// Var returns the id of variable at ith position
func (a *Assignment) Var(i int) int {
	return a.varlist[i]
}

// Variables returns the slice of variables
func (a *Assignment) Variables() []int {
	return a.varlist
}

// Value returns the value assigned to the variable on the ith position
func (a *Assignment) Value(i int) int {
	return a.values[i]
}

// Values returns the slice of values
func (a *Assignment) Values() []int {
	return a.values
}

// Index calculates the corresponding index given a stride
func (a *Assignment) Index(stride map[int]int) int {
	x := 0
	for i, v := range a.values {
		x += v * stride[a.varlist[i]]
	}
	return x
}

// Consistent returns true if the current assignment is consistent with a given valoration
func (a *Assignment) Consistent(values []int) bool {
	for i, v := range a.varlist {
		if v >= len(values) || values[v] == -1 {
			continue
		}
		if values[v] != a.values[i] {
			return false
		}
	}
	return true
}
