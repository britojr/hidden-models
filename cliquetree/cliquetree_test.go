package cliquetree

import (
	"bytes"
	"fmt"
	"reflect"
	"sort"
	"strings"
	"testing"

	"github.com/britojr/kbn/assignment"
	"github.com/britojr/kbn/factor"
	"github.com/britojr/kbn/floats"
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
			if !floats.AlmostEqual(u, v) {
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
			{.25, .35, .35, .05},
			{.20, .22, .40, .18},
		},
		result: [][]float64{
			{.15, .21, .14, .02},
			{.12, .088, .24, .072},
		},
	}, {
		cliques: [][]int{{0}, {1}, {0, 1, 2}, {2, 3}, {2, 4}},
		adj:     [][]int{{2}, {2}, {0, 1, 3, 4}, {2}, {2}},
		cardin:  []int{2, 2, 2, 2, 2},
		values: [][]float64{
			{.999, .001},
			{.998, .002},
			{.999, .06, .71, .05, .001, .94, .29, .95},
			{.95, .10, .05, .90},
			{.99, .30, .01, .70},
		},
		result: [][]float64{
			{.999, .001},
			{.998, .002},
			{.9960050, .0000599, .0014186, .0000001, .0009970, .0009381, .0005794, .0000019},
			{.9476094, .0002516, .0498741, .0022648},
			{.9875088, .0007549, .0099748, .0017615},
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
			for j, x := range tt.result[i] {
				if !floats.AlmostEqual(x, c.Calibrated(i).Values()[j], 1e-5) {
					t.Errorf("wrong values for clique %v, want %v, got %v", tt.cliques[i], tt.result[i], c.Calibrated(i).Values())
				}
			}
			pa := c.Parents()[i]
			if pa != -1 {
				v1 := c.Calibrated(i).SumOut(c.varin[i]).Values()
				v2 := c.Calibrated(pa).SumOut(c.varout[i]).Values()
				for j, x := range c.CalibratedSepSet(i).Values() {
					if !floats.AlmostEqual(x, v1[j]) {
						t.Errorf("wrong values for sepset of clique %v, want %v, got %v", tt.cliques[i], x, v1[j])
					}
					if !floats.AlmostEqual(x, v2[j]) {
						t.Errorf("wrong values for sepset of clique %v, want %v, got %v", tt.cliques[i], x, v2[j])
					}
				}
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

func TestProbOfEvidence(t *testing.T) {
	cases := []struct {
		cliques, adj [][]int
		cardin       []int
		values       [][]float64
		result       []struct {
			evidence []int
			prob     float64
		}
	}{{
		cliques: [][]int{{0}, {1}, {0, 1, 2}, {2, 3}, {2, 4}},
		adj:     [][]int{{2}, {2}, {0, 1, 3, 4}, {2}, {2}},
		cardin:  []int{2, 2, 2, 2, 2},
		values: [][]float64{
			{.999, .001},
			{.998, .002},
			{.999, .06, .71, .05, .001, .94, .29, .95},
			{.95, .10, .05, .90},
			{.99, .30, .01, .70},
		},
		result: []struct {
			evidence []int
			prob     float64
		}{{
			evidence: []int{0, 0, 0},
			prob:     .996004998,
		}, {
			evidence: []int{},
			prob:     1,
		}, {
			evidence: []int{0},
			prob:     .99899999,
		}, {
			evidence: []int{1, 1, 1},
			prob:     .0000019,
		}, {
			evidence: []int{1, 0, 1, 1},
			prob:     .000844308,
		}, {
			evidence: []int{1, 0, 1},
			prob:     .00093812,
		}, {
			evidence: []int{0, 0, 1},
			prob:     .000997002,
		}, {
			evidence: []int{0, 1, 1},
			prob:     .00057942,
		}, {
			evidence: []int{0, 1, 1, 1},
			prob:     .000521478,
		}, {
			evidence: []int{0, 1, 1, 1, 1},
			prob:     .0003650346,
		}},
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
		for _, r := range tt.result {
			got := c.ProbOfEvidence(r.evidence)
			if !floats.AlmostEqual(r.prob, got, 1e-7) {
				t.Errorf("wrong prob of evidence %v, want %v, got %v", r.evidence, r.prob, got)
			}
		}
	}
}

func TestConditionalProb(t *testing.T) {
	cases := []struct {
		cliques, adj [][]int
		cardin       []int
		values       [][]float64
		result       []struct {
			evidence []int
			prob     [][]float64
		}
	}{{
		cliques: [][]int{{0}, {1}, {0, 1, 2}, {2, 3}, {2, 4}},
		adj:     [][]int{{2}, {2}, {0, 1, 3, 4}, {2}, {2}},
		cardin:  []int{2, 2, 2, 2, 2},
		values: [][]float64{
			{.999, .001},
			{.998, .002},
			{.999, .06, .71, .05, .001, .94, .29, .95},
			{.95, .10, .05, .90},
			{.99, .30, .01, .70},
		},
		result: []struct {
			evidence []int
			prob     [][]float64
		}{{
			evidence: []int{},
			prob: [][]float64{
				{.999, .001},
				{.998, .002},
				{.9975, .0025},
				{.9479, .0521},
				{.9883, .0117},
			},
		}, {
			evidence: []int{0, 0, 0},
			prob: [][]float64{
				{1, 0},
				{1, 0},
				{1, 0},
				{.9491, .0509},
				{.9893, .0107},
			},
		}, {
			evidence: []int{0, 1, 1, 1},
			prob: [][]float64{
				{1, 0},
				{0, 1},
				{0, 1},
				{0, 1},
				{.3, .7},
			},
		}, {
			evidence: []int{0, 1, 1, 1, 1},
			prob: [][]float64{
				{1, 0},
				{0, 1},
				{0, 1},
				{0, 1},
				{0, 1},
			},
		}},
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
		c.StorePotentials()
		for _, r := range tt.result {
			c.ReduceByEvidence(r.evidence)
			c.UpDownCalibration()
			for i := 0; i < c.Size(); i++ {
				f := c.Calibrated(i).Normalize()
				for _, v := range c.Calibrated(i).Variables() {
					if v != i {
						f = f.SumOutOne(v)
					}
				}
				for j := range r.prob[i] {
					if !floats.AlmostEqual(r.prob[i][j], f.Values()[j], 1e-3) {
						t.Errorf("wrong probabilities for evid %v,  want %v, got %v", r.evidence, r.prob[i], f.Values())
					}
				}
			}
			c.RecoverPotentials()
		}
	}
}

func TestSaveOnLibdaiFormat(t *testing.T) {
	cases := []struct {
		cliques, adj [][]int
		cardin       []int
		values       [][]float64
		result       string
	}{{
		cliques: [][]int{{0, 1}, {1, 2}},
		adj:     [][]int{{1}, {0}},
		cardin:  []int{2, 2, 2},
		values: [][]float64{
			{.25, .35, .35, .05},
			{.20, .22, .40, .18},
		},
		result: "2\n\n" +
			"2\n" +
			"0 1 \n" +
			"2 2 \n" +
			"4\n" +
			fmt.Sprintf("%d     %.4f\n", 0, .25) +
			fmt.Sprintf("%d     %.4f\n", 1, .35) +
			fmt.Sprintf("%d     %.4f\n", 2, .35) +
			fmt.Sprintf("%d     %.4f\n", 3, .05) +
			"\n" +
			"2\n" +
			"1 2 \n" +
			"2 2 \n" +
			"4\n" +
			fmt.Sprintf("%d     %.4f\n", 0, .20) +
			fmt.Sprintf("%d     %.4f\n", 1, .22) +
			fmt.Sprintf("%d     %.4f\n", 2, .40) +
			fmt.Sprintf("%d     %.4f\n", 3, .18) +
			"\n",
	}, {
		cliques: [][]int{{0}, {1}, {0, 1, 2}, {2, 3}, {2, 4}},
		adj:     [][]int{{2}, {2}, {0, 1, 3, 4}, {2}, {2}},
		cardin:  []int{2, 2, 2, 2, 2},
		values: [][]float64{
			{.999, .001},
			{.998, .002},
			{.999, .06, .71, .05, .001, .94, .29, .95},
			{.95, .10, .05, .90},
			{.99, .30, .01, .70},
		},
		result: "5\n\n" +
			"1\n" +
			"0 \n" +
			"2 \n" +
			"2\n" +
			fmt.Sprintf("%d     %.4f\n", 0, .999) +
			fmt.Sprintf("%d     %.4f\n", 1, .001) +
			"\n" +
			"1\n" +
			"1 \n" +
			"2 \n" +
			"2\n" +
			fmt.Sprintf("%d     %.4f\n", 0, .998) +
			fmt.Sprintf("%d     %.4f\n", 1, .002) +
			"\n" +
			"3\n" +
			"0 1 2 \n" +
			"2 2 2 \n" +
			"8\n" +
			fmt.Sprintf("%d     %.4f\n", 0, .999) +
			fmt.Sprintf("%d     %.4f\n", 1, .06) +
			fmt.Sprintf("%d     %.4f\n", 2, .71) +
			fmt.Sprintf("%d     %.4f\n", 3, .05) +
			fmt.Sprintf("%d     %.4f\n", 4, .001) +
			fmt.Sprintf("%d     %.4f\n", 5, .94) +
			fmt.Sprintf("%d     %.4f\n", 6, .29) +
			fmt.Sprintf("%d     %.4f\n", 7, .95) +
			"\n" +
			"2\n" +
			"2 3 \n" +
			"2 2 \n" +
			"4\n" +
			fmt.Sprintf("%d     %.4f\n", 0, .95) +
			fmt.Sprintf("%d     %.4f\n", 1, .10) +
			fmt.Sprintf("%d     %.4f\n", 2, .05) +
			fmt.Sprintf("%d     %.4f\n", 3, .90) +
			"\n" +
			"2\n" +
			"2 4 \n" +
			"2 2 \n" +
			"4\n" +
			fmt.Sprintf("%d     %.4f\n", 0, .99) +
			fmt.Sprintf("%d     %.4f\n", 1, .30) +
			fmt.Sprintf("%d     %.4f\n", 2, .01) +
			fmt.Sprintf("%d     %.4f\n", 3, .70) +
			"\n",
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
		var b bytes.Buffer
		c.SaveOnLibdaiFormat(&b)
		got := b.String()
		if got != tt.result {
			for i := range tt.result {
				if got[i] != tt.result[i] {
					t.Errorf("Error on position %v, (%c)!=(%c)\nWant:\n[%v]\nGot:\n[%v]\n",
						i, got[i], tt.result[i], tt.result, got)
					break
				}
			}
		}
	}
}

func TestSaveOn(t *testing.T) {
	cases := []struct {
		cliques, adj [][]int
		cardin       []int
		values       [][]float64
		result       string
	}{{
		cliques: [][]int{{0, 1}, {1, 2}},
		adj:     [][]int{{1}, {0}},
		cardin:  []int{2, 2, 2},
		values: [][]float64{
			{.25, .35, .35, .05},
			{.20, .22, .40, .18},
		},
		result: "2\n" +
			"0 1 \n" +
			"1 2 \n" +
			"\n" +
			"1 \n" +
			"0 \n" +
			"\n" +
			"2 2 2 \n" +
			fmt.Sprintf("%.8f %.8f %.8f %.8f \n", .25, .35, .35, .05) +
			fmt.Sprintf("%.8f %.8f %.8f %.8f \n", .20, .22, .40, .18) +
			"\n",
	}, {
		cliques: [][]int{{0, 1}, {1, 2}},
		adj:     [][]int{{1}, {0}},
		result: "2\n" +
			"0 1 \n" +
			"1 2 \n" +
			"\n" +
			"1 \n" +
			"0 \n" +
			"\n" +
			"0 0 0 \n",
	}, {
		cliques: [][]int{{0}, {1}, {0, 1, 2}, {2, 3}, {2, 4}},
		adj:     [][]int{{2}, {2}, {0, 1, 3, 4}, {2}, {2}},
		cardin:  []int{2, 2, 2, 2, 2},
		values: [][]float64{
			{.999, .001},
			{.998, .002},
			{.999, .06, .71, .05, .001, .94, .29, .95},
			{.95, .10, .05, .90},
			{.99, .30, .01, .70},
		},
		result: "5\n" +
			"0 \n" +
			"1 \n" +
			"0 1 2 \n" +
			"2 3 \n" +
			"2 4 \n" +
			"\n" +
			"2 \n" +
			"2 \n" +
			"0 1 3 4 \n" +
			"2 \n" +
			"2 \n" +
			"\n" +
			"2 2 2 2 2 \n" +
			fmt.Sprintf("%.8f %.8f \n", .999, .001) +
			fmt.Sprintf("%.8f %.8f \n", .998, .002) +
			fmt.Sprintf("%.8f %.8f %.8f %.8f %.8f %.8f %.8f %.8f \n",
				.999, .06, .71, .05, .001, .94, .29, .95) +
			fmt.Sprintf("%.8f %.8f %.8f %.8f \n", .95, .10, .05, .90) +
			fmt.Sprintf("%.8f %.8f %.8f %.8f \n", .99, .30, .01, .70) +
			"\n",
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
		if len(tt.values) > 0 {
			err = c.SetAllPotentials(potentials)
		}
		if err != nil {
			t.Errorf(err.Error())
		}
		var b bytes.Buffer
		c.SaveOn(&b)
		got := b.String()
		if len(got) == 0 {
			t.Fatalf("Empty string")
		}
		if got != tt.result {
			for i := range got {
				if got[i] != tt.result[i] {
					t.Errorf("Error on position %v, (%c)!=(%c)\nWant:\n[%v]\nGot:\n[%v]\n",
						i, got[i], tt.result[i], tt.result, got)
					break
				}
			}
		}
	}
}

func TestLoadFrom(t *testing.T) {
	cases := []struct {
		cliques, adj [][]int
		cardin       []int
		values       [][]float64
		saved        string
	}{{
		cliques: [][]int{{0, 1}, {1, 2}},
		adj:     [][]int{{1}, {0}},
		cardin:  []int{2, 2, 2},
		values: [][]float64{
			{.25, .35, .35, .05},
			{.20, .22, .40, .18},
		},
		saved: "2\n" +
			"0 1 \n" +
			"1 2 \n" +
			"\n" +
			"1 \n" +
			"0 \n" +
			"\n" +
			"2 2 2 \n" +
			fmt.Sprintf("%.8f %.8f %.8f %.8f \n", .25, .35, .35, .05) +
			fmt.Sprintf("%.8f %.8f %.8f %.8f \n", .20, .22, .40, .18) +
			"\n",
	}, {
		cliques: [][]int{{0, 1}, {1, 2}},
		adj:     [][]int{{1}, {0}},
		cardin:  []int{0, 0, 0},
		values:  [][]float64{{}, {}},
		saved: "2\n" +
			"0 1 \n" +
			"1 2 \n" +
			"\n" +
			"1 \n" +
			"0 \n" +
			"\n" +
			"0 0 0 \n",
	}, {
		cliques: [][]int{{0}, {1}, {0, 1, 2}, {2, 3}, {2, 4}},
		adj:     [][]int{{2}, {2}, {0, 1, 3, 4}, {2}, {2}},
		cardin:  []int{2, 2, 2, 2, 2},
		values: [][]float64{
			{.999, .001},
			{.998, .002},
			{.999, .06, .71, .05, .001, .94, .29, .95},
			{.95, .10, .05, .90},
			{.99, .30, .01, .70},
		},
		saved: "5\n" +
			"0 \n" +
			"1 \n" +
			"0 1 2 \n" +
			"2 3 \n" +
			"2 4 \n" +
			"\n" +
			"2 \n" +
			"2 \n" +
			"0 1 3 4 \n" +
			"2 \n" +
			"2 \n" +
			"\n" +
			"2 2 2 2 2 \n" +
			fmt.Sprintf("%.8f %.8f \n", .999, .001) +
			fmt.Sprintf("%.8f %.8f \n", .998, .002) +
			fmt.Sprintf("%.8f %.8f %.8f %.8f %.8f %.8f %.8f %.8f \n",
				.999, .06, .71, .05, .001, .94, .29, .95) +
			fmt.Sprintf("%.8f %.8f %.8f %.8f \n", .95, .10, .05, .90) +
			fmt.Sprintf("%.8f %.8f %.8f %.8f \n", .99, .30, .01, .70) +
			"\n",
	}}
	for _, tt := range cases {
		c := LoadFrom(strings.NewReader(tt.saved))
		if c == nil {
			t.Fatalf("Nil clique tree returned")
		}
		if c.Size() != len(tt.cliques) {
			t.Errorf("wrong number of cliques, want %v, got %v", len(tt.cliques), c.Size())
		}
		if c.N() != len(tt.cardin) {
			t.Errorf("wrong cardinality, want %v, got %v", len(tt.cardin), c.N())
		}
		for i := 0; i < c.Size(); i++ {
			if !reflect.DeepEqual(c.Clique(i), tt.cliques[i]) {
				t.Errorf("wrong variables, want %v, got %v", tt.cliques[i], c.Clique(i))
			} else {
				for _, v := range tt.cliques[i] {
					if c.InitialPotential(i).Cardinality()[v] != tt.cardin[v] {
						t.Errorf("wrong cardinality, want %v, got %v", tt.cardin[v],
							c.InitialPotential(i).Cardinality()[v])
					}
				}
			}
			if !reflect.DeepEqual(tt.values[i], c.InitialPotential(i).Values()) {
				t.Errorf("wrong values, want %v, got %v", tt.values[i], c.InitialPotential(i).Values())
			}
			got := append([]int(nil), c.Neighbours(i)...)
			want := append([]int(nil), tt.adj[i]...)
			sort.Ints(got)
			sort.Ints(want)
			if !reflect.DeepEqual(want, got) {
				t.Errorf("wrong neighbours, want%v, got %v", want, got)
			}
		}
	}
}
