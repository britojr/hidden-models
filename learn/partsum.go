package learn

import (
	"fmt"
	"math"
	"os"
	"sort"
	"time"

	"github.com/britojr/kbn/cliquetree"
	"github.com/britojr/kbn/errchk"
	"github.com/britojr/kbn/mrf"
	"github.com/britojr/kbn/stats"
	"github.com/gonum/floats"
)

// PartsumCommand learns an approximation of the partition sum of a MRF
// using inference on a cliquetree
func PartsumCommand(
	dsfile string, delim, hdr uint,
	ctfile, mkfile, zfile string, discard float64,
) {
	data, _ := ExtractData(dsfile, delim, hdr)
	ct := LoadCliqueTree(ctfile)
	mk := LoadMRF(mkfile)

	start := time.Now()
	zs := estimatePartsum(ct, mk, data)
	elapsed := time.Since(start)

	zm := parsumStats(zs, discard)
	if len(zfile) > 0 {
		SavePartsum(zm, zfile)
	}
	fmt.Println(Sprintc(dsfile, ctfile, zfile, zm, discard, elapsed))
}

// SavePartsum saves the partsum estimates
func SavePartsum(zs []float64, fname string) {
	f, err := os.Create(fname)
	errchk.Check(err, fmt.Sprintf("Can't create file %v", fname))
	defer f.Close()
	fmt.Fprint(f, Sprintc(zs))
}

func estimatePartsum(ct *cliquetree.CliqueTree, mk *mrf.Mrf, data [][]int) []float64 {
	zs := make([]float64, len(data))
	for i, m := range data {
		p := ct.ProbOfEvidence(m)
		if p == 0 {
			panic(fmt.Sprintf("zero probability for evidence: %v", m))
		}
		zs[i] = mk.UnnormLogProb(m) - math.Log(p)
	}
	return zs
}

// parsumStats receives a slice of approximations of z
// and returns SD, Mean, Median, Mode, Min, Max
func parsumStats(zs []float64, d float64) []float64 {
	if d < 0 || d >= 0.5 {
		panic(fmt.Sprintf("invalid discard factor: %v", d))
	}
	a, b := int(float64(len(zs))*d), int(len(zs)+1-int(float64(len(zs))*d))
	sort.Float64s(zs)
	ws := zs[a:b]

	zm := []float64{
		stats.Stdev(ws), stats.Mean(ws), stats.Median(ws),
		stats.Mode(ws), floats.Min(ws), floats.Max(ws),
	}
	return zm
}
