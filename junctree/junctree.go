// Package junctree implements junctiontree
package junctree

import "github.com/britojr/tcc/characteristic"

// Node ...
type Node struct {
	Cliq, Sep []int
}

// JuncTree ...
type JuncTree struct {
	Nodes    []Node
	Children [][]int
}

// FromCharTree generates a junction tree from a characteristic tree and an inverse phi array
func FromCharTree(T *characteristic.Tree, iphi []int) *JuncTree {
	// add other root's children to the first root child
	children := characteristic.ChildrenList(T)
	first := children[0][0]
	children[first] = append(children[first], children[0][1:]...)

	// extract the clique list form the characteristic tree
	n := len(iphi)
	k := n - len(T.P) + 1
	K := characteristic.ExtractCliqueList(T, n, k)

	// visit the nodes in BFS order creating the clique with the relabled value
	// of each variable according to inverse phi
	index := 0
	mapIndex := make([]int, len(children))
	queue := []int{first}
	jt := new(JuncTree)
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
	return jt
}
