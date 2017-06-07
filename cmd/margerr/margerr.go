package main

import (
	"flag"
	"fmt"
	"math"
	"os"
	"path/filepath"
	"time"

	"github.com/britojr/kbn/learn"
	"github.com/britojr/kbn/utl"
	"github.com/britojr/kbn/utl/stats"
	"github.com/gonum/floats"
)

const (
	metrMSE = iota
	metrCrossEntropy
	metrInverse
	metrL1
	metrL2
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
	// flag.StringVar(&compfunc, "c", "mse", "compare function{mse|entropy}")

	// Parse and validate arguments
	flag.Parse()
	if len(exact) == 0 {
		fmt.Println("Missing exact marginals file")
		flag.PrintDefaults()
		os.Exit(1)
	}

	metrChoices = map[string]int{
		"mse":     metrMSE,
		"entropy": metrCrossEntropy,
		"inverse": metrInverse,
		"l1":      metrL1,
		"l2":      metrL2,
	}
}

func main() {
	marfs, _ := filepath.Glob("*.mar")
	fmt.Println(marfs)

	mp := utl.CreateFile(fmt.Sprintf("marginals_%v.txt", time.Now().Format(time.RFC3339)))
	defer mp.Close()
	cfuncs := []string{"mse", "entropy", "inverse", "l1", "l2"}
	fmt.Fprintln(mp, utl.Sprintc("marfile", cfuncs))

	for _, marf := range marfs {
		if marf != exact {
			d := []float64(nil)
			for _, v := range cfuncs {
				d = append(d, compare(exact, marf, metrChoices[v]))
			}
			fmt.Fprintln(mp, utl.Sprintc(marf, d))
		}
	}
}

func compare(exact, approx string, compfunc int) (d float64) {
	e, a := learn.LoadMarginals(exact), learn.LoadMarginals(approx)
	switch compfunc {
	case metrMSE:
		d = margMSE(e, a)
	case metrCrossEntropy:
		d = margCrossEntropy(e, a)
	case metrInverse:
		d = margCrossEntropy(a, e)
	case metrL1:
		d = margDistance(a, e, 1)
	case metrL2:
		d = margDistance(a, e, 2)
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
		c += crossEntropy(e[i], a[i])
	}
	return c / float64(len(e))
}

func margDistance(e, a [][]float64, l float64) (c float64) {
	for i := range e {
		c += floats.Distance(e[i], a[i], l)
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

func crossEntropy(xs, ys []float64) (c float64) {
	for i, v := range xs {
		// c -= v * math.Log(ys[i])
		if ys[i] != 0 {
			c -= v * math.Log(ys[i])
		}
	}
	return
}
