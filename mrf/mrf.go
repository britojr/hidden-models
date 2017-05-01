package mrf

import (
	"bufio"
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
	for i := range potentials {
		scanner.Scan()
		for j := range potentials[i].Values() {
			scanner.Scan()
			potentials[i].Values()[j] = utils.AtoF64(scanner.Text())
		}
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
