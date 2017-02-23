// Package em implements expectation-maximization algorithm
package em

import (
	"fmt"

	"github.com/britojr/kbn/assignment"
	"github.com/britojr/kbn/cliquetree"
	"github.com/britojr/kbn/factor"
	"github.com/britojr/kbn/filehandler"
)

var maxiterations = 1

// ExpectationMaximization ..
func ExpectationMaximization(ct *cliquetree.CliqueTree, ds *filehandler.DataSet) {
	// TODO: replace maxiterations for convergence test
	for i := 0; i < maxiterations; i++ {
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
		count[i] = factor.NewFactor(f.Variables(), f.Cardinality())
	}

	// calculate probability of every instance
	for _, m := range ds.Data() {
		fmt.Printf("line: %v\n", m)
		for i := 0; i < ct.Size(); i++ {
			fmt.Printf("%v\n", ct.GetPotential(i))
			break
		}
		for i := 0; i < ct.Size(); i++ {
			fmt.Printf("%v\n", ct.GetInitPotential(i))
			break
		}
		ct.RestrictByEvidence(m)
		for i := 0; i < ct.Size(); i++ {
			fmt.Printf("%v\n", ct.GetInitPotential(i))
			break
		}
		ct.UpDownCalibration()
		for i := 0; i < ct.Size(); i++ {
			fmt.Printf("%v\n", ct.Calibrated(i))
			break
		}
		for i := range count {
			f := ct.Calibrated(i)
			f.Normalize()
			assig := assignment.New(f.Variables(), f.Cardinality())
			for assig != nil {
				count[i].Add(assig, f.Get(assig))
				assig.Next()
			}
		}
	}

	return count
}
