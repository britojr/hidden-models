package likelihood

import (
	"fmt"
	"testing"

	"github.com/britojr/kbn/assignment"
	"github.com/britojr/kbn/utils"
)

type FakeCounter struct {
	cardin    []int
	numtuples int
	counts    map[string][]int
}

func (f FakeCounter) Count(assig *assignment.Assignment) (count int, ok bool) {
	panic("not implemented")
}
func (f FakeCounter) CountAssignments(varlist []int) []int {
	return f.counts[fmt.Sprint(varlist)]
}
func (f FakeCounter) Cardinality() []int {
	return f.cardin
}
func (f FakeCounter) NumTuples() int {
	return f.numtuples
}

func TestStructloglikelihood(t *testing.T) {
	cases := []struct {
		cliques [][]int
		sepsets [][]int
		counter FakeCounter
		result  float64
	}{
		{
			cliques: [][]int{{0, 1}, {1, 2}},
			sepsets: [][]int{{}, {1}},
			counter: FakeCounter{
				cardin:    []int{2, 2, 2},
				numtuples: 100,
				counts: map[string][]int{
					fmt.Sprint([]int{0, 1}): {25, 10, 35, 30},
					fmt.Sprint([]int{1, 2}): {40, 20, 10, 30},
					fmt.Sprint([]int{1}):    {60, 40},
				},
			},
			result: -191.2304,
		},
	}
	for _, tt := range cases {
		got := StructLog(tt.cliques, tt.sepsets, tt.counter)
		if !utils.FuzzyEqual(tt.result, got, 1e-4) {
			t.Errorf("want %v, got %v", tt.result, got)
		}
	}
}
