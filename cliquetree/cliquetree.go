package cliquetree

import (
	"github.com/britojr/kbn/factor"
	"github.com/britojr/kbn/utils"
	"github.com/willf/bitset"
)

// CliqueTree ..
type CliqueTree struct {
	nodes []node
}

type node struct {
	varset        *bitset.BitSet
	neighbours    []int
	initialPot    *factor.Factor
	calibratedPot *factor.Factor
}

// New ..
func New(n int) *CliqueTree {
	c := new(CliqueTree)
	c.nodes = make([]node, n)
	return c
}

// SetClique ..
func (c *CliqueTree) SetClique(i int, varlist []int) {
	c.nodes[i].varset = bitset.New(0)
	utils.SetFromSlice(c.nodes[i].varset, varlist)
}

// SetNeighbours ..
func (c *CliqueTree) SetNeighbours(i int, neighbours []int) {
	c.nodes[i].neighbours = neighbours
}

// SetPotential ..
func (c *CliqueTree) SetPotential(i int, potential *factor.Factor) {
	c.nodes[i].initialPot = potential
	c.nodes[i].calibratedPot = potential
}

// Calibrated ..
func (c *CliqueTree) Calibrated(i int) *factor.Factor {
	return c.nodes[i].calibratedPot
}

// IterativeCalibration ..
func (c *CliqueTree) IterativeCalibration() {

}
