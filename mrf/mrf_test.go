package mrf

import (
	"bytes"
	"fmt"
	"math"
	"reflect"
	"strings"
	"testing"

	"github.com/britojr/kbn/factor"
	"github.com/britojr/kbn/utl/floats"
)

func TestLoadFromUAI(t *testing.T) {
	cases := []struct {
		cardin  []int
		cliques [][]int
		values  [][]float64
		saved   string
	}{{
		cardin:  []int{2, 2, 2},
		cliques: [][]int{{0, 1}, {1, 2}},
		values: [][]float64{
			{.25, .35, .35, .05},
			{.20, .22, .40, .18},
		},
		saved: "MARKOV\n" +
			"3\n" +
			"2 2 2 \n" +
			"2\n" +
			"2 0 1 \n" +
			"2 1 2 \n" +
			"4\n" +
			fmt.Sprintf("%e\n%e\n%e\n%e\n\n", .25, .35, .35, .05) +
			"4\n" +
			fmt.Sprintf("%e\n%e\n%e\n%e\n\n", .20, .22, .40, .18),
	}, {
		cardin:  []int{2, 2, 2},
		cliques: [][]int{{0, 1}, {1, 2}},
		values: [][]float64{
			{.25, .35, .35, .05},
			{.20, .22, .40, .18},
		},
		saved: "MARKOV\n" +
			"3\n" +
			"2 2 2 \n" +
			"2\n" +
			"2 0 1 \n" +
			"2 1 2 \n" +
			"\n4" +
			fmt.Sprintf("\n%e %e %e %e \n", .25, .35, .35, .05) +
			"\n4" +
			fmt.Sprintf("\n%e %e %e %e \n", .20, .22, .40, .18),
	}, {
		cliques: [][]int{{0}, {1}, {0, 1, 2}, {2, 3}, {2, 4}},
		cardin:  []int{2, 2, 2, 2, 2},
		values: [][]float64{
			{.999, .001},
			{.998, .002},
			{.999, .06, .71, .05, .001, .94, .29, .95},
			{.95, .10, .05, .90},
			{.99, .30, .01, .70},
		},
		saved: "MARKOV\n" +
			"5\n" +
			"2 2 2 2 2 \n" +
			"5\n" +
			"1 0 \n" +
			"1 1 \n" +
			"3 0 1 2 \n" +
			"2 2 3 \n" +
			"2 2 4 \n" +
			"\n2" +
			fmt.Sprintf("\n%e %e \n", .999, .001) +
			"\n2" +
			fmt.Sprintf("\n%e %e \n", .998, .002) +
			"\n8" +
			fmt.Sprintf("\n%e %e %e %e %e %e %e %e \n",
				.999, .06, .71, .05, .001, .94, .29, .95) +
			"\n4" +
			fmt.Sprintf("\n%e %e %e %e \n", .95, .10, .05, .90) +
			"\n4" +
			fmt.Sprintf("\n%e %e %e %e \n", .99, .30, .01, .70),
	}}
	for _, tt := range cases {
		m := LoadFromUAI(strings.NewReader(tt.saved))
		if m == nil {
			t.Fatalf("Nil MRF returned")
		}
		if len(m.potentials) != len(tt.cliques) {
			t.Errorf("wrong number of potentials, want %v, got %v", len(tt.cliques), len(m.potentials))
		}
		if !reflect.DeepEqual(tt.cardin, m.cardin) {
			t.Errorf("wrong cardinality, want %v, got %v", tt.cardin, m.cardin)
		}
		for i, p := range m.potentials {
			if !reflect.DeepEqual(tt.cliques[i], p.Variables()) {
				t.Errorf("wrong variables, want %v, got %v", tt.cliques[i], p.Variables())
			}
			if !reflect.DeepEqual(tt.values[i], p.Values()) {
				t.Errorf("wrong values, want %v, got %v", tt.values[i], p.Values())
			}
		}
	}
}

