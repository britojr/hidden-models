package main

import (
	"flag"
	"fmt"
	"time"

	"github.com/britojr/kbn/cliquetree"
	"github.com/britojr/kbn/filehandler"
	"github.com/britojr/kbn/learn"
	"github.com/britojr/kbn/likelihood"
)

const (
	hiddencard = 2
)

// Define Flag variables
var (
	k          int     // treewidth
	numktrees  int     // number of k-trees to sample
	iterEM     int     // number of EM random restarts
	iterations int     // number of iterations of the whole process
	dsfile     string  // dataset file name
	delimiter  uint    // dataset file delimiter
	hdr        uint    // dataset file header type
	h          int     // number of hidden variables
	initpot    int     // type of initial potential
	check      bool    // validate cliquetree
	treefile   string  // file to save cliquetree
	epslon     float64 // minimum precision for EM convergence
	alpha      float64 // alpha parameter for dirichlet distribution
)

var (
	learner *learn.Learner
)

func parseFlags() {
	flag.IntVar(&k, "k", 5, "tree-width")
	flag.IntVar(&numktrees, "numk", 1, "number of ktrees samples")
	flag.IntVar(&iterEM, "iterem", 1, "number of EM iterations")
	flag.IntVar(&iterations, "iterations", 1, "number of iterations of the whole process")
	flag.StringVar(&dsfile, "f", "", "dataset file")
	flag.UintVar(&delimiter, "delimiter", ',', "field delimiter")
	flag.UintVar(&hdr, "hdr", 1, "1- name header, 2- cardinality header")
	flag.IntVar(&h, "h", 0, "hidden variables")
	flag.IntVar(&initpot, "initpot", 0,
		`		0- random values,
		1- empiric + dirichlet,
		2- empiric + random,
		3- empiric + uniform`)
	flag.BoolVar(&check, "check", false, "check tree")
	flag.StringVar(&treefile, "s", "", "saves the tree if informed a file name")
	flag.Float64Var(&epslon, "e", 0, "minimum precision for EM convergence")
	flag.Float64Var(&alpha, "a", 1, "alpha parameter for dirichlet distribution")

	// Parse and validate arguments
	flag.Parse()
	if len(dsfile) == 0 {
		fmt.Println("Please enter dataset file name.")
		return
	}
	fmt.Printf("Args: k=%v, h=%v, initpot=%v\n", k, h, initpot)
	fmt.Printf("eps=%v, numk=%v, iterEM=%v\n", epslon, numktrees, iterEM)
}

func main() {
	parseFlags()
	initializeLearner()
	// TODO: add here the MRF reading step

	ct, ll := learnStructureAndParamenters()
	for i := 1; i < iterations; i++ {
		currct, currll := learnStructureAndParamenters()
		if currll > ll {
			ct, ll = currct, currll
		}
	}
	fmt.Printf("Best LL: %v (%v)\n", ll, ct.Size())
}

func initializeLearner() {
	learner = learn.New(k, h, hiddencard, alpha)
	fmt.Printf("Loading dataset: %v\n", dsfile)
	start := time.Now()
	learner.LoadDataSet(dsfile, rune(delimiter), filehandler.HeaderFlags(hdr))
	elapsed := time.Since(start)
	fmt.Printf("Time: %v\n", elapsed)
}

func learnStructure() *cliquetree.CliqueTree {
	fmt.Println("Learning structure...")
	start := time.Now()
	ct := learn.RandomCliqueTree(learner.TotVar(), k)
	elapsed := time.Since(start)
	fmt.Printf("Time: %v\n", elapsed)

	ll := likelihood.StructLoglikelihood(ct.Cliques(), ct.SepSets(), learner.Counter())
	fmt.Printf("Structure LogLikelihood: %v\n", ll)
	return ct
}

func learnParameters(ct *cliquetree.CliqueTree) float64 {
	fmt.Println("Learning parameters...")
	start := time.Now()
	ll := learner.OptimizeParameters(ct, initpot, iterEM, epslon)
	elapsed := time.Since(start)
	fmt.Printf("Time: %v\n", elapsed)
	return ll
}

func learnStructureAndParamenters() (*cliquetree.CliqueTree, float64) {
	ct := learnStructure()
	ll := learnParameters(ct)
	fmt.Printf("Initial LL: %v\n", ll)
	for i := 1; i < iterations; i++ {
		currct := learnStructure()
		currll := learnParameters(currct)
		fmt.Printf("Current LL: %v\n", currll)
		if currll > ll {
			ct, ll = currct, currll
		}
	}

	if check {
		// TODO: remove this check
		learner.CheckTree(ct)
	}
	return ct, ll
}
