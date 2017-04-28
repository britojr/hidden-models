package learn

import (
	"fmt"
	"reflect"
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

func TestCreateRandomPortentials(t *testing.T) {
	cases := []struct {
		cliques [][]int
		cardin  []int
	}{
		{
			cliques: [][]int{{0, 1}, {1, 2}},
			cardin:  []int{2, 2, 2},
		},
	}
	for _, tt := range cases {
		faclist := CreateRandomPotentials(tt.cliques, tt.cardin)
		for _, f := range faclist {
			tot := utils.SliceSumFloat64(f.Values())
			if !utils.FuzzyEqual(tot, 1) {
				t.Errorf("random factor not normalized, sums to: %v", tot)
			}
			for _, v := range f.Values() {
				if v == 0 {
					t.Errorf("random factor has zero values: %v", f.Values())
				}
			}
		}
	}
}

func TestCreateUniformPortentials(t *testing.T) {
	fakeCounter := FakeCounter{
		cardin:    []int{2, 2, 2},
		numtuples: 100,
		counts: map[string][]int{
			fmt.Sprint([]int{0, 1}): {25, 10, 35, 30},
			fmt.Sprint([]int{1, 2}): {40, 20, 10, 30},
			fmt.Sprint([]int{1}):    {60, 40},
			fmt.Sprint([]int{0}):    {35, 65},
		},
	}
	cases := []struct {
		cliques [][]int
		cardin  []int
		numobs  int
		counter FakeCounter
		result  [][]float64
	}{
		{
			cliques: [][]int{{0, 1}, {1, 2}},
			cardin:  []int{2, 2, 2},
			numobs:  2,
			counter: fakeCounter,
			result:  [][]float64{{.25, .10, .35, .30}, {.60 / 2.0, .40 / 2.0, .60 / 2.0, .40 / 2.0}},
		},
		{
			cliques: [][]int{{0, 1}, {1, 2}},
			cardin:  []int{2, 2, 2},
			numobs:  3,
			counter: fakeCounter,
			result:  [][]float64{{.25, .10, .35, .30}, {.40, .20, .10, .30}},
		},
		{
			cliques: [][]int{{0, 1}, {1, 2}},
			cardin:  []int{2, 2, 2},
			numobs:  1,
			counter: fakeCounter,
			result:  [][]float64{{.35 / 2.0, .65 / 2.0, .35 / 2.0, .65 / 2.0}, {.25, .25, .25, .25}},
		},
	}
	for _, tt := range cases {
		faclist := CreateEmpiricPotentials(tt.counter, tt.cliques, tt.cardin, tt.numobs, EmpiricUniform)
		if len(faclist) != len(tt.result) {
			t.Errorf("wrong number of factors, expected %v, got %v", len(tt.result), len(faclist))
		}
		for i, f := range faclist {
			tot := utils.SliceSumFloat64(f.Values())
			if !utils.FuzzyEqual(tot, 1) {
				t.Errorf("uniform factor not normalized, sums to: %v", tot)
			}
			for _, v := range f.Values() {
				if v == 0 {
					t.Errorf("uniform factor has zero values: %v", f.Values())
				}
			}
			if !reflect.DeepEqual(tt.result[i], f.Values()) {
				t.Errorf("Wrong values, want %v, got %v", tt.result[i], f.Values())
			}
		}
	}
}
