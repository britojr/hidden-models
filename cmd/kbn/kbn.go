package main

import (
	"flag"
	"fmt"
	"log"
	"os"

	"github.com/britojr/kbn/dataset"
	"github.com/britojr/kbn/learn"
	"github.com/britojr/kbn/utl"
)

// Define subcommand names
const (
	structConst  = "struct"
	paramConst   = "param"
	partsumConst = "partsum"
)

// Define Flag variables
var (
	// common
	dsfile string // dataset file name
	delim  uint   // dataset file delimiter
	hdr    uint   // dataset file header type

	// struct command
	k         int    // treewidth
	h         int    // number of hidden variables
	nk        int    // number of k-trees to sample
	ctfileout string // cliquetree save file

	// param command
	ctfilein string  // cliquetree load file
	epslon   float64 // minimum precision for EM convergence
	iterem   int     // number of EM random restarts
	potdist  string  // initial potential distribution
	potmode  string  // mode to complete the initial potential distribution
	hcard    int     // cardinality of hiden variables
	alpha    float64 // alpha parameter for dirichlet distribution
	marfile  string  // save marginals file
	//ctfileout

	//partsum command
	//ctfilein
	mkfile  string  // markov random field uai file
	zfile   string  // file to save the log partsum
	discard float64 // discard factor
)

var (
	// Define subcommands
	structComm, paramComm, partsumComm *flag.FlagSet
	// Define choicemaps
	modeChoices, distChoices map[string]int
)

func init() {
	modeChoices = map[string]int{
		"independent": learn.ModeIndep,
		"conditional": learn.ModeCond,
		"full":        learn.ModeFull,
	}
	distChoices = map[string]int{
		"random":    learn.DistRandom,
		"uniform":   learn.DistUniform,
		"dirichlet": learn.DistDirichlet,
	}
}

func main() {
	parseFlags()

	// Verify that a subcommand has been provided
	// os.Arg[0] : main command, os.Arg[1] : subcommand
	if len(os.Args) < 2 {
		printDefaults()
		os.Exit(1)
	}

	switch os.Args[1] {
	case structConst:
		structComm.Parse(os.Args[2:])
		runStructComm()
	case paramConst:
		paramComm.Parse(os.Args[2:])
		runParamComm()
	case partsumConst:
		partsumComm.Parse(os.Args[2:])
		runPartsumComm()
	default:
		printDefaults()
		os.Exit(1)
	}
}

func runStructComm() {
	// Required Flags
	if dsfile == "" {
		fmt.Printf("\n error: missing dataset file\n\n")
		structComm.PrintDefaults()
		os.Exit(1)
	}

	log.Printf("d=%v, cs=%v, h=%v, k=%v\n",
		dsfile, ctfileout, h, k,
	)
	ds := dataset.NewFromFile(dsfile, rune(delim), dataset.HdrFlags(hdr))
	n := ds.NCols()
	sll, elapsed := learn.SampleStructure(ds, k, h, ctfileout)
	log.Println(utl.Sprintc(
		dsfile, ctfileout, n, k, h, sll, elapsed,
	))
}

func runParamComm() {
	// Required Flags
	if dsfile == "" || ctfilein == "" {
		fmt.Printf("\n error: missing dataset or structure file\n\n")
		paramComm.PrintDefaults()
		os.Exit(1)
	}
	var dist, mode int
	var ok bool
	if mode, ok = modeChoices[potmode]; !ok {
		fmt.Printf("\n error: invalid mode option\n\n")
		paramComm.PrintDefaults()
		os.Exit(1)
	}
	if dist, ok = distChoices[potdist]; !ok {
		fmt.Printf("\n error: invalid dist option\n\n")
		paramComm.PrintDefaults()
		os.Exit(1)
	}
	if potdist == "dirichlet" && alpha == 0 {
		fmt.Printf("\n error: missing alpha parameter\n\n")
		paramComm.PrintDefaults()
		os.Exit(1)
	}

	log.Printf(
		"a=%v, cl=%v, cs=%v, d=%v, dist=%v, e=%v, hc=%v, iterem=%v, mar=%v, mode=%v\n",
		alpha, ctfilein, ctfileout, dsfile, potdist, epslon, hcard, iterem, marfile, potmode,
	)
	learn.ParamCommand(
		dsfile, delim, hdr, ctfilein, ctfileout, marfile,
		hcard, alpha, epslon, iterem, dist, mode,
	)
}

