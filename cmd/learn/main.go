package main

import (
	"flag"
	"fmt"
	"time"

	"github.com/britojr/kbn/cliquetree"
	"github.com/britojr/kbn/filehandler"
	"github.com/britojr/kbn/learn"
)

func main() {
	// Define Flag variables
	var (
		k          int
		numktrees  int
		iterEM     int
		iterations int
		dsfile     string
		delimiter  uint
		hdr        uint
		h          int
		initpot    int
		check      bool
		treefile   string
		epslon     float64
	)
	flag.IntVar(&k, "k", 5, "tree-width")
	flag.IntVar(&numktrees, "numk", 1, "number of ktrees samples")
	flag.IntVar(&iterEM, "iterem", 1, "number of EM iterations")
	flag.IntVar(&iterations, "iterations", 1, "number of iterations of the whole process")
	flag.StringVar(&dsfile, "f", "", "dataset file")
	flag.UintVar(&delimiter, "delimiter", ',', "field delimiter")
	flag.UintVar(&hdr, "hdr", 1, "1- name header, 2- cardinality header")
	flag.IntVar(&h, "h", 0, "hidden variables")
	flag.IntVar(&initpot, "initpot", 1, "1- random values, 2- uniform values")
	flag.BoolVar(&check, "check", false, "check tree")
	flag.StringVar(&treefile, "s", "", "saves the tree if informed a file name")
	flag.Float64Var(&epslon, "e", 0, "minimum precision for EM convergence")

	// Parse and validate arguments
	flag.Parse()
	if len(dsfile) == 0 {
		fmt.Println("Please enter dataset file name.")
		return
	}
	fmt.Printf("Args: k=%v, h=%v, initpot=%v\n", k, h, initpot)
	fmt.Printf("eps=%v, numk=%v, iterEM=%v\n", epslon, numktrees, iterEM)

	learner := learn.New()
	learner.SetTreeWidth(k)
	learner.SetHiddenVars(h)
	learner.SetInitPot(initpot)
	learner.SetEpslon(epslon)

	fmt.Printf("Loading dataset: %v\n", dsfile)
	start := time.Now()
	learner.LoadDataSet(dsfile, rune(delimiter), filehandler.HeaderFlags(hdr))
	elapsed := time.Since(start)
	fmt.Printf("Time: %v\n", elapsed)

	ct, ll := learnStructureAndParamenters(learner, initpot, numktrees, iterEM, check)
	for i := 1; i < iterations; i++ {
		currct, currll := learnStructureAndParamenters(learner, initpot, numktrees, iterEM, check)
		if currll > ll {
			ct, ll = currct, currll
		}
	}
	fmt.Printf("Best LL: %v (%v)\n", ll, ct.Size())

	if len(treefile) > 0 {
		saveTree(ct, treefile)
	}
}

func learnStructureAndParamenters(learner *learn.Learner,
	initpot, numktrees, iterEM int, check bool) (*cliquetree.CliqueTree, float64) {
	fmt.Println("Learning structure...")

	start := time.Now()
	ct, ll := learner.GuessStructure(numktrees)
	elapsed := time.Since(start)
	fmt.Printf("Time: %v; Structure LogLikelihood: %v\n", elapsed, ll)

	fmt.Println("Learning parameters...")
	learner.InitializePotentials(ct, 2)
	fmt.Printf("Uniform LL: %v\n", learner.CalculateLikelihood(ct))
	learner.InitializePotentials(ct)
	fmt.Printf("Initial LL: %v\n", learner.CalculateLikelihood(ct))

	start = time.Now()
	learner.OptimizeParameters(ct, initpot, iterEM)
	elapsed = time.Since(start)
	fmt.Printf("Time: %v; CT: %v\n", elapsed, ct.Size())

	ll = learner.CalculateLikelihood(ct)
	fmt.Printf("Final LL: %v\n", ll)

	if check {
		learner.CheckTree(ct)
	}
	return ct, ll
}

func saveTree(ct *cliquetree.CliqueTree, treefile string) {
	learn.SaveCliqueTree(ct, treefile)
	learn.SaveMarginals(ct, treefile+"marg")
}
