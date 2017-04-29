package main

import (
	"flag"
	"fmt"
	"os"
	"time"

	"github.com/britojr/kbn/cliquetree"
	"github.com/britojr/kbn/filehandler"
	"github.com/britojr/kbn/learn"
	"github.com/britojr/kbn/likelihood"
	"github.com/britojr/kbn/utils"
)

// StepFlags type is used to store the flags indicating the steps that should be executed
// type StepFlags byte
const (
	// StructStep indicates execute structure learning step
	StructStep int = 1 << iota
	//ParamStep indicates parameter learning step
	ParamStep
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
	epslon     float64 // minimum precision for EM convergence
	alpha      float64 // alpha parameter for dirichlet distribution
	ctfile     string  // cliquetree file
	steps      int     // flags indicating what steps to execute
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
	flag.Float64Var(&epslon, "e", 0, "minimum precision for EM convergence")
	flag.Float64Var(&alpha, "a", 1, "alpha parameter for dirichlet distribution")
	flag.StringVar(&ctfile, "s", "", "cliquetree file")
	flag.IntVar(&steps, "steps", StructStep|ParamStep,
		`		step flags:
		1- structure learning,
		2- parameter learning`)

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

	var ct *cliquetree.CliqueTree
	var ll float64
	if steps&StructStep > 0 {
		ct, ll = learnStructureAndParamenters()
		fmt.Printf("Best LL: %v\n", ll)
		if len(ctfile) > 0 {
			f, err := os.Create(ctfile)
			utils.ErrCheck(err, fmt.Sprintf("Can't create file %v", ctfile))
			ct.SaveOn(f)
			f.Close()
		}
	} else {
		f, err := os.Open(ctfile)
		utils.ErrCheck(err, fmt.Sprintf("Can't open file %v", ctfile))
		ct = cliquetree.LoadFrom(f)
		f.Close()
		if steps&ParamStep > 0 {
			ll = learnParameters(ct)
			fmt.Printf("Best LL: %v\n", ll)
		}
	}

	// TODO: add here the MRF reading step
	// TODO: add inference step
	/*
		package mrf markovrf markrf
		func LoadFrom(r reader) mrf*
		func (m* mrf) mrfUnnormalidedMesures(l.dataset) []float64

		func (*c cliquetree) ProbOfAll(l.dataset) []float64
		func estimatePartitionFunction(phi, p []float64) float64
		// show obtained Z
	*/

	fmt.Printf("(%v)\n", ct.Size())
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

	// TODO: remove this structure likelihood
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

	// TODO: remove this check
	if check {
		learner.CheckTree(ct)
	}
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
	return ct, ll
}
