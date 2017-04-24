// Package em implements expectation-maximization algorithm
package em

import (
	"fmt"
	"math"

	"github.com/britojr/kbn/cliquetree"
	"github.com/britojr/kbn/factor"
	"github.com/britojr/kbn/filehandler"
	"github.com/britojr/kbn/likelihood"
)

// ExpectationMaximization ..
func ExpectationMaximization(ct *cliquetree.CliqueTree, ds filehandler.DataHandler, epslon float64) {
	diff := epslon * 10
	var llnew, llant float64
	llant = likelihood.Loglikelihood2(ct, ds)
	i := 0
	for ; diff >= epslon; i++ {
		fmt.Printf(".")
		newpot := expectationStep(ct, ds)
		for j := range newpot {
			if ct.Parents()[j] >= 0 {
				newpot[j] = newpot[j].Division(newpot[j].SumOut(ct.Varin(j)))
			} else {
				newpot[j].Normalize()
			}
		}
		ct.SetAllPotentials(newpot)
		llnew = likelihood.Loglikelihood2(ct, ds)
		diff = math.Abs((llnew - llant) / llant)
		llant = llnew
	}
	fmt.Printf("\nIterations: %v\n", i)
}

// expectationStep calculates the expected count of a list of observations and a cliquetree
func expectationStep(ct *cliquetree.CliqueTree, ds filehandler.DataHandler) []*factor.Factor {
	// initialize counter
	count := make([]*factor.Factor, ct.Size())
	for i := range count {
		count[i] = ct.InitialPotential(i).ClearCopy()
	}

	// calculate probability of every instance
	ct.StorePotentials()
	for _, m := range ds.Data() {
		ct.ReduceByEvidence(m)
		ct.UpDownCalibration()
		for i := range count {
			ct.Calibrated(i).Normalize()
			for j, v := range ct.Calibrated(i).Values() {
				count[i].Values()[j] += v
			}
		}
		ct.RecoverPotentials()
	}

	return count
}
