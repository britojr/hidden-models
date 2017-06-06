package main

import (
	"flag"
	"fmt"
	"os"
	"path/filepath"

	"github.com/britojr/kbn/learn"
	"github.com/britojr/kbn/utl/stats"
)

var (
	exact, compfunc string
)

func init() {
	parseFlags()
}

func parseFlags() {
	flag.StringVar(&exact, "e", "", "exact marginals file (required)")
	// flag.StringVar(&compfunc, "c", "mse", "compare function{mse|entropy}")

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
	fmt.Printf("marfile,mse,cross_entropy,variance\n")
	for _, marf := range marfs {
		if marf != exact {
			fmt.Printf("%v,%v,%v,%v",
				marf,
				compare(exact, marf, 0),
				compare(exact, marf, 1),
				compare(exact, marf, 1),
			)
		}
	}
}

func compare(exact, approx string, compfunc int) (d float64) {
	e, a := learn.LoadMarginals(exact), learn.LoadMarginals(approx)
	switch compfunc {
	case 0:
		d = margMSE(e, a)
	case 1:
		d = margCrossEntropy(e, a)
	case 2:
		d = margVariance(e, a)
	}
	return
}

func margMSE(e, a [][]float64) (mse float64) {
	for i := range e {
		mse += stats.MSE(e[i], a[i])
	}
	return mse / float64(len(e))
}

func margCrossEntropy(e, a [][]float64) (c float64) {
	for i := range e {
		c += stats.CrossEntropy(e[i], a[i])
	}
	return c / float64(len(e))
}

func margVariance(a, b [][]float64) float64 {
	c, d := 0, float64(0)
	for i := range a {
		for j, v := range a[i] {
			d += (v - b[i][j]) * (v - b[i][j])
			c++
		}
	}
	return d / float64(c)
}
