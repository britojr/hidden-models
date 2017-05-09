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
		got := NewFactor(w.varlist, w.cardin)
		got.SetValues(w.values)
		assig := assignment.New(w.varlist, w.cardin)
		i := 0
		for assig.Next() {
			v := got.Get(assig)
			if v != w.values[i] {
				t.Errorf("want(%v); got(%v)", w.values[i], v)
			}
			i++
		}
		if !reflect.DeepEqual(w.values, got.Values()) {
			t.Errorf("Wrong values, want %v, got %v", w.values, got.Values())
		}
		if !reflect.DeepEqual(got.Cardinality(), w.cardin) {
			t.Errorf("Wrong cardinality, want %v, got %v", w.cardin, got.Cardinality())
		}
	}
}

var testNewFactor = []struct {
	varlist []int
	cardin  []int
	values  []float64
}{
	{
		varlist: []int{0, 1},
		cardin:  []int{2, 2},
		values:  []float64{0, 0, 0, 0},
	},
	{
		varlist: []int{1, 3, 5},
		cardin:  []int{2, 2, 2, 2, 2, 2},
		values:  []float64{0, 0, 0, 0, 0, 0, 0, 0},
	},
	{
		varlist: []int{1},
		cardin:  []int{2, 2},
		values:  []float64{0, 0},
	},
	{
		varlist: []int{1},
		cardin:  []int{2, 1},
		values:  []float64{0},
	},
	{
		varlist: []int{},
		cardin:  []int{},
		values:  []float64{0},
	},
}

func TestNewFactor(t *testing.T) {
	for _, v := range testNewFactor {
		got := NewFactor(v.varlist, v.cardin)
		if !reflect.DeepEqual(v.varlist, got.Variables()) {
			t.Errorf("want %v, got %v", v.varlist, got.Variables())
		}
		if !reflect.DeepEqual(v.values, got.Values()) {
			t.Errorf("want %v, got %v", v.values, got.Values())
		}
	}
}

func TestVariables(t *testing.T) {
	for _, w := range tests {
		f := NewFactorValues(w.varlist, w.cardin, w.values)
		got := f.Variables()
		if !reflect.DeepEqual(w.varlist, got) {
			t.Errorf("want(%v); got(%v)", w.varlist, got)
		}
	}
}

func TestProduct(t *testing.T) {
	a := NewFactorValues(f01.varlist, f01.cardin, f01.values)
	b := NewFactorValues(f12.varlist, f12.cardin, f12.values)
	got := a.Product(b)
	assig := assignment.New(f012.varlist, f012.cardin)
	i := 0
	for assig.Next() {
		v := got.Get(assig)
		if !utils.FuzzyEqual(v, f012.values[i]) {
			t.Errorf("want(%v); got(%v)", f012.values[i], v)
		}
		i++
	}
}
func TestProduct2(t *testing.T) {
	cases := []struct {
		cardin []int
		varA   []int
		valA   []float64
		varB   []int
		valB   []float64
		varRes []int
		valRes []float64
	}{{
		cardin: []int{2, 2, 2, 2},
		varA:   []int{1, 0, 2},
		valA:   []float64{.2, .3, .4, .5, .6, .7, .8, .9},
		varB:   []int{2, 1},
		valB:   []float64{.11, .12, .13, .14},
		varRes: []int{0, 1, 2},
		valRes: []float64{.2 * .11, .4 * .11, .3 * .13, .5 * .13, .6 * .12, .8 * .12, .7 * .14, .9 * .14},
	}, {
		cardin: []int{2, 2, 2, 2},
		varA:   []int{1, 0, 2},
		valA:   []float64{.2, .3, .4, .5, .6, .7, .8, .9},
		varB:   []int{2, 1, 3},
		valB:   []float64{.11, .12, .13, .14, .15, .16, .17, .18},
		varRes: []int{0, 1, 2, 3},
		valRes: []float64{
			.2 * .11, .4 * .11, .3 * .13, .5 * .13, .6 * .12, .8 * .12, .7 * .14, .9 * .14,
			.2 * .15, .4 * .15, .3 * .17, .5 * .17, .6 * .16, .8 * .16, .7 * .18, .9 * .18,
		},
	}}
	for _, tt := range cases {
		a := NewFactorValues(tt.varA, tt.cardin, tt.valA)
		b := NewFactorValues(tt.varB, tt.cardin, tt.valB)
		got := a.Product(b)
		if !reflect.DeepEqual(tt.varRes, got.Variables()) {
			t.Errorf("Wrong variables, want %v, got %v", tt.varRes, got.Variables())
		}
		for i, v := range tt.valRes {
			if !utils.FuzzyEqual(v, got.Values()[i]) {
				t.Errorf("Wrong values, want %v, got %v", v, got.Values()[i])
			}
		}
	}
}

