package mrf

import (
	"bufio"
	"fmt"
	"io"
	"strings"

	"github.com/britojr/kbn/factor"
	"github.com/britojr/kbn/utils"
)

// Mrf markov random field
type Mrf struct {
	cardin     []int
	potentials []*factor.Factor
}

// LoadFromUAI creates a mrf from a reader in uai format
func LoadFromUAI(r io.Reader) *Mrf {
	scanner := bufio.NewScanner(r)
	scanner.Scan()
	scanner.Scan()
	// numvar := utils.Atoi(scanner.Text())
	scanner.Scan()
	cardin := utils.SliceAtoi(strings.Fields(scanner.Text()))
	scanner.Scan()
	potentials := make([]*factor.Factor, utils.Atoi(scanner.Text()))
	for i := range potentials {
		scanner.Scan()
		varlist := utils.SliceAtoi(strings.Fields(scanner.Text()))
		potentials[i] = factor.NewFactor(varlist[1:], cardin)
	}
	scanner.Scan()
	for i := range potentials {
		scanner.Scan()
		scanner.Scan()
		potentials[i].SetValues(utils.SliceAtoF64(strings.Fields(scanner.Text())))
		scanner.Scan()
	}
	return &Mrf{cardin, potentials}
}

// UnnormalizedMesure returns the "unnormalized probability" of given evidence
func (m *Mrf) UnnormalizedMesure(evid []int) float64 {
	q := float64(1)
	for _, f := range m.potentials {
		q *= f.GetEvidValue(evid)
	}
	return q
}

// Marginals return a map with all marginals normalized by given partition func z
func (m *Mrf) Marginals(z float64) map[int][]float64 {
	// TODO: fix marginals calculation
	ma := make(map[int][]float64)
	for _, p := range m.potentials {
		for j, v := range p.Variables() {
			if _, ok := ma[v]; !ok {
				q := p.SumOut(p.Variables()[:j])
				q = q.SumOut(p.Variables()[j+1:])
				ma[v] = q.Values()
				for k := range ma[v] {
					ma[v][k] /= z
				}
			}
		}
	}
	return ma
}

// Print prints all mrf values
func (m *Mrf) Print() {
	fmt.Println(m.cardin)
	for _, f := range m.potentials {
		fmt.Println(f.Variables(), f.Values())
	}
}
