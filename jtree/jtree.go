// Package jtree implements junctiontree
package jtree

import "github.com/britojr/tcc/characteristic"

// Node ...
type Node struct {
	cliq, sep []int
}

// JTree ...
type JTree struct {
	Nodes    []Node
	Children [][]int
}

// New generates a junction tree form a characteristic tree and an inverse phi array
func New(T *characteristic.Tree, iphi []int) *JTree {
	n := len(iphi)
	k := n - len(T.P) + 1
	children := characteristic.ChildrenList(T)
	first := children[0][0]
	children[first] = append(children[first], children[0][1:]...)
	jt := JTree{}
	K := characteristic.ExtractCliqueList(T, n, k)
	queue := []int{first}
	index := 0
	mapIndex := make([]int, len(children))
	for len(queue) > 0 {
		v := queue[0]
		mapIndex[v] = index
		queue = queue[1:]
		clique := make([]int, len(K[v])+1)
		clique[0] = iphi[v-1]
		for i := 0; i < len(K[v]); i++ {
			clique[i+1] = iphi[K[v][i]-1]
		}
		jt.Nodes = append(jt.Nodes, Node{clique, clique[1:]})
		index++
		queue = append(queue, children[v]...)
	}
	jt.Children = make([][]int, len(children)-1)
	for i := 1; i < len(children); i++ {
		jt.Children[mapIndex[i]] = make([]int, len(children[i]))
		for j := 0; j < len(children[i]); j++ {
			jt.Children[mapIndex[i]][j] = mapIndex[children[i][j]]
		}
	}
	return &jt
}
