package cliquetree

import (
	"testing"

	"github.com/britojr/kbn/assignment"
	"github.com/britojr/kbn/factor"
	"github.com/britojr/kbn/utils"
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
		p = append(p, factor.New(f.varlist, cardin, f.values))
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
		c.SetPotential(i, factor.New(f.varlist, cardin, f.values))
	}
	return c
}

func TestNew(t *testing.T) {
	c := New(len(factorList))
	for i, f := range factorList {
		c.SetClique(i, f.varlist)
		c.SetNeighbours(i, adjList[i])
		c.SetPotential(i, factor.New(f.varlist, cardin, f.values))
	}
}

func TestUpDownCalibration(t *testing.T) {
	c := initCliqueTree(factorList, adjList)
	c.UpDownCalibration()
	calculateCalibrated()
	for i, f := range cal {
		got := c.Calibrated(i)
		assig := assignment.New(f.Variables(), cardin)
		for assig != nil {
			u := f.Get(assig)
			v := got.Get(assig)
			if !utils.FuzzyEqual(u, v) {
				t.Errorf("F[%v][%v]: want(%v); got(%v)", i, assig, u, v)
			}
			assig.Next()
		}
	}
}

func TestIterativeCalibration(t *testing.T) {
	c := initCliqueTree(factorList, adjList)
	c.IterativeCalibration()
	calculateCalibrated()
	for i, f := range cal {
		got := c.Calibrated(i)
		assig := assignment.New(f.Variables(), cardin)
		for assig != nil {
			u := f.Get(assig)
			v := got.Get(assig)
			if !utils.FuzzyEqual(u, v) {
				t.Errorf("F[%v][%v]: want(%v); got(%v)", i, assig, u, v)
			}
			assig.Next()
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

func BenchmarkIterativeCalibration(b *testing.B) {
	ctrees := make([]*CliqueTree, 0)
	for _, bt := range benchTest {
		ctrees = append(ctrees, initCliqueTree(bt.fl, bt.adj))
	}
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		for _, c := range ctrees {
			c.IterativeCalibration()
		}
	}
}
