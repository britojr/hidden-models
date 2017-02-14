package factor

import (
	"testing"

	"github.com/britojr/kbn/assignment"
)

type testStruct struct {
	cardin  []int
	varlist []int
	values  []float64
}

var f01 = testStruct{
	cardin:  []int{3, 2, 2},
	varlist: []int{0, 1},
	values:  []float64{0.5, 0.8, 0.1, 0, 0.3, 0, 9},
}
var f12 = testStruct{
	cardin:  []int{3, 2, 2},
	varlist: []int{1, 2},
	values:  []float64{0.5, 0.7, 0.1, 0.2},
}
var f012 = testStruct{
	cardin:  []int{3, 2, 2},
	varlist: []int{0, 1, 2},
	values:  []float64{0.25, 0.35, 0.08, 0.16, 0.05, 0.07, 0, 0, 0.15, 0.21, 0.09, 0.18},
}

var f012_0 = testStruct{
	cardin:  []int{3, 2, 2},
	varlist: []int{1, 2},
	values:  []float64{0.45, 0.63, 0.17, 0.34},
}
var f012_1 = testStruct{
	cardin:  []int{3, 2, 2},
	varlist: []int{0, 2},
	values:  []float64{0.33, 0.51, 0.05, 0.07, 0.24, 0.39},
}
var f012_2 = testStruct{
	cardin:  []int{3, 2, 2},
	varlist: []int{0, 1},
	values:  []float64{0.60, 0.24, 0.12, 0, 0.36, 0.27},
}

var tests = []testStruct{
	f01,
	f12,
	f012,
}

var sumOutTests = []testStruct{
	f012_0,
	f012_1,
	f012_2,
}

func TestNew(t *testing.T) {
	for _, w := range tests {
		got := New(w.varlist, w.cardin, w.values)
		assig := assignment.New(w.varlist, w.cardin)
		i := 0
		for assig != nil {
			v := got.Get(assig)
			if v != w.values[i] {
				t.Errorf("want(%v); got(%v)", w.values[i], v)
			}
			assig.Next()
			i++
		}
	}
}

func TestProduct(t *testing.T) {
	a := New(f01.varlist, f01.cardin, f01.values)
	b := New(f12.varlist, f12.cardin, f12.values)
	got := a.Product(b)
	assig := assignment.New(f012.varlist, f012.cardin)
	i := 0
	for assig != nil {
		v := got.Get(assig)
		if v != f012.values[i] {
			t.Errorf("want(%v); got(%v)", f012.values[i], v)
		}
		assig.Next()
		i++
	}
}

func TestSumOut(t *testing.T) {
	a := New(f012.varlist, f012.cardin, f012.values)
	for i, w := range sumOutTests {
		got := a.SumOut(i)
		assig := assignment.New(w.varlist, w.cardin)
		i := 0
		for assig != nil {
			v := got.Get(assig)
			if v != w.values[i] {
				t.Errorf("want(%v); got(%v)", w.values[i], v)
			}
			assig.Next()
			i++
		}
	}
}
