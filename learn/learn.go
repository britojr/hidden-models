package learn

import (
	"fmt"
	"math"
	"os"

	"github.com/britojr/kbn/cliquetree"
	"github.com/britojr/kbn/conv"
	"github.com/britojr/kbn/counting"
	"github.com/britojr/kbn/counting/bitcounter"
	"github.com/britojr/kbn/em"
	"github.com/britojr/kbn/errchk"
	"github.com/britojr/kbn/factor"
	"github.com/britojr/kbn/likelihood"
	"github.com/britojr/kbn/list"
	"github.com/britojr/kbn/stats"
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
	data    [][]int
	counter *bitcounter.BitCounter
	n       int   // number of variables
	cardin  []int // cardinality slice
	// parameters
	k      int       // treewidth
	h      int       // number of hidden variables
	hcard  int       // cardinality of the hidden variables
	alphas []float64 // parameters for dirichlet distribution
}

// New creates new learner object with parameters
func New(data [][]int, cardin []int, k, h, hcard int, alpha ...float64) *Learner {
	l := new(Learner)
	l.k = k
	l.h = h
	l.hcard = hcard
	// create slice of alpha parameters
	if len(alpha) > 0 && alpha[0] > 0 {
		l.alphas = make([]float64, int(math.Pow(float64(hcard), float64(l.k+1))))
		for i := range l.alphas {
			l.alphas[i] = alpha[0]
		}
	}
	l.counter = bitcounter.NewBitCounter()
	l.counter.LoadFromData(data, cardin)
	l.n = len(cardin)
	// extend cardinality to hidden variables
	l.cardin = make([]int, l.n+l.h)
	copy(l.cardin, cardin)
	for i := l.n; i < len(l.cardin); i++ {
		l.cardin[i] = l.hcard
	}
	l.data = data
	return l
}

// Counter returns counter
func (l *Learner) Counter() counting.Counter {
	return l.counter
}

// Data returns dataset matrix
func (l *Learner) Data() [][]int {
	return l.data
}

// Cardinality returns cardinality slice
func (l *Learner) Cardinality() []int {
	return l.cardin
}

// TotVar returns total number of variables
func (l *Learner) TotVar() int {
	return l.n + l.h
}

// InitializePotentials initialize clique tree potentials
func (l *Learner) InitializePotentials(ct *cliquetree.CliqueTree, typePot, indePot int) {
	if typePot == FullRandom {
		ct.SetAllPotentials(CreateRandomPotentials(ct.Cliques(), l.cardin))
	} else {
		factors := CreateEmpiricPotentials(l.counter, ct.Cliques(), l.cardin, l.n, typePot, indePot, l.alphas...)
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
	typePot, indePot, iterations int, epslon float64) float64 {

	if iterations == 0 {
		l.InitializePotentials(ct, typePot, indePot)
		return l.CalculateLikelihood(ct)
	}
	// TODO: remove unnecessary LL calculations
	l.InitializePotentials(ct, typePot, indePot)
	fmt.Printf("LL before EM %v\n", l.CalculateLikelihood(ct))
	em.ExpectationMaximization(ct, l.data, epslon)
	bestll := l.CalculateLikelihood(ct)
	if iterations > 1 {
		fmt.Printf("curr LL %v\n", bestll)
		pot := make([]*factor.Factor, len(ct.Potentials()))
		copy(pot, ct.Potentials())
		for i := 1; i < iterations; i++ {
			l.InitializePotentials(ct, typePot, indePot)
			fmt.Printf("LL before EM %v\n", l.CalculateLikelihood(ct))
			em.ExpectationMaximization(ct, l.data, epslon)
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
	return likelihood.Loglikelihood(ct, l.data)
}

// GuessStructure tries a number of random structures and choses the best one and its log-likelihood
func (l *Learner) GuessStructure(iterations int) (*cliquetree.CliqueTree, float64) {
	bestStruct := cliquetree.NewRandom(l.n+l.h, l.k)
	bestScore := likelihood.StructLoglikelihood(bestStruct.Cliques(), bestStruct.SepSets(), l.counter)
	for i := 1; i < iterations; i++ {
		currStruct := cliquetree.NewRandom(l.n+l.h, l.k)
		currScore := likelihood.StructLoglikelihood(currStruct.Cliques(), currStruct.SepSets(), l.counter)
		if currScore > bestScore {
			bestScore = currScore
			bestStruct = currStruct
		}
	}
	return bestStruct, bestScore
}

// CreateEmpiricPotentials creates a list of clique tree potentials with counting
// for observed variables (empiric distribution), and expand uniformily or randomly for the hidden variables
func CreateEmpiricPotentials(counter counting.Counter, cliques [][]int, cardin []int,
	numobs, typePot, indePot int, alphas ...float64) []*factor.Factor {

	if typePot == EmpiricDirichlet && len(alphas) == 0 {
		panic("no parameters for dirichlet dirtributions")
	}
	factors := make([]*factor.Factor, len(cliques))
	for i := range factors {
		var observed, hidden []int
		if len(cardin) > numobs {
			observed, hidden = list.Split(cliques[i], numobs)
		} else {
			observed = cliques[i]
		}
		if len(observed) > 0 {
			values := conv.Sitof(counter.CountAssignments(observed))
			// factors[i] = P(observed)
			factors[i] = factor.NewFactorValues(observed, cardin, values).Normalize()
			if len(hidden) > 0 {
				if indePot != 0 {
					// g = P(hidden)
					g := factor.NewFactor(hidden, cardin)
					g.SetValues(proportionalValues(1, len(g.Values()), typePot, alphas))
					// P(observed, hidden) = P(observed) * P(hidden)
					factors[i] = factors[i].Product(g)
				} else {
					// g = P(hidden/observed)
					g := factor.NewFactor(cliques[i], cardin)
					lenobs := len(factors[i].Values())
					g.SetValues(proportionalValues(lenobs, len(g.Values())/lenobs, typePot, alphas))
					// P(observed, hidden) = P(observed) * P(hidden/observed)
					factors[i] = factors[i].Product(g)
				}
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
			stats.Dirichlet(alphas[:lenhidden], aux)
		case EmpiricRandom:
			stats.Random(aux)
		case EmpiricUniform:
			stats.Uniform(aux)
		default:
			stats.Uniform(aux)
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
	errchk.Check(err, "")
	defer f.Close()
	ct.SaveOn(f)
}

// LoadCliqueTree loads a clique tree from the given file
func LoadCliqueTree(fname string) *cliquetree.CliqueTree {
	f, err := os.Open(fname)
	errchk.Check(err, "")
	defer f.Close()
	return cliquetree.LoadFrom(f)
}
