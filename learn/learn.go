package learn

import (
	"fmt"
	"math"
	"os"

	"github.com/britojr/kbn/cliquetree"
	"github.com/britojr/kbn/counting"
	"github.com/britojr/kbn/counting/bitcounter"
	"github.com/britojr/kbn/em"
	"github.com/britojr/kbn/factor"
	"github.com/britojr/kbn/filehandler"
	"github.com/britojr/kbn/likelihood"
	"github.com/britojr/kbn/utils"
	"github.com/britojr/tcc/generator"
)

const (
	// FullRandom initialize potentials with random values
	FullRandom = iota
	// EmpiricDirichlet initialize potentials with empiric distribution extended by dirichlet
	EmpiricDirichlet
	// EmpiricRandom initialize potentials with empiric distribution extended by random
	EmpiricRandom
	// EmpiricUniform initialize potentials with empiric distribution extended by uniform
	EmpiricUniform
)

// Learner ..
type Learner struct {
	dataset *filehandler.DataSet
	counter *bitcounter.BitCounter
	n       int   // number of variables
	cardin  []int // cardinality slice
	// parameters
	k          int       // treewidth
	hidden     int       // number of hidden variables
	hiddencard int       // cardinality of the hidden variables
	alphas     []float64 // parameters for dirichlet distribution
}

// New creates new learner object with parameters
func New(k, hidden, hiddencard int, alpha ...float64) *Learner {
	l := new(Learner)
	l.k = k
	l.hidden = hidden
	l.hiddencard = hiddencard
	if len(alpha) > 0 {
		l.alphas = make([]float64, int(math.Pow(float64(hiddencard), float64(l.k+1))))
		for i := range l.alphas {
			l.alphas[i] = alpha[0]
		}
	}
	return l
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
	fmt.Printf("Variables: %v+%v, Instances: %v\n", l.n, l.hidden, len(l.dataset.Data()))
}

// Counter returns counter
func (l *Learner) Counter() counting.Counter {
	return l.counter
}

// Data returns dataset matrix
func (l *Learner) Data() [][]int {
	return l.dataset.Data()
}

// Cardinality returns cardinality slice
func (l *Learner) Cardinality() []int {
	return l.cardin
}

// TotVar returns total number of variables
func (l *Learner) TotVar() int {
	return l.n + l.hidden
}

// InitializePotentials initialize clique tree potentials
func (l *Learner) InitializePotentials(ct *cliquetree.CliqueTree, typePot int) {
	if typePot == FullRandom {
		ct.SetAllPotentials(CreateRandomPotentials(ct.Cliques(), l.cardin))
	} else {
		factors := CreateEmpiricPotentials(l.counter, ct.Cliques(), l.cardin, l.n, typePot, l.alphas...)
		for i := range factors {
			if len(ct.Varin(i)) != 0 {
				factors[i] = factors[i].Division(factors[i].SumOut(ct.Varin(i)))
			}
		}
		ct.SetAllPotentials(factors)
	}
}

// OptimizeParameters optimize the clique tree parameters
func (l *Learner) OptimizeParameters(ct *cliquetree.CliqueTree,
	typePot, iterations int, epslon float64) float64 {

	l.InitializePotentials(ct, typePot)
	em.ExpectationMaximization(ct, l.dataset.Data(), epslon)
	bestll := l.CalculateLikelihood(ct)
	if iterations > 1 {
		fmt.Printf("curr LL %v\n", bestll)
		pot := make([]*factor.Factor, len(ct.Potentials()))
		copy(pot, ct.Potentials())
		for i := 1; i < iterations; i++ {
			l.InitializePotentials(ct, typePot)
			em.ExpectationMaximization(ct, l.dataset.Data(), epslon)
			currll := l.CalculateLikelihood(ct)
			fmt.Printf("curr LL %v\n", currll)
			if currll > bestll {
				bestll = currll
				copy(pot, ct.Potentials())
			}
		}
		ct.SetAllPotentials(pot)
	}
	return bestll
}

