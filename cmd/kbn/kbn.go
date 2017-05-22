package main

import (
	"flag"
	"fmt"
	"log"
	"os"
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
	dsfile    string // dataset file name
	delimiter uint   // dataset file delimiter
	hdr       uint   // dataset file header type

	// struct command
	k         int    // treewidth
	h         int    // number of hidden variables
	hcard     int    // cardinality of hiden variables
	numktrees int    // number of k-trees to sample
	ctfileout string // cliquetree save file

	// param command
	ctfilein string  // cliquetree load file
	epslon   float64 // minimum precision for EM convergence
	iterem   int     // number of EM random restarts
	potdist  string  // initial potential distribution
	potmode  string  // mode to complete the initial potential distribution
	alpha    float64 // alpha parameter for dirichlet distribution
	marfile  string  // save marginals file
	//ctfileout

	//partsum command
	//ctfilein
	mkfile  string  // markov random field uai file
	zfile   string  // file to save the log partsum
	discard float64 // discard factor
)

// Define subcommands
var (
	structComm, paramComm, partsumComm *flag.FlagSet
)

func printDefaults() {
	fmt.Printf("Usage:\n\n")
	fmt.Printf("\tkbn <command> [arguments]\n\n")
	fmt.Printf("The commands are:\n\n")
	fmt.Printf("\t%v\n\t%v\n\t%v\n", structConst, paramConst, partsumConst)
	fmt.Println()
}

func main() {
	// Subcommands
	structComm = flag.NewFlagSet(structConst, flag.ExitOnError)
	paramComm = flag.NewFlagSet(paramConst, flag.ExitOnError)
	partsumComm = flag.NewFlagSet(partsumConst, flag.ExitOnError)

	// struct subcommand flags
	structComm.UintVar(&delimiter, "delim", ',', "field delimiter")
	structComm.UintVar(&hdr, "hdr", 4, "1- name header, 2- cardinality header,  4- name_card header")
	structComm.StringVar(&dsfile, "d", "", "dataset csv file (required)")
	structComm.StringVar(&ctfileout, "cs", "", "cliquetree save file")
	structComm.IntVar(&k, "k", 5, "treewidth of the structure")
	structComm.IntVar(&h, "h", 0, "number of hidden variables")
	structComm.IntVar(&hcard, "hc", 2, "cardinality of hidden variables")
	structComm.IntVar(&numktrees, "nk", 1, "number of ktrees samples")

	// param subcommand flags
	paramComm.UintVar(&delimiter, "delim", ',', "field delimiter")
	paramComm.UintVar(&hdr, "hdr", 4, "1- name header, 2- cardinality header,  4- name_card header")
	paramComm.StringVar(&dsfile, "d", "", "dataset csv file (required)")
	paramComm.StringVar(&ctfilein, "cl", "", "cliquetree load file (required)")
	paramComm.StringVar(&ctfileout, "cs", "", "cliquetree save file")
	paramComm.IntVar(&iterem, "iterem", 1, "number of EM iterations")
	paramComm.Float64Var(&epslon, "e", 1e-2, "minimum precision for EM convergence")
	paramComm.Float64Var(&alpha, "a", 1, "alpha parameter, required for --dist=dirichlet")
	paramComm.StringVar(&potdist, "dist", "uniform", "distribution {random|uniform|dirichlet} (required)")
	paramComm.StringVar(&potmode, "mode", "independent", "mode {independent|conditional|full} (required)")

	// partsum subcommand flags
	partsumComm.UintVar(&delimiter, "delim", ',', "field delimiter")
	partsumComm.UintVar(&hdr, "hdr", 4, "1- name header, 2- cardinality header,  4- name_card header")
	partsumComm.StringVar(&dsfile, "d", "", "dataset csv file (required)")
	partsumComm.StringVar(&ctfilein, "cl", "", "cliquetree load file (required)")
	partsumComm.StringVar(&mkfile, "mrf", "", "mrf load file (required)")
	partsumComm.StringVar(&zfile, "z", "", "file to save the partition sum")
	partsumComm.Float64Var(&discard, "dis", 0, "discard factor should be in [0,1)")

	// Verify that a subcommand has been provided
	// os.Arg[0] : main command
	// os.Arg[1] : subcommand
	if len(os.Args) < 2 {
		printDefaults()
		os.Exit(1)
	}

	// Switch on the subcommand
	// Parse the flags for appropriate FlagSet
	// FlagSet.Parse() requires a set of arguments to parse as input
	// os.Args[2:] will be all arguments starting after the subcommand at os.Args[1]
	switch os.Args[1] {
	case structConst:
		structComm.Parse(os.Args[2:])
	case paramConst:
		paramComm.Parse(os.Args[2:])
	case partsumConst:
		partsumComm.Parse(os.Args[2:])
	default:
		printDefaults()
		os.Exit(1)
	}

	// Check which subcommand was Parsed using the FlagSet.Parsed() function. Handle each case accordingly.
	// FlagSet.Parse() will evaluate to false if no flags were parsed (i.e. the user did not provide any flags)
	if structComm.Parsed() {
		// Required Flags
		if dsfile == "" {
			structComm.PrintDefaults()
			os.Exit(1)
		}
		// Print
		log.Printf("ds=%v, cs=%v, k=%v, h=%v, hcard=%v\n",
			dsfile, ctfileout, k, h, hcard,
		)
	}

	if paramComm.Parsed() {
		// Required Flags
		if dsfile == "" || ctfilein == "" {
			paramComm.PrintDefaults()
			os.Exit(1)
		}

		modeChoices := map[string]bool{"independent": true, "conditional": true, "full": true}
		if _, ok := modeChoices[potmode]; !ok {
			paramComm.PrintDefaults()
			os.Exit(1)
		}
		distChoices := map[string]bool{"random": true, "uniform": true, "dirichlet": true}
		if _, ok := distChoices[potdist]; !ok {
			paramComm.PrintDefaults()
			os.Exit(1)
		}
		if potdist == "dirichlet" && alpha == 0 {
			paramComm.PrintDefaults()
			os.Exit(1)
		}
		log.Printf("ds=%v, cl=%v, cst=%v, mode=%v, dist=%v, alpha=%v, eps=%v, iterem=%v\n",
			dsfile, ctfilein, ctfileout, potmode, potdist, alpha, epslon, iterem,
		)
	}

	if partsumComm.Parsed() {
		// Required Flags
		if dsfile == "" || ctfilein == "" || mkfile == "" {
			partsumComm.PrintDefaults()
			os.Exit(1)
		}
		if discard < 0 || discard >= 1 {
			partsumComm.PrintDefaults()
			os.Exit(1)
		}
		log.Printf("ds=%v, cl=%v, mrf=%v, zfile=%v, discard=%v\n",
			dsfile, ctfilein, mkfile, zfile, discard,
		)
	}
}
