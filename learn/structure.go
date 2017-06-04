package learn

import (
	"fmt"
	"io"
	"time"

	"github.com/britojr/kbn/cliquetree"
	"github.com/britojr/kbn/dataset"
	"github.com/britojr/kbn/likelihood"
	"github.com/britojr/kbn/utl"
	"github.com/britojr/kbn/utl/stats"
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

// SaveMarginas load a cliquetree to calculate and save its marginals
func SaveMarginas(ctfile, marfile string) {
	ct := loadCliqueTree(ctfile)
	saveCTMarginals(ct, -1, marfile)
}

// MarginalsMSE compares two marginals and return the Mean Squared Error
func MarginalsMSE(exact, approx string) (mse float64) {
	e, a := loadMarginals(exact), loadMarginals(approx)

	for i := range e {
		mse += stats.MSE(e[i], a[i])
	}
	return mse / float64(len(e))
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
	if obs > 0 {
		writeMarginals(f, ma[:obs])
	} else {
		writeMarginals(f, ma)
	}
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

func loadMarginals(fname string) [][]float64 {
	f := utl.OpenFile(fname)
	defer f.Close()
	return readMarginals(f)
}

func readMarginals(r io.Reader) (ma [][]float64) {
	fmt.Fscanln(r)
	var n int
	fmt.Fscanf(r, "%d", &n)
	ma = make([][]float64, n)
	for i := range ma {
		fmt.Fscanf(r, "%d", &n)
		ma[i] = make([]float64, n)
		for j := range ma[i] {
			fmt.Fscanf(r, "%f", &ma[i][j])
		}
	}
	return
}
