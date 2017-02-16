// Package junctree implements junctiontree
package junctree

import (
	"sort"

	"github.com/britojr/tcc/characteristic"
)

// Node of a junctree has a clique a Sepset (intersection with the parent clique)
type Node struct {
	Clique []int
	Sepset []int
}

// JuncTree ...
type JuncTree struct {
	Nodes    []Node
	Children [][]int
}

// FromCharTree generates a junction tree from a characteristic tree and an inverse phi array
func FromCharTree(T *characteristic.Tree, iphi []int) *JuncTree {
	jt := new(JuncTree)
	n := len(iphi)
	k := n - len(T.P) + 1

	// create children matrix
	children := make([][]int, len(T.P))
	for i := 0; i < len(T.P); i++ {
		if T.P[i] != -1 {
			children[T.P[i]] = append(children[T.P[i]], i)
		}
	}
	jt.Children = children

	jt.Nodes = make([]Node, len(children))
	jt.Nodes[0].Clique = make([]int, k)
	// Initialize auxiliar (not relabled) clique matrix
	K := make([][]int, n-k+1)
	K[0] = make([]int, k)
	for i := 0; i < k; i++ {
		K[0][i] = n - (k - i) + 1
		jt.Nodes[0].Clique[i] = iphi[K[0][i]-1]
	}

	// Visit T in BFS order, starting with the children of the root
	queue := make([]int, n)
	start, end := 0, 0
	for i := 0; i < len(children[0]); i++ {
		queue[end] = children[0][i]
		end++
	}
	for start != end {
		v := queue[start]
		start++
		// update unlabled clique K
		for i := 0; i < len(K[T.P[v]]); i++ {
			if i != T.L[v] {
				K[v] = append(K[v], K[T.P[v]][i])
			}
		}
		if T.P[v] != 0 {
			K[v] = append(K[v], T.P[v])
			sort.Ints(K[v])
		}

		// enqueue the children of v
		for i := 0; i < len(children[v]); i++ {
			queue[end] = children[v][i]
			end++
		}

		// create the relabled clique
		clique := make([]int, len(K[v])+1)
		clique[0] = iphi[v-1]
		for i := 0; i < len(K[v]); i++ {
			clique[i+1] = iphi[K[v][i]-1]
		}
		jt.Nodes[v].Clique = clique
		jt.Nodes[v].Sepset = clique[1:]
	}

	return jt
}
