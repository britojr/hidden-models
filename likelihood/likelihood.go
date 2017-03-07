package likelihood

import (
	"math"

	"github.com/britojr/kbn/cliquetree"
	"github.com/britojr/kbn/utils"
)

// loglikelihood1 calculates the log-likelihood weighted by counts
func loglikelihood1(ct *cliquetree.CliqueTree, counter utils.Counter, numobs int) (ll float64) {
	for i, clique := range ct.Cliques() {
		var observed, hidden []int
		if len(counter.Cardinality()) > numobs {
			observed, hidden = utils.SliceSplit(clique, numobs)
		} else {
			observed = clique
		}
		if len(observed) > 0 {
			count := counter.CountAssignments(observed)
			f := ct.Calibrated(i).SumOut(hidden)
			for j, v := range count {
				ll += float64(v) * math.Log(f.Values()[j])
			}
		}
	}
	for i, sepset := range ct.SepSets() {
		var observed, hidden []int
		if len(counter.Cardinality()) > numobs {
			observed, hidden = utils.SliceSplit(sepset, numobs)
		} else {
			observed = sepset
		}
		if len(observed) > 0 {
			count := counter.CountAssignments(observed)
			f := ct.Calibrated(i).SumOut(hidden)
			for j, v := range count {
				ll -= float64(v) * math.Log(f.Values()[j])
			}
		}
	}
	return
}

// loglikelihood2 calculates the log-likelihood line by line
func loglikelihood2(cliques, sepsets [][]int, counter utils.Counter) (ll float64) {
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

// StructLog calculates the log-likelihood
func StructLog(cliques, sepsets [][]int, counter utils.Counter) (ll float64) {
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

func sumLogCount(varlist []int, counter utils.Counter) (ll float64) {
	values := counter.CountAssignments(varlist)
	for _, v := range values {
		if v != 0 {
			ll += float64(v) * math.Log(float64(v))
		}
	}
	return ll
}
