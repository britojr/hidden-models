package main

import (
	"flag"
	"fmt"
	"math"
	"os"
	"sort"
	"time"

	"github.com/britojr/kbn/cliquetree"
	"github.com/britojr/kbn/factor"
	"github.com/britojr/kbn/filehandler"
	"github.com/britojr/kbn/learn"
	"github.com/britojr/kbn/likelihood"
	"github.com/britojr/kbn/mrf"
	"github.com/britojr/kbn/utils"
)

// StepFlags type is used to store the flags indicating the steps that should be executed
// type StepFlags byte
const (
	// StructStep indicates execute structure learning step
	StructStep int = 1 << iota
	//ParamStep indicates parameter learning step
	ParamStep
	//InferStep indicates inference step
	InferStep
)

const (
	hiddencard = 2
)

// Define Flag variables
var (
	k          int     // treewidth
	numktrees  int     // number of k-trees to sample
	iterEM     int     // number of EM random restarts
	iterations int     // number of iterations of the whole process
	dsfile     string  // dataset file name
	delimiter  uint    // dataset file delimiter
	hdr        uint    // dataset file header type
	h          int     // number of hidden variables
	initpot    int     // type of initial potential
	check      bool    // validate cliquetree
	epslon     float64 // minimum precision for EM convergence
	alpha      float64 // alpha parameter for dirichlet distribution
	ctfile     string  // cliquetree file
	mkfile     string  // markov random field uai file
	steps      int     // flags indicating what steps to execute
)

var (
	learner *learn.Learner
)

func parseFlags() {
	flag.IntVar(&k, "k", 5, "tree-width")
	flag.IntVar(&numktrees, "numk", 1, "number of ktrees samples")
	flag.IntVar(&iterEM, "iterem", 1, "number of EM iterations")
	flag.IntVar(&iterations, "iterations", 1, "number of iterations of the whole process")
	flag.StringVar(&dsfile, "f", "", "dataset file")
	flag.UintVar(&delimiter, "delimiter", ',', "field delimiter")
	flag.UintVar(&hdr, "hdr", 1, "1- name header, 2- cardinality header")
	flag.IntVar(&h, "h", 0, "hidden variables")
	flag.IntVar(&initpot, "initpot", 0,
		`		0- random values,
		1- empiric + dirichlet,
		2- empiric + random,
		3- empiric + uniform`)
	flag.BoolVar(&check, "check", false, "check tree")
	flag.Float64Var(&epslon, "e", 1e-2, "minimum precision for EM convergence")
	flag.Float64Var(&alpha, "a", 0.5, "alpha parameter for dirichlet distribution")
	flag.StringVar(&ctfile, "c", "", "cliquetree file")
	flag.StringVar(&mkfile, "m", "", "MRF file")
	flag.IntVar(&steps, "steps", StructStep|ParamStep,
		`		step flags:
		1- structure learning,
		2- parameter learning,
		3- inference step`)

	// Parse and validate arguments
	flag.Parse()
	if len(dsfile) == 0 {
		fmt.Println("Please enter dataset file name.")
		return
	}
	fmt.Printf("Args: k=%v, h=%v, initpot=%v\n", k, h, initpot)
	fmt.Printf("eps=%v, alph=%v, numk=%v, iterEM=%v\n", epslon, alpha, numktrees, iterEM)
}

func main() {
	parseFlags()
	initializeLearner()

	var ct *cliquetree.CliqueTree
	var ll float64
	if steps&StructStep > 0 {
		ct, ll = learnStructureAndParamenters()
		fmt.Printf("Best LL: %v\n", ll)
		if len(ctfile) > 0 {
			f, err := os.Create(ctfile)
			utils.ErrCheck(err, fmt.Sprintf("Can't create file %v", ctfile))
			ct.SaveOn(f)
			f.Close()
		}
	} else {
		if len(ctfile) == 0 {
			fmt.Println("Inform a valid clique tree filename")
			return
		}
		f, err := os.Open(ctfile)
		utils.ErrCheck(err, fmt.Sprintf("Can't open file %v", ctfile))
		ct = cliquetree.LoadFrom(f)
		f.Close()
		if steps&ParamStep > 0 {
			ll = learnParameters(ct)
			fmt.Printf("Best LL: %v\n", ll)
		}
	}
	if steps&InferStep > 0 {
		if len(mkfile) == 0 {
			fmt.Println("Please inform a valid MRF file")
			return
		}
		inferenceStep(ct)
	}
}

