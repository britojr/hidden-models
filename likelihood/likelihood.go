package likelihood

import (
	"fmt"
	"math"

	"github.com/britojr/kbn/cliquetree"
	"github.com/britojr/kbn/filehandler"
	"github.com/britojr/kbn/utils"
)

// Loglikelihood calculates the log-likelihood line by line
func Loglikelihood(ct *cliquetree.CliqueTree, ds filehandler.DataHandler) (ll float64) {
	for _, m := range ds.Data() {
		v := ct.ProbOfEvidence(m)
		if v != 0 {
			ll += math.Log(v)
		} else {
			// an evidence should never have zero probability
			panic(fmt.Sprintf("zero probability for evid: %v", m))
		}
	}
	return
}

// StructLoglikelihood calculates the log-likelihood based on the counting of observed variables
func StructLoglikelihood(cliques, sepsets [][]int, counter utils.Counter) (ll float64) {
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
