package cliquetree

import (
	"reflect"
	"sort"
	"testing"

	"github.com/britojr/kbn/assignment"
	"github.com/britojr/kbn/factor"
	"github.com/britojr/kbn/utils"
	"github.com/britojr/tcc/characteristic"
)

type factorStruct struct {
	varlist []int
	values  []float64
}

//                 A, B, C, D, E, F, G, H, I, J, K, L
var cardin = []int{2, 2, 2, 2, 2, 2, 2, 2, 2, 2, 2, 2}
var fAB = factorStruct{
	varlist: []int{0, 1},
	values:  []float64{0.5, 0.1, 0.7, 0.2},
}
var fABC = factorStruct{
	varlist: []int{0, 1, 2},
	values:  []float64{0.5, 0.1, 0.3, 0.8, 0.0, 0.9, 0.6, 0.7},
}
var fABD = factorStruct{
	varlist: []int{0, 1, 3},
	values:  []float64{0.2, 0.5, 0.5, 0.8, 0.1, 0.9, 0.3, 0.7},
}
var fABE = factorStruct{
	varlist: []int{0, 1, 4},
	values:  []float64{0.9, 0.3, 0.12, 0.48, 0.19, 0.98, 0.1, 0.7},
}
var fBDF = factorStruct{
	varlist: []int{1, 3, 5},
	values:  []float64{0.8, 0.1, 0.23, 0.85, 0.5, 0.45, 0.76, 0.12},
}
var fADG = factorStruct{
	varlist: []int{0, 3, 6},
	values:  []float64{0.2, 0.75, 0.66, 0.2, 0.41, 0.19, 0.3, 0.7},
}

var fBFH = factorStruct{
	varlist: []int{1, 5, 7},
	values:  []float64{0.2, 0.75, 0.66, 0.2, 0.41, 0.19, 0.3, 0.7},
}
var fBFL = factorStruct{
	varlist: []int{1, 5, 11},
	values:  []float64{0.2, 0.75, 0.66, 0.2, 0.41, 0.19, 0.3, 0.7},
}
var fFHL = factorStruct{
	varlist: []int{5, 7, 11},
	values:  []float64{0.2, 0.75, 0.66, 0.2, 0.41, 0.19, 0.3, 0.7},
}
var fAGI = factorStruct{
	varlist: []int{0, 6, 8},
	values:  []float64{0.2, 0.75, 0.66, 0.2, 0.41, 0.19, 0.3, 0.7},
}
var fADI = factorStruct{
	varlist: []int{0, 3, 8},
	values:  []float64{0.2, 0.75, 0.66, 0.2, 0.41, 0.19, 0.3, 0.7},
}
var fADGI = factorStruct{
	varlist: []int{0, 3, 6, 8},
	values:  []float64{0.2, 0.75, 0.66, 0.2, 0.41, 0.19, 0.3, 0.7, 0.2, 0.75, 0.66, 0.2, 0.41, 0.19, 0.3, 0.7},
}
var fADGK = factorStruct{
	varlist: []int{0, 3, 6, 10},
	values:  []float64{0.2, 0.75, 0.66, 0.2, 0.41, 0.19, 0.3, 0.7, 0.2, 0.75, 0.66, 0.2, 0.41, 0.19, 0.3, 0.7},
}
var fAGIK = factorStruct{
	varlist: []int{0, 6, 8, 10},
	values:  []float64{0.2, 0.75, 0.66, 0.2, 0.41, 0.19, 0.3, 0.7, 0.2, 0.75, 0.66, 0.2, 0.41, 0.19, 0.3, 0.7},
}

var factorList = []factorStruct{fAB, fABC, fABD, fABE, fBDF, fADG}
var varinX = [][]int{[]int(nil), []int{2}, []int{3}, []int{4}, []int{5}, []int{6}}
var varoutX = [][]int{[]int(nil), []int(nil), []int(nil), []int(nil), []int{0}, []int{1}}
var adjList = [][]int{
	[]int{1, 2, 3},
	[]int{0},
	[]int{0, 4, 5},
	[]int{0},
	[]int{2},
	[]int{2},
}

var cal []*factor.Factor

