package main

import (
	"flag"
	"fmt"
	"os"
	"path/filepath"
	"time"

	"github.com/britojr/kbn/learn"
	"github.com/britojr/kbn/utl"
	"github.com/britojr/kbn/utl/errchk"
)

var (
	exact       string
	metrChoices map[string]int
)

func init() {
	parseFlags()
}

func parseFlags() {
	flag.StringVar(&exact, "e", "", "exact marginals file (required)")

	// Parse and validate arguments
	flag.Parse()
	if len(exact) == 0 {
		fmt.Println("Missing exact marginals file")
		flag.PrintDefaults()
		os.Exit(1)
	}
}

func main() {
	marfs, _ := filepath.Glob("*.mar")
	fmt.Println(marfs)

	mp := utl.CreateFile(fmt.Sprintf("marginals_%v.txt", time.Now().Format(time.RFC3339)))
	defer mp.Close()
	cfuncs := []string{"mse", "entropy", "l1", "l2", "abs", "hel"}
	fmt.Fprintln(mp, utl.Sprintc("marfile", cfuncs))

	for _, marf := range marfs {
		if marf != exact {
			d := []float64(nil)
			for _, v := range cfuncs {
				distanfunc, err := learn.ValidDistanceFunc(v)
				errchk.Check(err, "")
				d = append(d, learn.CompareMarginals(exact, marf, distanfunc))
			}
			fmt.Fprintln(mp, utl.Sprintc(marf, d))
		}
	}
}
