package learn

import (
	"fmt"
	"log"
	"time"

	"github.com/britojr/kbn/cliquetree"
	"github.com/britojr/kbn/dataset"
	"github.com/britojr/kbn/em"
	"github.com/britojr/kbn/factor"
	"github.com/britojr/kbn/likelihood"
	"github.com/britojr/kbn/list"
	"github.com/britojr/kbn/utl/conv"
	"github.com/britojr/kbn/utl/stats"
)

// DependenceMode specifies the dependence mode to use on factor initialilazion
type DependenceMode string

// define possible modes of potential initialilazion
const (
	ModeIndep DependenceMode = "indep" // empiric times distribution independently: p(x,y) = p(x)*p(y)
	ModeCond                 = "cond"  // empiric times distribution conditionaly: p(x,y) = p(x)*p(y|x)
	ModeFull                 = "full"  // distribution only
)

var depmodes = map[string]DependenceMode{
	string(ModeIndep): ModeIndep,
	string(ModeCond):  ModeCond,
	string(ModeFull):  ModeFull,
}

// String returns the distance functons names
func (d DependenceMode) String() string { return string(d) }

// ValidDependenceMode returns a dependence mode option from a string
func ValidDependenceMode(a string) (d DependenceMode, err error) {
	var ok bool
	d, ok = depmodes[a]
	if !ok {
		err = fmt.Errorf("invalid dependence mode string: %v", a)
	}
	return
}

// Distribution specifies the type of distribution to use on factor initialilazion
type Distribution string

// define available distributions used in potential initialilazion
const (
	DistRandom    Distribution = "rand"
	DistUniform                = "unif"
	DistDirichlet              = "dirichlet"
)

var distributions = map[string]Distribution{
	string(DistRandom):    DistRandom,
	string(DistUniform):   DistUniform,
	string(DistDirichlet): DistDirichlet,
}

// String returns the distance functons names
func (d Distribution) String() string { return string(d) }

// ValidDistribution returns a distribution option from a string
func ValidDistribution(a string) (d Distribution, err error) {
	var ok bool
	d, ok = distributions[a]
	if !ok {
		err = fmt.Errorf("invalid distribution string: %v", a)
	}
	return
}

// Counter describes an object that can count occurrences of assignments in a dataset
type Counter interface {
	CountAssignments([]int) []int
}

// Parameters learns the parameters of a cliquetree structure based on a dataset
// the learned structure is saved in the optional output file
func Parameters(
	ds *dataset.Dataset, ctin, ctout, marfile string, hc int,
	alpha, epslon float64, potdist Distribution, potmode DependenceMode, skipEM bool,
) (float64, time.Duration) {
	ct := LoadCliqueTree(ctin)
	log.Printf("Successfully read cliquetree\n")
	cardin := extendCardin(ds.Cardin(), ct.N(), hc)

	start := time.Now()
	ll := learnParameters(
		ct, ds, cardin, ds.NCols(), alpha, epslon, potdist, potmode, skipEM,
	)
	elapsed := time.Since(start)
	log.Printf("Learned parameters in %v\n", elapsed)

	if len(ctout) > 0 {
		SaveCliqueTree(ct, ctout)
		log.Printf("Saved cliquetree in %v\n", ctout)
	}
	if len(marfile) > 0 {
		saveCTMarginals(ct, ds.NCols(), marfile)
		log.Printf("Saved marginals in %v\n", marfile)
	}
	return ll, elapsed
}

func learnParameters(
	ct *cliquetree.CliqueTree, ds *dataset.Dataset, cardin []int, n int,
	alpha, epslon float64, potdist Distribution, potmode DependenceMode, skipEM bool,
) (ll float64) {
	initializePotentials(ct, ds, cardin, n, potdist, potmode, alpha)
	if skipEM {
		return
	}
	// no need to run EM if there is no latent variables
	if ct.N() > ds.NCols() || potmode == ModeFull {
		ll = em.ExpectationMaximization(ct, ds.Data(), epslon)
	} else {
		// if there are no latent variables, calculate LL just using counting
		ll = likelihood.StructLoglikelihood(ct.Cliques(), ct.SepSets(), ds)
	}
	return
}

func initializePotentials(
	ct *cliquetree.CliqueTree, ds *dataset.Dataset, cardin []int, n int,
	potdist Distribution, potmode DependenceMode, alpha float64,
) {
	if potmode == ModeFull {
		ct.SetAllPotentials(createRandomPotentials(ct.Cliques(), cardin, potdist, alpha))
	} else {
		factors := createEmpiricPotentials(ds, ct.Cliques(), cardin, n, potdist, potmode, alpha)
		for i := range factors {
			if len(ct.Varin(i)) != 0 {
				factors[i] = factors[i].Division(factors[i].SumOut(ct.Varin(i)))
			}
		}
		ct.SetAllPotentials(factors)
	}
}

func createEmpiricPotentials(
	ds Counter, cliques [][]int, cardin []int,
	numobs int, potdist Distribution, potmode DependenceMode, alpha float64,
) []*factor.Factor {
	factors := make([]*factor.Factor, len(cliques))
	for i := range factors {
		observed, hidden := separate(numobs, len(cardin), cliques[i])
		if len(observed) > 0 {
			values := conv.Sitof(ds.CountAssignments(observed))
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

func conditionalValues(lenobs, lenhidden int, potdist Distribution, alpha float64) []float64 {
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
	cliques [][]int, cardin []int, dist Distribution, alpha float64,
) []*factor.Factor {
	factors := make([]*factor.Factor, len(cliques))
	for i := range factors {
		factors[i] = factor.NewFactor(cliques[i], cardin)
		initRandomPotential(factors[i].Values(), dist, alpha)
	}
	return factors
}

func initRandomPotential(values []float64, dist Distribution, alpha float64) {
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
