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
	)
	flag.IntVar(&k, "k", 5, "tree-width")
	flag.IntVar(&iterations, "it", 100, "number of iterations/samples")
	flag.StringVar(&dsfile, "f", "", "dataset file")
	flag.UintVar(&delimiter, "delimiter", ',', "field delimiter")
	flag.UintVar(&hdr, "hdr", 1, "name header 1, cardinality header 2")

	// Parse and validate arguments
	flag.Parse()
	if len(dsfile) == 0 {
		fmt.Println("Please enter dataset file name.")
		return
	}
	fmt.Printf("Args: it=%v, k=%v\n", iterations, k)

	learner := learn.New()
	learner.SetIterations(iterations)
	learner.SetTreeWidth(k)

	fmt.Printf("Loading dataset: %v\n", dsfile)
	learner.LoadDataSet(dsfile, rune(delimiter), filehandler.HeaderFlags(hdr))

	fmt.Println("Learning junction tree...")
	start := time.Now()
	_, ll := learner.BestJuncTree()
	elapsed := time.Since(start)

	fmt.Printf("Time: %v, LogLikelihood: %v\n", elapsed, ll)
}
