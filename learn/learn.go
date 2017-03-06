package learn

import (
	"fmt"
	"math"

	"github.com/britojr/kbn/assignment"
	"github.com/britojr/kbn/cliquetree"
	"github.com/britojr/kbn/counting/bitcounter"
	"github.com/britojr/kbn/em"
	"github.com/britojr/kbn/factor"
	"github.com/britojr/kbn/filehandler"
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
	hidden     int   // number of hidden variables
	hiddencard int   // default cardinality of the hidden variables
	cardin     []int // cardinality slice
	norm       bool
	initpot    int
}

// New ..
func New() *Learner {
	l := new(Learner)
	l.iterations = 100
	l.treewidth = 3
	l.hiddencard = 2
	l.norm = true
	l.initpot = 1
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

// SetNorm ..
func (l *Learner) SetNorm(norm bool) {
	l.norm = norm
}

// SetInitPot ..
func (l *Learner) SetInitPot(initpot int) {
	l.initpot = initpot
}

// LoadDataSet ..
func (l *Learner) LoadDataSet(dsfile string, delimiter rune, dsHdrlns filehandler.HeaderFlags) {
	l.dataset = filehandler.NewDataSet(dsfile, delimiter, dsHdrlns)
	l.dataset.Read()
	l.counter = bitcounter.NewBitCounter()
	l.counter.LoadFromData(l.dataset.Data(), l.dataset.Cardinality())
	l.n = len(l.dataset.Cardinality())
	// extend cardinality to hidden variables
	l.cardin = make([]int, l.n+l.hidden)
	copy(l.cardin, l.dataset.Cardinality())
	for i := l.n; i < len(l.cardin); i++ {
		l.cardin[i] = l.hiddencard
	}
	fmt.Printf("Variables: %v, Instances: %v\n", l.n, len(l.dataset.Data()))
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
	score := Structloglikelihood(ct.Cliques(), ct.SepSets(), l.counter)
	return ct, score
}

// OptimizeParameters optimize the clique tree parameters
func (l *Learner) OptimizeParameters(ct *cliquetree.CliqueTree) {
	// initialize clique tree potentials
	if l.initpot == 1 {
		ct.SetAllPotentials(CreateRandomPortentials(ct.Cliques(), l.cardin))
	} else {
		ct.SetAllPotentials(CreateUniformPortentials(ct.Cliques(), l.cardin, l.n, l.counter))
	}

	fmt.Printf("Initial param: %v (%v)=0\n", ct.BkpPotential(0).Values()[0], ct.BkpPotential(0).Variables())
	// call EM until convergence
	em.ExpectationMaximization(ct, l.dataset, l.norm)

	// check resulting parameters TODO: remove
	// check if they are uniform
	l.checkUniform(ct)
	// check if after summing out the hidden variables they are the same as initial count
	l.checkWithInitialCount(ct)
}

func (l *Learner) checkUniform(ct *cliquetree.CliqueTree) {
	fmt.Println("checkUniform")
	uniform := CreateUniformPortentials(ct.Cliques(), l.cardin, l.n, l.counter)
	fmt.Printf("Uniform param: %v (%v)=0\n", uniform[0].Values()[0], uniform[0].Variables())
	diff, i, j, err := factor.MaxDifference(uniform, ct.BkpPotentialList())
	utils.ErrCheck(err, "")
	fmt.Printf("f[%v][%v]=%v; g[%v][%v]=%v\n", i, j, uniform[i].Values()[j], i, j, ct.BkpPotential(i).Values()[j])
	if diff > 0 {
		fmt.Printf(" > Not uniform: maxdiff = %v\n", diff)
		if diff > 1e-6 {
			fmt.Println(" !! Significant difference!")
		}
	} else {
		fmt.Printf(" > Is uniform: maxdiff = %v\n", diff)
	}
}

func (l *Learner) checkWithInitialCount(ct *cliquetree.CliqueTree) {
	fmt.Println("checkWithInitialCount")
	initialCount := make([]*factor.Factor, ct.Size())
	sumOutHidden := make([]*factor.Factor, ct.Size())
	for i := range initialCount {
		var observed, hidden []int
		if l.hidden > 0 {
			observed, hidden = utils.SliceSplit(ct.Clique(i), l.n)
		} else {
			observed = ct.Clique(i)
		}
		if len(observed) > 0 {
			values := utils.SliceItoF64(l.counter.CountAssignments(observed))
			sumOutHidden[i] = ct.BkpPotential(i)
			if len(hidden) > 0 {
				sumOutHidden[i] = sumOutHidden[i].SumOut(hidden)
			}
			initialCount[i] = factor.NewFactorValues(observed, l.cardin, values)
			initialCount[i].Normalize()
			sumOutHidden[i].Normalize()
		}
	}

	fmt.Printf("IniCount param: %v (%v)=0\n", initialCount[0].Values()[0], initialCount[0].Variables())
	fmt.Printf("sumOut param: %v (%v)=0\n", sumOutHidden[0].Values()[0], sumOutHidden[0].Variables())
	diff, i, j, err := factor.MaxDifference(initialCount, sumOutHidden)
	utils.ErrCheck(err, "")
	fmt.Printf("f[%v][%v]=%v; g[%v][%v]=%v\n", i, j, initialCount[i].Values()[j], i, j, sumOutHidden[i].Values()[j])
	if diff > 0 {
		fmt.Printf(" > Different from initial counting: maxdiff = %v\n", diff)
		if diff > 1e-6 {
			fmt.Println(" >> Significant difference!")
		}
	} else {
		fmt.Printf(" > Exactly the initial counting: maxdiff = %v\n", diff)
	}
}

// Counter returns a counting value for an assignment
type Counter interface {
	Count(assig *assignment.Assignment) (count int, ok bool)
	CountAssignments(varlist []int) []int
	Cardinality() []int
	NumTuples() int
}

// Structloglikelihood calculates the log-likelihood
func Structloglikelihood(cliques, sepsets [][]int, counter Counter) (ll float64) {
	// for each node adds the count of every attribution of the clique and
	// subtracts the count of every attribution of the sepset
	for i := range cliques {
		ll += sumLogCount(cliques[i], counter)
	}
	for i := range sepsets {
		ll -= sumLogCount(sepsets[i], counter)
	}
	ll -= float64(counter.NumTuples()) * math.Log(float64(counter.NumTuples()))
	return
}

func sumLogCount(varlist []int, counter Counter) (ll float64) {
	values := counter.CountAssignments(varlist)
	for _, v := range values {
		if v != 0 {
			ll += float64(v) * math.Log(float64(v))
		}
	}
	return ll
}

// CreateUniformPortentials creates a list of clique tree potentials with uniform values for the hidden variables
func CreateUniformPortentials(cliques [][]int, cardin []int,
	numobs int, counter Counter) []*factor.Factor {

	factors := make([]*factor.Factor, len(cliques))
	for i := range factors {
		var observed, hidden []int
		if len(cardin) > numobs {
			observed, hidden = utils.SliceSplit(cliques[i], numobs)
		} else {
			observed = cliques[i]
		}
		if len(observed) > 0 {
			values := utils.SliceItoF64(counter.CountAssignments(observed))
			factors[i] = factor.NewFactorValues(observed, cardin, values)
			if len(hidden) > 0 {
				g := factor.NewFactor(hidden, cardin)
				g.SetUniform()
				factors[i] = factors[i].Product(g)
			}
			factors[i].Normalize()
		} else {
			factors[i] = factor.NewFactor(hidden, cardin)
			factors[i].SetUniform()
			fmt.Printf(" >> Factor %v: has only hidden variables\n", i)
		}
	}
	return factors
}

// CreateRandomPortentials creates a list of clique potentials with random values
func CreateRandomPortentials(cliques [][]int, cardin []int) []*factor.Factor {
	factors := make([]*factor.Factor, len(cliques))
	for i := range factors {
		factors[i] = factor.NewFactor(cliques[i], cardin).SetRandom()
	}
	return factors
}
