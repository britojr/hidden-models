package learn

import (
	"fmt"
	"math"

	"github.com/britojr/kbn/cliquetree"
	"github.com/britojr/kbn/counting/bitcounter"
	"github.com/britojr/kbn/em"
	"github.com/britojr/kbn/factor"
	"github.com/britojr/kbn/filehandler"
	"github.com/britojr/kbn/junctree"
	"github.com/britojr/kbn/utils"
	"github.com/britojr/tcc/generator"
)

// Learner ..
type Learner struct {
	//parameters
	iterations int
	treewidth  int
	n          int // number of variables
	dataset    *filehandler.DataSet
	counter    *bitcounter.BitCounter
	hidden     int // number of hidden variables
	hiddencard int // default cardinality of the hidden variables
}

// New ..
func New() *Learner {
	l := new(Learner)
	l.iterations = 100
	l.treewidth = 3
	l.hiddencard = 2
	return l
}

// SetTreeWidth ..
func (l *Learner) SetTreeWidth(k int) {
	l.treewidth = k
}

// SetIterations ..
func (l *Learner) SetIterations(it int) {
	l.iterations = it
}

// SetHiddenVars ..
func (l *Learner) SetHiddenVars(h int) {
	l.hidden = h
}

// LoadDataSet ..
func (l *Learner) LoadDataSet(dsfile string, delimiter rune, dsHdrlns filehandler.HeaderFlags) {
	l.dataset = filehandler.NewDataSet(dsfile, delimiter, dsHdrlns)
	l.dataset.Read()
	l.counter = bitcounter.NewBitCounter()
	l.counter.LoadFromData(l.dataset.Data(), l.dataset.Cardinality())
	l.n = len(l.dataset.Cardinality())
}

// BestJuncTree ..
func (l *Learner) BestJuncTree() (*junctree.JuncTree, float64) {
	bestStruct, bestScore := l.newRandomStruct()
	for i := 1; i < l.iterations; i++ {
		currStruct, currScore := l.newRandomStruct()
		if currScore > bestScore {
			bestScore = currScore
			bestStruct = currStruct
		}
	}
	return bestStruct, bestScore
}

// calcLL calculates the loglikelihood of a list of cliques
func (l *Learner) calcLL(nodelist []junctree.Node) (ll float64) {
	// for each node adds the count of every attribution of the clique and
	// subtracts the count of every attribution of the Sepset
	for _, node := range nodelist {
		values := l.counter.GetOccurrences(node.Clique)
		for _, v := range values {
			if v != 0 {
				ll += float64(v) * math.Log(float64(v))
			}
		}
		values = l.counter.GetOccurrences(node.Sepset)
		for _, v := range values {
			if v != 0 {
				ll -= float64(v) * math.Log(float64(v))
			}
		}
	}
	ll -= float64(l.dataset.Size()) * math.Log(float64(l.dataset.Size()))
	return
}

func (l *Learner) newRandomStruct() (*junctree.JuncTree, float64) {
	T, iphi, err := generator.RandomCharTree(l.n+l.hidden, l.treewidth)
	utils.ErrCheck(err, "")
	jt := junctree.FromCharTree(T, iphi)
	score := l.calcLL(jt.Nodes)
	return jt, score
}

// OptimizeParameters ..
func (l *Learner) OptimizeParameters(jt *junctree.JuncTree) *cliquetree.CliqueTree {
	// extend cardinality to hidden variables
	cardin := make([]int, l.n+l.hidden)
	copy(cardin, l.dataset.Cardinality())
	for i := l.n; i < len(cardin); i++ {
		cardin[i] = l.hiddencard
	}

	// initialize clique tree TODO: fix redundant code merge junctree on clique tree
	ct := cliquetree.New(len(jt.Nodes))
	for i, n := range jt.Nodes {
		ct.SetClique(i, n.Clique)
		ct.SetNeighbours(i, jt.Adj[i])
	}
	temp := l.counter.GetOccurrences(ct.Clique(0))
	tempf := utils.SliceItoF64(temp)
	fmt.Printf("occur for clique(%v): %v\n", ct.Clique(0), temp)
	fmt.Printf("float: %v\n", tempf)
	var tot float64
	for _, v := range tempf {
		tot += v
	}
	for i := range tempf {
		tempf[i] /= tot
	}
	fmt.Printf("norm: %v\n", tempf)

	// initialize clique tree potentials
	for i := 0; i < ct.Size(); i++ {
		values := utils.SliceItoF64(l.counter.GetOccurrences(ct.Clique(i)))
		if l.hidden > 0 {
			// TODO: change this to avoid this new allocations
			observed, hidden := utils.SliceSplit(ct.Clique(i), l.n)
			f := factor.New(observed, cardin, values)
			g := factor.NewUniform(hidden, cardin)
			ct.SetPotential(i, f.Product(g))
		} else {
			ct.SetPotential(i, factor.New(ct.Clique(i), cardin, values))
		}
	}
	// call EM until convergence
	em.ExpectationMaximization(ct, l.dataset)
	// return learned structure
	return ct
}
