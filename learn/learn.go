package learn

import (
	"fmt"
	"sort"

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

	var keys []int
	for k := range ma {
		keys = append(keys, k)
	}
	fmt.Fprintf(f, "MAR\n")
	fmt.Fprintf(f, "%d ", obs)
	sort.Ints(keys)
	for i := 0; i < obs; i++ {
		fmt.Fprintf(f, "%d ", len(ma[keys[i]]))
		for _, v := range ma[keys[i]] {
			fmt.Fprintf(f, "%.5f ", v)
		}
	}
}
