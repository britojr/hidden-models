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

var send, receive []*factor.Factor

// New ..
func New(n int) *CliqueTree {
	c := new(CliqueTree)
	c.nodes = make([]node, n)
	return c
}

// SetClique ..
func (c *CliqueTree) SetClique(i int, varlist []int) {
	// TODO: give hint size for the bitset
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

// UpDownCalibration ..
func (c *CliqueTree) UpDownCalibration() {
	send = make([]*factor.Factor, len(c.nodes))
	receive = make([]*factor.Factor, len(c.nodes))
	n := len(c.nodes[0].neighbours)
	p := make([]*factor.Factor, n+1)
	q := make([]*factor.Factor, n-1)
	msg := make([]*factor.Factor, n)

	p[0] = c.nodes[0].initialPot
	for i, ch := range c.nodes[0].neighbours {
		msg[i] = c.upwardmessage(ch, 0)
		p[i+1] = p[i].Product(msg[i])
	}
	c.nodes[0].calibratedPot = p[n]

	q[n-2] = msg[n-1]
	for i := n - 3; i >= 0; i-- {
		q[i] = q[i+1].Product(msg[i+1])
	}

	for i := 0; i < n; i++ {
		m := p[i]
		if i != n-1 {
			m = m.Product(q[i])
		}
		ch := c.nodes[0].neighbours[i]
		diff := utils.SetSubtract(c.nodes[0].varset, c.nodes[ch].varset)
		for _, x := range diff {
			m = m.SumOut(x)
		}
		c.downwardmessage(0, ch, m)
	}
}

func (c *CliqueTree) downwardmessage(pa, j int, pm *factor.Factor) {
	n := len(c.nodes[j].neighbours)
	p := make([]*factor.Factor, n+1)
	q := make([]*factor.Factor, n-1)
	msg := make([]*factor.Factor, n)

	p[0] = c.nodes[j].initialPot
	for i, ch := range c.nodes[j].neighbours {
		if ch != pa {
			msg[i] = send[ch]
		} else {
			msg[i] = pm
		}
		p[i+1] = p[i].Product(msg[i])
	}
	c.nodes[j].calibratedPot = p[n]

	if n > 1 {
		q[n-2] = msg[n-1]
		for i := n - 3; i >= 0; i-- {
			q[i] = q[i+1].Product(msg[i+1])
		}
	}

	for i := 0; i < n; i++ {
		ch := c.nodes[j].neighbours[i]
		if ch == pa {
			continue
		}
		m := p[i]
		if i != n-1 {
			m = m.Product(q[i])
		}
		diff := utils.SetSubtract(c.nodes[j].varset, c.nodes[ch].varset)
		for _, x := range diff {
			m = m.SumOut(x)
		}
		c.downwardmessage(j, ch, m)
	}
}

func (c *CliqueTree) upwardmessage(i, pa int) *factor.Factor {
	msg := c.nodes[i].initialPot
	if len(c.nodes[i].neighbours) > 1 {
		for _, ne := range c.nodes[i].neighbours {
			if ne != pa {
				msg = msg.Product(c.upwardmessage(ne, i))
			}
		}
	}
	diff := utils.SetSubtract(c.nodes[i].varset, c.nodes[pa].varset)
	for _, x := range diff {
		msg = msg.SumOut(x)
	}
	send[i] = msg
	return msg
}
