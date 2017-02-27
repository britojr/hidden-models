package cliquetree

import (
	"sort"

	"github.com/britojr/kbn/factor"
	"github.com/britojr/tcc/characteristic"
)

// CliqueTree ..
type CliqueTree struct {
	nodes         []node
	send, receive []*factor.Factor
	prev, post    [][]*factor.Factor
	parent        []int
}

type node struct {
	varlist       []int          // wich variables participate on this clique
	neighbours    []int          // cliques that are adjacent to this one
	origPot       *factor.Factor // original clique potential
	currPot       *factor.Factor // initial clique potential for calibration
	calibratedPot *factor.Factor // calibrated potential
	sepset        []int
}

// New ..
func New(n int) *CliqueTree {
	c := new(CliqueTree)
	c.nodes = make([]node, n)
	return c
}

// Size returns the number of cliques
func (c *CliqueTree) Size() int {
	return len(c.nodes)
}

// SetClique ..
func (c *CliqueTree) SetClique(i int, varlist []int) {
	c.nodes[i].varlist = varlist
}

// SetSepSet ..
func (c *CliqueTree) SetSepSet(i int, varlist []int) {
	c.nodes[i].sepset = varlist
}

// Clique ..
func (c *CliqueTree) Clique(i int) []int {
	return c.nodes[i].varlist
}

// SepSet ..
func (c *CliqueTree) SepSet(i int) []int {
	return c.nodes[i].sepset
}

// SetNeighbours ..
func (c *CliqueTree) SetNeighbours(i int, neighbours []int) {
	c.nodes[i].neighbours = neighbours
}

// SetPotential ..
func (c *CliqueTree) SetPotential(i int, potential *factor.Factor) {
	c.nodes[i].origPot = potential
	c.nodes[i].currPot = potential
}

// SetAllPotentials ..
func (c *CliqueTree) SetAllPotentials(potentials []*factor.Factor) {
	for i, potential := range potentials {
		c.nodes[i].origPot = potential
		c.nodes[i].currPot = potential
	}
}

// GetBkpPotential ..
func (c *CliqueTree) GetBkpPotential(i int) *factor.Factor {
	return c.nodes[i].origPot
}

// GetPotential ..
func (c *CliqueTree) GetPotential(i int) *factor.Factor {
	return c.nodes[i].currPot
}

// RestrictByEvidence applies an evidence tuple to each potential on the clique tree
func (c *CliqueTree) RestrictByEvidence(evidence []int) {
	for i := range c.nodes {
		c.nodes[i].currPot = c.nodes[i].origPot.Restrict(evidence)
	}
}

// Calibrated ..
func (c *CliqueTree) Calibrated(i int) *factor.Factor {
	if c.nodes[i].calibratedPot == nil {
		panic("Clique tree wasn't calibrated")
	}
	return c.nodes[i].calibratedPot
}

// UpDownCalibration ..
func (c *CliqueTree) UpDownCalibration() {
	c.send = make([]*factor.Factor, len(c.nodes))
	c.receive = make([]*factor.Factor, len(c.nodes))
	c.prev = make([][]*factor.Factor, len(c.nodes))
	c.post = make([][]*factor.Factor, len(c.nodes))
	root := 0

	c.upwardmessage(root, -1)
	c.downwardmessage(-1, root)
}

func (c *CliqueTree) upwardmessage(v, pa int) {
	c.prev[v] = make([]*factor.Factor, 1, len(c.nodes[v].neighbours)+1)
	c.prev[v][0] = c.nodes[v].currPot
	if len(c.nodes[v].neighbours) > 1 {
		for _, ne := range c.nodes[v].neighbours {
			if ne != pa {
				c.upwardmessage(ne, v)
				c.prev[v] = append(c.prev[v], c.send[ne].Product(c.prev[v][len(c.prev[v])-1]))
			}
		}
	}
	if pa != -1 {
		c.send[v] = c.prev[v][len(c.prev[v])-1].Marginalize(c.nodes[pa].varlist)
	}
}

func (c *CliqueTree) downwardmessage(pa, v int) {
	c.nodes[v].calibratedPot = c.prev[v][len(c.prev[v])-1]
	n := len(c.nodes[v].neighbours)
	if pa != -1 {
		c.nodes[v].calibratedPot = c.nodes[v].calibratedPot.Product(c.receive[v])
		n--
	}
	if len(c.nodes[v].neighbours) == 1 && pa != -1 {
		return
	}

	c.post[v] = make([]*factor.Factor, n)
	i := len(c.post[v]) - 1
	c.post[v][i] = c.receive[v]
	i--
	for k := len(c.nodes[v].neighbours) - 1; k >= 0 && i >= 0; k-- {
		ch := c.nodes[v].neighbours[k]
		if ch == pa {
			continue
		}
		c.post[v][i] = c.send[ch]
		if c.post[v][i+1] != nil {
			c.post[v][i] = c.post[v][i].Product(c.post[v][i+1])
		}
		i--
	}

	k := 0
	for _, ch := range c.nodes[v].neighbours {
		if ch == pa {
			continue
		}
		msg := c.prev[v][k]
		if c.post[v][k] != nil {
			msg = msg.Product(c.post[v][k])
		}
		c.receive[ch] = msg.Marginalize(c.nodes[ch].varlist)
		c.downwardmessage(v, ch)
		k++
	}
}

// FromCharTree generates a clique tree from a characteristic tree and an inverse phi array
func FromCharTree(T *characteristic.Tree, iphi []int) *CliqueTree {
	// determine number of variables (n) and treewidth (k)
	n := len(iphi)
	k := n - len(T.P) + 1

	// create children matrix
	children := make([][]int, len(T.P))
	for i := 0; i < len(T.P); i++ {
		if T.P[i] != -1 {
			children[T.P[i]] = append(children[T.P[i]], i)
		}
	}

	// create relabled cliques list
	cliques := make([][]int, len(children))
	cliques[0] = make([]int, k)
	// Initialize auxiliar (not relabled) clique matrix
	K := make([][]int, n-k+1)
	K[0] = make([]int, k)
	for i := 0; i < k; i++ {
		K[0][i] = n - (k - i) + 1
		cliques[0][i] = iphi[K[0][i]-1]
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
		cliques[v] = make([]int, len(K[v])+1)
		cliques[v][0] = iphi[v-1]
		for i := 0; i < len(K[v]); i++ {
			cliques[v][i+1] = iphi[K[v][i]-1]
		}
	}

	// create new clique tree
	c := New(len(children))
	// initialize root clique
	c.SetClique(0, cliques[0])
	c.SetNeighbours(0, children[0])

	for i := 1; i < c.Size(); i++ {
		// set cliques and sepset
		c.SetClique(i, cliques[i])
		c.SetSepSet(i, cliques[i][1:])
		// set adjacency list as children plus parent
		children[i] = append(children[i], T.P[i])
		c.SetNeighbours(i, children[i])

	}
	// set parents slice
	c.parent = T.P

	return c
}