var benchTest = []struct {
	fl  []factorStruct
	adj [][]int
}{
	{factorList, adjList},
	{
		[]factorStruct{
			fABD, fBFH, fBFL, fFHL, fAGI, fADI, fADG,
			fBDF, fAB, fABC, fADGI, fABE, fADGK, fAGIK},
		[][]int{
			[]int{8, 7, 6},
			[]int{7, 2, 3},
			[]int{1},
			[]int{1},
			[]int{6},
			[]int{6},
			[]int{4, 5, 10, 0},
			[]int{0, 1},
			[]int{0, 9, 11},
			[]int{8},
			[]int{12, 13, 6},
			[]int{8},
			[]int{10},
			[]int{10},
		},
	},
}

func calculateCalibrated() {
	cal = make([]*factor.Factor, len(factorList))
	p := make([]*factor.Factor, 0)
	for _, f := range factorList {
		p = append(p, factor.NewFactorValues(f.varlist, cardin, f.values))
	}
	//AB
	cal[0] = p[0].Product(p[1].SumOutOne(2)).
		Product(p[2].Product(p[4].SumOutOne(5).Product(p[5].SumOutOne(6))).SumOutOne(3)).
		Product(p[3].SumOutOne(4))
	//ABC
	cal[1] = p[1].Product(p[0]).
		Product(p[2].Product(p[4].SumOutOne(5).Product(p[5].SumOutOne(6))).SumOutOne(3)).
		Product(p[3].SumOutOne(4))
	//ABD
	cal[2] = p[2].Product(p[4].SumOutOne(5)).
		Product(p[5].SumOutOne(6)).
		Product(p[0].Product(p[1].SumOutOne(2).Product(p[3].SumOutOne(4))))
	//ABE
	cal[3] = p[3].Product(p[0]).
		Product(p[1].SumOutOne(2)).
		Product(p[2].Product(p[4].SumOutOne(5).Product(p[5].SumOutOne(6))).SumOutOne(3))
	//BDF
	cal[4] = p[4].
		Product(p[2].Product(p[5].SumOutOne(6).Product(p[0]).Product(p[1].SumOutOne(2)).Product(p[3].SumOutOne(4))).SumOutOne(0))
	//ADG
	cal[5] = p[5].
		Product(p[2].Product(p[4].SumOutOne(5).Product(p[0]).Product(p[1].SumOutOne(2)).Product(p[3].SumOutOne(4))).SumOutOne(1))
}

func initCliqueTree(factorList []factorStruct, adjList [][]int) *CliqueTree {
	c := New(len(factorList))
	c.varin = varinX
	c.varout = varoutX
	for i, f := range factorList {
		c.SetClique(i, f.varlist)
		c.SetNeighbours(i, adjList[i])
		c.SetPotential(i, factor.NewFactorValues(f.varlist, cardin, f.values))
	}
	return c
}

func TestNew(t *testing.T) {
	c := New(len(factorList))
	for i, f := range factorList {
		c.SetClique(i, f.varlist)
		c.SetNeighbours(i, adjList[i])
		c.SetPotential(i, factor.NewFactorValues(f.varlist, cardin, f.values))
	}
}

func TestNewStructure(t *testing.T) {
	cases := []struct {
		cliques, adj, sepsets, varin, varout [][]int
		parents                              []int
		err                                  error
	}{{
		cliques: [][]int{{0, 1}, {0, 2}, {2, 3}, {2, 4}},
		adj:     [][]int{{1}, {0, 2, 3}, {1}, {1}},
		sepsets: [][]int{nil, {0}, {2}, {2}},
		parents: []int{-1, 0, 1, 1},
		varin:   [][]int{nil, {2}, {3}, {4}},
		varout:  [][]int{nil, {1}, {0}, {0}},
	}}
	for _, tt := range cases {
		c, err := NewStructure(tt.cliques, tt.adj)
		if tt.err != err {
			t.Errorf("wrong err, want %v, got %v", tt.err, err)
		}
		if tt.err == nil {
			for i := range tt.cliques {
				sort.Ints(tt.cliques[i])
				sort.Ints(tt.sepsets[i])
				if !reflect.DeepEqual(tt.cliques[i], c.Clique(i)) {
					t.Errorf("wrong clique, want %v, got %v", tt.cliques[i], c.Clique(i))
				}
				if !reflect.DeepEqual(tt.sepsets[i], c.SepSet(i)) {
					t.Errorf("wrong sepset, want %v, got %v", tt.sepsets[i], c.SepSet(i))
				}
				if !reflect.DeepEqual(tt.adj[i], c.Neighbours(i)) {
					t.Errorf("wrong adj, want %v, got %v", tt.adj[i], c.Neighbours(i))
				}
				if !reflect.DeepEqual(tt.varin[i], c.varin[i]) {
					t.Errorf("wrong varin, want %v, got %v", tt.varin[i], c.varin[i])
				}
				if !reflect.DeepEqual(tt.varout[i], c.varout[i]) {
					t.Errorf("wrong varout, want %v, got %v", tt.varout[i], c.varout[i])
				}
			}
			if !reflect.DeepEqual(tt.parents, c.Parents()) {
				t.Errorf("wrong parents, want %v, got %v", tt.parents, c.Parents())
			}
		}
	}
}

