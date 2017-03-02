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
		norm       bool
		initpot    int
	)
	flag.IntVar(&k, "k", 5, "tree-width")
	flag.IntVar(&iterations, "it", 100, "number of iterations/samples")
	flag.StringVar(&dsfile, "f", "", "dataset file")
	flag.UintVar(&delimiter, "delimiter", ',', "field delimiter")
	flag.UintVar(&hdr, "hdr", 1, "1- name header, 2- cardinality header")
	flag.IntVar(&h, "h", 0, "hidden variables")
	flag.BoolVar(&norm, "norm", true, "normalize potentials")
	flag.IntVar(&initpot, "initpot", 1, "1- random values, 2- uniform values")

	// Parse and validate arguments
	flag.Parse()
	if len(dsfile) == 0 {
		fmt.Println("Please enter dataset file name.")
		return
	}
	fmt.Printf("Args: it=%v, k=%v, h=%v\n", iterations, k, h)
	fmt.Printf("Args: norm=%v, initpot=%v\n", norm, initpot)

	learner := learn.New()
	learner.SetIterations(iterations)
	learner.SetTreeWidth(k)
	learner.SetHiddenVars(h)
	learner.SetNorm(norm)
	learner.SetInitPot(initpot)

	fmt.Printf("Loading dataset: %v\n", dsfile)
	start := time.Now()
	learner.LoadDataSet(dsfile, rune(delimiter), filehandler.HeaderFlags(hdr))
	elapsed := time.Since(start)
	fmt.Printf("Time: %v\n", elapsed)

	fmt.Println("Learning structure...")
	start = time.Now()
	ct, ll := learner.GuessStructure()
	elapsed = time.Since(start)
	fmt.Printf("Time: %v; LogLikelihood: %v\n", elapsed, ll)

	fmt.Println("Learning parameters...")
	start = time.Now()
	learner.OptimizeParameters(ct)
	elapsed = time.Since(start)
	fmt.Printf("Time: %v; CT: %v\n", elapsed, ct.Size())

}
