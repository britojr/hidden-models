package learn

import (
	"fmt"
	"os"
	"time"

	"github.com/britojr/kbn/cliquetree"
	"github.com/britojr/kbn/counting/bitcounter"
	"github.com/britojr/kbn/errchk"
	"github.com/britojr/kbn/likelihood"
)

// StructureCommand learns a cliquetree structure corresponding to the given dataset
// with the specified treewidth and additional latent variables
// the learned structure is saved in the given file
func StructureCommand(dsfile, ctfile string, k, h, hc, nk int, delim, hdr uint) {
	data, dscardin := ExtractData(dsfile, delim, hdr)
	n := len(dscardin)

	start := time.Now()
	ct := cliquetree.NewRandom(n+h, k)
	counter := bitcounter.NewBitCounter()
	counter.LoadFromData(data, dscardin)
	sll := likelihood.StructLoglikelihood(ct.Cliques(), ct.SepSets(), counter)
	elapsed := time.Since(start)

	fmt.Printf("%v,%v,%v,%v,%v,%v,%v\n", dsfile, ctfile, n, k, h, sll, elapsed)

	if len(ctfile) > 0 {
		f, err := os.Create(ctfile)
		errchk.Check(err, fmt.Sprintf("Can't create file %v", ctfile))
		ct.SaveOn(f)
		f.Close()
	}
}
