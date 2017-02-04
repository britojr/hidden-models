package learn

import (
	"math"

	"github.com/britojr/playgo/counting/bitcounter"
	"github.com/britojr/playgo/filehandler"
	"github.com/britojr/playgo/junctree"
	"github.com/britojr/playgo/utils"
	"github.com/britojr/tcc/generator"
)

// Learner ..
type Learner struct {
	//parameters
	iterations int
	treewidth  int
	n          int
	dataset    *filehandler.DataSet
	counter    *bitcounter.BitCounter
}

// New ..
func New() *Learner {
	l := new(Learner)
	l.iterations = 100
	l.treewidth = 3
	return l
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

// SetTreeWidth ..
func (l *Learner) SetTreeWidth(k int) {
	l.treewidth = k
}

// SetIterations ..
func (l *Learner) SetIterations(n int) {
	l.iterations = n
}

// calcLL calculates the loglikelihood of a junctree for the data
func (l *Learner) calcLL(jt *junctree.JuncTree) (ll float64) {
	// for each node adds the count of every attribution of the clique and
	// subtracts the count of every attribution of the separator
	for _, node := range jt.Nodes {
		r := l.counter.Marginalize(node.Cliq...)
		next := r.ValueIteratorNonZero()
		v := next()
		for v != nil {
			ll += float64(*v) * math.Log(float64(*v))
			v = next()
		}
		next = r.SumOut(node.Cliq[0]).ValueIteratorNonZero()
		v = next()
		for v != nil {
			ll -= float64(*v) * math.Log(float64(*v))
			v = next()
		}
	}
	ll -= float64(l.dataset.Size()) * math.Log(float64(l.dataset.Size()))
	return
}

func (l *Learner) newRandomStruct() (*junctree.JuncTree, float64) {
	T, iphi, err := generator.RandomCharTree(l.n, l.treewidth)
	utils.ErrCheck(err, "")
	jt := junctree.FromCharTree(T, iphi)
	score := l.calcLL(jt)
	return jt, score
}
