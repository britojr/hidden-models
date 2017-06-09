// generate a cliquetree file format from the output of libdai
package main

import (
	"bufio"
	"fmt"
	"io"
	"os"
	"strings"

	"github.com/britojr/kbn/cliquetree"
	"github.com/britojr/kbn/utl"
	"github.com/britojr/kbn/utl/conv"
	"github.com/britojr/kbn/utl/errchk"
)

const (
	cliqueConst = "VarElim result"
	edgesConst  = "Spanning tree"
)

func main() {
	if len(os.Args) < 3 {
		fmt.Printf("\nUsage: %v <lidaifile.out> <filename.ct0>\n\n", os.Args[0])
		os.Exit(1)
	}
	r := utl.OpenFile(os.Args[1])
	w := utl.CreateFile(os.Args[2])
	defer r.Close()
	defer w.Close()
	c := parseToCT(r)
	c.SaveOn(w)
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
