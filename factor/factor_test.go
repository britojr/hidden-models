package factor

import (
	"errors"
	"math"
	"reflect"
	"testing"

	"github.com/britojr/kbn/assignment"
	"github.com/britojr/kbn/utils"
)

type testStruct struct {
	cardin  []int
	varlist []int
	values  []float64
}

var f01 = testStruct{
	cardin:  []int{3, 2, 2},
	varlist: []int{0, 1},
	values:  []float64{0.5, 0.1, 0.3, 0.8, 0.0, 0.9},
}
var f12 = testStruct{
	cardin:  []int{3, 2, 2},
	varlist: []int{1, 2},
	values:  []float64{0.5, 0.1, 0.7, 0.2},
}
var f012 = testStruct{
	cardin:  []int{3, 2, 2},
	varlist: []int{0, 1, 2},
	values:  []float64{0.25, 0.05, 0.15, 0.08, 0.00, 0.09, 0.35, 0.07, 0.21, 0.16, 0.00, 0.18},
}

var f012_0 = testStruct{
	cardin:  []int{3, 2, 2},
	varlist: []int{1, 2},
	values:  []float64{0.45, 0.17, 0.63, 0.34},
}
var f012_1 = testStruct{
	cardin:  []int{3, 2, 2},
	varlist: []int{0, 2},
	values:  []float64{0.33, 0.05, 0.24, 0.51, 0.07, 0.39},
}
var f012_2 = testStruct{
	cardin:  []int{3, 2, 2},
	varlist: []int{0, 1},
	values:  []float64{0.60, 0.12, 0.36, 0.24, 0.0, 0.27},
}

var tests = []testStruct{
	f01,
	f12,
	f012,
}

var sumOutOneTests = []testStruct{
	f012_0,
	f012_1,
	f012_2,
}

func TestNew(t *testing.T) {
	for _, w := range tests {
		got := New(w.varlist, w.cardin, w.values)
		assig := assignment.New(w.varlist, w.cardin)
		i := 0
		for {
			v := got.Get(assig)
			if v != w.values[i] {
				t.Errorf("want(%v); got(%v)", w.values[i], v)
			}
			if hasnext := assig.Next(); !hasnext {
				break
			}
			i++
		}
	}
}

func TestVariables(t *testing.T) {
	for _, w := range tests {
		f := New(w.varlist, w.cardin, w.values)
		got := f.Variables()
		if !reflect.DeepEqual(w.varlist, got) {
			t.Errorf("want(%v); got(%v)", w.varlist, got)
		}
	}
}

func TestProduct(t *testing.T) {
	a := New(f01.varlist, f01.cardin, f01.values)
	b := New(f12.varlist, f12.cardin, f12.values)
	got := a.Product(b)
	assig := assignment.New(f012.varlist, f012.cardin)
	i := 0
	for {
		v := got.Get(assig)
		if !utils.FuzzyEqual(v, f012.values[i]) {
			t.Errorf("want(%v); got(%v)", f012.values[i], v)
		}
		if hasnext := assig.Next(); !hasnext {
			break
		}
		i++
	}
}

func TestSumOutOne(t *testing.T) {
	a := New(f012.varlist, f012.cardin, f012.values)
	for i, w := range sumOutOneTests {
		got := a.SumOutOne(i)
		assig := assignment.New(w.varlist, w.cardin)
		i := 0
		for {
			v := got.Get(assig)
			if !utils.FuzzyEqual(v, w.values[i]) {
				t.Errorf("want(%v); got(%v)", w.values[i], v)
			}
			if hasnext := assig.Next(); !hasnext {
				break
			}
			i++
		}
	}
}

var testRestrict = []struct {
	varlist    []int
	cardin     []int
	values     []float64
	evid       []int
	restricted []float64
}{
	{
		varlist:    []int{0, 1},
		cardin:     []int{2, 2},
		values:     []float64{0.1, 0.2, 0.3, 0.4},
		evid:       []int{1, 0},
		restricted: []float64{0.0, 0.2, 0.0, 0.0},
	},
	{
		varlist:    []int{1, 3},
		cardin:     []int{2, 2, 2, 2},
		values:     []float64{0.1, 0.2, 0.3, 0.4},
		evid:       []int{0, 1, 0},
		restricted: []float64{0.0, 0.2, 0.0, 0.4},
	},
}

func TestRestrict(t *testing.T) {
	for _, v := range testRestrict {
		f := New(v.varlist, v.cardin, v.values)
		f = f.Restrict(v.evid)
		if !reflect.DeepEqual(v.restricted, f.values) {
			t.Errorf("want %v, got %v", v.restricted, f.values)
		}
	}
}

