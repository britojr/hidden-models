package cliquetree

import (
	"bufio"
	"fmt"
	"io"
	"reflect"
	"sort"
	"strings"

	"github.com/britojr/kbn/factor"
	"github.com/britojr/kbn/utils"
	"github.com/britojr/tcc/characteristic"
)

// CliqueTree ..
type CliqueTree struct {
	cliques [][]int // wich variables participate on this clique
	sepsets [][]int // sepsets for each node (intersection with the parent clique)
	varin   [][]int // the difference between clique and parent
	varout  [][]int // the difference between clique and parent

	neighbours [][]int // cliques that are adjacent to this one, including parent
	parent     []int   // the parent of each node

	initialPotStored    []*factor.Factor // original clique potential
	initialPot          []*factor.Factor // initial clique potential for calibration
	calibratedPot       []*factor.Factor // calibrated potential
	calibratedPotSepSet []*factor.Factor // calibrated potential for the sepset
	calibratedPotStored []*factor.Factor // auxiliar for storing calibrated potentials

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
	c.initialPotStored = make([]*factor.Factor, n)
	c.initialPot = make([]*factor.Factor, n)
	c.calibratedPot = make([]*factor.Factor, n)
	c.calibratedPotSepSet = make([]*factor.Factor, n)
	return c
}

// NewStructure creates a new clique tree structure
func NewStructure(cliques, adj [][]int) (*CliqueTree, error) {
	c := new(CliqueTree)
	n := len(cliques)
	if len(adj) != n {
		return nil, fmt.Errorf("wrong size for adjacency list: %v", len(adj))
	}
	c.cliques = make([][]int, n)
	c.neighbours = make([][]int, n)
	for i := range cliques {
		c.cliques[i] = append([]int(nil), cliques[i]...)
		sort.Ints(c.cliques[i])
		c.neighbours[i] = append([]int(nil), adj[i]...)
	}
	c.bfsOrder(0)
	c.sepsets = make([][]int, n)
	c.varin, c.varout = make([][]int, n), make([][]int, n)
	for i := 1; i < n; i++ {
		c.sepsets[i], c.varin[i], c.varout[i] = utils.OrderedSliceDiff(c.cliques[c.parent[i]], c.cliques[i])
	}
	return c, nil
}

