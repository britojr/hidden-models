package cliquetree

import "github.com/britojr/kbn/factor"

// CliqueTree ..
type CliqueTree struct {
	nodes []node
}

type node struct {
	varlist       []int          // wich variables participate on this clique
	neighbours    []int          // cliques that are adjacent to this one
	initialPot    *factor.Factor // initial clique potential
	calibratedPot *factor.Factor // calibrated potential
}

var send, receive []*factor.Factor
var prev, post [][]*factor.Factor
var parent []int
var children [][]int

// New ..
func New(n int) *CliqueTree {
	c := new(CliqueTree)
	c.nodes = make([]node, n)
	return c
}

// SetClique ..
func (c *CliqueTree) SetClique(i int, varlist []int) {
	c.nodes[i].varlist = varlist
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
	prev = make([][]*factor.Factor, len(c.nodes))
	post = make([][]*factor.Factor, len(c.nodes))
	root := 0

	c.upwardmessage(root, -1)
	c.downwardmessage(-1, root)
}

func (c *CliqueTree) upwardmessage(v, pa int) {
	prev[v] = make([]*factor.Factor, 1, len(c.nodes[v].neighbours)+1)
	prev[v][0] = c.nodes[v].initialPot
	if len(c.nodes[v].neighbours) > 1 {
		for _, ne := range c.nodes[v].neighbours {
			if ne != pa {
				c.upwardmessage(ne, v)
				prev[v] = append(prev[v], send[ne].Product(prev[v][len(prev[v])-1]))
			}
		}
	}
	if pa != -1 {
		send[v] = prev[v][len(prev[v])-1].Marginalize(c.nodes[pa].varlist)
	}
}

func (c *CliqueTree) downwardmessage(pa, v int) {
	c.nodes[v].calibratedPot = prev[v][len(prev[v])-1]
	n := len(c.nodes[v].neighbours)
	if pa != -1 {
		c.nodes[v].calibratedPot = c.nodes[v].calibratedPot.Product(receive[v])
		n--
	}
	if len(c.nodes[v].neighbours) == 1 && pa != -1 {
		return
	}

	post[v] = make([]*factor.Factor, n)
	i := len(post[v]) - 1
	post[v][i] = receive[v]
	i--
	for k := len(c.nodes[v].neighbours) - 1; k >= 0 && i >= 0; k-- {
		ch := c.nodes[v].neighbours[k]
		if ch == pa {
			continue
		}
		post[v][i] = send[ch]
		if post[v][i+1] != nil {
			post[v][i] = post[v][i].Product(post[v][i+1])
		}
		i--
	}

	k := 0
	for _, ch := range c.nodes[v].neighbours {
		if ch == pa {
			continue
		}
		msg := prev[v][k]
		if post[v][k] != nil {
			msg = msg.Product(post[v][k])
		}
		receive[ch] = msg.Marginalize(c.nodes[ch].varlist)
		c.downwardmessage(v, ch)
		k++
	}
}

// IterativeCalibration ..
func (c *CliqueTree) IterativeCalibration() {
	send = make([]*factor.Factor, len(c.nodes))
	receive = make([]*factor.Factor, len(c.nodes))
	prev = make([][]*factor.Factor, len(c.nodes))
	post = make([][]*factor.Factor, len(c.nodes))
	root := 0
	order := c.bfsOrder(root)
	for i := len(order) - 1; i >= 0; i-- {
		c.upmessage(order[i])
	}
	for _, v := range order {
		c.downmessage(v)
	}
}

func (c *CliqueTree) upmessage(v int) {
	prev[v] = make([]*factor.Factor, 1, len(c.nodes[v].neighbours)+1)
	prev[v][0] = c.nodes[v].initialPot
	for _, ch := range children[v] {
		prev[v] = append(prev[v], send[ch].Product(prev[v][len(prev[v])-1]))
	}
	if parent[v] != -1 {
		send[v] = prev[v][len(prev[v])-1].Marginalize(c.nodes[parent[v]].varlist)
	}
}

func (c *CliqueTree) downmessage(v int) {
	c.nodes[v].calibratedPot = prev[v][len(prev[v])-1]
	if parent[v] != -1 {
		c.nodes[v].calibratedPot = c.nodes[v].calibratedPot.Product(receive[v])
	}
	if len(children[v]) == 0 {
		return
	}
	post[v] = make([]*factor.Factor, len(children[v]))
	i := len(post[v]) - 1
	post[v][i] = receive[v]
	i--
	for ; i >= 0; i-- {
		ch := children[v][i+1]
		post[v][i] = send[ch]
		if post[v][i+1] != nil {
			post[v][i] = post[v][i].Product(post[v][i+1])
		}
	}
	for k, ch := range children[v] {
		msg := prev[v][k]
		if post[v][k] != nil {
			msg = msg.Product(post[v][k])
		}
		receive[ch] = msg.Marginalize(c.nodes[ch].varlist)
	}
}

func (c *CliqueTree) bfsOrder(root int) []int {
	parent = make([]int, len(c.nodes))
	children = make([][]int, len(c.nodes))
	visit := make([]bool, len(c.nodes))
	queue := make([]int, len(c.nodes))
	start, end := 0, 0
	parent[root] = -1
	visit[root] = true
	queue[end] = root
	end++
	for start < end {
		v := queue[start]
		children[v] = make([]int, 0)
		start++
		for _, ne := range c.nodes[v].neighbours {
			if !visit[ne] {
				children[v] = append(children[v], ne)
				parent[ne] = v
				visit[ne] = true
				queue[end] = ne
				end++
			}
		}
	}
	return queue
}