func TestUnnormalizedProb(t *testing.T) {
	cases := []struct {
		cardin []int
		pot    []*factor.Factor
		result []struct {
			evid []int
			prob float64
		}
	}{{
		cardin: []int{2, 2, 2},
		pot: []*factor.Factor{
			factor.NewFactorValues([]int{0, 1}, []int{2, 2, 2}, []float64{.25, .35, .35, .05}),
			factor.NewFactorValues([]int{1, 2}, []int{2, 2, 2}, []float64{.20, .22, .40, .18}),
		},
		result: []struct {
			evid []int
			prob float64
		}{
			{[]int{0, 0, 0}, .25 * .20},
			{[]int{1, 1, 1}, .05 * .18},
			{[]int{0, 1, 0}, .35 * .22},
		},
	}, {
		cardin: []int{2, 2, 2, 2, 2},
		pot: []*factor.Factor{
			factor.NewFactorValues([]int{0}, []int{2, 2, 2, 2, 2}, []float64{.999, .001}),
			factor.NewFactorValues([]int{1}, []int{2, 2, 2, 2, 2}, []float64{.998, .002}),
			factor.NewFactorValues([]int{0, 1, 2}, []int{2, 2, 2, 2, 2},
				[]float64{.999, .06, .71, .05, .001, .94, .29, .95}),
			factor.NewFactorValues([]int{2, 3}, []int{2, 2, 2, 2, 2}, []float64{.95, .10, .05, .90}),
			factor.NewFactorValues([]int{2, 4}, []int{2, 2, 2, 2, 2}, []float64{.99, .30, .01, .70}),
		},
		result: []struct {
			evid []int
			prob float64
		}{
			{[]int{0, 0, 0, 0, 0}, .999 * .998 * .999 * .95 * .99},
			{[]int{1, 1, 1, 1, 1}, .001 * .002 * .95 * .90 * .70},
			{[]int{0, 1, 1, 0, 0}, .999 * .002 * .29 * .10 * .30},
		},
	}}
	for _, tt := range cases {
		m := &Mrf{tt.cardin, tt.pot}
		for _, r := range tt.result {
			got := m.UnnormalizedProb(r.evid)
			if !floats.AlmostEqual(r.prob, got) {
				t.Errorf("wrong value, want %v, got %v", r.prob, got)
			}
		}
	}
}

func TestUnnormLogProb(t *testing.T) {
	cases := []struct {
		cardin []int
		pot    []*factor.Factor
		result []struct {
			evid []int
			prob float64
		}
	}{{
		cardin: []int{2, 2, 2},
		pot: []*factor.Factor{
			factor.NewFactorValues([]int{0, 1}, []int{2, 2, 2}, []float64{.25, .35, .35, .05}),
			factor.NewFactorValues([]int{1, 2}, []int{2, 2, 2}, []float64{.20, .22, .40, .18}),
		},
		result: []struct {
			evid []int
			prob float64
		}{
			{[]int{0, 0, 0}, math.Log(.25) + math.Log(.20)},
			{[]int{1, 1, 1}, math.Log(.05 * .18)},
			{[]int{0, 1, 0}, math.Log(.35) + math.Log(.22)},
		},
	}, {
		cardin: []int{2, 2, 2, 2, 2},
		pot: []*factor.Factor{
			factor.NewFactorValues([]int{0}, []int{2, 2, 2, 2, 2}, []float64{.999, .001}),
			factor.NewFactorValues([]int{1}, []int{2, 2, 2, 2, 2}, []float64{.998, .002}),
			factor.NewFactorValues([]int{0, 1, 2}, []int{2, 2, 2, 2, 2},
				[]float64{.999, .06, .71, .05, .001, .94, .29, .95}),
			factor.NewFactorValues([]int{2, 3}, []int{2, 2, 2, 2, 2}, []float64{.95, .10, .05, .90}),
			factor.NewFactorValues([]int{2, 4}, []int{2, 2, 2, 2, 2}, []float64{.99, .30, .01, .70}),
		},
		result: []struct {
			evid []int
			prob float64
		}{
			{[]int{0, 0, 0, 0, 0},
				math.Log(.999) + math.Log(.998) + math.Log(.999) + math.Log(.95) + math.Log(.99)},
			{[]int{1, 1, 1, 1, 1},
				math.Log(.001) + math.Log(.002) + math.Log(.95) + math.Log(.90) + math.Log(.70)},
			{[]int{0, 1, 1, 0, 0},
				math.Log(.999) + math.Log(.002) + math.Log(.29) + math.Log(.10) + math.Log(.30)},
		},
	}}
	for _, tt := range cases {
		m := &Mrf{tt.cardin, tt.pot}
		for _, r := range tt.result {
			got := m.UnnormLogProb(r.evid)
			if !floats.AlmostEqual(r.prob, got) {
				t.Errorf("%v != %v, for evid %v", r.prob, got, r.evid)
			}
		}
	}
}

