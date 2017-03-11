package utils

import (
	"reflect"
	"sort"
	"testing"

	"github.com/willf/bitset"
)

var testFuzzyEqual = []struct {
	a, b  float64
	equal bool
}{
	{0.1, 0.2, false},
	{0, 0, true},
	{0.0002, 0.0002, true},
	{1.0005, 1.0005 + (epslon / 2.0), true},
	{0.0005, 0.0005 + epslon, false},
}

func TestFuzzyEqual(t *testing.T) {
	for _, v := range testFuzzyEqual {
		got := FuzzyEqual(v.a, v.b)
		if got != v.equal {
			t.Errorf("%v == %v : got %v, want %v", v.a, v.b, got, v.equal)
		}
	}
}

var testSliceSplit = []struct {
	slice []int
	n     int
	a, b  []int
}{
	{[]int{3, 4, 8, 9, 1, 6, 2, 0}, 6, []int{3, 4, 1, 2, 0}, []int{8, 9, 6}},
	{[]int{3, 4, 8, 9, 1, 6, 2, 0}, 0, []int{}, []int{3, 4, 8, 9, 1, 6, 2, 0}},
	{[]int{3, 4, 8}, 9, []int{3, 4, 8}, []int{}},
	{[]int{8}, 8, []int{}, []int{8}},
	{[]int{}, 8, []int{}, []int{}},
}

func TestSliceSplit(t *testing.T) {
	for _, v := range testSliceSplit {
		a, b := SliceSplit(v.slice, v.n)
		if !reflect.DeepEqual(a, v.a) {
			t.Errorf("got %v, want %v", a, v.a)
		}
		if !reflect.DeepEqual(b, v.b) {
			t.Errorf("got %v, want %v", b, v.b)
		}
	}

}

var testSliceUnion = []struct {
	a, b, res []int
}{
	{[]int{}, []int{}, []int{}},
	{[]int{}, []int{1}, []int{1}},
	{[]int{1}, []int{}, []int{1}},
	{[]int{1}, []int{1}, []int{1}},
	{[]int{2}, []int{1}, []int{1, 2}},
	{[]int{6, 4, 2, 8}, []int{8, 9, 3, 1, 2}, []int{1, 2, 3, 4, 6, 8, 9}},
}

func TestSliceUnion(t *testing.T) {
	for _, v := range testSliceUnion {
		got := SliceUnion(v.a, v.b)
		sort.Ints(got)
		if !reflect.DeepEqual(got, v.res) {
			t.Errorf("got %v want %v", got, v.res)
		}
	}
}

var testSliceDifference = []struct {
	a, b, res []int
}{
	{[]int{}, []int{}, []int{}},
	{[]int{}, []int{1}, []int{}},
	{[]int{1}, []int{}, []int{1}},
	{[]int{1}, []int{1}, []int{}},
	{[]int{2}, []int{1}, []int{2}},
	{[]int{6, 4, 2, 8}, []int{8, 9, 3, 1, 2}, []int{4, 6}},
	{[]int{5, 7}, []int{7, 5, 6}, []int{}},
	{[]int{5, 7}, []int{1, 2, 3}, []int{5, 7}},
}

func TestSliceDifference(t *testing.T) {
	for _, v := range testSliceDifference {
		got := SliceDifference(v.a, v.b)
		sort.Ints(got)
		if !reflect.DeepEqual(got, v.res) {
			t.Errorf("got %v want %v", got, v.res)
		}
	}
}

var testNormalizeSlice = []struct {
	values, normalized []float64
}{
	{
		[]float64{0.15, 0.25, 0.35, 0.25},
		[]float64{0.15, 0.25, 0.35, 0.25},
	},
	{
		[]float64{15, 25, 35, 25},
		[]float64{0.15, 0.25, 0.35, 0.25},
	},
	{
		[]float64{10, 20, 30, 40, 50, 60, 70, 80},
		[]float64{1.0 / 36, 2.0 / 36, 3.0 / 36, 4.0 / 36, 5.0 / 36, 6.0 / 36, 7.0 / 36, 8.0 / 36},
	},
	{
		[]float64{0.15},
		[]float64{1},
	},
	// {
	// 	[]float64{},
	// 	[]float64{},
	// },
	// {
	// 	[]float64{0, 0, 0},
	// 	[]float64{0, 0, 0},
	// },
}

func TestNormalizeSlice(t *testing.T) {
	for _, v := range testNormalizeSlice {
		NormalizeSlice(v.values)
		if !reflect.DeepEqual(v.values, v.normalized) {
			t.Errorf("want %v, got %v", v.normalized, v.values)
		}
	}
}

var testNormalizeIntSlice = []struct {
	values     []int
	normalized []float64
}{
	{
		[]int{15, 25, 35, 25},
		[]float64{0.15, 0.25, 0.35, 0.25},
	},
	{
		[]int{10, 20, 30, 40, 50, 60, 70, 80},
		[]float64{1.0 / 36, 2.0 / 36, 3.0 / 36, 4.0 / 36, 5.0 / 36, 6.0 / 36, 7.0 / 36, 8.0 / 36},
	},
	{
		[]int{15},
		[]float64{1},
	},
	{
		[]int{},
		[]float64{},
	},
	{
		[]int{0, 0, 0},
		[]float64{0, 0, 0},
	},
}

func TestNormalizeIntSlice(t *testing.T) {
	for _, v := range testNormalizeIntSlice {
		got := NormalizeIntSlice(v.values)
		if !reflect.DeepEqual(v.normalized, got) {
			t.Errorf("want %v, got %v", v.normalized, got)
		}
	}
}

var testListIntersection = []struct {
	list   [][]int
	result []int
}{
	{
		[][]int{
			{3, 5, 7},
			{3, 4, 5, 0},
			{2, 1, 5, 0, 3},
		},
		[]int{3, 5},
	},
	{
		[][]int{
			{1, 9, 7},
			{7, 3},
			{8, 8, 7},
			{9, 7},
		},
		[]int{7},
	},
	{
		[][]int{
			{7, 3},
		},
		[]int{3, 7},
	},
	{
		[][]int{},
		[]int{},
	},
}

func TestListIntersection(t *testing.T) {
	for _, v := range testListIntersection {
		setlist := make([]*bitset.BitSet, len(v.list))
		for i := range v.list {
			setlist[i] = NewBitSet()
			SetSlice(setlist[i], v.list[i])
		}
		b := ListIntersection(setlist)
		got := SliceFromBitSet(b)
		if !reflect.DeepEqual(v.result, got) {
			t.Errorf("want %v,  got %v", v.result, got)
		}
	}
}

var testSliceSumFloat64 = []struct {
	values []float64
	sum    float64
}{
	{
		[]float64{5, 5},
		10,
	},
	{
		[]float64{1.5, 3.5, 0.5},
		5.5,
	},
	{
		[]float64{},
		0,
	},
	{
		[]float64(nil),
		0,
	},
}

func TestSliceSumFloat64(t *testing.T) {
	for _, v := range testSliceSumFloat64 {
		got := SliceSumFloat64(v.values)
		if v.sum != got {
			t.Errorf("want %v, got %v", v.sum, got)
		}
	}
}
