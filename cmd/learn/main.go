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
	)
	flag.IntVar(&k, "k", 5, "tree-width")
	flag.IntVar(&iterations, "it", 100, "number of iterations/samples")
	flag.StringVar(&dsfile, "f", "", "dataset file")
	flag.UintVar(&delimiter, "delimiter", ',', "field delimiter")
	flag.UintVar(&hdr, "hdr", 1, "name header 1, cardinality header 2")
	flag.IntVar(&h, "h", 0, "hidden variables")

	// Parse and validate arguments
	flag.Parse()
	if len(dsfile) == 0 {
		fmt.Println("Please enter dataset file name.")
		return
	}
	fmt.Printf("Args: it=%v, k=%v, h=%v\n", iterations, k, h)

	learner := learn.New()
	learner.SetIterations(iterations)
	learner.SetTreeWidth(k)
	learner.SetHiddenVars(h)

	fmt.Printf("Loading dataset: %v\n", dsfile)
	start := time.Now()
	learner.LoadDataSet(dsfile, rune(delimiter), filehandler.HeaderFlags(hdr))
	elapsed := time.Since(start)
	fmt.Printf("Time: %v\n", elapsed)

	fmt.Println("Learning junction tree...")
	start = time.Now()
	jt, ll := learner.BestJuncTree()
	elapsed = time.Since(start)
	fmt.Printf("Time: %v; LogLikelihood: %v\n", elapsed, ll)

	fmt.Println("Learning parameters...")
	start = time.Now()
	ct := learner.OptimizeParameters(jt)
	elapsed = time.Since(start)
	fmt.Printf("Time: %v; CT: %v\n", elapsed, ct.Size())

}
