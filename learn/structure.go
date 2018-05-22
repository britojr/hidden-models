package learn

import (
	"fmt"
	"io"
	"time"

	"github.com/britojr/kbn/cliquetree"
	"github.com/britojr/kbn/dataset"
	"github.com/britojr/kbn/likelihood"
	"github.com/britojr/utl/ioutl"
	"github.com/britojr/utl/stats"
)

// DistanceFunc specifies a distance function
type DistanceFunc string

// defines the available distance functions
const (
	MSE          DistanceFunc = "mse"
	CrossEntropy              = "entropy"
	L1Distance                = "l1"
	L2Distance                = "l2"
	MaxAbsError               = "max-abs"
	MeanAbsError              = "mean-abs"
	Hellinger                 = "hel"
)

var distances = map[string]DistanceFunc{
	"mse":      MSE,
	"entropy":  CrossEntropy,
	"l1":       L1Distance,
	"l2":       L2Distance,
	"max-abs":  MaxAbsError,
	"mean-abs": MeanAbsError,
	"hel":      Hellinger,
}

// String returns the distance functons names
func (d DistanceFunc) String() string { return string(d) }

// ValidDistanceFunc returns a distance function option from a string
func ValidDistanceFunc(a string) (d DistanceFunc, err error) {
	var ok bool
	d, ok = distances[a]
	if !ok {
		err = fmt.Errorf("invalid distance function string: %v", a)
	}
	return
}

// SampleStructure samples a cliquetree structure with limited treewidth
// with number of variables corresponding to the given dataset plus latent variables
func SampleStructure(ds *dataset.Dataset, k, h int, ctfile string) (float64, time.Duration) {
	start := time.Now()
	ct := cliquetree.NewRandom(ds.NCols()+h, k)
	sll := likelihood.StructLoglikelihood(ct.Cliques(), ct.SepSets(), ds)
	elapsed := time.Since(start)
	if len(ctfile) > 0 {
		SaveCliqueTree(ct, ctfile)
	}
	return sll, elapsed
}

// CompareMarginals compares two marginals and return a difference
func CompareMarginals(exact, approx string, dsfunc DistanceFunc) (d float64) {
	e, a := LoadMarginals(exact), LoadMarginals(approx)
	switch dsfunc {
	case MSE:
		d = stats.MatMSE(e, a)
	case CrossEntropy:
		d = stats.MatCrossEntropy(e, a)
	case L1Distance:
		d = stats.MatDistance(e, a, 1)
	case L2Distance:
		d = stats.MatDistance(e, a, 2)
	case MaxAbsError:
		d = stats.MatMaxAbsErr(e, a)
	case MeanAbsError:
		d = stats.MatMeanAbsErr(e, a)
	case Hellinger:
		d = stats.MatHellDist(e, a)
	default:
		panic(fmt.Sprintf("invalid distance function: %v", dsfunc))
	}
	return
}

// SaveMarginas load a cliquetree to calculate and save its marginals
func SaveMarginas(ctfile, marfile string) {
	ct := LoadCliqueTree(ctfile)
	saveCTMarginals(ct, -1, marfile)
}

// LoadMarginals read a MAR file and returns a slice of floats
func LoadMarginals(fname string) [][]float64 {
	f := ioutl.OpenFile(fname)
	defer f.Close()
	return readMarginals(f)
}

// LoadCliqueTree from a file name
func LoadCliqueTree(fname string) *cliquetree.CliqueTree {
	f := ioutl.OpenFile(fname)
	defer f.Close()
	return cliquetree.LoadFrom(f)
}

// SaveCliqueTree on a file
func SaveCliqueTree(ct *cliquetree.CliqueTree, fname string) {
	f := ioutl.CreateFile(fname)
	defer f.Close()
	ct.SaveOn(f)
}

func saveCTMarginals(ct *cliquetree.CliqueTree, obs int, fname string) {
	f := ioutl.CreateFile(fname)
	defer f.Close()
	ma := ct.Marginals()
	if obs > 0 {
		writeMarginals(f, ma[:obs])
	} else {
		writeMarginals(f, ma)
	}
}

func writeMarginals(w io.Writer, ma [][]float64) {
	fmt.Fprintf(w, "MAR\n")
	fmt.Fprintf(w, "%d ", len(ma))
	for i := range ma {
		fmt.Fprintf(w, "%d ", len(ma[i]))
		for _, v := range ma[i] {
			fmt.Fprintf(w, "%v ", v)
		}
	}
}

func readMarginals(r io.Reader) (ma [][]float64) {
	var mar string
	fmt.Fscanln(r, &mar)
	var n int
	fmt.Fscanf(r, "%d", &n)
	ma = make([][]float64, n)
	for i := range ma {
		fmt.Fscanf(r, "%d", &n)
		ma[i] = make([]float64, n)
		for j := range ma[i] {
			fmt.Fscanf(r, "%f", &ma[i][j])
		}
	}
	return
}
