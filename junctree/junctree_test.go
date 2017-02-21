package junctree

import (
	"reflect"
	"testing"

	"github.com/britojr/tcc/characteristic"
)

var jt = JuncTree{
	Nodes: []Node{
		{Clique: []int{1, 2, 8}, Sepset: []int(nil)},
		{Clique: []int{0, 4, 7, 1}, Sepset: []int{4, 7, 1}},
		{Clique: []int{10, 1, 2, 8}, Sepset: []int{1, 2, 8}},
		{Clique: []int{9, 1, 2, 8}, Sepset: []int{1, 2, 8}},
		{Clique: []int{3, 10, 2, 8}, Sepset: []int{10, 2, 8}},
		{Clique: []int{4, 7, 1, 2}, Sepset: []int{7, 1, 2}},
		{Clique: []int{5, 7, 1, 8}, Sepset: []int{7, 1, 8}},
		{Clique: []int{6, 0, 4, 7}, Sepset: []int{0, 4, 7}},
		{Clique: []int{7, 1, 2, 8}, Sepset: []int{1, 2, 8}},
	},
	/*[][]int{
		{2, 3, 8},
		{7},
		{4},
		[]int(nil),
		[]int(nil),
		{1},
		[]int(nil),
		[]int(nil),
		{5, 6},
	},*/
	Adj: [][]int{
		{2, 3, 8},
		{7, 5},
		{4, 0},
		{0},
		{2},
		{1, 8},
		{8},
		{1},
		{5, 6, 0},
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
	for i := range got.Nodes {
		if !reflect.DeepEqual(got.Nodes[i].Clique, want.Nodes[i].Clique) {
			t.Errorf("Clique[%v]; Got: %v; Want: %v", i, got.Nodes[i].Clique, want.Nodes[i].Clique)
		}
		if !reflect.DeepEqual(got.Nodes[i].Sepset, want.Nodes[i].Sepset) {
			t.Errorf("Sepset[%v]; Got: %v; Want: %v", i, got.Nodes[i].Sepset, want.Nodes[i].Sepset)
		}
		if !reflect.DeepEqual(got.Adj[i], want.Adj[i]) {
			t.Errorf("Adj[%v]; Got: %v; Want: %v", i, got.Adj[i], want.Adj[i])
		}
	}
}
