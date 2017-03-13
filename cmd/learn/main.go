package main

import (
	"flag"
	"fmt"
	"time"

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

	fmt.Println("Learning structure...")
	start = time.Now()
	ct, ll := learner.GuessStructure(iterations)
	elapsed = time.Since(start)
	fmt.Printf("Time: %v; Structure LogLikelihood: %v\n", elapsed, ll)

	fmt.Println("Learning parameters...")
	learner.InitializePotentials(ct, 2)
	fmt.Printf("Uniform LL: %v\n", learner.CalculateLikelihood(ct))
	learner.InitializePotentials(ct)
	fmt.Printf("Initial LL: %v\n", learner.CalculateLikelihood(ct))
	start = time.Now()
	learner.OptimizeParameters(ct)
	elapsed = time.Since(start)
	fmt.Printf("Time: %v; CT: %v\n", elapsed, ct.Size())
	fmt.Printf("Final LL: %v\n", learner.CalculateLikelihood(ct))

	if check {
		learner.CheckTree(ct)
	}

}