func runPartsumComm() {
	// Required Flags
	if dsfile == "" || ctfilein == "" || mkfile == "" {
		fmt.Printf("\n error: missing dataset/structure/MRF files\n\n")
		partsumComm.PrintDefaults()
		os.Exit(1)
	}
	if discard < 0 || discard >= .5 {
		fmt.Printf("\n error: invalid dircard factor\n\n")
		partsumComm.PrintDefaults()
		os.Exit(1)
	}

	log.Printf("cl=%v, d=%v, dis=%v, m=%v, zf=%v, \n",
		ctfilein, dsfile, discard, mkfile, zfile,
	)
	learn.PartsumCommand(
		dsfile, delim, hdr, ctfilein, mkfile, zfile, discard,
	)
}

func parseFlags() {
	// Subcommands
	structComm = flag.NewFlagSet(structConst, flag.ExitOnError)
	paramComm = flag.NewFlagSet(paramConst, flag.ExitOnError)
	partsumComm = flag.NewFlagSet(partsumConst, flag.ExitOnError)

	// struct subcommand flags
	structComm.UintVar(&delim, "delim", ',', "field delimiter")
	structComm.UintVar(&hdr, "hdr", 4, "1- name header, 2- cardinality header,  4- name_card header")
	structComm.StringVar(&dsfile, "d", "", "dataset csv file (required)")
	structComm.StringVar(&ctfileout, "cs", "", "cliquetree save file")
	structComm.IntVar(&k, "k", 3, "treewidth of the structure")
	structComm.IntVar(&h, "h", 0, "number of hidden variables")
	structComm.IntVar(&nk, "nk", 1, "number of ktrees samples")

	// param subcommand flags
	paramComm.UintVar(&delim, "delim", ',', "field delimiter")
	paramComm.UintVar(&hdr, "hdr", 4, "1- name header, 2- cardinality header,  4- name_card header")
	paramComm.StringVar(&dsfile, "d", "", "dataset csv file (required)")
	paramComm.StringVar(&ctfilein, "cl", "", "cliquetree load file (required)")
	paramComm.StringVar(&ctfileout, "cs", "", "cliquetree save file")
	paramComm.StringVar(&marfile, "mar", "", "cliquetree marginals save file")
	paramComm.IntVar(&iterem, "iterem", 1, "number of EM iterations")
	paramComm.Float64Var(&epslon, "e", 1e-2, "minimum precision for EM convergence")
	paramComm.Float64Var(&alpha, "a", 1, "alpha parameter, required for --dist=dirichlet")
	paramComm.StringVar(&potdist, "dist", "uniform", "distribution {random|uniform|dirichlet} (required)")
	paramComm.IntVar(&hcard, "hc", 2, "cardinality of hidden variables")
	paramComm.StringVar(&potmode, "mode", "independent", "mode {independent|conditional|full} (required)")

	// partsum subcommand flags
	partsumComm.UintVar(&delim, "delim", ',', "field delimiter")
	partsumComm.UintVar(&hdr, "hdr", 4, "1- name header, 2- cardinality header,  4- name_card header")
	partsumComm.StringVar(&dsfile, "d", "", "dataset csv file (required)")
	partsumComm.StringVar(&ctfilein, "cl", "", "cliquetree load file (required)")
	partsumComm.StringVar(&mkfile, "m", "", "mrf load file (required)")
	partsumComm.StringVar(&zfile, "z", "", "file to save the partition sum")
	partsumComm.Float64Var(&discard, "dis", 0, "discard factor should be in [0,0.5)")
}

func printDefaults() {
	fmt.Printf("Usage:\n\n")
	fmt.Printf("\tkbn <command> [arguments]\n\n")
	fmt.Printf("The commands are:\n\n")
	fmt.Printf("\t%v\n\t%v\n\t%v\n", structConst, paramConst, partsumConst)
	fmt.Println()
}
