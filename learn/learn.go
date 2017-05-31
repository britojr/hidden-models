package learn

import (
	"fmt"
	"io"

	"github.com/britojr/kbn/cliquetree"
	"github.com/britojr/kbn/mrf"
	"github.com/britojr/kbn/utl"
)

func saveCliqueTree(ct *cliquetree.CliqueTree, fname string) {
	f := utl.CreateFile(fname)
	defer f.Close()
	ct.SaveOn(f)
}

func loadCliqueTree(fname string) *cliquetree.CliqueTree {
	f := utl.OpenFile(fname)
	defer f.Close()
	return cliquetree.LoadFrom(f)
}

func loadMRF(fname string) *mrf.Mrf {
	f := utl.OpenFile(fname)
	defer f.Close()
	return mrf.LoadFromUAI(f)
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
