// run experiments in batch
package main

import (
	"flag"
	"fmt"
	"io/ioutil"
	"os"
	"path"
	"path/filepath"
	"strings"
	"time"

	yaml "gopkg.in/yaml.v2"

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

const (
	cFileDelim   = "file_delim"
	cFileHeader  = "file_header"
	cStructBlk   = "struct_blk"
	cTreewidth   = "treewidth"
	cHNum        = "hnum"
	cHProp       = "hprop"
	cRepeat      = "repeat"
	cParamsBlk   = "params_blk"
	cAlpha       = "alpha"
	cDist        = "dist"
	cMode        = "mode"
	cHCard       = "hcard"
	cEMThreshold = "em_threshold"
	cMaxIterEM   = "em_max_iterations"
	cPartsumBlk  = "partsum_blk"
	cDiscards    = "discards"
)

// Define parameters defaults
var (
	delim         = uint(',')
	hdr           = uint(4)
	nk            = 0
	maxIterEM     = 0
	defaultHCard  = []int{2}
	defaultEpslon = 1e-3

	parMap map[interface{}]interface{}
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
	if v, ok := parMap[cStructBlk]; ok {
		structBlks := v.([]interface{})
		for _, it := range structBlks {
			blk := it.(map[interface{}]interface{})
			var h, k int
			repeat := 1
			if v, ok := blk[cRepeat]; ok {
				repeat = v.(int)
				fmt.Printf("%v : '%v'\n", cRepeat, repeat)
			}
			if v, ok := blk[cTreewidth]; ok {
				k = v.(int)
				fmt.Printf("%v : '%v'\n", cTreewidth, k)
			}
			if v, ok := blk[cHProp]; ok {
				h = int(convToF64(v) * float64(n))
				fmt.Printf("%v : '%v'\n", cHProp, h)
			}
			if v, ok := blk[cHNum]; ok {
				h = v.(int)
				fmt.Printf("%v : '%v'\n", cHNum, h)
			}
			for i := 1; i <= repeat; i++ {
				ctfi := structSaveName(csvf, k, h, i)
				structureCommand(csvf, delim, hdr, ctfi, k, h, nk)
			}
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
		if v, ok := parMap[cParamsBlk]; ok {
			for _, it := range v.([]interface{}) {
				blk := it.(map[interface{}]interface{})
				var alpha float64
				var potdist, potmode string
				repeat := 1
				epslon := defaultEpslon
				hc := defaultHCard
				if v, ok := blk[cRepeat]; ok {
					repeat = v.(int)
					fmt.Printf("%v : '%v'\n", cRepeat, repeat)
				}
				if v, ok := blk[cAlpha]; ok {
					alpha = convToF64(v)
					fmt.Printf("%v : '%v'\n", cAlpha, alpha)
				}
				if v, ok := blk[cEMThreshold]; ok {
					epslon = convToF64(v)
					fmt.Printf("%v : '%v'\n", cEMThreshold, epslon)
				}
				if v, ok := blk[cMaxIterEM]; ok {
					maxIterEM = v.(int)
					fmt.Printf("%v : '%v'\n", cMaxIterEM, maxIterEM)
				}
				if v, ok := blk[cDist]; ok {
					potdist = v.(string)
					fmt.Printf("%v : '%v'\n", cDist, potdist)
				}
				if v, ok := blk[cMode]; ok {
					potmode = v.(string)
					fmt.Printf("%v : '%v'\n", cMode, potmode)
				}
				if v, ok := blk[cHCard]; ok {
					hc = make([]int, len(v.([]interface{})))
					for i, c := range v.([]interface{}) {
						hc[i] = c.(int)
					}
					fmt.Printf("%v : '%v'\n", cHCard, hc)
				}
				for i := 1; i <= repeat; i++ {
					ctfo, marf := paramSaveNames(ctfi, alpha, potdist, potmode, i)
					paramCommand(
						csvf, delim, hdr,
						ctfi, ctfo, marf, hc,
						alpha, epslon, maxIterEM, potdist, potmode,
					)
				}
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
		if v, ok := parMap[cParamsBlk]; ok {
			for _, it := range v.([]interface{}) {
				blk := it.(map[interface{}]interface{})
				var dis float64
				if v, ok := blk[cDiscards]; ok {
					dis = convToF64(v)
					fmt.Printf("%v : '%v'\n", cDiscards, dis)
				}
				mkfile := mrfFileName(csvf)
				zfile := partsumSaveName(ctfi, dis)
				partsumCommand(
					csvf, delim, hdr, ctfi, mkfile, zfile, dis,
				)
			}
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
	dsfile string, delim, hdr uint, ctin, ctout, marfile string, hc []int,
	alpha, epslon float64, iterem int, potdist, potmode string,
) {
	skipEM := false
	ds := dataset.NewFromFile(dsfile, rune(delim), dataset.HdrFlags(hdr))
	mode, _ := learn.ValidDependenceMode(potmode)
	dist, _ := learn.ValidDistribution(potdist)
	ll, elapsed := learn.Parameters(
		ds, ctin, ctout, marfile, hc, alpha, epslon, maxIterEM, dist, mode, skipEM,
	)
	fmt.Fprintln(paramfp, ioutl.Sprintc(
		dsfile, ctin, ctout, ll, elapsed, alpha, epslon, potdist, potmode, maxIterEM,
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
	parMap = make(map[interface{}]interface{})
	data, err := ioutil.ReadFile(argfile)
	errchk.Check(err, "")
	errchk.Check(yaml.Unmarshal([]byte(data), &parMap), "")
	if v, ok := parMap[cFileDelim]; ok {
		delim = uint(v.(string)[0])
		fmt.Printf("%v : '%c'\n", cFileDelim, delim)
	}
	if v, ok := parMap[cFileHeader]; ok {
		hdr = uint(v.(int))
		fmt.Printf("%v : '%v'\n", cFileHeader, hdr)
	}
}

func convToF64(v interface{}) float64 {
	switch v.(type) {
	case int:
		return float64(v.(int))
	case float64:
		return v.(float64)
	}
	return 0
}
