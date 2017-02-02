package counting

import (
	"reflect"
	"testing"
)

var cardinality = []int{2, 3, 2, 4}
var dataset = [][]int{
	{0, 0, 0, 0}, //               = 0
	{0, 2, 0, 2}, //2*2 + 2*12     = 28
	{1, 0, 1, 3}, //1+ 1*6 + 3*12  = 43
	{0, 0, 1, 1}, //6 + 12         = 18
	{0, 2, 0, 1}, //2*2 + 12       = 16
	{1, 1, 1, 1}, //1 +2 +6 +12    = 21
	{1, 2, 0, 3}, //1 +2*2 + 3*12  = 41
	{0, 1, 1, 0}, //2+6            = 8
	{0, 0, 1, 1}, //1*6 +1*12      = 18
	{0, 0, 0, 0}, //               = 0
	{0, 0, 1, 1}, //1*6 +1*12      = 18
}
var sp = SparseTable{
	strideMap: map[int]int{
		0: 1,
		1: 2,
		2: 6,
		3: 12,
	},
	countMap: map[int]int{
		0:  2,
		8:  1,
		18: 3,
		16: 1,
		21: 1,
		28: 1,
		41: 1,
		43: 1,
	},
	varOrdering: []int{0, 1, 2, 3},
	cardinality: map[int]int{
		0: 2,
		1: 3,
		2: 2,
		3: 4,
	},
}

// { 0,  0}, = 0
// { 0,  0}, = 0
// { 1,  0}, = 1
// { 0,  1}, = 3
// { 0,  1}, = 3
// { 0,  1}, = 3
// { 1,  1}, = 4
// { 2,  1}, = 5
// { 2,  2}, = 8
// { 0,  3}, = 9
// { 2,  3}, = 11
var sp13 = SparseTable{
	strideMap: map[int]int{
		1: 1,
		3: 3,
	},
	countMap: map[int]int{
		0:  2,
		1:  1,
		3:  3,
		4:  1,
		5:  1,
		8:  1,
		9:  1,
		11: 1,
	},
}

// { 0, 0, 0}, = 0
// { 0, 0, 0}, = 0
// { 1, 1, 0}, = 4
// { 2, 0, 1}, = 8
// { 0, 1, 1}, = 9
// { 0, 1, 1}, = 9
// { 0, 1, 1}, = 9
// { 1, 1, 1}, = 10
// { 2, 0, 2}, = 14
// { 2, 0, 3}, = 20
// { 0, 1, 3}, = 21

var spReduc = SparseTable{
	strideMap: map[int]int{
		1: 1,
		2: 3,
		3: 6,
	},
	countMap: map[int]int{
		0:  2,
		4:  1,
		8:  1,
		9:  3,
		10: 1,
		14: 1,
		20: 1,
		21: 1,
	},
}

func TestLoadFromData(t *testing.T) {
	want := &sp
	got := NewSparse()
	got.LoadFromData(dataset, cardinality)
	if !reflect.DeepEqual(want, got) {
		t.Errorf("got %v; want %v", got, want)
	}
}

/*func TestMarginalize(t *testing.T) {
	want := &sp13
	got := sp.Marginalize(1, 3)
	if !reflect.DeepEqual(want, got) {
		t.Errorf("got %v; want %v", got, want)
	}
}*/

func TestReduce(t *testing.T) {
	want := &spReduc
	got := sp.Reduce()
	if !reflect.DeepEqual(want, got) {
		t.Errorf("got %v; want %v", got, want)
	}
}
