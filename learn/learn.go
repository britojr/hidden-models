package learn

import (
	"math"

	"github.com/britojr/kbn/counting/bitcounter"
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
}

// New ..
func New() *Learner {
	l := new(Learner)
	l.iterations = 100
	l.treewidth = 3
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
	// subtracts the count of every attribution of the separator
	for _, node := range nodelist {
		values := l.counter.GetOccurrences(node.Clique)
		for _, v := range values {
			if v != 0 {
				ll += float64(v) * math.Log(float64(v))
			}
		}
		values = l.counter.GetOccurrences(node.Separator)
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
	T, iphi, err := generator.RandomCharTree(l.n, l.treewidth)
	utils.ErrCheck(err, "")
	jt := junctree.FromCharTree(T, iphi)
	score := l.calcLL(jt.Nodes)
	return jt, score
}
