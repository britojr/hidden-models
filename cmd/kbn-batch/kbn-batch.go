/*
run experiments in batch
*/
package main

import (
	"flag"
	"fmt"
	"os"
	"path"
	"path/filepath"
	"strings"
	"time"

	"github.com/britojr/kbn/dataset"
	"github.com/britojr/kbn/learn"
	"github.com/britojr/utl/errchk"
	"github.com/britojr/utl/ioutl"
)

const (
	structStep int = 1 << iota
	paramStep
	partsumStep
)

type structArg struct {
	k    int
	hf   float64
	iter int
}
type paramArg struct {
	alpha            float64
	potdist, potmode string
	iter             int
}

// Define parameters defaults
var (
	delim  = uint(',')
	hdr    = uint(4)
	nk     = 0
	hc     = 2
	iterem = 0
	epslon = 0.01

	structArgs = []structArg{}
	paramArgs  = []paramArg{}
	discards   = []float64{}
	//
	argfile                      string
	step                         int
	structfn, paramfn, partsumfn string
	structfp, paramfp, partsumfp *os.File
)

func main() {
	// log.SetFlags(0)
	// log.SetOutput(ioutil.Discard)
	parseFlags()
	readParameters(argfile)

	fmt.Println("reading files:")
	csvfs, _ := filepath.Glob("*.csv")
	fmt.Println(csvfs)

	if step&structStep > 0 {
		batchStruct(csvfs)
	}
	if step&paramStep > 0 {
		batchParam(csvfs)
	}
	if step&partsumStep > 0 {
		batchPartsum(csvfs)
	}
}

func batchStruct(csvfs []string) {
	fmt.Println("running struct...")
	var err error
	t := time.Now().Format(time.RFC3339)
	structfn = fmt.Sprintf("structure_%v.txt", t)
	structfp, err = os.Create(structfn)
	errchk.Check(err, fmt.Sprintf("%v", structfn))
	fmt.Fprintf(structfp, "dsfile,ctfile,n,k,h,sll,elapsed\n")
	defer structfp.Close()

	for _, csvf := range csvfs {
		generateStructs(csvf)
	}
}

func generateStructs(csvf string) {
	ds := dataset.NewFromFile(csvf, rune(delim), dataset.HdrFlags(hdr))
	n := ds.NCols()
	for _, it := range structArgs {
		h := int(it.hf * float64(n))
		for i := 1; i <= it.iter; i++ {
			ctfi := structSaveName(csvf, it.k, h, i)
			structureCommand(csvf, delim, hdr, ctfi, it.k, h, nk)
		}
	}
}

func batchParam(csvfs []string) {
	fmt.Println("running param...")
	var err error
	t := time.Now().Format(time.RFC3339)
	paramfn = fmt.Sprintf("parameters_%v.txt", t)
	paramfp, err = os.Create(paramfn)
	errchk.Check(err, fmt.Sprintf("%v", paramfn))
	fmt.Fprintf(paramfp,
		"dsfile,ctin,ctout,ll,elapsed,alpha,epslon,potdist,potmode,iterem\n",
	)
	defer paramfp.Close()

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
			for i := 1; i <= it.iter; i++ {
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
	fmt.Println("running partsum...")
	var err error
	t := time.Now().Format(time.RFC3339)
	partsumfn = fmt.Sprintf("partsum_%v.txt", t)
	partsumfp, err = os.Create(partsumfn)
	errchk.Check(err, fmt.Sprintf("%v", partsumfn))
	fmt.Fprintf(partsumfp,
		"dsfile,ctfile,zfile,sd, mean, median, mode, min, max,discard,elapsed\n",
	)
	defer partsumfp.Close()

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
			zfile := partsumSaveName(ctfi, dis)
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

func paramSaveNames(ctfi string, alpha float64, potdist, potmode string, i int) (string, string) {
	// return ".ctp", ".ctp.mar"
	name := strings.TrimSuffix(ctfi, path.Ext(ctfi))
	name = fmt.Sprintf("%v_(%v_%v_%v_%v)", name, alpha, potdist, potmode, i)
	return fmt.Sprintf("%v.ctp", name), fmt.Sprintf("%v.mar", name)
}

func partsumSaveName(ctfi string, dis float64) string {
	// return ".z"
	return fmt.Sprintf("%v_(%v).sum", ctfi, dis)
}

func structureCommand(
	dsfile string, delim, hdr uint, ctfile string, k, h, nk int,
) {
	ds := dataset.NewFromFile(dsfile, rune(delim), dataset.HdrFlags(hdr))
	n := ds.NCols()
	sll, elapsed := learn.SampleStructure(ds, k, h, ctfile)
	fmt.Fprintln(structfp,
		ioutl.Sprintc(dsfile, ctfile, n, k, h, sll, elapsed),
	)
}

func paramCommand(
	dsfile string, delim, hdr uint, ctin, ctout, marfile string, hc int,
	alpha, epslon float64, iterem int, potdist, potmode string,
) {
	skipEM := false
	ds := dataset.NewFromFile(dsfile, rune(delim), dataset.HdrFlags(hdr))
	mode, _ := learn.ValidDependenceMode(potmode)
	dist, _ := learn.ValidDistribution(potdist)
	ll, elapsed := learn.Parameters(
		ds, ctin, ctout, marfile, hc, alpha, epslon, dist, mode, skipEM,
	)
	fmt.Fprintln(paramfp, ioutl.Sprintc(
		dsfile, ctin, ctout, ll, elapsed, alpha, epslon, potdist, potmode, iterem,
	))
}

func partsumCommand(
	dsfile string, delim, hdr uint,
	ctfile, mkfile, zfile string, discard float64,
) {
	ds := dataset.NewFromFile(dsfile, rune(delim), dataset.HdrFlags(hdr))
	zm, elapsed := learn.PartitionSum(ds, ctfile, mkfile, zfile, discard)
	fmt.Fprintln(partsumfp,
		ioutl.Sprintc(dsfile, ctfile, zfile, zm, discard, elapsed),
	)
}

func parseFlags() {
	flag.StringVar(&argfile, "a", "", "parameters file")
	flag.IntVar(&step, "s", 7, "1- structure, 2- parameters,  4- partition sum")

	// Parse and validate arguments
	flag.Parse()
	if len(argfile) == 0 {
		fmt.Println("Missing parameters file name")
		flag.PrintDefaults()
		os.Exit(1)
	}
}

func readParameters(argfile string) {
	r, err := os.Open(argfile)
	errchk.Check(err, fmt.Sprintf("Can't open file %v", argfile))
	defer r.Close()
	var (
		nst, npr, nps    int
		k, iter          int
		potdist, potmode string
		hf, alpha, dis   float64
	)
	fmt.Fscanf(r, "%d ", &nst)
	for i := 0; i < nst; i++ {
		// k int, hf float64, iter int
		fmt.Fscanf(r, "%d %f %d", &k, &hf, &iter)
		structArgs = append(structArgs, structArg{k, hf, iter})
	}
	fmt.Fscanf(r, "%d", &npr)
	for i := 0; i < npr; i++ {
		// alpha float64, potdist, potmode int, iter int
		fmt.Fscanf(r, "%f %d %d %d", &alpha, &potdist, &potmode, &iter)
		paramArgs = append(paramArgs, paramArg{alpha, potdist, potmode, iter})
	}
	fmt.Fscanf(r, "%d", &nps)
	for i := 0; i < nps; i++ {
		// dis float64
		fmt.Fscanf(r, "%f", &dis)
		discards = append(discards, dis)
	}
	fmt.Println(structArgs)
	fmt.Println(paramArgs)
	fmt.Println(discards)
}