func (c *CliqueTree) bfsOrder(root int) []int {
	c.parent = make([]int, len(c.cliques))
	visit := make([]bool, len(c.cliques))
	queue := make([]int, len(c.cliques))
	start, end := 0, 0
	c.parent[root] = -1
	visit[root] = true
	queue[end] = root
	end++
	for start < end {
		v := queue[start]
		start++
		for _, ne := range c.neighbours[v] {
			if !visit[ne] {
				c.parent[ne] = v
				visit[ne] = true
				queue[end] = ne
				end++
			}
		}
	}
	return queue
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

// Varin returns the variables that are on the ith clique and not on its parent
func (c *CliqueTree) Varin(i int) []int {
	return c.varin[i]
}

// SetNeighbours ..
func (c *CliqueTree) SetNeighbours(i int, neighbours []int) {
	c.neighbours[i] = neighbours
}

// Neighbours ..
func (c *CliqueTree) Neighbours(i int) []int {
	return c.neighbours[i]
}

// Parents ..
func (c *CliqueTree) Parents() []int {
	return c.parent
}

// SetPotential ..
func (c *CliqueTree) SetPotential(i int, potential *factor.Factor) {
	c.initialPot[i] = potential
}

// SetAllPotentials ..
func (c *CliqueTree) SetAllPotentials(potentials []*factor.Factor) error {
	c.initialPot = append([]*factor.Factor(nil), potentials...)
	// check potentials scope
	for i := range c.cliques {
		if i >= len(c.initialPot) || c.initialPot[i] == nil {
			return fmt.Errorf("no factor for clique %v", c.cliques[i])
		}
		if !reflect.DeepEqual(c.initialPot[i].Variables(), c.cliques[i]) {
			return fmt.Errorf("wrong scope, clique %v has factor %v", c.cliques[i], c.initialPot[i].Variables())
		}
	}
	return nil
}

// InitialPotential ..
func (c *CliqueTree) InitialPotential(i int) *factor.Factor {
	return c.initialPot[i]
}

// Potentials ..
func (c *CliqueTree) Potentials() []*factor.Factor {
	return c.initialPot
}

// ReduceByEvidence applies an evidence tuple to each potential on the clique tree
func (c *CliqueTree) ReduceByEvidence(evidence []int) {
	for i := range c.initialPot {
		c.initialPot[i] = c.initialPot[i].Reduce(evidence)
	}
}

// StorePotentials stores the initial potential values in order to b able to recover them later
func (c *CliqueTree) StorePotentials() {
	c.initialPotStored = append([]*factor.Factor(nil), c.initialPot...)
}

// RecoverPotentials recover intial potentials previously stored
func (c *CliqueTree) RecoverPotentials() {
	c.initialPot = append([]*factor.Factor(nil), c.initialPotStored...)
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
	return c.calibratedPotSepSet[i]
}

// SetCalibrated ..
func (c *CliqueTree) SetCalibrated(i int, f *factor.Factor) {
	c.calibratedPot[i] = f
}

// SetCalibratedSepSet ..
func (c *CliqueTree) SetCalibratedSepSet(i int, f *factor.Factor) {
	c.calibratedPotSepSet[i] = f
}

// LoadCalibration ..
func (c *CliqueTree) LoadCalibration() {
	for i := range c.calibratedPot {
		c.calibratedPot[i] = c.initialPot[i]
	}
}

// StoreCalibration stores the calibrated values in order to retract them later
func (c *CliqueTree) StoreCalibration() {
	c.calibratedPotStored = append([]*factor.Factor(nil), c.calibratedPot...)
}

// RecoverCalibration recover calibration previously stored
func (c *CliqueTree) RecoverCalibration() {
	c.calibratedPot = append([]*factor.Factor(nil), c.calibratedPotStored...)
}

// ProbOfEvidence ..
func (c *CliqueTree) ProbOfEvidence(evid []int) float64 {
	root := 0
	send := make([]*factor.Factor, c.Size())
	c.upwardreduction(root, -1, evid, send)
	// summout all variables of the resulting (calibrated) root factor
	return utils.SliceSumFloat64(send[root].Values())
}

func (c *CliqueTree) upwardreduction(v, pa int, evid []int, send []*factor.Factor) {
	prev := c.initialPot[v].Reduce(evid)
	for _, ne := range c.neighbours[v] {
		if ne != pa {
			c.upwardreduction(ne, v, evid, send)
			prev = prev.Product(send[ne])
		}
	}
	send[v] = prev.SumOut(c.varin[v])
}

// UpDownCalibration ..
func (c *CliqueTree) UpDownCalibration() {
	// -------------------------------------------------------------------------
	// send[i] contains the message that ith node sends up to its parent
	// receive[i] contains the message that ith node receives from his parent
	// -------------------------------------------------------------------------
	c.send = make([]*factor.Factor, c.Size())
	c.receive = make([]*factor.Factor, c.Size())
	// -------------------------------------------------------------------------
	// post[i][j] contains the product of every message that node i received
	// from its j+1 children to the last children
	// prev[i][j] contains the product of node i initial potential and
	// every message that node i received from its fist children to the j-1 children
	// So the message to be sent from i to j will be the product of prev and post
	// -------------------------------------------------------------------------
	c.prev = make([][]*factor.Factor, c.Size())
	c.post = make([][]*factor.Factor, c.Size())

	c.calibratedPot = make([]*factor.Factor, c.Size())
	c.calibratedPotSepSet = make([]*factor.Factor, c.Size())
	root := 0

	c.upwardmessage(root, -1)
	c.downwardmessage(-1, root)
}

func (c *CliqueTree) upwardmessage(v, pa int) {
	c.prev[v] = make([]*factor.Factor, 1, len(c.neighbours[v])+1)
	c.prev[v][0] = c.initialPot[v]
	for _, ne := range c.neighbours[v] {
		if ne != pa {
			c.upwardmessage(ne, v)
			c.prev[v] = append(c.prev[v], c.send[ne].Product(c.prev[v][len(c.prev[v])-1]))
		}
	}
	if pa != -1 {
		c.send[v] = c.prev[v][len(c.prev[v])-1].SumOut(c.varin[v])
	}
}

func (c *CliqueTree) downwardmessage(pa, v int) {
	c.calibratedPot[v] = c.prev[v][len(c.prev[v])-1]
	n := len(c.neighbours[v])
	if pa != -1 {
		c.calibratedPot[v] = c.calibratedPot[v].Product(c.receive[v])
		n--
		// calculate calibrated sepset
		c.calibratedPotSepSet[v] = c.calibratedPot[v].SumOut(c.varin[v])
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
		c.receive[ch] = msg.SumOut(c.varout[ch])
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
	varout := make([][]int, n)
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
			} else {
				varout[v] = append(varout[v], iphi[K[T.P[v]][i]-1])
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
	c.varin = make([][]int, n)
	c.varout = varout
	// initialize root clique
	sort.Ints(cliques[0])
	c.SetClique(0, cliques[0])
	c.SetNeighbours(0, children[0])

	for i := 1; i < c.Size(); i++ {
		// set cliques and sepset
		c.varin[i] = append(c.varin[i], cliques[i][0])
		c.SetSepSet(i, append([]int(nil), cliques[i][1:]...))
		sort.Ints(c.SepSet(i))
		sort.Ints(cliques[i])
		c.SetClique(i, cliques[i])
		// set adjacency list as children plus parent
		children[i] = append(children[i], T.P[i])
		c.SetNeighbours(i, children[i])

	}
	// set parents slice
	c.parent = T.P

	return c
}

// Marginals return a map with all marginals
func (c CliqueTree) Marginals() map[int][]float64 {
	c.UpDownCalibration()
	m := make(map[int][]float64)
	for i := 0; i < c.Size(); i++ {
		for j, v := range c.Calibrated(i).Variables() {
			if _, ok := m[v]; !ok {
				f := c.Calibrated(i).SumOut(c.Calibrated(i).Variables()[:j])
				f = f.SumOut(c.Calibrated(i).Variables()[j+1:])
				m[v] = f.Values()
			}
		}
	}
	return m
}

// SaveOnLibdaiFormat saves a clique tree in libDAI factor graph format on the given writer
func (c *CliqueTree) SaveOnLibdaiFormat(f io.Writer) {
	// number of potentials
	fmt.Fprintf(f, "%d\n", c.Size())
	fmt.Fprintln(f)
	for i := 0; i < c.Size(); i++ {
		// number of variables
		fmt.Fprintf(f, "%d\n", len(c.InitialPotential(i).Variables()))
		// variables
		for _, v := range c.InitialPotential(i).Variables() {
			fmt.Fprintf(f, "%d ", v)
		}
		fmt.Fprintln(f)
		// cardinalities
		for _, v := range c.InitialPotential(i).Variables() {
			fmt.Fprintf(f, "%d ", c.InitialPotential(i).Cardinality()[v])
		}
		fmt.Fprintln(f)
		// number of factor values
		fmt.Fprintf(f, "%d\n", len(c.InitialPotential(i).Values()))
		// factor values
		for j, v := range c.InitialPotential(i).Values() {
			fmt.Fprintf(f, "%d     %.4f\n", j, v)
		}
		fmt.Fprintln(f)
	}
}

// SaveOn saves a clique on the given writer
func (c *CliqueTree) SaveOn(w io.Writer) {
	// number of cliques
	fmt.Fprintf(w, "%d\n", c.Size())
	// cliques
	for i := 0; i < c.Size(); i++ {
		for _, v := range c.Clique(i) {
			fmt.Fprintf(w, "%d ", v)
		}
		fmt.Fprintln(w)
	}
	fmt.Fprintln(w)
	// adjacency
	for i := 0; i < c.Size(); i++ {
		for _, v := range c.Neighbours(i) {
			fmt.Fprintf(w, "%d ", v)
		}
		fmt.Fprintln(w)
	}
	fmt.Fprintln(w)
	// cardinality of all variables
	cardin := c.InitialPotential(0).Cardinality()
	for _, v := range cardin {
		fmt.Fprintf(w, "%d ", v)
	}
	fmt.Fprintln(w)
	// factor values
	for i := 0; i < c.Size(); i++ {
		for _, v := range c.InitialPotential(i).Values() {
			fmt.Fprintf(w, "%.8f ", v)
		}
		fmt.Fprintln(w)
	}
	fmt.Fprintln(w)
}

// LoadFrom loads a clique tree from the given reader
func LoadFrom(r io.Reader) *CliqueTree {
	scanner := bufio.NewScanner(r)
	scanner.Scan()
	size := utils.Atoi(scanner.Text())
	cliques := make([][]int, size)
	for i := 0; i < size; i++ {
		scanner.Scan()
		cliques[i] = append(cliques[i], utils.SliceAtoi(strings.Fields(scanner.Text()))...)
	}
	scanner.Scan()
	adj := make([][]int, size)
	for i := 0; i < size; i++ {
		scanner.Scan()
		adj[i] = append(adj[i], utils.SliceAtoi(strings.Fields(scanner.Text()))...)
	}
	scanner.Scan()
	scanner.Scan()
	cardin := utils.SliceAtoi(strings.Fields(scanner.Text()))
	potentials := make([]*factor.Factor, size)
	for i := range potentials {
		scanner.Scan()
		values := utils.SliceAtoF64(strings.Fields(scanner.Text()))
		potentials[i] = factor.NewFactorValues(cliques[i], cardin, values)
	}

	c, err := NewStructure(cliques, adj)
	utils.ErrCheck(err, "")
	err = c.SetAllPotentials(potentials)
	utils.ErrCheck(err, "")

	return c
}