func inferenceStep(ct *cliquetree.CliqueTree) {
	// read MRF
	fmt.Printf("Loading MRF file: %v\n", mkfile)
	f, err := os.Open(mkfile)
	utils.ErrCheck(err, fmt.Sprintf("Can't create file %v", mkfile))
	mk := mrf.LoadFromUAI(f)
	f.Close()
	if mk == nil {
		fmt.Printf("an error occurred while loading file %v\n", mkfile)
		return
	}

	// inference step
	fmt.Println("Calculating partition function...")
	start := time.Now()
	z := estimatePartitionFunction(ct, mk, learner.Data())
	elapsed := time.Since(start)
	fmt.Printf("Time: %v\n", elapsed)
	fmt.Printf("Partition function (Log): %.8f\n", math.Log(z))
}

func estimatePartitionFunction(ct *cliquetree.CliqueTree, mk *mrf.Mrf, data [][]int) float64 {
	var z, p, phi float64
	for _, m := range data {
		p = ct.ProbOfEvidence(m)
		if p != 0 {
			phi = mk.UnnormalizedMesure(m)
			z += phi / p
		} else {
			panic(fmt.Sprintf("zero probability for evid: %v", m))
		}
	}
	return z / float64(len(data))
}

func initializeLearner() {
	learner = learn.New(k, h, hiddencard, alpha)
	fmt.Printf("Loading dataset: %v\n", dsfile)
	start := time.Now()
	learner.LoadDataSet(dsfile, rune(delimiter), filehandler.HeaderFlags(hdr))
	elapsed := time.Since(start)
	fmt.Printf("Time: %v\n", elapsed)
}

func learnStructure() *cliquetree.CliqueTree {
	fmt.Println("Learning structure...")
	start := time.Now()
	ct := learn.RandomCliqueTree(learner.TotVar(), k)
	elapsed := time.Since(start)
	fmt.Printf("Time: %v\n", elapsed)

	// TODO: remove this structure likelihood
	ll := likelihood.StructLoglikelihood(ct.Cliques(), ct.SepSets(), learner.Counter())
	fmt.Printf("Structure LogLikelihood: %v\n", ll)
	return ct
}

func learnParameters(ct *cliquetree.CliqueTree) float64 {
	fmt.Println("Learning parameters...")
	start := time.Now()
	ll := learner.OptimizeParameters(ct, initpot, iterEM, epslon)
	elapsed := time.Since(start)
	fmt.Printf("Time: %v\n", elapsed)

	// TODO: remove this check
	if check {
		CheckTree(ct)
	}
	return ll
}

func learnStructureAndParamenters() (*cliquetree.CliqueTree, float64) {
	ct := learnStructure()
	ll := learnParameters(ct)
	fmt.Printf("Initial LL: %v\n", ll)
	for i := 1; i < iterations; i++ {
		currct := learnStructure()
		currll := learnParameters(currct)
		fmt.Printf("Current LL: %v\n", currll)
		if currll > ll {
			ct, ll = currct, currll
		}
	}
	return ct, ll
}

// =============================================================================

// SaveMarginals saves all marginals of a cliquetree
func SaveMarginals(ct *cliquetree.CliqueTree, ll float64, fname string) {
	f, err := os.Create(fname)
	utils.ErrCheck(err, "")
	defer f.Close()
	m := ct.Marginals()

	var keys []int
	for k := range m {
		keys = append(keys, k)
	}
	sort.Ints(keys)
	for _, k := range keys {
		fmt.Fprintf(f, "{%d} %v\n", k, m[k])
	}
	fmt.Fprintf(f, "LL=%v\n", ll)
}

// CheckTree ..
func CheckTree(ct *cliquetree.CliqueTree) {
	ct.UpDownCalibration()
	// check if they are uniform
	checkUniform(ct)
	// check if after summing out the hidden variables they are the same as initial count
	checkWithInitialCount(ct)
	// check if tre have any zero factor
	checkCliqueTree(ct)
}

