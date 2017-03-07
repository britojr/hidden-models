package cliquetree

import (
	"sort"

	"github.com/britojr/kbn/factor"
	"github.com/britojr/kbn/utils"
	"github.com/britojr/tcc/characteristic"
)

// CliqueTree ..
type CliqueTree struct {
	cliques [][]int // wich variables participate on this clique
	sepsets [][]int // sepsets for each node (intersection with the parent clique)

	neighbours [][]int // cliques that are adjacent to this one, including parent
	parent     []int   // the parent of each node

	origPot       []*factor.Factor // original clique potential
	currPot       []*factor.Factor // initial clique potential for calibration
	calibratedPot []*factor.Factor // calibrated potential

	// auxiliar for message passing, send to parent and receive from parent
	send, receive []*factor.Factor
	// axiliar to reduce (memoize) number of factor multiplications
	prev, post [][]*factor.Factor
}

// New ..
func New(n int) *CliqueTree {
	c := new(CliqueTree)
	c.cliques = make([][]int, n)
	c.sepsets = make([][]int, n)
	c.neighbours = make([][]int, n)
	c.parent = make([]int, n)
	c.origPot = make([]*factor.Factor, n)
	c.currPot = make([]*factor.Factor, n)
	c.calibratedPot = make([]*factor.Factor, n)
	return c
}

// Size returns the number of cliques
func (c *CliqueTree) Size() int {
	return len(c.cliques)
}

// SetClique ..
func (c *CliqueTree) SetClique(i int, varlist []int) {
	c.cliques[i] = varlist
}

// SetSepSet ..
func (c *CliqueTree) SetSepSet(i int, varlist []int) {
	c.sepsets[i] = varlist
}

// Clique returns the ith clique
func (c *CliqueTree) Clique(i int) []int {
	return c.cliques[i]
}

// Cliques returns the complete clique list ..
func (c *CliqueTree) Cliques() [][]int {
	return c.cliques
}

// SepSet returns the ith sepset
func (c *CliqueTree) SepSet(i int) []int {
	return c.sepsets[i]
}

// SepSets returns complete sepset list
func (c *CliqueTree) SepSets() [][]int {
	return c.sepsets
}

// SetNeighbours ..
func (c *CliqueTree) SetNeighbours(i int, neighbours []int) {
	c.neighbours[i] = neighbours
}

// SetPotential ..
func (c *CliqueTree) SetPotential(i int, potential *factor.Factor) {
	c.origPot[i] = potential
	c.currPot[i] = potential
}

// SetAllPotentials ..
func (c *CliqueTree) SetAllPotentials(potentials []*factor.Factor) {
	c.origPot = append([]*factor.Factor(nil), potentials...)
	c.currPot = append([]*factor.Factor(nil), potentials...)
}

// BkpPotentialList returns a list with all the original potentials
func (c *CliqueTree) BkpPotentialList() []*factor.Factor {
	return c.origPot
}

// BkpPotential ..
func (c *CliqueTree) BkpPotential(i int) *factor.Factor {
	return c.origPot[i]
}

// CurrPotential ..
func (c *CliqueTree) CurrPotential(i int) *factor.Factor {
	return c.currPot[i]
}

// ReduceByEvidence applies an evidence tuple to each potential on the clique tree
func (c *CliqueTree) ReduceByEvidence(evidence []int) {
	for i := range c.currPot {
		c.currPot[i] = c.origPot[i].Reduce(evidence)
	}
}

// Calibrated returns the calibrated potential for the ith clique
func (c *CliqueTree) Calibrated(i int) *factor.Factor {
	if c.calibratedPot[i] == nil {
		panic("Clique tree wasn't calibrated")
	}
	return c.calibratedPot[i]
}

// CalibratedSepSet returns the calibrated potential for the sepset of the ith clique
func (c *CliqueTree) CalibratedSepSet(i int) *factor.Factor {
	if c.calibratedPot[i] == nil {
		panic("Clique tree wasn't calibrated")
	}
	if c.parent[i] < 0 {
		return nil
	}
	diff := utils.SliceDifference(c.Clique(i), c.Clique(c.parent[i]))
	return c.calibratedPot[i].SumOut(diff)
}

// LoadCalibration ..
func (c *CliqueTree) LoadCalibration() {
	for i := range c.calibratedPot {
		c.calibratedPot[i] = c.currPot[i]
	}
}

// UpDownCalibration ..
func (c *CliqueTree) UpDownCalibration() {
	c.send = make([]*factor.Factor, c.Size())
	c.receive = make([]*factor.Factor, c.Size())
	// post[i][j] contains the product of every message that node i received
	// from its j+1 children to the last children
	// prev[i][j] contains the product of node i initial potential and
	// every message that node i received from its fist children to the j-1 children
	// So the message to be sent from i to j will be the product of prev and post
	c.prev = make([][]*factor.Factor, c.Size())
	c.post = make([][]*factor.Factor, c.Size())
	root := 0

	c.upwardmessage(root, -1)
	c.downwardmessage(-1, root)
}

func (c *CliqueTree) upwardmessage(v, pa int) {
	c.prev[v] = make([]*factor.Factor, 1, len(c.neighbours[v])+1)
	c.prev[v][0] = c.currPot[v]
	if len(c.neighbours[v]) > 1 {
		for _, ne := range c.neighbours[v] {
			if ne != pa {
				c.upwardmessage(ne, v)
				c.prev[v] = append(c.prev[v], c.send[ne].Product(c.prev[v][len(c.prev[v])-1]))
			}
		}
	}
	if pa != -1 {
		c.send[v] = c.prev[v][len(c.prev[v])-1].Marginalize(c.cliques[pa])
	}
}

func (c *CliqueTree) downwardmessage(pa, v int) {
	c.calibratedPot[v] = c.prev[v][len(c.prev[v])-1]
	n := len(c.neighbours[v])
	if pa != -1 {
		c.calibratedPot[v] = c.calibratedPot[v].Product(c.receive[v])
		n--
	}
	if len(c.neighbours[v]) == 1 && pa != -1 {
		return
	}

	c.post[v] = make([]*factor.Factor, n)
	i := len(c.post[v]) - 1
	c.post[v][i] = c.receive[v]
	i--
	for k := len(c.neighbours[v]) - 1; k >= 0 && i >= 0; k-- {
		ch := c.neighbours[v][k]
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
	for _, ch := range c.neighbours[v] {
		if ch == pa {
			continue
		}
		msg := c.prev[v][k]
		if c.post[v][k] != nil {
			msg = msg.Product(c.post[v][k])
		}
		c.receive[ch] = msg.Marginalize(c.cliques[ch])
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
