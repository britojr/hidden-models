// Package em implements expectation-maximization algorithm
package em

import (
	"fmt"

	"github.com/britojr/kbn/assignment"
	"github.com/britojr/kbn/cliquetree"
	"github.com/britojr/kbn/factor"
	"github.com/britojr/kbn/filehandler"
)

var maxiterations = 10

// ExpectationMaximization ..
func ExpectationMaximization(ct *cliquetree.CliqueTree, ds *filehandler.DataSet) {
	// TODO: replace maxiterations for convergence test
	for i := 0; i < maxiterations; i++ {
		fmt.Printf("expect-step %v:\n", i)
		fmt.Printf("%v\n", ct.GetPotential(0))
		newpot := expectationStep(ct, ds)
		ct.SetAllPotentials(newpot)
	}
}

// expectationStep ..
func expectationStep(ct *cliquetree.CliqueTree, ds *filehandler.DataSet) []*factor.Factor {
	// initialize counter
	count := make([]*factor.Factor, ct.Size())
	for i := range count {
		f := ct.GetPotential(i)
		count[i] = factor.NewZeroes(f.Variables(), f.Cardinality())
	}

	// calculate probability of every instance
	for _, m := range ds.Data() {
		ct.RestrictByEvidence(m)
		ct.UpDownCalibration()
		for i := range count {
			f := ct.Calibrated(i)
			f.Normalize() // TODO: check if this is really needed
			assig := assignment.New(f.Variables(), f.Cardinality())
			for assig != nil {
				count[i].Add(assig, f.Get(assig))
				assig.Next()
			}
		}
	}

	return count
}