func TestOrderedSliceDiff(t *testing.T) {
	cases := []struct {
		a, b, inter, in, out []int
	}{{
		a:     []int{2, 3, 4},
		b:     []int{2, 4, 5},
		inter: []int{2, 4},
		in:    []int{5},
		out:   []int{3},
	}, {
		a:     []int{5, 6, 7},
		b:     []int{2, 4, 5},
		inter: []int{5},
		in:    []int{2, 4},
		out:   []int{6, 7},
	}}
	for _, tt := range cases {
		inter, in, out := orderedSliceDiff(tt.a, tt.b)
		if !reflect.DeepEqual(tt.inter, inter) {
			t.Errorf("wrong inter, want %v, got %v", tt.inter, inter)
		}
		if !reflect.DeepEqual(tt.in, in) {
			t.Errorf("wrong in, want %v, got %v", tt.in, in)
		}
		if !reflect.DeepEqual(tt.out, out) {
			t.Errorf("wrong out, want %v, got %v", tt.out, out)
		}
	}
}

func TestUpDownCalibration(t *testing.T) {
	c := initCliqueTree(factorList, adjList)
	c.UpDownCalibration()
	calculateCalibrated()
	for i, f := range cal {
		got := c.Calibrated(i)
		assig := assignment.New(f.Variables(), cardin)
		for assig.Next() {
			u := f.Get(assig)
			v := got.Get(assig)
			if !utils.FuzzyEqual(u, v) {
				t.Errorf("F[%v][%v]: want(%v); got(%v)", i, assig, u, v)
			}
		}
	}
}

func TestUpDownCalibration2(t *testing.T) {
	cases := []struct {
		cliques, adj [][]int
		cardin       []int
		values       [][]float64
		result       [][]float64
	}{{
		cliques: [][]int{{0, 1}, {1, 2}},
		adj:     [][]int{{1}, {0}},
		cardin:  []int{2, 2, 2},
		values: [][]float64{
			{.25, .35, .35, .5},
			{.20, .40, .22, .18},
		},
		result: [][]float64{
			{.15, .14, .21, .02},
			{.12, .24, .088, .072},
		},
	}}
	for _, tt := range cases {
		c, err := NewStructure(tt.cliques, tt.adj)
		if err != nil {
			t.Errorf(err.Error())
		}
		potentials := make([]*factor.Factor, len(tt.values))
		for i, v := range tt.values {
			potentials[i] = factor.NewFactorValues(tt.cliques[i], tt.cardin, v)
		}
		err = c.SetAllPotentials(potentials)
		if err != nil {
			t.Errorf(err.Error())
		}
		c.UpDownCalibration()
		for i := range tt.cliques {
			if !reflect.DeepEqual(tt.result[i], c.Calibrated(i).Values()) {
				t.Errorf("wrong values for clique %v, want %v, got %v", tt.cliques[i], tt.result[i], c.Calibrated(i).Values())
			}
		}
	}

}

func BenchmarkUpDownCalibration(b *testing.B) {
	ctrees := make([]*CliqueTree, 0)
	for _, bt := range benchTest {
		ctrees = append(ctrees, initCliqueTree(bt.fl, bt.adj))
	}
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		for _, c := range ctrees {
			c.UpDownCalibration()
		}
	}
}

