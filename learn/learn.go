package learn

import (
	"fmt"
	"os"
	"sort"

	"github.com/britojr/kbn/cliquetree"
	"github.com/britojr/kbn/counting/bitcounter"
	"github.com/britojr/kbn/em"
	"github.com/britojr/kbn/factor"
	"github.com/britojr/kbn/filehandler"
	"github.com/britojr/kbn/likelihood"
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
	initpot    int
	epslon     float64 // epslon for EM convergence
}

// New ..
func New() *Learner {
	l := new(Learner)
	l.treewidth = 3
	l.hiddencard = 2
	l.initpot = 1
	l.epslon = 1e-6
	return l
}

// SetTreeWidth ..
func (l *Learner) SetTreeWidth(k int) {
	l.treewidth = k
}

// SetHiddenVars ..
func (l *Learner) SetHiddenVars(h int) {
	l.hidden = h
}

// SetInitPot ..
func (l *Learner) SetInitPot(initpot int) {
	l.initpot = initpot
}

// SetEpslon ..
func (l *Learner) SetEpslon(epslon float64) {
	if epslon != 0 {
		l.epslon = epslon
	}
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
	fmt.Printf("Tot var: %v\n", len(l.cardin))
	fmt.Printf("Variables: %v, Instances: %v\n", l.n, len(l.dataset.Data()))
}

// GuessStructure tries a number of random structures and choses the best one and its log-likelihood
func (l *Learner) GuessStructure(iterations int) (*cliquetree.CliqueTree, float64) {
	bestStruct, bestScore := l.randomStruct()
	for i := 1; i < iterations; i++ {
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
	score := likelihood.StructLog(ct.Cliques(), ct.SepSets(), l.counter)
	return ct, score
}

// InitializePotentials initialize clique tree potentials
func (l *Learner) InitializePotentials(ct *cliquetree.CliqueTree, initpot ...int) {
	aux := l.initpot
	if len(initpot) > 0 {
		aux = initpot[0]
	}
	switch aux {
	case 2:
		ct.SetAllPotentials(CreateEmpiricPotentials(ct.Cliques(), l.cardin, l.n, l.counter, true))
	case 3:
		ct.SetAllPotentials(CreateEmpiricPotentials(ct.Cliques(), l.cardin, l.n, l.counter, false))
	default:
		// case 1:
		ct.SetAllPotentials(CreateRandomPotentials(ct.Cliques(), l.cardin))
	}
}

// OptimizeParameters optimize the clique tree parameters
func (l *Learner) OptimizeParameters(ct *cliquetree.CliqueTree, initpot, iterations int) {
	l.InitializePotentials(ct, initpot)
	em.ExpectationMaximization(ct, l.dataset, l.epslon)
	if iterations > 1 {
		ll := l.CalculateLikelihood(ct)
		fmt.Printf("curr LL %v\n", ll)
		pot := make([]*factor.Factor, len(ct.Potentials()))
		copy(pot, ct.Potentials())
		for i := 1; i < iterations; i++ {
			l.InitializePotentials(ct, initpot)
			em.ExpectationMaximization(ct, l.dataset, l.epslon)
			currll := l.CalculateLikelihood(ct)
			fmt.Printf("curr LL %v\n", currll)
			if currll > ll {
				ll = currll
				copy(pot, ct.Potentials())
			}
		}
		ct.SetAllPotentials(pot)
	}
}

// CalculateLikelihood calculates the likelihood of a clique tree
func (l *Learner) CalculateLikelihood(ct *cliquetree.CliqueTree) float64 {
	ct.UpDownCalibration()
	return likelihood.Loglikelihood2(ct, l.dataset, l.n)
}

// CreateEmpiricPotentials creates a list of clique tree potentials with counting
// for observed variables (empiric distribution), and expand uniformily or randomly for the hidden variables
func CreateEmpiricPotentials(cliques [][]int, cardin []int,
	numobs int, counter utils.Counter, uniform bool) []*factor.Factor {

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
				if uniform {
					g.SetUniform()
				} else {
					g.SetRandom()
				}
				factors[i] = factors[i].Product(g)
			}
			factors[i].Normalize()
		} else {
			factors[i] = factor.NewFactor(hidden, cardin)
			factors[i].SetUniform()
		}
	}
	return factors
}

// CreateRandomPotentials creates a list of clique potentials with random values
func CreateRandomPotentials(cliques [][]int, cardin []int) []*factor.Factor {
	factors := make([]*factor.Factor, len(cliques))
	for i := range factors {
		factors[i] = factor.NewFactor(cliques[i], cardin).SetRandom()
	}
	return factors
}

