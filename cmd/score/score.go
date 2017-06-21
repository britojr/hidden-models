// calculates cliquetree score
package main

import (
	"fmt"
	"os"
	"time"

	"github.com/britojr/kbn/dataset"
	"github.com/britojr/kbn/learn"
	"github.com/britojr/kbn/likelihood"
)

func main() {
	if len(os.Args) < 3 {
		fmt.Printf("\nUsage: %v <cliquetree.ct0> <dataset.csv>\n\n", os.Args[0])
		os.Exit(1)
	}

	c := learn.LoadCliqueTree(os.Args[1])
	ds := dataset.NewFromFile(os.Args[2], rune(','), dataset.HdrNameCard)
	start := time.Now()
	ll := likelihood.StructLoglikelihood(c.Cliques(), c.SepSets(), ds)
	elapsed := time.Since(start)
	fmt.Printf("LL: %v, time: %v\n", ll, elapsed)
}
