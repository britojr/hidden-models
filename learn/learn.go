package learn

import (
	"fmt"
	"log"
	"os"
	"sort"
	"strings"
	"time"

	"github.com/britojr/kbn/cliquetree"
	"github.com/britojr/kbn/errchk"
	"github.com/britojr/kbn/filehandler"
	"github.com/britojr/kbn/mrf"
)

// SaveCliqueTree saves a clique tree on the given file
func SaveCliqueTree(ct *cliquetree.CliqueTree, fname string) {
	f, err := os.Create(fname)
	errchk.Check(err, fmt.Sprintf("Can't create file %v", fname))
	defer f.Close()
	ct.SaveOn(f)
}

// LoadCliqueTree loads a clique tree from the given file
func LoadCliqueTree(fname string) *cliquetree.CliqueTree {
	f, err := os.Open(fname)
	errchk.Check(err, fmt.Sprintf("Can't open file %v", fname))
	defer f.Close()
	return cliquetree.LoadFrom(f)
}

// LoadMRF loads a MRF from the given file
func LoadMRF(fname string) *mrf.Mrf {
	f, err := os.Open(fname)
	errchk.Check(err, fmt.Sprintf("Can't open file %v", fname))
	defer f.Close()
	return mrf.LoadFromUAI(f)
}

// ExtractData reads a dataset and returns is data and cardinality
func ExtractData(dsfile string, delimiter, hdr uint) (data [][]int, cardin []int) {
	log.Printf("Loading dataset: %v\n", dsfile)
	start := time.Now()
	dataset := filehandler.NewDataSet(dsfile, rune(delimiter), filehandler.HeaderFlags(hdr))
	dataset.Read()
	elapsed := time.Since(start)
	log.Printf("Time: %v\n", elapsed)
	return dataset.Data(), dataset.Cardinality()
}

// Sprintc returns the default formats in a comma-separated string
func Sprintc(a ...interface{}) string {
	s := fmt.Sprintln(a...)
	s = strings.Trim(s, "\n")
	s = strings.Replace(s, " ", ",", -1)
	s = strings.Replace(s, "[", "", -1)
	s = strings.Replace(s, "]", "", -1)
	return s
}

// SaveCTMarginals saves marginals of observed variables of a clique tree in UAI format
func SaveCTMarginals(ct *cliquetree.CliqueTree, obs int, fname string) {
	f, err := os.Create(fname)
	errchk.Check(err, "")
	defer f.Close()
	ma := ct.Marginals()

	var keys []int
	for k := range ma {
		keys = append(keys, k)
	}
	fmt.Fprintf(f, "MAR\n")
	fmt.Fprintf(f, "%d ", obs)
	sort.Ints(keys)
	for i := 0; i < obs; i++ {
		fmt.Fprintf(f, "%d ", len(ma[keys[i]]))
		for _, v := range ma[keys[i]] {
			fmt.Fprintf(f, "%.5f ", v)
		}
	}
}
