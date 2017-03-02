// Package em implements expectation-maximization algorithm
package em

import (
	"fmt"

	"github.com/britojr/kbn/cliquetree"
	"github.com/britojr/kbn/factor"
	"github.com/britojr/kbn/filehandler"
)

var maxiterations = 5

const epslon = 1e-6

// ExpectationMaximization ..
func ExpectationMaximization(ct *cliquetree.CliqueTree, ds *filehandler.DataSet, norm bool) {
	// TODO: replace maxiterations for convergence test
	diff := epslon + 1
	for i := 1; i <= maxiterations || diff >= epslon; i++ {
		fmt.Printf("Iteration: %v\n", i)
		newpot := expectationStep(ct, ds)
		if norm {
			for j := range newpot {
				newpot[j].Normalize()
			}
		}
		diff = factor.MaxDifference(ct.BkpPotentialList(), newpot)
		ct.SetAllPotentials(newpot)
	}
}

// expectationStep ..
func expectationStep(ct *cliquetree.CliqueTree, ds *filehandler.DataSet) []*factor.Factor {
	// initialize counter
	count := make([]*factor.Factor, ct.Size())
	for i := range count {
		count[i] = ct.GetPotential(i).ClearCopy()
	}

	// calculate probability of every instance
	for _, m := range ds.Data() {
		ct.RestrictByEvidence(m)
		ct.UpDownCalibration()
		// ct.LoadCalibration()
		for i := range count {
			ct.Calibrated(i).Normalize()
			for j, v := range ct.Calibrated(i).Values() {
				count[i].Values()[j] += v
			}
		}
	}

	return count
}