// CalculateLikelihood calculates the likelihood of a clique tree
func (l *Learner) CalculateLikelihood(ct *cliquetree.CliqueTree) float64 {
	ct.UpDownCalibration()
	return likelihood.Loglikelihood(ct, l.dataset.Data())
}

// GuessStructure tries a number of random structures and choses the best one and its log-likelihood
func (l *Learner) GuessStructure(iterations int) (*cliquetree.CliqueTree, float64) {
	bestStruct := RandomCliqueTree(l.n+l.hidden, l.k)
	bestScore := likelihood.StructLoglikelihood(bestStruct.Cliques(), bestStruct.SepSets(), l.counter)
	for i := 1; i < iterations; i++ {
		currStruct := RandomCliqueTree(l.n+l.hidden, l.k)
		currScore := likelihood.StructLoglikelihood(currStruct.Cliques(), currStruct.SepSets(), l.counter)
		if currScore > bestScore {
			bestScore = currScore
			bestStruct = currStruct
		}
	}
	return bestStruct, bestScore
}

// RandomCliqueTree creates a new cliquetree from a randomized chartree
func RandomCliqueTree(n, k int) *cliquetree.CliqueTree {
	T, iphi, err := generator.RandomCharTree(n, k)
	utils.ErrCheck(err, "")
	ct := cliquetree.FromCharTree(T, iphi)
	return ct
}

// CreateEmpiricPotentials creates a list of clique tree potentials with counting
// for observed variables (empiric distribution), and expand uniformily or randomly for the hidden variables
func CreateEmpiricPotentials(counter counting.Counter, cliques [][]int, cardin []int,
	numobs, typePot int, alphas ...float64) []*factor.Factor {

	if typePot == EmpiricDirichlet && len(alphas) == 0 {
		panic("no parameters for dirichlet dirtributions")
	}
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
			// factors[i] = P(observed)
			factors[i] = factor.NewFactorValues(observed, cardin, values).Normalize()
			if len(hidden) > 0 {
				// g = P(hidden/observed)
				g := factor.NewFactor(cliques[i], cardin)
				lenobs := len(factors[i].Values())
				g.SetValues(proportionalValues(lenobs, len(g.Values())/lenobs, typePot, alphas))
				// P(observed, hidden) = P(observed) * P(hidden/observed)
				factors[i] = factors[i].Product(g)
			}
		} else {
			factors[i] = factor.NewFactor(cliques[i], cardin)
			factors[i].SetValues(proportionalValues(1, len(factors[i].Values()), typePot, alphas))
		}
	}
	return factors
}

func proportionalValues(lenobs, lenhidden, typePot int, alphas []float64) []float64 {
	values := make([]float64, lenobs*lenhidden)
	aux := make([]float64, lenhidden)
	for i := 0; i < lenobs; i++ {
		switch typePot {
		case EmpiricDirichlet:
			utils.Dirichlet(alphas[:lenhidden], aux)
		case EmpiricRandom:
			utils.Random(aux)
		case EmpiricUniform:
			utils.Uniform(aux)
		}
		for j := 0; j < lenhidden; j++ {
			values[i+(j*lenobs)] = aux[j]
		}
	}
	return values
}

// CreateRandomPotentials creates a list of clique potentials with random values
func CreateRandomPotentials(cliques [][]int, cardin []int) []*factor.Factor {
	factors := make([]*factor.Factor, len(cliques))
	for i := range factors {
		factors[i] = factor.NewFactor(cliques[i], cardin).SetRandom()
	}
	return factors
}

// SaveCliqueTree saves a clique tree on the given file
func SaveCliqueTree(ct *cliquetree.CliqueTree, fname string) {
	f, err := os.Create(fname)
	utils.ErrCheck(err, "")
	defer f.Close()
	ct.SaveOn(f)
}

// LoadCliqueTree loads a clique tree from the given file
func LoadCliqueTree(fname string) *cliquetree.CliqueTree {
	f, err := os.Open(fname)
	utils.ErrCheck(err, "")
	defer f.Close()
	return cliquetree.LoadFrom(f)
}