// CheckTree ..
func (l *Learner) CheckTree(ct *cliquetree.CliqueTree) {
	ct.UpDownCalibration()
	// check if they are uniform
	l.checkUniform(ct)
	// check if after summing out the hidden variables they are the same as initial count
	l.checkWithInitialCount(ct)

	checkCliqueTree(ct)
}

func (l *Learner) checkUniform(ct *cliquetree.CliqueTree) {
	fmt.Println("checkUniform")
	uniform := CreateEmpiricPotentials(ct.Cliques(), l.cardin, l.n, l.counter, true)
	fmt.Printf("Uniform param: %v (%v)=0\n", uniform[0].Values()[0], uniform[0].Variables())
	calibrated := make([]*factor.Factor, ct.Size())
	for i := range calibrated {
		calibrated[i] = ct.Calibrated(i)
	}
	diff, i, j, err := factor.MaxDifference(uniform, calibrated)
	utils.ErrCheck(err, "")
	fmt.Printf("f[%v][%v]=%v; g[%v][%v]=%v\n", i, j, uniform[i].Values()[j], i, j, ct.Calibrated(i).Values()[j])
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
			// sumOutHidden[i] = ct.InitialPotential(i)
			sumOutHidden[i] = ct.Calibrated(i)
			if len(hidden) > 0 {
				sumOutHidden[i] = sumOutHidden[i].SumOut(hidden)
			}
			initialCount[i] = factor.NewFactorValues(observed, l.cardin, values)
			initialCount[i].Normalize()
			sumOutHidden[i].Normalize()
		}
	}

	if initialCount[0] != nil {
		fmt.Printf("IniCount param: %v (%v)=0\n", initialCount[0].Values()[0], initialCount[0].Variables())
		fmt.Printf("sumOut param: %v (%v)=0\n", sumOutHidden[0].Values()[0], sumOutHidden[0].Variables())
	}
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

// checkCliqueTree ..
func checkCliqueTree(ct *cliquetree.CliqueTree) {
	printTree := func(f *factor.Factor) {
		fmt.Printf("(%v)\n", f.Variables())
		fmt.Println("tree:")
		for i := 0; i < ct.Size(); i++ {
			fmt.Printf("node %v: neighb: %v clique: %v septset: %v parent: %v\n",
				i, ct.Neighbours(i), ct.Clique(i), ct.SepSet(i), ct.Parents()[i])
		}
		fmt.Println("original potentials:")
		for i := 0; i < ct.Size(); i++ {
			fmt.Printf("node %v:\n var: %v\n values: %v\n",
				i, ct.InitialPotential(i).Variables(), ct.InitialPotential(i).Values())
		}
	}

	for i := range ct.Potentials() {
		f := ct.InitialPotential(i)
		sum := 0.0
		for _, v := range f.Values() {
			sum += v
		}
		if utils.FuzzyEqual(sum, 0) {
			printTree(f)
			panic("original zero factor")
		}
		f = ct.Calibrated(i)
		sum = 0.0
		for _, v := range f.Values() {
			sum += v
		}
		if utils.FuzzyEqual(sum, 0) {
			printTree(f)
			panic("original zero factor")
		}
	}
}

// SaveCliqueTree saves a clique tree in libDAI factor graph format
func SaveCliqueTree(ct *cliquetree.CliqueTree, fname string) {
	f, err := os.Create(fname)
	utils.ErrCheck(err, "")
	defer f.Close()
	// number of potentials
	fmt.Fprintf(f, "%d\n", ct.Size())
	fmt.Fprintln(f)
	for i := 0; i < ct.Size(); i++ {
		fmt.Fprintf(f, "%d\n", len(ct.Calibrated(i).Variables()))
		for _, v := range ct.Calibrated(i).Variables() {
			fmt.Fprintf(f, "%d ", v)
		}
		fmt.Fprintln(f)

		for _, v := range ct.Calibrated(i).Variables() {
			fmt.Fprintf(f, "%d ", ct.Calibrated(i).Cardinality()[v])
		}
		fmt.Fprintln(f)

		fmt.Fprintf(f, "%d\n", len(ct.Calibrated(i).Values()))
		for j, v := range ct.InitialPotential(i).Values() {
			fmt.Fprintf(f, "%d     %.4f\n", j, v)
		}
		fmt.Fprintln(f)
	}
}

// SaveMarginals saves all marginals of a cliquetree
func SaveMarginals(ct *cliquetree.CliqueTree, ll float64, fname string) {
	f, err := os.Create(fname)
	utils.ErrCheck(err, "")
	defer f.Close()
	m := ct.Marginals()

	var keys []int
	for k := range m {
		keys = append(keys, k)
	}
	sort.Ints(keys)
	for _, k := range keys {
		fmt.Fprintf(f, "{%d} %v\n", k, m[k])
	}
	fmt.Fprintf(f, "LL=%v\n", ll)
}