func checkUniform(ct *cliquetree.CliqueTree) {
	fmt.Println("checking uniform...")
	uniform := learn.CreateEmpiricPotentials(learner.Counter(),
		ct.Cliques(), learner.Cardinality(), learner.TotVar()-h, learn.EmpiricUniform)
	// normalize uniform using tree direction
	for i := range uniform {
		if len(ct.Varin(i)) != 0 {
			uniform[i] = uniform[i].Division(uniform[i].SumOut(ct.Varin(i)))
		}
	}
	fmt.Printf("Uniform param: %v (%v)=0\n", uniform[0].Values()[0], uniform[0].Variables())
	calibrated := make([]*factor.Factor, ct.Size())
	for i := range calibrated {
		// calibrated[i] = ct.Calibrated(i)
		calibrated[i] = ct.InitialPotential(i)
	}
	diff, i, j, err := factor.MaxDifference(uniform, calibrated)
	utils.ErrCheck(err, "")
	fmt.Printf("f[%v][%v]=%v; g[%v][%v]=%v\n", i, j, uniform[i].Values()[j], i, j, calibrated[i].Values()[j])
	fmt.Printf(" maxdiff = %v\n", diff)
	if diff > 1e-6 {
		fmt.Println(" >>> Not uniform")
	}
}

func checkWithInitialCount(ct *cliquetree.CliqueTree) {
	fmt.Println("checking count...")
	initialCount := make([]*factor.Factor, ct.Size())
	sumOutHidden := make([]*factor.Factor, ct.Size())
	for i := range initialCount {
		var observed, hidden []int
		if h > 0 {
			observed, hidden = utils.SliceSplit(ct.Clique(i), learner.TotVar()-h)
		} else {
			observed = ct.Clique(i)
		}
		if len(observed) > 0 {
			values := utils.SliceItoF64(learner.Counter().CountAssignments(observed))
			// sumOutHidden[i] = ct.InitialPotential(i)
			sumOutHidden[i] = ct.Calibrated(i)
			if len(hidden) > 0 {
				sumOutHidden[i] = sumOutHidden[i].SumOut(hidden)
			}
			initialCount[i] = factor.NewFactorValues(observed, learner.Cardinality(), values)
			initialCount[i].Normalize()
			// sumOutHidden[i].Normalize()
		}
	}

	if initialCount[0] != nil {
		fmt.Printf("IniCount param: %v (%v)=0\n", initialCount[0].Values()[0], initialCount[0].Variables())
		fmt.Printf("sumOut param: %v (%v)=0\n", sumOutHidden[0].Values()[0], sumOutHidden[0].Variables())
	}
	diff, i, j, err := factor.MaxDifference(initialCount, sumOutHidden)
	utils.ErrCheck(err, "")
	fmt.Printf("f[%v][%v]=%v; g[%v][%v]=%v\n", i, j, initialCount[i].Values()[j], i, j, sumOutHidden[i].Values()[j])
	fmt.Printf(" maxdiff = %v\n", diff)
	if diff > 1e-6 {
		fmt.Println(" >> Different from initial counting")
	}
}

// checkCliqueTree ..
func checkCliqueTree(ct *cliquetree.CliqueTree) {
	printTree := func(f *factor.Factor) {
		fmt.Printf("(%v)\n", f.Variables())
		fmt.Println("tree:")
		for i := 0; i < ct.Size(); i++ {
			fmt.Printf("node %v: neighb: %v clique: %v septset: %v parent: %v\n",
				i, ct.Neighbours(i), ct.Clique(i), ct.SepSet(i), ct.Parents()[i])
		}
		fmt.Println("original potentials:")
		for i := 0; i < ct.Size(); i++ {
			fmt.Printf("node %v:\n var: %v\n values: %v\n",
				i, ct.InitialPotential(i).Variables(), ct.InitialPotential(i).Values())
		}
	}

	for i := range ct.Potentials() {
		f := ct.InitialPotential(i)
		sum := 0.0
		for _, v := range f.Values() {
			sum += v
		}
		if utils.FuzzyEqual(sum, 0) {
			printTree(f)
			panic("original zero factor")
		}
		f = ct.Calibrated(i)
		sum = 0.0
		for _, v := range f.Values() {
			sum += v
		}
		if utils.FuzzyEqual(sum, 0) {
			printTree(f)
			panic("original zero factor")
		}
	}
}
