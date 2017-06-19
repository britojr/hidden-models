package learn

import (
	"fmt"
	"reflect"
	"testing"

	"github.com/britojr/kbn/factor"
	"github.com/britojr/kbn/utl/floats"
)

type FakeCounter struct {
	counts map[string][]int
}

func (f FakeCounter) CountAssignments(varlist []int) []int {
	return f.counts[fmt.Sprint(varlist)]
}

func TestCreateEmpiricPotentials(t *testing.T) {
	fakeCounter := FakeCounter{
		counts: map[string][]int{
			fmt.Sprint([]int{0, 1, 2}): {15, 10, 5, 25, 5, 20, 15, 5},
			fmt.Sprint([]int{0, 1}):    {20, 30, 20, 30},
			fmt.Sprint([]int{0, 2}):    {20, 35, 20, 25},
			fmt.Sprint([]int{1, 2}):    {25, 30, 25, 20},
			fmt.Sprint([]int{0}):       {40, 60},
			fmt.Sprint([]int{1}):       {50, 50},
			fmt.Sprint([]int{2}):       {55, 45},
		},
	}
	cases := []struct {
		cliques [][]int
		cardin  []int
		numobs  int
		counter FakeCounter
		result  [][]float64
	}{{
		cliques: [][]int{{0, 1}, {1, 2}},
		cardin:  []int{2, 2, 2},
		numobs:  2,
		counter: fakeCounter,
		result:  [][]float64{{.20, .30, .20, .30}, {.25, .25, .25, .25}},
	}, {
		cliques: [][]int{{0, 1}, {1, 2}},
		cardin:  []int{2, 2, 2},
		numobs:  3,
		counter: fakeCounter,
		result:  [][]float64{{.20, .30, .20, .30}, {.25, .30, .25, .20}},
	}, {
		cliques: [][]int{{0, 1}, {1, 2}},
		cardin:  []int{2, 2, 2},
		numobs:  1,
		counter: fakeCounter,
		result:  [][]float64{{.20, .30, .20, .30}, {.25, .25, .25, .25}},
	}}

	// test conditional uniform
	for _, tt := range cases {
		faclist := createEmpiricPotentials(tt.counter, tt.cliques, tt.cardin, tt.numobs, DistUniform, ModeCond, 0)
		if len(faclist) != len(tt.result) {
			t.Errorf("wrong number of factors, expected %v, got %v", len(tt.result), len(faclist))
		}
		for i, f := range faclist {
			tot := floats.Sum(f.Values())
			if !floats.AlmostEqual(tot, 1) {
				t.Errorf("uniform factor not normalized, sums to: %v", tot)
			}
			for _, v := range f.Values() {
				if v == 0 {
					t.Errorf("uniform factor has zero values: %v", f.Values())
				}
			}
			if !reflect.DeepEqual(tt.result[i], f.Values()) {
				t.Errorf("Wrong values, want %v, got %v", tt.result[i], f.Values())
			}
		}
	}

	// test random independent
	for _, tt := range cases {
		faclist := createEmpiricPotentials(tt.counter, tt.cliques, tt.cardin, tt.numobs, DistRandom, ModeIndep, 0)
		if len(faclist) != len(tt.result) {
			t.Errorf("wrong number of factors, expected %v, got %v", len(tt.result), len(faclist))
		}
		for _, f := range faclist {
			tot := floats.Sum(f.Values())
			if !floats.AlmostEqual(tot, 1) {
				t.Errorf("random factor not normalized, sums to: %v", tot)
			}
			for _, v := range f.Values() {
				if v == 0 {
					t.Errorf("random factor has zero values: %v", f.Values())
				}
			}
			hid := []int(nil)
			for j := range tt.cardin {
				if j >= tt.numobs {
					hid = append(hid, j)
				}
			}
			q := f.SumOut(hid)
			stridobs := len(q.Values())

			for i := 0; i < len(f.Values()); i += stridobs {
				prop := q.Values()[0] / f.Values()[i]
				for j := 1; j < stridobs; j++ {
					prop2 := q.Values()[j] / f.Values()[i+j]
					if !floats.AlmostEqual(prop, prop2, 1e-9) {
						t.Errorf("wrong proportion %v, %v", prop, prop2)
					}
				}
			}
		}
	}
}

func TestCreateRandomPortentials(t *testing.T) {
	cases := []struct {
		cliques [][]int
		cardin  []int
	}{{
		cliques: [][]int{{0, 1}, {1, 2}},
		cardin:  []int{2, 2, 2},
	}}
	for _, tt := range cases {
		faclist := createRandomPotentials(tt.cliques, tt.cardin, DistRandom, 0)
		for _, f := range faclist {
			tot := floats.Sum(f.Values())
			if !floats.AlmostEqual(tot, 1) {
				t.Errorf("random factor not normalized, sums to: %v", tot)
			}
			for _, v := range f.Values() {
				if v == 0 {
					t.Errorf("random factor has zero values: %v", f.Values())
				}
			}
		}
	}
}

func TestConditionalValues(t *testing.T) {
	cases := []struct {
		varlist, hidden, cardin []int
		lenobs, lenhidden, dist int
		alpha                   float64
	}{
		{[]int{0, 1}, []int{0, 1}, []int{2, 2}, 1, 4, DistUniform, 0},
		{[]int{0, 1}, []int{1}, []int{2, 2}, 2, 2, DistUniform, 0},
		{[]int{0, 1}, []int{}, []int{2, 2}, 4, 1, DistUniform, 0},
		{[]int{0, 1, 2}, []int{1, 2}, []int{2, 2, 2}, 2, 4, DistUniform, 0},
		{[]int{0, 1, 2}, []int{2}, []int{2, 2, 2}, 4, 2, DistUniform, 0},
		{[]int{0, 1}, []int{0, 1}, []int{2, 2}, 1, 4, DistRandom, 0},
		{[]int{0, 1}, []int{1}, []int{2, 2}, 2, 2, DistRandom, 0},
		{[]int{0, 1}, []int{}, []int{2, 2}, 4, 1, DistRandom, 0},
		{[]int{0, 1, 2}, []int{1, 2}, []int{2, 2, 2}, 2, 4, DistRandom, 0},
		{[]int{0, 1, 2}, []int{2}, []int{2, 2, 2}, 4, 2, DistRandom, 0},
		{[]int{0, 1}, []int{0, 1}, []int{2, 2}, 1, 4, DistDirichlet, 1.5},
		{[]int{0, 1}, []int{1}, []int{2, 2}, 2, 2, DistDirichlet, 0.3},
		{[]int{0, 1}, []int{}, []int{2, 2}, 4, 1, DistDirichlet, 0.5},
		{[]int{0, 1, 2}, []int{1, 2}, []int{2, 2, 2}, 2, 4, DistDirichlet, 0.3},
		{[]int{0, 1, 2}, []int{2}, []int{2, 2, 2}, 4, 2, DistDirichlet, 1.8},
	}
	for _, tt := range cases {
		values := conditionalValues(tt.lenobs, tt.lenhidden, tt.dist, tt.alpha)
		g := factor.NewFactorValues(tt.varlist, tt.cardin, values).SumOut(tt.hidden)
		if tt.lenobs != len(g.Values()) {
			t.Errorf("wrong size, want %v got %v", tt.lenobs, len(g.Values()))
		}
		for _, v := range g.Values() {
			if !floats.AlmostEqual(v, float64(1)) {
				t.Errorf("wrong value, want 1.0, got %v (typePot %v) val=%v", v, tt.dist, g.Values())
				break
			}
		}
	}
}
