package cliquetree

import (
	"reflect"
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

func TestUpDownCalibration(t *testing.T) {
	c := initCliqueTree(factorList, adjList)
	c.UpDownCalibration()
	calculateCalibrated()
	for i, f := range cal {
		got := c.Calibrated(i)
		assig := assignment.New(f.Variables(), cardin)
		for {
			u := f.Get(assig)
			v := got.Get(assig)
			if !utils.FuzzyEqual(u, v) {
				t.Errorf("F[%v][%v]: want(%v); got(%v)", i, assig, u, v)
			}
			if hasnext := assig.Next(); !hasnext {
				break
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
}{
	{
		iphi: []int{0, 10, 9, 3, 4, 5, 6, 7, 1, 2, 8},
		chartree: characteristic.Tree{
			P: []int{-1, 5, 0, 0, 2, 8, 8, 1, 0},
			L: []int{-1, 2, -1, -1, 0, 2, 1, 2, -1},
		},
		cliques: [][]int{
			{1, 2, 8},
			{0, 4, 7, 1},
			{10, 1, 2, 8},
			{9, 1, 2, 8},
			{3, 10, 2, 8},
			{4, 7, 1, 2},
			{5, 7, 1, 8},
			{6, 0, 4, 7},
			{7, 1, 2, 8},
		},
		sepsets: [][]int{
			[]int(nil),
			{4, 7, 1},
			{1, 2, 8},
			{1, 2, 8},
			{10, 2, 8},
			{7, 1, 2},
			{7, 1, 8},
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
			if !reflect.DeepEqual(got.nodes[i].neighbours, v.adj[i]) {
				t.Errorf("Adj[%v]; Got: %v; Want: %v", i, got.nodes[i].neighbours, v.adj[i])
			}
			if got.parent[i] != v.parent[i] {
				t.Errorf("parent[%v]; Got: %v; Want: %v", i, got.parent[i], v.parent[i])
			}
		}
	}
}
