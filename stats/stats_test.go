package stats

import (
	"reflect"
	"testing"

	"github.com/britojr/kbn/floats"
)

func TestMean(t *testing.T) {
	cases := []struct {
		xs   []float64
		mean float64
	}{
		{[]float64{2, 2, 2}, 2},
		{[]float64{1, 2, 3}, 2},
		{[]float64{5, 4, 1, 2, 3, 6}, 3.5},
		{[]float64{12, 12, 12, 12, 13013}, 2612.2},
	}
	for _, tt := range cases {
		got := Mean(tt.xs)
		if !floats.AlmostEqual(tt.mean, got) {
			t.Errorf("wrong value,  want %v, got %v", tt.mean, got)
		}
	}
}
func TestMedian(t *testing.T) {
	cases := []struct {
		xs   []float64
		mean float64
	}{
		{[]float64{1, 2, 3}, 2},
		{[]float64{2, 2, 2}, 2},
		{[]float64{5, 4, 1, 2, 3, 6}, 3.5},
		{[]float64{3, 1, 7}, 3},
	}
	for _, tt := range cases {
		got := Median(tt.xs)
		if !floats.AlmostEqual(tt.mean, got) {
			t.Errorf("wrong value,  want %v, got %v", tt.mean, got)
		}
	}
}

func TestVariance(t *testing.T) {
	cases := []struct {
		xs   []float64
		want float64
	}{
		{[]float64{2, 2, 2}, 0},
		// {[]float64{1, 2, 3}, 2.0 / 3.0},
		// {[]float64{5, 4, 1, 2, 3, 6}, 2.916666667},
		// {[]float64{12, 12, 12, 12, 13013}, 27044160.16},
	}
	for _, tt := range cases {
		got := Variance(tt.xs)
		if !floats.AlmostEqual(tt.want, got, 1e-6) {
			t.Errorf("wrong value,  want %v, got %v", tt.want, got)
		}
	}
}

func TestStdev(t *testing.T) {
	cases := []struct {
		xs []float64
		sd float64
	}{
		{[]float64{2, 2, 2}, 0},
		// {[]float64{1, 2, 3}, 0.816496581},
		// {[]float64{5, 4, 1, 2, 3, 6}, 1.707825128},
		// {[]float64{12, 12, 12, 12, 13013}, math.Sqrt(27044160.16)},
	}
	for _, tt := range cases {
		got := Stdev(tt.xs)
		if !floats.AlmostEqual(tt.sd, got, 1e-6) {
			t.Errorf("wrong value,  want %v, got %v", tt.sd, got)
		}
	}
}

func TestDirichlet(t *testing.T) {
	cases := []struct {
		alphas []float64
	}{
		{[]float64{3.2, 3.2, 3.2, 3.2}},
		{[]float64{0.1, 0.1, 0.1, 0.1, 0.1, 0.1, 0.1, 0.1}},
		{[]float64{0.01, 0.01}},
		{[]float64{5}},
	}
	for _, tt := range cases {
		values := make([]float64, len(tt.alphas))
		Dirichlet(tt.alphas, values)
		if len(tt.alphas) != len(values) {
			t.Errorf("wrong size, want %v, got %v", len(tt.alphas), len(values))
		}
		if len(tt.alphas) != 0 && !floats.AlmostEqual(1, floats.Sum(values)) {
			t.Errorf("not normalized %v", values)
		}
	}

	// test different outcomes
	alphas := []float64{0.7, 0.7, 0.7, 0.7, 0.7, 0.7, 0.7, 0.7}
	a, b := make([]float64, len(alphas)), make([]float64, len(alphas))
	Dirichlet(alphas, a)
	Dirichlet(alphas, b)
	count := 0
	for i := range alphas {
		if floats.AlmostEqual(a[i], b[i]) {
			count++
		}
	}
	if count == len(alphas) {
		t.Errorf("Sampled the same distribution:\n%v\n%v", a, b)
	}
}

func TestNormalize(t *testing.T) {
	cases := []struct {
		values, normalized []float64
	}{{
		[]float64{0.15, 0.25, 0.35, 0.25},
		[]float64{0.15, 0.25, 0.35, 0.25},
	}, {
		[]float64{15, 25, 35, 25},
		[]float64{0.15, 0.25, 0.35, 0.25},
	}, {
		[]float64{10, 20, 30, 40, 50, 60, 70, 80},
		[]float64{1.0 / 36, 2.0 / 36, 3.0 / 36, 4.0 / 36, 5.0 / 36, 6.0 / 36, 7.0 / 36, 8.0 / 36},
	}, {
		[]float64{0.15},
		[]float64{1},
	}}
	for _, tt := range cases {
		Normalize(tt.values)
		if !reflect.DeepEqual(tt.values, tt.normalized) {
			t.Errorf("want %v, got %v", tt.normalized, tt.values)
		}
	}
}