func TestSumOutOne(t *testing.T) {
	a := NewFactorValues(f012.varlist, f012.cardin, f012.values)
	for i, w := range sumOutOneTests {
		got := a.SumOutOne(i)
		assig := assignment.New(w.varlist, w.cardin)
		i := 0
		for assig.Next() {
			v := got.Get(assig)
			if !utils.FuzzyEqual(v, w.values[i]) {
				t.Errorf("want(%v); got(%v)", w.values[i], v)
			}
			i++
		}
	}
}

var testReduce = []struct {
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
	{
		varlist:    []int{3, 1},
		cardin:     []int{2, 2, 2, 2},
		values:     []float64{0.1, 0.3, 0.2, 0.4},
		evid:       []int{0, 1, 0},
		restricted: []float64{0.0, 0.0, 0.2, 0.4},
	},
}

func TestReduce(t *testing.T) {
	for _, v := range testReduce {
		f := NewFactorValues(v.varlist, v.cardin, v.values)
		g := f.Reduce(v.evid)
		if !reflect.DeepEqual(v.restricted, g.values) {
			t.Errorf("want %v, got %v", v.restricted, g.values)
		}
		if !reflect.DeepEqual(v.values, f.values) {
			t.Errorf("Original values changed, want %v, got %v", v.values, f.values)
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
	// {
	// 	varlist:    []int{1},
	// 	cardin:     []int{2, 2},
	// 	values:     []float64{},
	// 	normalized: []float64{},
	// },
	// {
	// 	varlist:    []int{1},
	// 	cardin:     []int{2, 2},
	// 	values:     []float64{0, 0.0},
	// 	normalized: []float64{0, 0},
	// },
}

func TestNormalize(t *testing.T) {
	for _, v := range testNormalize {
		f := NewFactorValues(v.varlist, v.cardin, v.values)
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
				f[i] = NewFactorValues(v.varlist, v.cardin, v.alist[i])
			}
			if len(v.blist[i]) > 0 {
				g[i] = NewFactorValues(v.varlist, v.cardin, v.blist[i])
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

var testSetUniform = []struct {
	f       *Factor
	uniform []float64
}{
	{
		f:       NewFactorValues([]int{0, 1}, []int{2, 2}, []float64{10, 20, 30, 40}),
		uniform: []float64{0.25, 0.25, 0.25, 0.25},
	},
	{
		f:       NewFactor([]int{1, 2, 3}, []int{2, 2, 2, 2, 2}),
		uniform: []float64{0.125, 0.125, 0.125, 0.125, 0.125, 0.125, 0.125, 0.125},
	},
	{
		f:       NewFactor([]int{1}, []int{2, 2}),
		uniform: []float64{0.5, 0.5},
	},
	{
		f:       NewFactorValues([]int{1}, []int{2, 2}, []float64{}),
		uniform: []float64{},
	},
	{
		f:       NewFactor([]int{}, []int{}),
		uniform: []float64{1},
	},
}

func TestSetUniform(t *testing.T) {
	for _, v := range testSetUniform {
		got := v.f.SetUniform()
		if !reflect.DeepEqual(v.uniform, got.Values()) {
			t.Errorf("want %v, got %v", v.uniform, got.Values())
		}
	}
}

var testSetRandom = []struct {
	f    *Factor
	size int
}{
	{
		f:    NewFactorValues([]int{0, 1}, []int{2, 2}, []float64{10, 20, 30, 40}),
		size: 4,
	},
	{
		f:    NewFactor([]int{1, 2, 3}, []int{2, 2, 2, 2, 2}),
		size: 8,
	},
	{
		f:    NewFactor([]int{1}, []int{2, 2}),
		size: 2,
	},
	// {
	// 	f:    NewFactorValues([]int{1}, []int{2, 2}, []float64{}),
	// 	size: 0,
	// },
	{
		f:    NewFactor([]int{}, []int{}),
		size: 1,
	},
}

func TestSetRandom(t *testing.T) {
	for _, v := range testSetRandom {
		got := v.f.SetRandom()
		if v.size != len(got.Values()) {
			t.Errorf("want %v, got %v", v.size, got.Values())
		}
		if v.size != 0 && !utils.FuzzyEqual(1, utils.SliceSumFloat64(got.Values())) {
			t.Errorf("not normalized, sum %v", utils.SliceSumFloat64(got.Values()))
		}
	}

	// test different outcomes
	f := NewFactor([]int{1, 2, 3}, []int{2, 2, 2, 2, 2})
	f.SetRandom()
	values := append([]float64(nil), f.values...)
	f.SetRandom()
	count := 0
	for i := range values {
		if utils.FuzzyEqual(values[i], f.values[i]) {
			count++
		}
	}
	if count == len(values) {
		t.Errorf("Sampled the same distribution:\n%v\n%v", f.values, f.SetRandom().values)
	}
}

func TestSetDirichlet(t *testing.T) {
	cases := []struct {
		f      *Factor
		size   int
		alphas []float64
	}{{
		f:      NewFactorValues([]int{0, 1}, []int{2, 2}, []float64{10, 20, 30, 40}),
		size:   4,
		alphas: []float64{3.2, 3.2, 3.2, 3.2},
	}, {
		f:      NewFactor([]int{1, 2, 3}, []int{2, 2, 2, 2, 2}),
		size:   8,
		alphas: []float64{0.1, 0.1, 0.1, 0.1, 0.1, 0.1, 0.1, 0.1},
	}, {
		f:      NewFactor([]int{1}, []int{2, 2}),
		size:   2,
		alphas: []float64{0.01, 0.01},
	}, {
		f:      NewFactor([]int{}, []int{}),
		size:   1,
		alphas: []float64{5},
	}}
	for _, tt := range cases {
		got := tt.f.SetDirichlet(tt.alphas)
		if tt.size != len(got.Values()) {
			t.Errorf("want %v, got %v", tt.size, got.Values())
		}
		if tt.size != 0 && !utils.FuzzyEqual(1, utils.SliceSumFloat64(got.Values())) {
			t.Errorf("not normalized, sum %v", utils.SliceSumFloat64(got.Values()))
		}
	}

	// test different outcomes
	f := NewFactor([]int{1, 2, 3}, []int{2, 2, 2, 2, 2})
	alphas := []float64{0.7, 0.7, 0.7, 0.7, 0.7, 0.7, 0.7, 0.7}
	f.SetDirichlet(alphas)
	values := append([]float64(nil), f.values...)
	f.SetDirichlet(alphas)
	count := 0
	for i := range values {
		if utils.FuzzyEqual(values[i], f.values[i]) {
			count++
		}
	}
	if count == len(values) {
		t.Errorf("Sampled the same distribution:\n%v\n%v", f.values, f.SetDirichlet(alphas).values)
	}
}

var testSumOut = []struct {
	f       *Factor
	varlist []int
	r       *Factor
}{
	{
		f: NewFactorValues([]int{0, 1}, []int{2, 2}, []float64{10, 20, 30, 40}),
		r: NewFactorValues([]int{0, 1}, []int{2, 2}, []float64{10, 20, 30, 40}),
	},
	{
		f:       NewFactorValues([]int{0, 1}, []int{2, 2}, []float64{10, 20, 30, 40}),
		varlist: []int{2},
		r:       NewFactorValues([]int{0, 1}, []int{2, 2}, []float64{10, 20, 30, 40}),
	},
	{
		f:       NewFactorValues([]int{0, 1}, []int{2, 2}, []float64{10, 20, 30, 40}),
		varlist: []int{0},
		r:       NewFactorValues([]int{1}, []int{2, 2}, []float64{30, 70}),
	},
	{
		f:       NewFactorValues([]int{1, 0}, []int{2, 2}, []float64{10, 30, 20, 40}),
		varlist: []int{0},
		r:       NewFactorValues([]int{1}, []int{2, 2}, []float64{30, 70}),
	},
	{
		f:       NewFactorValues([]int{0, 1}, []int{2, 2}, []float64{10, 20, 30, 40}),
		varlist: []int{1},
		r:       NewFactorValues([]int{0}, []int{2, 2}, []float64{40, 60}),
	},
	{
		f:       NewFactorValues([]int{1}, []int{2, 2}, []float64{10, 20}),
		varlist: []int{1},
		r:       NewFactorValues([]int{}, []int{2, 2}, []float64{30}),
	},
	{
		f:       NewFactorValues([]int{0}, []int{1}, []float64{15}),
		varlist: []int{0},
		r:       NewFactorValues([]int{}, []int{1}, []float64{15}),
	},
	{
		f: NewFactorValues([]int{1, 2, 3}, []int{2, 2, 2, 2, 2},
			[]float64{1, 2, 3, 4, 5, 6, 7, 8}),
		varlist: []int{3, 1},
		r:       NewFactorValues([]int{2}, []int{2, 2, 2, 2, 2}, []float64{14, 22}),
	},
	{
		f: NewFactorValues([]int{3, 2, 1}, []int{2, 2, 2, 2, 2},
			[]float64{1, 5, 3, 7, 2, 6, 4, 8}),
		varlist: []int{3, 1},
		r:       NewFactorValues([]int{2}, []int{2, 2, 2, 2, 2}, []float64{14, 22}),
	},
	{
		f: NewFactorValues([]int{1, 2, 3}, []int{2, 2, 2, 2, 2},
			[]float64{1, 2, 3, 4, 5, 6, 7, 8}),
		varlist: []int{2},
		r:       NewFactorValues([]int{1, 3}, []int{2, 2, 2, 2, 2}, []float64{1 + 3, 2 + 4, 5 + 7, 6 + 8}),
	},
	{
		f: NewFactorValues([]int{0, 1}, []int{2, 2, 2},
			[]float64{.25, .35, .35, .05}),
		varlist: []int{0},
		r:       NewFactorValues([]int{1}, []int{2, 2, 2}, []float64{.6, .4}),
	},
	{
		f: NewFactorValues([]int{1, 2}, []int{2, 2, 2},
			[]float64{.20, .22, .40, .18}),
		varlist: []int{2},
		r:       NewFactorValues([]int{1}, []int{2, 2, 2}, []float64{.6, .4}),
	},
}

func TestSumOut(t *testing.T) {
	for _, v := range testSumOut {
		got := v.f.SumOut(v.varlist)
		if !reflect.DeepEqual(v.r.Variables(), got.Variables()) {
			t.Errorf("want %v, got %v", v.r.Variables(), got.Variables())
		}
		for i, x := range v.r.Values() {
			if !utils.FuzzyEqual(x, got.Values()[i]) {
				t.Errorf("want %v, got %v", v.r.Values(), got.Values())
			}
		}
	}
}

func TestDivision(t *testing.T) {
	cases := []struct {
		f      *Factor
		g      *Factor
		result *Factor
	}{{
		f:      NewFactorValues([]int{0, 1}, []int{2, 2}, []float64{10, 20, 30, 40}),
		g:      NewFactorValues([]int{0, 1}, []int{2, 2}, []float64{5, 5, 30, 5}),
		result: NewFactorValues([]int{0, 1}, []int{2, 2}, []float64{2, 4, 1, 8}),
	}, {
		f:      NewFactorValues([]int{0, 1}, []int{2, 2}, []float64{10, 20, 30, 60}),
		g:      NewFactorValues([]int{1}, []int{2, 2}, []float64{5, 30}),
		result: NewFactorValues([]int{0, 1}, []int{2, 2}, []float64{2, 4, 1, 2}),
	}, {
		f:      NewFactorValues([]int{1, 2}, []int{2, 2, 2}, []float64{20, 22, 40, 18}),
		g:      NewFactorValues([]int{0, 1}, []int{2, 2, 2}, []float64{25, 35, 35, 5}),
		result: NewFactorValues([]int{1, 2}, []int{2, 2, 2}, []float64{.333333, .55, .666666, .45}),
	}}
	for _, tt := range cases {
		got := tt.f.Division(tt.g)
		if !reflect.DeepEqual(tt.result.Variables(), got.Variables()) {
			t.Errorf("want %v, got %v", tt.result.Variables(), got.Variables())
		}
		for j, v := range tt.result.Values() {
			if !utils.FuzzyEqual(v, got.Values()[j], 1e-6) {
				t.Errorf("want %v, got %v", tt.result.Values(), got.Values())
			}
		}
	}
	for _, tt := range cases {
		fvalues := append([]float64(nil), tt.f.Values()...)
		gvalues := append([]float64(nil), tt.g.Values()...)
		tt.f.Division(tt.g)
		if !reflect.DeepEqual(tt.f.Values(), fvalues) {
			t.Errorf("factor changed want %v, got %v", tt.f.Values(), fvalues)
		}
		if !reflect.DeepEqual(tt.g.Values(), gvalues) {
			t.Errorf("factor changed want %v, got %v", tt.g.Values(), gvalues)
		}
	}
}

func TestGetEvidValue(t *testing.T) {
	cases := []struct {
		f      *Factor
		evid   []int
		result float64
	}{{
		NewFactorValues([]int{0, 1}, []int{2, 2}, []float64{10, 20, 30, 40}),
		[]int{1, 1, 1, 1}, 40,
	}, {
		NewFactorValues([]int{1}, []int{2, 2}, []float64{5, 30}),
		[]int{0, 1, 0, 0}, 30,
	}, {
		NewFactorValues([]int{1, 2}, []int{2, 2, 2}, []float64{20, 22, 40, 18}),
		[]int{0, 1, 0, 0}, 22,
	}}
	for _, tt := range cases {
		got := tt.f.GetEvidValue(tt.evid)
		if tt.result != got {
			t.Errorf("wrong value, want %v, got %v", tt.result, got)
		}
	}
}