var testNormalize = []struct {
	varlist    []int
	cardin     []int
	values     []float64
	normalized []float64
}{
	{
		varlist:    []int{0, 1},
		cardin:     []int{2, 2},
		values:     []float64{10, 20, 30, 40},
		normalized: []float64{0.1, 0.2, 0.3, 0.4},
	},
	{
		varlist:    []int{1, 3, 5},
		cardin:     []int{2, 2, 2, 2, 2, 2},
		values:     []float64{10, 20, 30, 40, 50, 60, 70, 80},
		normalized: []float64{1.0 / 36, 2.0 / 36, 3.0 / 36, 4.0 / 36, 5.0 / 36, 6.0 / 36, 7.0 / 36, 8.0 / 36},
	},
	{
		varlist:    []int{0, 1},
		cardin:     []int{2, 2},
		values:     []float64{0.15, 0.25, 0.35, 0.25},
		normalized: []float64{0.15, 0.25, 0.35, 0.25},
	},
	{
		varlist:    []int{1},
		cardin:     []int{2, 2},
		values:     []float64{0.15},
		normalized: []float64{1},
	},
	{
		varlist:    []int{1},
		cardin:     []int{2, 2},
		values:     []float64{},
		normalized: []float64{},
	},
	{
		varlist:    []int{1},
		cardin:     []int{2, 2},
		values:     []float64{0, 0.0},
		normalized: []float64{0, 0},
	},
}

func TestNormalize(t *testing.T) {
	for _, v := range testNormalize {
		f := New(v.varlist, v.cardin, v.values)
		f.Normalize()
		if !reflect.DeepEqual(v.normalized, f.values) {
			t.Errorf("want %v, got %v", v.normalized, f.values)
		}
	}
}

var testMaxDifference = []struct {
	varlist []int
	cardin  []int
	alist   [][]float64
	blist   [][]float64
	maxdiff float64
	err     error
}{
	{
		varlist: []int{1, 3, 5},
		cardin:  []int{2, 2, 2, 2, 2, 2},
		alist: [][]float64{
			{11, 21, 31, 41, 51, 61, 71, 81},
			{10, 20, 30, 40, 50, 60, 70, 80},
		},
		blist: [][]float64{
			{11, 21, 31, 41, 51, 61, 71, 81},
			{10, 20, 30, 40, 50, 60, 70, 80},
		},
		maxdiff: 0,
	},
	{
		varlist: []int{1, 3, 5},
		cardin:  []int{2, 2, 2, 2, 2, 2},
		alist: [][]float64{
			{11, 21, 31, 41, 51, 61, 71, 81},
			{10, 20, 30, 40, 50, 60, 72, 80},
		},
		blist: [][]float64{
			{10, 20, 30, 40, 50, 60, 70, 80},
			{10, 20, 30, 40, 50, 60, 70, 80},
		},
		maxdiff: 2,
	},
	{
		varlist: []int{1, 3, 5},
		cardin:  []int{2, 2, 2, 2, 2, 2},
		alist: [][]float64{
			{11, 21, 31, 0.8, 0.005},
			{10, 20, 30, 40, 50, 60, 70, 80},
			{11, 21, 31, 0.8, 0.005},
		},
		blist: [][]float64{
			{11, 21, 31, 0.8, 0.005},
			{10, 20, 30, 40, 50, 60, 70, 80},
			{11, 21, 31, 0.8, 0.004},
		},
		maxdiff: 0.001,
	},
	{
		varlist: []int{1, 3, 5},
		cardin:  []int{2, 2, 2, 2, 2, 2},
		alist: [][]float64{
			{11, 21, 31, 0.8, 0.005},
			{10, 20, 30, 40, 50, 60, 70, 80},
			{},
			{11, 21, 31, 0.8, 0.005},
		},
		blist: [][]float64{
			{11, 21, 31, 0.8, 0.005},
			{10, 20, 30, 40, 50, 60, 70, 80},
			{},
			{11, 21, 31, 0.8, 0.004},
		},
		maxdiff: 0.001,
	},
	{
		varlist: []int{1, 3, 5},
		cardin:  []int{2, 2, 2, 2, 2, 2},
		alist: [][]float64{
			{11, 21, 31, 0.8, 0.005},
			{10, 20, 30, 40, 50, 60, 70, 80},
			{},
			{11, 21, 31, 0.8, 0.005},
		},
		blist: [][]float64{
			{11, 21, 31, 0.8, 0.005},
			{10, 20, 30, 40, 50, 60, 70, 80},
			{11, 21, 31, 0.8, 0.004},
			{11, 21, 31, 0.8, 0.004},
		},
		err: errors.New("incompatible list of factors"),
	},
}

func TestMaxDifference(t *testing.T) {
	for _, v := range testMaxDifference {
		f := make([]*Factor, len(v.alist))
		g := make([]*Factor, len(v.alist))
		for i := range f {
			if len(v.alist[i]) > 0 {
				f[i] = New(v.varlist, v.cardin, v.alist[i])
			}
			if len(v.blist[i]) > 0 {
				g[i] = New(v.varlist, v.cardin, v.blist[i])
			}

		}
		got, i, j, err := MaxDifference(f, g)
		if (v.err != nil && err == nil) || (v.err == nil && err != nil) {
			t.Errorf("want %v, got %v", v.err, err)
		}
		if err == nil {
			if v.maxdiff != got {
				t.Errorf("want %v, got %v", v.maxdiff, got)
			}
			if math.Abs(f[i].Values()[j]-g[i].Values()[j]) != got {
				t.Errorf("want %v, got %v", math.Abs(f[i].Values()[j]-g[i].Values()[j]), got)
			}
		}
	}
}
