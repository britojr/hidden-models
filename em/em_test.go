package em

import (
	"fmt"
	"testing"

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
