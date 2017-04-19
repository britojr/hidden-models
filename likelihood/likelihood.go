package likelihood

import (
	"fmt"
	"math"

	"github.com/britojr/kbn/cliquetree"
	"github.com/britojr/kbn/filehandler"
	"github.com/britojr/kbn/utils"
)

// Loglikelihood1 calculates the log-likelihood weighted by counts
func Loglikelihood1(ct *cliquetree.CliqueTree, counter utils.Counter, numobs int) (ll float64) {
	for i, clique := range ct.Cliques() {
		// TODO: how to avoid this split?
		observed, hidden := utils.SliceSplit(clique, numobs)
		if len(observed) > 0 {
			count := counter.CountAssignments(observed)
			for j, v := range ct.Calibrated(i).SumOut(hidden).Values() {
				if v != 0 {
					ll += float64(count[j]) * math.Log(v)
				}
			}
		}
	}
	for i, sepset := range ct.SepSets() {
		observed, hidden := utils.SliceSplit(sepset, numobs)
		if len(observed) > 0 {
			count := counter.CountAssignments(observed)
			for j, v := range ct.CalibratedSepSet(i).SumOut(hidden).Values() {
				if v != 0 {
					ll -= float64(count[j]) * math.Log(v)
				}
			}
		}
	}
	return
}

// Loglikelihood2 calculates the log-likelihood line by line
func Loglikelihood2(ct *cliquetree.CliqueTree, ds filehandler.DataHandler) (ll float64) {
	for _, m := range ds.Data() {
		v := ct.ProbOfEvidence(m)
		if v != 0 {
			ll += math.Log(v)
		} else {
			panic(fmt.Sprintf("zero probability for evid: %v", m))
		}
	}
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
