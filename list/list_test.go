package list

import (
	"reflect"
	"sort"
	"testing"

	"github.com/willf/bitset"
)

func TestSplit(t *testing.T) {
	cases := []struct {
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
	for _, v := range cases {
		a, b := Split(v.slice, v.n)
		if !reflect.DeepEqual(a, v.a) {
			t.Errorf("got %v, want %v", a, v.a)
		}
		if !reflect.DeepEqual(b, v.b) {
			t.Errorf("got %v, want %v", b, v.b)
		}
	}
}

func TestOrderedDiff(t *testing.T) {
	cases := []struct {
		a, b, inter, in, out []int
	}{{
		a:     []int{2, 3, 4},
		b:     []int{2, 4, 5},
		inter: []int{2, 4},
		in:    []int{5},
		out:   []int{3},
	}, {
		a:     []int{5, 6, 7},
		b:     []int{2, 4, 5},
		inter: []int{5},
		in:    []int{2, 4},
		out:   []int{6, 7},
	}}
	for _, tt := range cases {
		inter, in, out := OrderedDiff(tt.a, tt.b)
		if !reflect.DeepEqual(tt.inter, inter) {
			t.Errorf("wrong inter, want %v, got %v", tt.inter, inter)
		}
		if !reflect.DeepEqual(tt.in, in) {
			t.Errorf("wrong in, want %v, got %v", tt.in, in)
		}
		if !reflect.DeepEqual(tt.out, out) {
			t.Errorf("wrong out, want %v, got %v", tt.out, out)
		}
	}
}

func TestUnion(t *testing.T) {
	cases := []struct {
		a, b, res []int
	}{
		{[]int{}, []int{}, []int{}},
		{[]int{}, []int{1}, []int{1}},
		{[]int{1}, []int{}, []int{1}},
		{[]int{1}, []int{1}, []int{1}},
		{[]int{2}, []int{1}, []int{1, 2}},
		{[]int{6, 4, 2, 8}, []int{8, 9, 3, 1, 2}, []int{1, 2, 3, 4, 6, 8, 9}},
	}
	for _, v := range cases {
		got := Union(v.a, v.b)
		sort.Ints(got)
		if !reflect.DeepEqual(got, v.res) {
			t.Errorf("got %v want %v", got, v.res)
		}
	}
}

func TestDifference(t *testing.T) {
	cases := []struct {
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
	for _, v := range cases {
		got := Difference(v.a, v.b)
		sort.Ints(got)
		if !reflect.DeepEqual(got, v.res) {
			t.Errorf("got %v want %v", got, v.res)
		}
	}
}

func TestIntersectionBits(t *testing.T) {
	cases := []struct {
		lst    [][]int
		result []int
	}{{
		[][]int{{3, 5, 7},
			{3, 4, 5, 0},
			{2, 1, 5, 0, 3},
		},
		[]int{3, 5},
	}, {
		[][]int{{1, 9, 7},
			{7, 3},
			{8, 8, 7},
			{9, 7},
		},
		[]int{7},
	}, {
		[][]int{{7, 3}},
		[]int{3, 7},
	}, {
		[][]int{},
		[]int{},
	}}
	for _, v := range cases {
		setlist := make([]*bitset.BitSet, len(v.lst))
		for i := range v.lst {
			setlist[i] = NewBitSet()
			SetBits(setlist[i], v.lst[i])
		}
		b := IntersectionBits(setlist)
		got := FromBitSet(b)
		if !reflect.DeepEqual(v.result, got) {
			t.Errorf("want %v,  got %v", v.result, got)
		}
	}
}
