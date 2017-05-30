package learn

import (
	"fmt"
	"time"

	"github.com/britojr/kbn/cliquetree"
	"github.com/britojr/kbn/dataset"
	"github.com/britojr/kbn/likelihood"
)

// StructureCommand learns a cliquetree structure corresponding to the given dataset
// with the specified treewidth and additional latent variables
// the learned structure is saved in the given file
func StructureCommand(
	dsfile string, delim, hdr uint, ctfile string, k, h, nk int,
) {
	n, sll, elapsed := StructureCommandValues(
		dsfile, delim, hdr, ctfile, k, h, nk,
	)
	fmt.Println(Sprintc(dsfile, ctfile, n, k, h, sll, elapsed))
}

// StructureCommandValues learns a cliquetree structure corresponding to the given dataset
// with the specified treewidth and additional latent variables
// the learned structure is saved in the given file
func StructureCommandValues(
	dsfile string, delim, hdr uint, ctfile string, k, h, nk int,
) (int, float64, time.Duration) {
	ds := dataset.NewFromFile(dsfile, rune(delim), dataset.HdrFlags(hdr))
	n := ds.NCols()

	start := time.Now()
	ct := cliquetree.NewRandom(n+h, k)
	sll := likelihood.StructLoglikelihood(ct.Cliques(), ct.SepSets(), ds)
	elapsed := time.Since(start)

	if len(ctfile) > 0 {
		SaveCliqueTree(ct, ctfile)
	}
	return n, sll, elapsed
}
