package mrf

import (
	"fmt"
	"reflect"
	"strings"
	"testing"

	"github.com/britojr/kbn/factor"
	"github.com/britojr/kbn/utils"
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
			"\n4" +
			fmt.Sprintf("\n%.6f %.6f %.6f %.6f \n", .25, .35, .35, .05) +
			"\n4" +
			fmt.Sprintf("\n%.6f %.6f %.6f %.6f \n", .20, .22, .40, .18),
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
			fmt.Sprintf("\n%.6f %.6f \n", .999, .001) +
			"\n2" +
			fmt.Sprintf("\n%.6f %.6f \n", .998, .002) +
			"\n8" +
			fmt.Sprintf("\n%.6f %.6f %.6f %.6f %.6f %.6f %.6f %.6f \n",
				.999, .06, .71, .05, .001, .94, .29, .95) +
			"\n4" +
			fmt.Sprintf("\n%.6f %.6f %.6f %.6f \n", .95, .10, .05, .90) +
			"\n4" +
			fmt.Sprintf("\n%.6f %.6f %.6f %.6f \n", .99, .30, .01, .70),
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

func TestUnnormalizedMesure(t *testing.T) {
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
			got := m.UnnormalizedMesure(r.evid)
			if !utils.FuzzyEqual(r.prob, got) {
				t.Errorf("wrong value, want %v, got %v", r.prob, got)
			}
		}
	}
}
