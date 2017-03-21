package em

import (
	"fmt"
	"testing"

	"github.com/britojr/kbn/cliquetree"
	"github.com/britojr/kbn/factor"
	"github.com/britojr/kbn/utils"
)

func TestCheckFactorDiff(t *testing.T) {
	cases := []struct {
		alist     [][]float64
		blist     [][]float64
		threshold float64
		diff      float64
		err       error
	}{{
		alist:     [][]float64{},
		blist:     [][]float64{},
		threshold: 0,
		diff:      0,
	}, {
		alist: [][]float64{
			{11, 21},
			{10, 20},
		},
		blist: [][]float64{
			{11, 21},
			{10, 20},
		},
		threshold: 10,
		diff:      0,
	}, {
		alist: [][]float64{
			{11, 21},
		},
		blist: [][]float64{
			{11, 21},
			{10, 20},
		},
		err: fmt.Errorf("missing potential"),
	}, {
		alist: [][]float64{
			{11, 21},
			{11, 21},
		},
		blist: [][]float64{
			{11, 21},
			{10},
		},
		err: fmt.Errorf("incompatible slices"),
	}, {
		alist: [][]float64{
			{11, 21},
			{10, 20},
		},
		blist: [][]float64{
			{10, 21},
			{10, 22},
		},
		threshold: 10,
		diff:      2,
	}, {
		alist: [][]float64{
			{11, 21},
			{10, 20},
		},
		blist: [][]float64{
			{10, 21},
			{10, 22},
		},
		threshold: 1,
		diff:      1,
	}, {
		alist: [][]float64{
			{11, 21},
			{10, 20},
		},
		blist: [][]float64{
			{10, 21},
			{4, 31},
		},
		threshold: 5,
		diff:      6,
	}}
	for _, tt := range cases {
		fs := make([]*factor.Factor, len(tt.alist))
		for i := range tt.alist {
			fs[i] = factor.NewFactor([]int{}, []int{})
			fs[i].SetValues(tt.alist[i])
		}
		gs := make([]*factor.Factor, len(tt.blist))
		for i := range tt.blist {
			gs[i] = factor.NewFactor([]int{}, []int{})
			gs[i].SetValues(tt.blist[i])
		}

		diff, err := checkFactorDiff(fs, gs, tt.threshold)
		if (tt.err != nil && err == nil) || (tt.err == nil && err != nil) {
			t.Errorf("different err,  want %v, got %v", tt.err, err)
		}
		if tt.err == nil {
			if !utils.FuzzyEqual(tt.diff, diff, 1e-4) {
				t.Errorf("wrong max diff, want %v, got %v", tt.diff, diff)
			}
		}
	}
}

type FakeDataHandler struct {
	data [][]int
}

func (f FakeDataHandler) Data() [][]int {
	return f.data
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

func TestExpectationStep(t *testing.T) {
	cases := []struct {
		cliques, adj [][]int
		cardin       []int
		values       [][]float64
		ds           FakeDataHandler
		result       [][]float64
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
		ds: FakeDataHandler{
			data: [][]int{
				{0, 1, 1, 0, 1},
				{0, 1, 1, 0, 1},
				{1, 1, 1, 1, 1},
				{0, 1, 1, 0, 1},
			},
		},
		result: [][]float64{
			{3, 1},
			{0, 4},
			{0, 0, 0, 0, 0, 0, 3, 1},
			{0, 3, 0, 1},
			{0, 0, 0, 4},
		},
	}}
	for _, tt := range cases {
		c, err := initiCliqueTree(tt.cliques, tt.adj, tt.cardin, tt.values)
		if err != nil {
			t.Errorf(err.Error())
		}
		got := expectationStep(c, tt.ds)
		if len(got) != len(tt.result) {
			t.Errorf("wrong number of factors, want %v, got %v", len(got), len(tt.result))
		}
		for i := range tt.result {
			for j := range tt.result[i] {
				if !utils.FuzzyEqual(tt.result[i][j], got[i].Values()[j]) {
					t.Errorf("wrong counting, want %v, got %v", tt.result[i], got[i].Values())
				}
			}
		}
	}
}
