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
		iterations int
		dsfile     string
		delimiter  uint
		hdr        uint
		h          int
		initpot    int
		check      bool
	)
	flag.IntVar(&k, "k", 5, "tree-width")
	flag.IntVar(&iterations, "it", 1, "number of iterations/samples")
	flag.StringVar(&dsfile, "f", "", "dataset file")
	flag.UintVar(&delimiter, "delimiter", ',', "field delimiter")
	flag.UintVar(&hdr, "hdr", 1, "1- name header, 2- cardinality header")
	flag.IntVar(&h, "h", 0, "hidden variables")
	flag.IntVar(&initpot, "initpot", 1, "1- random values, 2- uniform values")
	flag.BoolVar(&check, "check", false, "check tree")

	// Parse and validate arguments
	flag.Parse()
	if len(dsfile) == 0 {
		fmt.Println("Please enter dataset file name.")
		return
	}
	fmt.Printf("Args: it=%v, k=%v, h=%v, initpot=%v\n", iterations, k, h, initpot)

	learner := learn.New()
	learner.SetIterations(iterations)
	learner.SetTreeWidth(k)
	learner.SetHiddenVars(h)
	learner.SetInitPot(initpot)

	fmt.Printf("Loading dataset: %v\n", dsfile)
	start := time.Now()
	learner.LoadDataSet(dsfile, rune(delimiter), filehandler.HeaderFlags(hdr))
	elapsed := time.Since(start)
	fmt.Printf("Time: %v\n", elapsed)

	ct, ll := learnStructureAndParamenters(learner, check)
	for i := 1; i < iterations; i++ {
		currct, currll := learnStructureAndParamenters(learner, check)
		if currll > ll {
			ct, ll = currct, currll
		}
	}
	fmt.Printf("Best LL: %v (%v)\n", ll, ct.Size())
}

func learnStructureAndParamenters(learner *learn.Learner, check bool) (*cliquetree.CliqueTree, float64) {
	fmt.Println("Learning structure...")
	start := time.Now()
	ct, ll := learner.GuessStructure(1)
	elapsed := time.Since(start)
	fmt.Printf("Time: %v; Structure LogLikelihood: %v\n", elapsed, ll)

	fmt.Println("Learning parameters...")
	learner.InitializePotentials(ct, 2)
	fmt.Printf("Uniform LL: %v\n", learner.CalculateLikelihood(ct, 1))
	fmt.Printf("Uniform LL: %v\n", learner.CalculateLikelihood(ct, 2))
	learner.InitializePotentials(ct)
	fmt.Printf("Initial LL: %v\n", learner.CalculateLikelihood(ct, 1))
	fmt.Printf("Initial LL: %v\n", learner.CalculateLikelihood(ct, 2))
	start = time.Now()
	learner.OptimizeParameters(ct)
	elapsed = time.Since(start)
	fmt.Printf("Time: %v; CT: %v\n", elapsed, ct.Size())
	ll = learner.CalculateLikelihood(ct, 1)
	fmt.Printf("Final LL: %v\n", ll)
	ll = learner.CalculateLikelihood(ct, 2)
	fmt.Printf("Final LL: %v\n", ll)

	if check {
		learner.CheckTree(ct)
	}
	return ct, ll
}
