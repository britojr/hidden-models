// calculates cliquetree score
package main

import (
	"fmt"
	"os"

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
	ll := likelihood.StructLoglikelihood(c.Cliques(), c.SepSets(), ds)
	fmt.Println(ll)
}