func TestSaveOnLibdaiFormat(t *testing.T) {
	cases := []struct {
		cliques [][]int
		cardin  []int
		values  [][]float64
		result  string
	}{{
		cliques: [][]int{{0, 1}, {1, 2}},
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
			fmt.Sprintf("%d     %e\n", 0, .25) +
			fmt.Sprintf("%d     %e\n", 1, .35) +
			fmt.Sprintf("%d     %e\n", 2, .35) +
			fmt.Sprintf("%d     %e\n", 3, .05) +
			"\n" +
			"2\n" +
			"1 2 \n" +
			"2 2 \n" +
			"4\n" +
			fmt.Sprintf("%d     %e\n", 0, .20) +
			fmt.Sprintf("%d     %e\n", 1, .22) +
			fmt.Sprintf("%d     %e\n", 2, .40) +
			fmt.Sprintf("%d     %e\n", 3, .18) +
			"\n",
	}, {
		cliques: [][]int{{0}, {1}, {0, 1, 2}, {2, 3}, {2, 4}},
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
			fmt.Sprintf("%d     %e\n", 0, .999) +
			fmt.Sprintf("%d     %e\n", 1, .001) +
			"\n" +
			"1\n" +
			"1 \n" +
			"2 \n" +
			"2\n" +
			fmt.Sprintf("%d     %e\n", 0, .998) +
			fmt.Sprintf("%d     %e\n", 1, .002) +
			"\n" +
			"3\n" +
			"0 1 2 \n" +
			"2 2 2 \n" +
			"8\n" +
			fmt.Sprintf("%d     %e\n", 0, .999) +
			fmt.Sprintf("%d     %e\n", 1, .06) +
			fmt.Sprintf("%d     %e\n", 2, .71) +
			fmt.Sprintf("%d     %e\n", 3, .05) +
			fmt.Sprintf("%d     %e\n", 4, .001) +
			fmt.Sprintf("%d     %e\n", 5, .94) +
			fmt.Sprintf("%d     %e\n", 6, .29) +
			fmt.Sprintf("%d     %e\n", 7, .95) +
			"\n" +
			"2\n" +
			"2 3 \n" +
			"2 2 \n" +
			"4\n" +
			fmt.Sprintf("%d     %e\n", 0, .95) +
			fmt.Sprintf("%d     %e\n", 1, .10) +
			fmt.Sprintf("%d     %e\n", 2, .05) +
			fmt.Sprintf("%d     %e\n", 3, .90) +
			"\n" +
			"2\n" +
			"2 4 \n" +
			"2 2 \n" +
			"4\n" +
			fmt.Sprintf("%d     %e\n", 0, .99) +
			fmt.Sprintf("%d     %e\n", 1, .30) +
			fmt.Sprintf("%d     %e\n", 2, .01) +
			fmt.Sprintf("%d     %e\n", 3, .70) +
			"\n",
	}}
	for _, tt := range cases {
		potentials := make([]*factor.Factor, len(tt.values))
		for i, v := range tt.values {
			potentials[i] = factor.NewFactorValues(tt.cliques[i], tt.cardin, v)
		}
		m := &Mrf{tt.cardin, potentials}
		var b bytes.Buffer
		m.SaveOnLibdaiFormat(&b)
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
