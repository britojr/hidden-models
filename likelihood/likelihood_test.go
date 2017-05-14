package likelihood

import (
	"fmt"
	"testing"

	"github.com/britojr/kbn/assignment"
	"github.com/britojr/kbn/cliquetree"
	"github.com/britojr/kbn/factor"
	"github.com/britojr/kbn/floats"
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

func initiCliqueTree(cliques, adj [][]int, cardin []int, values [][]float64) (*cliquetree.CliqueTree, error) {
	c, err := cliquetree.NewStructure(cliques, adj)
	if err != nil {
		return nil, err
	}
	potentials := make([]*factor.Factor, len(values))
	for i, v := range values {
		potentials[i] = factor.NewFactorValues(cliques[i], cardin, v)
	}
	err = c.SetAllPotentials(potentials)
	if err != nil {
		return nil, err
	}
	return c, nil
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
		got := StructLoglikelihood(tt.cliques, tt.sepsets, tt.counter)
		if !floats.AlmostEqual(tt.result, got, 1e-4) {
			t.Errorf("want %v, got %v", tt.result, got)
		}
	}
}

func TestLoglikelihood(t *testing.T) {
	cases := []struct {
		cliques, adj [][]int
		cardin       []int
		values       [][]float64
		result       []struct {
			data [][]int
			ll   float64
		}
	}{{
		cliques: [][]int{{0}, {1}, {0, 1, 2}, {2, 3}, {2, 4}},
		adj:     [][]int{{2}, {2}, {0, 1, 3, 4}, {2}, {2}},
		cardin:  []int{2, 2, 2, 2, 2},
		values: [][]float64{
			{.999, .001},
			{.998, .002},
			{.999, .06, .71, .05, .001, .94, .29, .95},
			{.95, .10, .05, .90},
			{.99, .30, .01, .70},
		},
		result: []struct {
			data [][]int
			ll   float64
		}{{
			data: [][]int{
				{0, 1, 1, 0, 1},
				{0, 1, 1, 0, 1},
				{1, 1, 1, 1, 1},
				{0, 1, 1, 0, 1},
			},
			ll: -43.97392118,
		}},
	}}
	for _, tt := range cases {
		c, err := initiCliqueTree(tt.cliques, tt.adj, tt.cardin, tt.values)
		if err != nil {
			t.Errorf(err.Error())
		}
		c.UpDownCalibration()
		for i := range tt.result {
			got := Loglikelihood(c, tt.result[i].data)
			if !floats.AlmostEqual(tt.result[i].ll, got, 1e-7) {
				t.Errorf("wrong ll2, want %v, got %v", tt.result[i].ll, got)
			}
		}
	}
}
