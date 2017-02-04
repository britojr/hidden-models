package main

import (
	"flag"
	"fmt"
	"time"

	"github.com/britojr/playgo/filehandler"
	"github.com/britojr/playgo/learn"
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
	// Print arguments
	fmt.Printf("Args: it=%v\n", iterations)
	fmt.Printf("DataSet: %v\n", dsfile)

	learner := learn.New()
	learner.SetIterations(iterations)
	learner.SetTreeWidth(k)
	learner.LoadDataSet(dsfile, rune(delimiter), filehandler.HeaderFlags(hdr))

	start := time.Now()
	_, ll := learner.BestJuncTree()
	elapsed := time.Since(start)

	fmt.Printf("Time: %v, LogLikelihood: %v\n", elapsed, ll)
}
