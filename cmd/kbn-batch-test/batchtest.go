/*
run experiments in batch

*/
package main

import (
	"fmt"
	"os"
	"path"
	"path/filepath"
	"strings"

	"github.com/britojr/kbn/errchk"
	"github.com/britojr/kbn/learn"
)

// Define parameters
// defaults
var (
	delim  = uint(',')
	hdr    = uint(4)
	nk     = 0
	hc     = 2
	iterem = 0
	epslon = 0.01

	// struct
	structArgs = []struct {
		k    int
		hf   float64
		iter int
	}{
		{3, 0.0, 1},
		{3, 0.2, 1},
		// {3, 0.0, 3},
		// {3, 0.1, 3},
		// {3, 0.3, 3},
		// {5, 0.0, 3},
		// {5, 0.1, 3},
		// {5, 0.3, 3},
		// {7, 0.0, 3},
		// {7, 0.1, 3},
		// {7, 0.3, 3},
	}

	// param
	paramArgs = []struct {
		alpha            float64
		potdist, potmode int
		iter             int
	}{
		{1, learn.DistRandom, learn.ModeFull, 2},
		{1, learn.DistUniform, learn.ModeFull, 1},
	}

	// partsum
	discards = []float64{0, 0.1, 0.2}
)

var (
	structfn                     = "struct.txt"
	paramfn                      = "param.txt"
	partsumfn                    = "partsum.txt"
	structfp, paramfp, partsumfp *os.File
)

func main() {
	// log.SetFlags(0)
	// log.SetOutput(ioutil.Discard)
	var err error
	fmt.Println("reading files:")
	csvfs, _ := filepath.Glob("*.csv")
	fmt.Println(csvfs)

	fmt.Println("running struct...")
	structfp, err = os.Create(structfn)
	errchk.Check(err, fmt.Sprintf("%v", structfn))
	fmt.Fprintf(structfp, "dsfile,ctfile,n,k,h,sll,elapsed\n")
	defer structfp.Close()
	batchStruct(csvfs)

	fmt.Println("running param...")
	paramfp, err = os.Create(paramfn)
	errchk.Check(err, fmt.Sprintf("%v", paramfn))
	fmt.Fprintf(paramfp,
		"dsfile,ctin,ctout,ll,elapsed,alpha,epslon,potdist,potmode,iterem\n",
	)
	defer paramfp.Close()
	batchParam(csvfs)

	fmt.Println("running partsum...")
	partsumfp, err = os.Create(partsumfn)
	errchk.Check(err, fmt.Sprintf("%v", partsumfn))
	fmt.Fprintf(partsumfp, "dsfile,ctfile,zfile,zm,discard,elapsed\n")
	defer partsumfp.Close()
	batchPartsum(csvfs)
}

func batchStruct(csvfs []string) {
	for _, csvf := range csvfs {
		generateStructs(csvf)
	}
}

func generateStructs(csvf string) {
	_, dscardin := learn.ExtractData(csvf, delim, hdr)
	n := len(dscardin)
	for _, it := range structArgs {
		h := int(it.hf * float64(n))
		for i := 1; i <= it.iter; i++ {
			ctfi := structSaveName(csvf, it.k, h, i)
			structureCommand(csvf, delim, hdr, ctfi, it.k, h, nk)
		}
	}
}

func batchParam(csvfs []string) {
	for _, csvf := range csvfs {
		name := strings.TrimSuffix(csvf, path.Ext(csvf))
		name = strings.TrimSuffix(name, path.Ext(name))
		ctfis, _ := filepath.Glob(fmt.Sprintf("%v*.ct0", name))
		generateParams(csvf, ctfis)
	}
}

func generateParams(csvf string, ctfis []string) {
	for _, ctfi := range ctfis {
		for _, it := range paramArgs {
			for i := 1; i < it.iter; i++ {
				ctfo, marf := paramSaveNames(ctfi, it.alpha, it.potdist, it.potmode, i)
				paramCommand(
					csvf, delim, hdr,
					ctfi, ctfo, marf, hc,
					it.alpha, epslon, iterem, it.potdist, it.potmode,
				)
			}
		}
	}
}

func batchPartsum(csvfs []string) {
	for _, csvf := range csvfs {
		name := strings.TrimSuffix(csvf, path.Ext(csvf))
		name = strings.TrimSuffix(name, path.Ext(name))
		ctfis, _ := filepath.Glob(fmt.Sprintf("%v*.ctp", name))
		generatePartsums(csvf, ctfis)
	}
}

func generatePartsums(csvf string, ctfis []string) {
	for _, ctfi := range ctfis {
		for _, dis := range discards {
			mkfile := mrfFileName(csvf)
			zfile := partsumSaveName(mkfile, dis)
			partsumCommand(
				csvf, delim, hdr, ctfi, mkfile, zfile, dis,
			)
		}
	}
}

func mrfFileName(csvf string) string {
	// return ".uai"
	return strings.TrimSuffix(csvf, path.Ext(csvf))
}

func structSaveName(csvf string, k, h, i int) string {
	// return ".ct0"
	name := strings.TrimSuffix(csvf, path.Ext(csvf))
	name = strings.TrimSuffix(name, path.Ext(name))
	return fmt.Sprintf("%v_(%v_%v_%v).ct0", name, k, h, i)
}

func paramSaveNames(ctfi string, alpha float64, potdist, potmode, i int) (string, string) {
	// return ".ctp", ".ctp.mar"
	name := strings.TrimSuffix(ctfi, path.Ext(ctfi))
	name = fmt.Sprintf("%v_(%v_%v_%v_%v)", name, alpha, potdist, potmode, i)
	return fmt.Sprintf("%v.ctp", name), fmt.Sprintf("%v.mar", name)
}

func partsumSaveName(mkfile string, dis float64) string {
	// return ".z"
	return fmt.Sprintf("%v_(%v).sum", mkfile, dis)
}

func structureCommand(
	dsfile string, delim, hdr uint, ctfile string, k, h, nk int,
) {
	n, sll, elapsed := learn.StructureCommandValues(
		dsfile, delim, hdr, ctfile, k, h, nk,
	)
	fmt.Fprintln(structfp,
		learn.Sprintc(dsfile, ctfile, n, k, h, sll, elapsed),
	)
}

func paramCommand(
	dsfile string, delim, hdr uint, ctin, ctout, marfile string, hc int,
	alpha, epslon float64, iterem, potdist, potmode int,
) {
	ll, elapsed := learn.ParamCommandValues(
		dsfile, delim, hdr, ctin, ctout, marfile, hc,
		alpha, epslon, iterem, potdist, potmode,
	)
	fmt.Fprintln(paramfp, learn.Sprintc(
		dsfile, ctin, ctout, ll, elapsed, alpha, epslon, potdist, potmode, iterem,
	))
}

func partsumCommand(
	dsfile string, delim, hdr uint,
	ctfile, mkfile, zfile string, discard float64,
) {
	zm, elapsed := learn.PartsumCommandValues(
		dsfile, delim, hdr, ctfile, mkfile, zfile, discard,
	)
	fmt.Fprintln(partsumfp,
		learn.Sprintc(dsfile, ctfile, zfile, zm, discard, elapsed),
	)
}
