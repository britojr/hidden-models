// Package em implements expectation-maximization algorithm
package em

import (
	"fmt"
	"log"
	"math"

	"github.com/britojr/kbn/cliquetree"
	"github.com/britojr/kbn/factor"
	"github.com/britojr/kbn/likelihood"
)

// ExpectationMaximization runs EM algorithm for a cliquetree tree and returns
// the loglikelihood after convergence
func ExpectationMaximization(ct *cliquetree.CliqueTree, data [][]int, epslon float64) float64 {
	diff := epslon * 10
	var llnew, llant float64
	llant = likelihood.Loglikelihood(ct, data)
	i := 0
	for ; diff >= epslon; i++ {
		fmt.Printf(".")
		newpot := expectationStep(ct, data)
		maximizationStep(ct, newpot)
		ct.SetAllPotentials(newpot)
		llnew = likelihood.Loglikelihood(ct, data)
		diff = math.Abs((llnew - llant) / llant)
		llant = llnew
	}
	log.Printf("\nEM Iterations: %v\n", i)
	return llnew
}

// expectationStep calculates the expected count of a list of observations and a cliquetree
func expectationStep(ct *cliquetree.CliqueTree, data [][]int) []*factor.Factor {
	// initialize counter
	count := make([]*factor.Factor, ct.Size())
	for i := range count {
		count[i] = ct.InitialPotential(i).ClearCopy()
	}

	// calculate probability of every instance
	ct.StorePotentials()
	for _, m := range data {
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

func maximizationStep(ct *cliquetree.CliqueTree, newpot []*factor.Factor) {
	for j := range newpot {
		if ct.Parents()[j] >= 0 {
			newpot[j] = newpot[j].Division(newpot[j].SumOut(ct.Varin(j)))
		} else {
			newpot[j].Normalize()
		}
	}
}
