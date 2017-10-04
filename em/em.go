// Package em implements expectation-maximization algorithm
package em

import (
	"log"
	"math"

	"github.com/britojr/kbn/cliquetree"
	"github.com/britojr/kbn/factor"
	"github.com/britojr/utl/floats"
)

// ExpectationMaximization runs EM algorithm for a cliquetree tree and returns
// the loglikelihood after convergence
func ExpectationMaximization(ct *cliquetree.CliqueTree, data [][]int, epslon float64) float64 {
	log.Printf("Start EM\n")
	var llnew, llant float64
	newpot, llant := expectationStep(ct, data)
	maximizationStep(ct, newpot)
	ct.SetAllPotentials(newpot)
	diff := epslon * 10
	i := 0
	for ; diff >= epslon; i++ {
		newpot, llnew = expectationStep(ct, data)
		maximizationStep(ct, newpot)
		ct.SetAllPotentials(newpot)
		diff = math.Abs((llnew - llant) / llant)
		llant = llnew
		log.Printf("\tdiff: %v\n", diff)
	}
	log.Printf("EM Iterations: %v\n", i)
	return llnew
}

// expectationStep calculates the expected count of a list of observations and a cliquetree
func expectationStep(ct *cliquetree.CliqueTree, data [][]int) ([]*factor.Factor, float64) {
	// initialize counter
	count := make([]*factor.Factor, ct.Size())
	for i := range count {
		count[i] = ct.InitialPotential(i).ClearCopy()
	}

	// calculate probability of every instance
	ct.StorePotentials()
	var ll float64
	for _, m := range data {
		ct.ReduceByEvidence(m)
		ct.UpDownCalibration()
		// accumulate the log-probability to return loglikelihood
		ll += lprob(ct.Calibrated(0).Values())
		for i := range count {
			ct.Calibrated(i).Normalize()
			for j, v := range ct.Calibrated(i).Values() {
				count[i].Values()[j] += v
			}
		}
		ct.RecoverPotentials()
	}
	return count, ll
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

func lprob(values []float64) float64 {
	p := floats.Sum(values)
	if p == 0 {
		panic("invalid log(0)")
	}
	return math.Log(p)
}
