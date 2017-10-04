package learn

import (
	"fmt"
	"math"
	"os"
	"sort"
	"time"

	"github.com/britojr/kbn/cliquetree"
	"github.com/britojr/kbn/dataset"
	"github.com/britojr/kbn/mrf"
	"github.com/britojr/utl/errchk"
	"github.com/britojr/utl/floats"
	"github.com/britojr/utl/ioutl"
	"github.com/britojr/utl/stats"
)

// PartitionSum calculates an approximation of the partition sum of a MRF
// using inference on a cliquetree
func PartitionSum(
	ds *dataset.Dataset, ctfile, mkfile, zfile string, discard float64,
) ([]float64, time.Duration) {
	ct := LoadCliqueTree(ctfile)
	mk := loadMRF(mkfile)

	start := time.Now()
	zs := estimatePartsum(ct, mk, ds.Data())
	elapsed := time.Since(start)

	zm := partsumStats(zs, discard)
	if len(zfile) > 0 {
		savePartsum(zm, zfile)
	}
	return zm, elapsed
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

// partsumStats receives a slice of approximations of z
// and returns SD, Mean, Median, Mode, Min, Max
func partsumStats(zs []float64, d float64) []float64 {
	if d < 0 || d >= 0.5 {
		panic(fmt.Sprintf("invalid discard factor: %v", d))
	}
	a, b := int(float64(len(zs))*d), int(len(zs)-int(float64(len(zs))*d)-1)
	sort.Float64s(zs)
	ws := zs[a:b]

	zm := []float64{
		stats.Stdev(ws), stats.Mean(ws), stats.Median(ws),
		stats.Mode(ws), floats.Min(ws), floats.Max(ws),
	}
	return zm
}

func loadMRF(fname string) *mrf.Mrf {
	f := ioutl.OpenFile(fname)
	defer f.Close()
	return mrf.LoadFromUAI(f)
}

func savePartsum(zs []float64, fname string) {
	f, err := os.Create(fname)
	errchk.Check(err, fmt.Sprintf("Can't create file %v", fname))
	defer f.Close()
	fmt.Fprint(f, ioutl.Sprintc(zs))
}
