package main

import (
	"flag"
	"fmt"
	"io/ioutil"
	"log"
	"os"

	"github.com/britojr/kbn/dataset"
	"github.com/britojr/kbn/learn"
	"github.com/britojr/utl/ioutl"
)

// Define subcommand names
const (
	structConst  = "struct"
	paramConst   = "param"
	partsumConst = "partsum"
	marginConst  = "marginals"
	margerrConst = "margerr"
)

// Define Flag variables
var (
	// common
	dsfile  string // dataset file name
	delim   uint   // dataset file delimiter
	hdr     uint   // dataset file header type
	verbose bool   // verbose mode

	// struct command
	k         int    // treewidth
	h         int    // number of hidden variables
	nk        int    // number of k-trees to sample
	ctfileout string // cliquetree save file

	// param command
	ctfilein string  // cliquetree load file
	epslon   float64 // minimum precision for EM convergence
	skipEM   bool    // if should skip EM
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

	// margerr command
	exactmar string // correct maginals file
	compmode string // type of function used to compare two marginals

	// Define subcommands
	structComm, paramComm, partsumComm, marginComm, margerrComm *flag.FlagSet
)

func init() {
	initSubcommands()
}

func main() {
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
	case marginConst:
		marginComm.Parse(os.Args[2:])
		runMarginComm()
	case margerrConst:
		margerrComm.Parse(os.Args[2:])
		runMargerrComm()
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
	if !verbose {
		log.SetFlags(0)
		log.SetOutput(ioutil.Discard)
	}

	fmt.Printf("d=%v, cs=%v, h=%v, k=%v\n",
		dsfile, ctfileout, h, k,
	)
	ds := dataset.NewFromFile(dsfile, rune(delim), dataset.HdrFlags(hdr))
	n := ds.NCols()
	sll, elapsed := learn.SampleStructure(ds, k, h, ctfileout)
	fmt.Println(ioutl.Sprintc(
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
	if !verbose {
		log.SetFlags(0)
		log.SetOutput(ioutil.Discard)
	}
	mode, err := learn.ValidDependenceMode(potmode)
	if err != nil {
		fmt.Printf("\n error: invalid mode option\n\n")
		paramComm.PrintDefaults()
		os.Exit(1)
	}
	dist, err := learn.ValidDistribution(potdist)
	if err != nil {
		fmt.Printf("\n error: invalid dist option\n\n")
		paramComm.PrintDefaults()
		os.Exit(1)
	}
	if potdist == "dirichlet" && alpha == 0 {
		fmt.Printf("\n error: missing alpha parameter\n\n")
		paramComm.PrintDefaults()
		os.Exit(1)
	}

	fmt.Printf(
		"a=%v, cl=%v, cs=%v, d=%v, dist=%v, e=%v, hc=%v, mar=%v, mode=%v, skipem=%v\n",
		alpha, ctfilein, ctfileout, dsfile, potdist, epslon, hcard, marfile, potmode, skipEM,
	)
	ds := dataset.NewFromFile(dsfile, rune(delim), dataset.HdrFlags(hdr))
	ll, elapsed := learn.Parameters(
		ds, ctfilein, ctfileout, marfile, hcard, alpha, epslon, dist, mode, skipEM,
	)
	fmt.Println(ioutl.Sprintc(
		dsfile, ctfilein, ctfileout, ll, elapsed, alpha, epslon, potdist, potmode, skipEM,
	))
}

func runPartsumComm() {
	// Required Flags
	if dsfile == "" || ctfilein == "" || mkfile == "" {
		fmt.Printf("\n error: missing dataset/structure/MRF files\n\n")
		partsumComm.PrintDefaults()
		os.Exit(1)
	}
	if !verbose {
		log.SetFlags(0)
		log.SetOutput(ioutil.Discard)
	}
	if discard < 0 || discard >= .5 {
		fmt.Printf("\n error: invalid dircard factor\n\n")
		partsumComm.PrintDefaults()
		os.Exit(1)
	}

	fmt.Printf("cl=%v, d=%v, dis=%v, m=%v, zf=%v\n",
		ctfilein, dsfile, discard, mkfile, zfile,
	)
	ds := dataset.NewFromFile(dsfile, rune(delim), dataset.HdrFlags(hdr))
	zm, elapsed := learn.PartitionSum(ds, ctfilein, mkfile, zfile, discard)
	fmt.Println(ioutl.Sprintc(dsfile, ctfilein, zfile, zm, discard, elapsed))
}

func runMarginComm() {
	// Required Flags
	if ctfilein == "" || marfile == "" {
		fmt.Printf("\n error: missing cliquetree or marginals files\n\n")
		marginComm.PrintDefaults()
		os.Exit(1)
	}
	if !verbose {
		log.SetFlags(0)
		log.SetOutput(ioutil.Discard)
	}

	fmt.Printf("c=%v, m=%v\n", ctfilein, marfile)
	learn.SaveMarginas(ctfilein, marfile)
	fmt.Println(ioutl.Sprintc(ctfilein, marfile))
}

func runMargerrComm() {
	// Required Flags
	if exactmar == "" || marfile == "" {
		fmt.Printf("\n error: inform two marginal files to compare\n\n")
		margerrComm.PrintDefaults()
		os.Exit(1)
	}
	if !verbose {
		log.SetFlags(0)
		log.SetOutput(ioutil.Discard)
	}
	comp, err := learn.ValidDistanceFunc(compmode)
	if err != nil {
		fmt.Printf("\n error: invalid compare function option\n\n")
		margerrComm.PrintDefaults()
		os.Exit(1)
	}

	fmt.Printf("e=%v, a=%v, c=%v\n", exactmar, marfile, comp)
	dif := learn.CompareMarginals(exactmar, marfile, comp)
	fmt.Println(ioutl.Sprintc(exactmar, marfile, dif))
}

func initSubcommands() {
	// Subcommands
	structComm = flag.NewFlagSet(structConst, flag.ExitOnError)
	paramComm = flag.NewFlagSet(paramConst, flag.ExitOnError)
	partsumComm = flag.NewFlagSet(partsumConst, flag.ExitOnError)
	marginComm = flag.NewFlagSet(marginConst, flag.ExitOnError)
	margerrComm = flag.NewFlagSet(margerrConst, flag.ExitOnError)

	// struct subcommand flags
	structComm.UintVar(&delim, "delim", ',', "field delimiter")
	structComm.UintVar(&hdr, "hdr", 4, "1- name header, 2- cardinality header,  4- name_card header")
	structComm.StringVar(&dsfile, "d", "", "dataset csv file (required)")
	structComm.StringVar(&ctfileout, "cs", "", "cliquetree save file")
	structComm.IntVar(&k, "k", 3, "treewidth of the structure")
	structComm.IntVar(&h, "h", 0, "number of hidden variables")
	structComm.IntVar(&nk, "nk", 1, "number of ktrees samples")
	structComm.BoolVar(&verbose, "v", false, "prints detailed steps")

	// param subcommand flags
	paramComm.UintVar(&delim, "delim", ',', "field delimiter")
	paramComm.UintVar(&hdr, "hdr", 4, "1- name header, 2- cardinality header,  4- name_card header")
	paramComm.StringVar(&dsfile, "d", "", "dataset csv file (required)")
	paramComm.StringVar(&ctfilein, "cl", "", "cliquetree load file (required)")
	paramComm.StringVar(&ctfileout, "cs", "", "cliquetree save file")
	paramComm.StringVar(&marfile, "mar", "", "cliquetree marginals save file")
	paramComm.BoolVar(&skipEM, "skipem", false, "if should skip EM")
	paramComm.Float64Var(&epslon, "e", 1e-2, "minimum precision for EM convergence")
	paramComm.Float64Var(&alpha, "a", 1, "alpha parameter, required for --dist=dirichlet")
	paramComm.StringVar(&potdist, "dist", string(learn.DistUniform),
		fmt.Sprintf("distribution {%v|%v|%v} (required)", learn.DistRandom, learn.DistUniform, learn.DistDirichlet),
	)
	paramComm.IntVar(&hcard, "hc", 2, "cardinality of hidden variables")
	paramComm.StringVar(&potmode, "mode", string(learn.ModeIndep),
		fmt.Sprintf("mode {%v|%v|%v} (required)", learn.ModeIndep, learn.ModeCond, learn.ModeFull),
	)
	paramComm.BoolVar(&verbose, "v", false, "prints detailed steps")

	// partsum subcommand flags
	partsumComm.UintVar(&delim, "delim", ',', "field delimiter")
	partsumComm.UintVar(&hdr, "hdr", 4, "1- name header, 2- cardinality header,  4- name_card header")
	partsumComm.StringVar(&dsfile, "d", "", "dataset csv file (required)")
	partsumComm.StringVar(&ctfilein, "cl", "", "cliquetree load file (required)")
	partsumComm.StringVar(&mkfile, "m", "", "mrf load file (required)")
	partsumComm.StringVar(&zfile, "z", "", "file to save the partition sum")
	partsumComm.Float64Var(&discard, "dis", 0, "discard factor should be in [0,0.5)")
	partsumComm.BoolVar(&verbose, "v", false, "prints detailed steps")

	// margin subcommand flags
	marginComm.StringVar(&ctfilein, "c", "", "cliquetree load file (required)")
	marginComm.StringVar(&marfile, "m", "", "cliquetree marginals save file (required)")

	// margerr subcommand flags
	margerrComm.StringVar(&exactmar, "e", "", "exact marginals file (required)")
	margerrComm.StringVar(&marfile, "a", "", "approximation marginals file (required)")
	margerrComm.StringVar(&compmode, "c", "abs", "compare function {mse|entropy|l1|l2|abs|hel}")
}

func printDefaults() {
	fmt.Printf("Usage:\n\n")
	fmt.Printf("\tkbn <command> [arguments]\n\n")
	fmt.Printf("The commands are:\n\n")
	fmt.Printf("\t%v\n\t%v\n\t%v\n\t%v\n\t%v\n",
		structConst, paramConst, partsumConst, marginConst, margerrConst,
	)
	fmt.Println()
}