var testFromCharTree = []struct {
	iphi             []int
	chartree         characteristic.Tree
	cliques, sepsets [][]int
	adj              [][]int
	parent           []int
	varin            [][]int
	varout           [][]int
}{
	{
		iphi: []int{0, 10, 9, 3, 4, 5, 6, 7, 1, 2, 8},
		chartree: characteristic.Tree{
			P: []int{-1, 5, 0, 0, 2, 8, 8, 1, 0},
			L: []int{-1, 2, -1, -1, 0, 2, 1, 2, -1},
		},
		cliques: [][]int{
			{1, 2, 8},
			{0, 1, 4, 7},
			{1, 2, 8, 10},
			{1, 2, 8, 9},
			{2, 3, 8, 10},
			{1, 2, 4, 7},
			{1, 5, 7, 8},
			{0, 4, 6, 7},
			{1, 2, 7, 8},
		},
		sepsets: [][]int{
			[]int(nil),
			{1, 4, 7},
			{1, 2, 8},
			{1, 2, 8},
			{2, 8, 10},
			{1, 2, 7},
			{1, 7, 8},
			{0, 4, 7},
			{1, 2, 8},
		},
		adj: [][]int{
			{2, 3, 8},
			{7, 5},
			{4, 0},
			{0},
			{2},
			{1, 8},
			{8},
			{1},
			{5, 6, 0},
		},
		parent: []int{-1, 5, 0, 0, 2, 8, 8, 1, 0},
		varin:  [][]int{[]int(nil), []int{0}, []int{10}, []int{9}, []int{3}, []int{4}, []int{5}, []int{6}, []int{7}},
		varout: [][]int{[]int(nil), []int{2}, []int(nil), []int(nil), []int{1}, []int{8}, []int{2}, []int{1}, []int(nil)},
	},
}

func TestFromCharTree(t *testing.T) {
	for _, v := range testFromCharTree {
		got := FromCharTree(&v.chartree, v.iphi)
		for i := 0; i < got.Size(); i++ {
			if !reflect.DeepEqual(got.Clique(i), v.cliques[i]) {
				t.Errorf("Clique[%v]; Got: %v; Want: %v", i, got.Clique(i), v.cliques[i])
			}
			if !reflect.DeepEqual(got.SepSet(i), v.sepsets[i]) {
				t.Errorf("Sepset[%v]; Got: %v; Want: %v", i, got.SepSet(i), v.sepsets[i])
			}
			if !reflect.DeepEqual(got.neighbours[i], v.adj[i]) {
				t.Errorf("Adj[%v]; Got: %v; Want: %v", i, got.neighbours[i], v.adj[i])
			}
			if got.parent[i] != v.parent[i] {
				t.Errorf("parent[%v]; Got: %v; Want: %v", i, got.parent[i], v.parent[i])
			}
			if !reflect.DeepEqual(got.varin[i], v.varin[i]) {
				t.Errorf("varin[%v]; Got: %v; Want: %v", i, got.varin[i], v.varin[i])
			}
			if !reflect.DeepEqual(got.varout[i], v.varout[i]) {
				t.Errorf("varout[%v]; Got: %v; Want: %v", i, got.varout[i], v.varout[i])
			}
		}
	}
}

func TestReduceByEvidence(t *testing.T) {
	cases := []struct {
		n          int
		potentials []*factor.Factor
		part       []struct {
			evidence []int
			reduced  [][]float64
		}
	}{
		{
			n: 2,
			potentials: []*factor.Factor{
				factor.NewFactorValues([]int{0, 1}, []int{2, 2, 2}, []float64{.25, .10, .35, .30}),
				factor.NewFactorValues([]int{1, 2}, []int{2, 2, 2}, []float64{.40, .20, .10, .30}),
			},
			part: []struct {
				evidence []int
				reduced  [][]float64
			}{
				{
					[]int{0, 1},
					[][]float64{{0, 0, .35, 0}, {0, .20, 0, .30}},
				},
				{
					[]int{1, 0},
					[][]float64{{0, .10, 0, 0}, {.40, 0, .10, 0}},
				},
				{
					[]int{1, 0, 1},
					[][]float64{{0, .10, 0, 0}, {0, 0, .10, 0}},
				},
			},
		},
	}
	for _, tt := range cases {
		c := New(tt.n)
		c.SetAllPotentials(tt.potentials)
		for k := range tt.part {
			c.StorePotentials()
			c.ReduceByEvidence(tt.part[k].evidence)
			for i, f := range c.initialPotStored {
				if !reflect.DeepEqual(tt.potentials[i].Values(), f.Values()) {
					t.Errorf("Original potential changed, want %v, got %v", tt.potentials[i].Values(), f.Values())
				}
				if !reflect.DeepEqual(tt.part[k].reduced[i], c.InitialPotential(i).Values()) {
					t.Errorf("Wrong reduction, want %v, got %v", tt.part[k].reduced[i], c.InitialPotential(i).Values())
				}
			}
			c.RecoverPotentials()
			for i := range tt.potentials {
				if !reflect.DeepEqual(tt.potentials[i].Values(), c.InitialPotential(i).Values()) {
					t.Errorf("Wrong recover, want %v, got %v", tt.potentials[i].Values(), c.InitialPotential(i).Values())
				}
			}
		}
	}
}
