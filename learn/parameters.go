package learn

import (
	"time"

	"github.com/britojr/kbn/cliquetree"
	"github.com/britojr/kbn/conv"
	"github.com/britojr/kbn/counting"
	"github.com/britojr/kbn/counting/bitcounter"
	"github.com/britojr/kbn/em"
	"github.com/britojr/kbn/factor"
	"github.com/britojr/kbn/list"
	"github.com/britojr/kbn/stats"
)

// define possible modes of potential initialilazion
const (
	ModeIndep = iota // empiric times distribution independently: p(x,y) = p(x)*p(y)
	ModeCond         // empiric times distribution conditionaly: p(x,y) = p(x)*p(y|x)
	ModeFull         // distribution only
)

// define distributions used in potential initialilazion
const (
	DistRandom = iota
	DistUniform
	DistDirichlet
)

// ParamCommand learns the parameters of a cliquetree structure based on the dataset
// the learned structure is saved in the output file
func ParamCommand(
	dsfile string, delim, hdr uint, ctin, ctout, marfile string, hc int,
	alpha, epslon float64, iterem, potdist, potmode int,
) {
	data, dscardin := ExtractData(dsfile, delim, hdr)
	n := len(dscardin)
	counter := bitcounter.NewBitCounter()
	counter.LoadFromData(data, dscardin)

	ct := LoadCliqueTree(ctin)
	cardin := extendCardin(dscardin, ct.N(), hc)

	start := time.Now()
	ll := learnParameters(
		ct, counter, data, cardin, n,
		alpha, epslon, potdist, potmode, iterem,
	)
	elapsed := time.Since(start)

	if len(ctout) > 0 {
		SaveCliqueTree(ct, ctout)
	}
	if len(marfile) > 0 {
		SaveCTMarginals(ct, n, marfile)
	}
	Printcln(dsfile, ctin, ctout, ll, elapsed, alpha, epslon, potdist, potmode, iterem)
}

func learnParameters(
	ct *cliquetree.CliqueTree, counter counting.Counter, data [][]int, cardin []int, n int,
	alpha, epslon float64, potdist, potmode, iter int,
) (ll float64) {
	initializePotentials(ct, counter, cardin, n, potdist, potmode, alpha)
	ll = em.ExpectationMaximization(ct, data, epslon)
	return
}

func initializePotentials(
	ct *cliquetree.CliqueTree, counter counting.Counter, cardin []int, n int,
	potdist, potmode int, alpha float64,
) {
	if potmode == ModeFull {
		ct.SetAllPotentials(createRandomPotentials(ct.Cliques(), cardin, potdist, alpha))
	} else {
		factors := createEmpiricPotentials(counter, ct.Cliques(), cardin, n, potdist, potmode, alpha)
		for i := range factors {
			if len(ct.Varin(i)) != 0 {
				factors[i] = factors[i].Division(factors[i].SumOut(ct.Varin(i)))
			}
		}
		ct.SetAllPotentials(factors)
	}
}

func createEmpiricPotentials(
	counter counting.Counter, cliques [][]int, cardin []int,
	numobs, potdist, potmode int, alpha float64,
) []*factor.Factor {
	factors := make([]*factor.Factor, len(cliques))
	for i := range factors {
		observed, hidden := separate(numobs, len(cardin), cliques[i])
		if len(observed) > 0 {
			values := conv.Sitof(counter.CountAssignments(observed))
			// factors[i] = P(observed)
			factors[i] = factor.NewFactorValues(observed, cardin, values).Normalize()
			if len(hidden) > 0 {
				if potmode == ModeIndep {
					// g = P(hidden)
					g := factor.NewFactor(hidden, cardin)
					initRandomPotential(g.Values(), potdist, alpha)
					// P(observed, hidden) = P(observed) * P(hidden)
					factors[i] = factors[i].Product(g)
				} else {
					// g = P(hidden/observed)
					g := factor.NewFactor(cliques[i], cardin)
					lenobs := len(factors[i].Values())
					g.SetValues(conditionalValues(lenobs, len(g.Values())/lenobs, potdist, alpha))
					// P(observed, hidden) = P(observed) * P(hidden/observed)
					factors[i] = factors[i].Product(g)
				}
			}
		} else {
			factors[i] = factor.NewFactor(cliques[i], cardin)
			initRandomPotential(factors[i].Values(), potdist, alpha)
		}
	}
	return factors
}

func conditionalValues(lenobs, lenhidden, potdist int, alpha float64) []float64 {
	values := make([]float64, lenobs*lenhidden)
	aux := make([]float64, lenhidden)
	for i := 0; i < lenobs; i++ {
		initRandomPotential(aux, potdist, alpha)
		for j := 0; j < lenhidden; j++ {
			values[i+(j*lenobs)] = aux[j]
		}
	}
	return values
}

func createRandomPotentials(
	cliques [][]int, cardin []int, dist int, alpha float64,
) []*factor.Factor {
	factors := make([]*factor.Factor, len(cliques))
	for i := range factors {
		factors[i] = factor.NewFactor(cliques[i], cardin)
		initRandomPotential(factors[i].Values(), dist, alpha)
	}
	return factors
}

func initRandomPotential(values []float64, dist int, alpha float64) {
	switch dist {
	case DistRandom:
		stats.Random(values)
	case DistUniform:
		stats.Uniform(values)
	case DistDirichlet:
		stats.Dirichlet1(alpha, values)
	}
}

func separate(n, t int, varlist []int) (observed, hidden []int) {
	if t > n {
		return list.Split(varlist, n)
	}
	return varlist, nil
}

// extendCardin extends cardinality to add hidden variables
func extendCardin(dscardin []int, t, hc int) []int {
	cardin := make([]int, t)
	copy(cardin, dscardin)
	for i := len(dscardin); i < len(cardin); i++ {
		cardin[i] = hc
	}
	return cardin
}
