package learn

import (
	"time"

	"github.com/britojr/kbn/cliquetree"
	"github.com/britojr/kbn/dataset"
	"github.com/britojr/kbn/likelihood"
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
