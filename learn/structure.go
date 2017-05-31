package learn

import (
	"fmt"
	"io"
	"time"

	"github.com/britojr/kbn/cliquetree"
	"github.com/britojr/kbn/dataset"
	"github.com/britojr/kbn/likelihood"
	"github.com/britojr/kbn/utl"
)

// SampleStructure samples a cliquetree structure with limited treewidth
// with number of variables corresponding to the given dataset plus latent variables
func SampleStructure(ds *dataset.Dataset, k, h int, ctfile string) (float64, time.Duration) {
	start := time.Now()
	ct := cliquetree.NewRandom(ds.NCols()+h, k)
	sll := likelihood.StructLoglikelihood(ct.Cliques(), ct.SepSets(), ds)
	elapsed := time.Since(start)
	if len(ctfile) > 0 {
		saveCliqueTree(ct, ctfile)
	}
	return sll, elapsed
}

func loadCliqueTree(fname string) *cliquetree.CliqueTree {
	f := utl.OpenFile(fname)
	defer f.Close()
	return cliquetree.LoadFrom(f)
}

func saveCliqueTree(ct *cliquetree.CliqueTree, fname string) {
	f := utl.CreateFile(fname)
	defer f.Close()
	ct.SaveOn(f)
}

func saveCTMarginals(ct *cliquetree.CliqueTree, obs int, fname string) {
	f := utl.CreateFile(fname)
	defer f.Close()
	ma := ct.Marginals()
	writeMarginals(f, ma[:obs])
}

func writeMarginals(w io.Writer, ma [][]float64) {
	fmt.Fprintf(w, "MAR\n")
	fmt.Fprintf(w, "%d ", len(ma))
	for i := range ma {
		fmt.Fprintf(w, "%d ", len(ma[i]))
		for _, v := range ma[i] {
			fmt.Fprintf(w, "%.5f ", v)
		}
	}
}
