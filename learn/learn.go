package learn

import (
	"fmt"
	"math"
	"reflect"

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

// GuessStructure tries a number of random structures and choses the best one and its log-likelihood
func (l *Learner) GuessStructure() (*cliquetree.CliqueTree, float64) {
	bestStruct, bestScore := l.randomStruct()
	for i := 1; i < l.iterations; i++ {
		currStruct, currScore := l.randomStruct()
		if currScore > bestScore {
			bestScore = currScore
			bestStruct = currStruct
		}
	}
	return bestStruct, bestScore
}

// creates a new cliquetree from a randomized chartree and calculates its log-likelihood
func (l *Learner) randomStruct() (*cliquetree.CliqueTree, float64) {
	T, iphi, err := generator.RandomCharTree(l.n+l.hidden, l.treewidth)
	utils.ErrCheck(err, "")
	ct := cliquetree.FromCharTree(T, iphi)
	score := l.loglikelihood(ct)
	return ct, score
}

// loglikelihood calculates the log-likelihood of a list of cliques
func (l *Learner) loglikelihood(ct *cliquetree.CliqueTree) (ll float64) {
	// for each node adds the count of every attribution of the clique and
	// subtracts the count of every attribution of the sepset
	for i := 0; i < ct.Size(); i++ {
		values := l.counter.GetOccurrences(ct.Clique(i))
		for _, v := range values {
			if v != 0 {
				ll += float64(v) * math.Log(float64(v))
			}
		}
		values = l.counter.GetOccurrences(ct.SepSet(i))
		for _, v := range values {
			if v != 0 {
				ll -= float64(v) * math.Log(float64(v))
			}
		}
	}
	ll -= float64(l.dataset.Size()) * math.Log(float64(l.dataset.Size()))
	return
}

// CreateUniformPortentials creates a list of clique tree potentials with uniform values for the hidden variables
func (l *Learner) CreateUniformPortentials(ct *cliquetree.CliqueTree, cardin []int) []*factor.Factor {
	factors := make([]*factor.Factor, ct.Size())
	for i := range factors {
		values := utils.SliceItoF64(l.counter.GetOccurrences(ct.Clique(i)))
		var observed, hidden []int
		if l.hidden > 0 {
			// TODO: change this to avoid this new allocations
			observed, hidden = utils.SliceSplit(ct.Clique(i), l.n)
		}
		if len(hidden) > 0 {
			g := factor.NewFactor(hidden, cardin)
			g.SetUniform()
			factors[i] = factor.New(observed, cardin, values).Product(g)
		} else {
			factors[i] = factor.New(ct.Clique(i), cardin, values)
		}
	}
	return factors
}

// CreateRandomPortentials creates a list of clique tree potentials with random values
func (l *Learner) CreateRandomPortentials(ct *cliquetree.CliqueTree, cardin []int) []*factor.Factor {
	factors := make([]*factor.Factor, ct.Size())
	for i := range factors {
		factors[i] = factor.NewFactor(ct.Clique(i), cardin)
		factors[i].SetRandom()
	}
	return factors
}

// OptimizeParameters optimize the clique tree parameters
func (l *Learner) OptimizeParameters(ct *cliquetree.CliqueTree) {
	// extend cardinality to hidden variables
	cardin := make([]int, l.n+l.hidden)
	copy(cardin, l.dataset.Cardinality())
	for i := l.n; i < len(cardin); i++ {
		cardin[i] = l.hiddencard
	}

	// initialize clique tree potentials
	// ct.SetAllPotentials(l.CreateUniformPortentials(ct, cardin))
	ct.SetAllPotentials(l.CreateRandomPortentials(ct, cardin))

	//TODO: remove
	// count := make([]*factor.Factor, ct.Size())
	// for i := range count {
	// 	count[i] = ct.GetPotential(i).Clone()
	// }
	count := l.CreateUniformPortentials(ct, cardin)
	// fmt.Println("Initial clique tree")
	// for i := 0; i < ct.Size(); i++ {
	// 	fmt.Printf("%v\n", ct.GetInitPotential(i))
	// }
	// fmt.Println("==========================================")

	// call EM until convergence
	em.ExpectationMaximization(ct, l.dataset)

	//TODO: remove
	// fmt.Println("==========================================")
	// fmt.Println("Clique tree post EM")
	// for i := 0; i < ct.Size(); i++ {
	// 	fmt.Printf("%v\n", ct.GetInitPotential(i))
	// }
	for i := range count {
		if !reflect.DeepEqual(count[i], ct.GetPotential(i)) {
			fmt.Printf("diff >>>>>>>>>>>>>>>>>>>:\n%v\n%v\n", count[i], ct.GetPotential(i))
			break
		}
	}
}

// TODO: remove bellow

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
