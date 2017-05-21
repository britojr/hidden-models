package floats

import (
	"reflect"
	"testing"
)

func TestAlmostEqual(t *testing.T) {
	cases := []struct {
		a, b  float64
		equal bool
		eps   []float64
	}{
		{0.1, 0.2, false, nil},
		{0, 0, true, nil},
		{0.0002, 0.0002, true, nil},
		{1.0005, 1.0005 + (epslon / 2.0), true, nil},
		{0.0005, 0.0005 + epslon, false, nil},
		{1.00050, 1.00051, false, []float64{1e-5}},
		{1.00050, 1.00051, true, []float64{1e-4}},
	}
	for _, tt := range cases {
		got := AlmostEqual(tt.a, tt.b, tt.eps...)
		if got != tt.equal {
			t.Errorf("%v == %v : got %v, want %v", tt.a, tt.b, got, tt.equal)
		}
	}
}

func TestMax(t *testing.T) {
	cases := []struct {
		xs     []float64
		result float64
	}{
		{[]float64{1, 2, 3}, 3},
		{[]float64{2, 2, 2}, 2},
		{[]float64{2, 3, 6, 5, 4, 1}, 6},
	}
	for _, tt := range cases {
		got := Max(tt.xs)
		if !AlmostEqual(tt.result, got) {
			t.Errorf("wrong value,  want %v, got %v", tt.result, got)
		}
	}
}

func TestMin(t *testing.T) {
	cases := []struct {
		xs     []float64
		result float64
	}{
		{[]float64{1, 2, 3}, 1},
		{[]float64{2, 2, 2}, 2},
		{[]float64{2, 3, 6, 5, 4, 1}, 1},
	}
	for _, tt := range cases {
		got := Min(tt.xs)
		if !AlmostEqual(tt.result, got) {
			t.Errorf("wrong value,  want %v, got %v", tt.result, got)
		}
	}
}

func TestSum(t *testing.T) {
	cases := []struct {
		values []float64
		sum    float64
	}{
		{[]float64{5, 5}, 10},
		{[]float64{1.5, 3.5, 0.5}, 5.5},
		{[]float64{}, 0},
		{[]float64(nil), 0},
	}
	for _, tt := range cases {
		got := Sum(tt.values)
		if tt.sum != got {
			t.Errorf("want %v, got %v", tt.sum, got)
		}
	}
}

func TestNormalize(t *testing.T) {
	cases := []struct {
		values []float64
		r      []float64
	}{
		{[]float64{5, 5}, []float64{.5, .5}},
		{[]float64{3, 1}, []float64{.75, .25}},
		{[]float64{1.5, 3.5, 1.0}, []float64{1.5 / 6, 3.5 / 6, 1.0 / 6}},
	}
	for _, tt := range cases {
		Normalize(tt.values)
		if !reflect.DeepEqual(tt.r, tt.values) {
			t.Errorf("want %v, got %v", tt.r, tt.values)
		}
	}
}
