// Package jtree implements junctiontree
package jtree

import (
	"sort"

	"github.com/britojr/tcc/characteristic"
)

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
	n := len(iphi)
	k := n - len(T.P) + 1

	// Create children vector from T.P.
	children := make([][]int, n)
	for i := 0; i < len(T.P); i++ {
		if T.P[i] != -1 {
			children[T.P[i]] = append(children[T.P[i]], i)
		}
	}

	// Visit T in BFS order, starting with the children of R.
	K := make([][]int, n)
	queue := make([]int, n)
	m := make([]bool, n)
	start, end := 0, 0
	for i := 0; i < len(children[0]); i++ {
		m[children[0][i]] = true
		queue[end] = children[0][i]
		end++
	}
	for start != end {
		v := queue[start]
		start++
		if T.P[v] == 0 {
			for i := n - k + 1; i <= n; i++ {
				K[v] = append(K[v], i)
			}
		} else {
			for i := 0; i < len(K[T.P[v]]); i++ {
				if i != T.L[v] {
					K[v] = append(K[v], K[T.P[v]][i])
				}
			}
			K[v] = append(K[v], T.P[v])
			sort.Ints(K[v])
		}
		// for i := 0; i < len(K[v]); i++ {
		// 	u := K[v][i]
		// 	// adj[u-1] = append(adj[u-1], v-1)
		// 	adj[v-1] = append(adj[v-1], u-1)
		// }
		for i := 0; i < len(children[v]); i++ {
			if !m[children[v][i]] {
				m[children[v][i]] = true
				queue[end] = children[v][i]
				end++
			}
		}
	}

	return &JTree{nil, nil}
}
