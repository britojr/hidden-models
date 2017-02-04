package jtree

import (
	"reflect"
	"testing"

	"github.com/britojr/tcc/characteristic"
)

var jt = JTree{
	[]Node{
		{[]int{10, 1, 2, 8}, []int{1, 2, 8}},
		{[]int{3, 10, 2, 8}, []int{10, 2, 8}},
		{[]int{9, 1, 2, 8}, []int{1, 2, 8}},
		{[]int{7, 1, 2, 8}, []int{1, 2, 8}},
		{[]int{4, 7, 1, 2}, []int{7, 1, 2}},
		{[]int{5, 7, 1, 8}, []int{7, 1, 8}},
		{[]int{0, 4, 7, 1}, []int{4, 7, 1}},
		{[]int{6, 0, 4, 7}, []int{0, 4, 7}},
	},
	[][]int{
		{1, 2, 3},
		{}, // (nil),
		{}, // (nil),
		{4, 5},
		{6},
		{}, // (nil),
		{7},
		{}, // (nil),
	},
}

var iphi = []int{0, 10, 9, 3, 4, 5, 6, 7, 1, 2, 8}
var T = characteristic.Tree{
	P: []int{-1, 5, 0, 0, 2, 8, 8, 1, 0},
	L: []int{-1, 2, -1, -1, 0, 2, 1, 2, -1},
}

func TestFromCharTree(t *testing.T) {
	want := &jt
	got := FromCharTree(&T, iphi)
	if !reflect.DeepEqual(got, want) {
		t.Errorf("Got = %v; want %v", got, want)
	}
}
