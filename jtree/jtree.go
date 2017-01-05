// Package jtree implements junctiontree
package jtree

import "github.com/britojr/tcc/characteristic"

type jTreeNode struct {
	cliq, sep []int
}

// JTree ...
type JTree struct {
	Nodes []jTreeNode
	P     []int
}

// Generate generates a junction tree form a characteristic tree and an inverse phi array
func Generate(T *characteristic.Tree, iphi []int) *JTree {
	return &JTree{nil, nil}
}
