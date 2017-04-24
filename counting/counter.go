package counting

import "github.com/britojr/kbn/assignment"

// Counter returns a counting value for an assignment
type Counter interface {
	Count(assig *assignment.Assignment) (count int, ok bool)
	CountAssignments(varlist []int) []int
	Cardinality() []int
	NumTuples() int
}
