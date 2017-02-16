package cliquetree

import (
	"github.com/britojr/kbn/factor"
	"github.com/willf/bitset"
)

// CliqueTree ..
type CliqueTree []struct {
	varset        *bitset.BitSet
	neighbours    []int
	initialPot    *factor.Factor
	calibratedPot *factor.Factor
}

// New ..
func New(n int) *CliqueTree {

}

// SetClique ..
func SetClique(i int, varlist []int) {

}

// SetNeighbours ..
func SetNeighbours(i int, neighbours []int) {

}

// SetPotential ..
func SetPotential(i int, potential *factor.Factor) {

}
