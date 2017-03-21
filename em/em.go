// Package em implements expectation-maximization algorithm
package em

import (
	"fmt"
	"math"

	"github.com/britojr/kbn/cliquetree"
	"github.com/britojr/kbn/factor"
	"github.com/britojr/kbn/filehandler"
	"github.com/britojr/kbn/utils"
)

const epslon = 1e-10

// ExpectationMaximization ..
func ExpectationMaximization(ct *cliquetree.CliqueTree, ds *filehandler.DataSet) {
	diff := epslon * 10
	var err error
	for i := 1; diff >= epslon; i++ {
		fmt.Printf("Iteration: %v\n", i)
		// TODO: check what is to be done for the maximization step
		newpot := expectationStep(ct, ds)
		for j := range newpot {
			newpot[j].Normalize()
		}
		diff, err = checkFactorDiff(ct, newpot, diff)
		utils.ErrCheck(err, "")
		ct.SetAllPotentials(newpot)
	}
}

// expectationStep ..
func expectationStep(ct *cliquetree.CliqueTree, ds *filehandler.DataSet) []*factor.Factor {
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

func checkFactorDiff(ct *cliquetree.CliqueTree, fs []*factor.Factor, threshold float64) (float64, error) {
	if ct.Size() != len(fs) {
		return 0, fmt.Errorf("missing potentials %v x %v", ct.Size(), len(fs))
	}
	var diff float64
	for i := 0; i < ct.Size(); i++ {
		d, err := maxDiff(ct.InitialPotential(i).Values(), fs[i].Values(), threshold)
		if err != nil {
			return 0, err
		}
		if d > diff {
			diff = d
			if diff >= threshold {
				return diff, nil
			}
		}
	}
	return diff, nil
}

// MaxDiff calculates the max difference of two slices,
// if the difference is already bigger than threshold, stops the calculation
func maxDiff(a, b []float64, threshold float64) (float64, error) {
	if len(a) != len(b) {
		return 0, fmt.Errorf("slices of different sizes %v x %v", len(a), len(b))
	}
	var diff float64
	for i := range a {
		if d := math.Abs(a[i] - b[i]); d > diff {
			diff = d
			if diff >= threshold {
				return diff, nil
			}
		}
	}
	return diff, nil
}
