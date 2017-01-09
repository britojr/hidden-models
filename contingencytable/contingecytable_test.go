package contingecytable

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
var sp = Sparse{
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
}

func TestLoadFromData(t *testing.T) {
	want := &sp
	got := LoadFromData(dataset, cardinality)
	if !reflect.DeepEqual(want, got) {
		t.Errorf("got %v; want %v", got, want)
	}
}
