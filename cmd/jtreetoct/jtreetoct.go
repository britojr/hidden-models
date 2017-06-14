// generate a cliquetree file format from the output of libdai
package main

import (
	"bufio"
	"fmt"
	"io"
	"math"
	"os"
	"strings"

	"github.com/britojr/kbn/cliquetree"
	"github.com/britojr/kbn/factor"
	"github.com/britojr/kbn/utl"
	"github.com/britojr/kbn/utl/conv"
	"github.com/britojr/kbn/utl/errchk"
)

const (
	cliqueConst = "VarElim result"
	edgesConst  = "Spanning tree"
	valsConst   = "QaValues"
)

func main() {
	if len(os.Args) < 3 {
		fmt.Printf("\nUsage: %v <lidaifile.out> <filename.ct0> [vals.out]\n\n", os.Args[0])
		os.Exit(1)
	}
	r := utl.OpenFile(os.Args[1])
	w := utl.CreateFile(os.Args[2])
	defer r.Close()
	defer w.Close()
	c := parseToCT(r)
	if len(os.Args) > 3 {
		v := utl.OpenFile(os.Args[3])
		defer v.Close()
		parseValues(v, c)
	}
	c.SaveOn(w)
}

func parseValues(r io.Reader, c *cliquetree.CliqueTree) {
	cardin := make([]int, c.N())
	for i := range cardin {
		cardin[i] = 2
	}
	pot := make([]*factor.Factor, len(c.Cliques()))
	num := int(math.Pow(2, 21))
	var aux rune
	for i := 0; i < len(c.Cliques()); i++ {
		fmt.Fscanf(r, "%c", &aux)
		fmt.Printf("(%c)\n", aux)
		pot[i] = factor.NewFactor(c.Clique(i), cardin)
		if len(pot[i].Values()) < num {
			fmt.Println(len(pot[i].Values()))
			panic("invalid size")
		}
		for j := range pot[i].Values() {
			// fmt.Fscanf(r, "%f", &pot[i].Values()[j])
			_, err := fmt.Fscanf(r, "%f", &pot[i].Values()[j])
			errchk.Check(err, "")
			if j < 2 || j >= num-2 {
				fmt.Println(pot[i].Values()[j])
			}
		}
		fmt.Fscanf(r, "%c", &aux)
		fmt.Printf("(%c)\n", aux)
	}
	c.SetAllPotentials(pot)
}

func parseToCT(r io.Reader) *cliquetree.CliqueTree {
	m := findLines(r, []string{cliqueConst, edgesConst})
	cliques := parseClique(m[cliqueConst])
	adj := parseAdj(m[edgesConst], len(cliques))
	c, err := cliquetree.NewStructure(cliques, adj)
	errchk.Check(err, "")
	return c
}

func findLines(r io.Reader, prefixes []string) map[string]string {
	found := make(map[string]string)
	scanner := bufio.NewScanner(r)
	for scanner.Scan() {
		for _, prefix := range prefixes {
			if strings.HasPrefix(scanner.Text(), prefix) {
				found[prefix] = scanner.Text()
				break
			}
		}
	}
	return found
}

func parseClique(s string) [][]int {
	r := strings.Trim(s, cliqueConst+": (){}")
	r = strings.Replace(r, ",", "", -1)
	r = strings.Replace(r, "}", "", -1)
	r = strings.Replace(r, "x", "", -1)
	rs := strings.FieldsFunc(r, func(c rune) bool {
		return c == '{'
	})
	m := make([][]int, len(rs))
	for i, r := range rs {
		m[i] = conv.Satoi(strings.Fields(r))
	}
	return m
}

func parseAdj(s string, n int) [][]int {
	r := strings.Trim(s, edgesConst+": ()")
	r = strings.Replace(r, ",", "", -1)
	r = strings.Replace(r, ")", "", -1)
	r = strings.Replace(r, "->", " ", -1)
	rs := strings.FieldsFunc(r, func(c rune) bool {
		return c == '('
	})
	m := make([][]int, n)
	for _, r := range rs {
		e := conv.Satoi(strings.Fields(r))
		m[e[0]] = append(m[e[0]], e[1])
		m[e[1]] = append(m[e[1]], e[0])
	}
	return m
}
