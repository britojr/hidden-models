package learn

import (
	"bytes"
	"reflect"
	"strings"
	"testing"

	"github.com/britojr/kbn/utl/floats"
)

func TestWriteMarginals(t *testing.T) {
	cases := []struct {
		values  [][]float64
		content string
	}{{
		[][]float64{
			{.9, .1},
			{.1, .7, .2},
			{.8, .200013},
			{.5, .5},
		},
		`MAR
4 2 0.90000 0.10000 3 0.10000 0.70000 0.20000 2 0.80000 0.20001 2 0.50000 0.50000`,
	}}
	for _, tt := range cases {
		var b bytes.Buffer
		writeMarginals(&b, tt.values)
		got := b.String()
		if got != tt.content {
			for i := range tt.content {
				if got[i] != tt.content[i] {
					t.Errorf("Error on position %v, (%c)!=(%c)\nWant:\n[%v]\nGot:\n[%v]\n",
						i, got[i], tt.content[i], tt.content, got)
					break
				}
			}
		}
	}
}

func TestReadMarginals(t *testing.T) {
	cases := []struct {
		values  [][]float64
		content string
	}{{
		[][]float64{
			{.9, .1},
			{.1, .7, .2},
			{.8, .20001},
			{.5, .5},
		},
		`MAR
4 2 0.90000 0.10000 3 0.10000 0.70000 0.20000 2 0.80000 0.20001 2 0.50000 0.50000`,
	}}
	for _, tt := range cases {
		got := readMarginals(strings.NewReader(tt.content))
		if !reflect.DeepEqual(tt.values, got) {
			t.Errorf("Error on reading marginas\n%v\n%v\n", tt.values, got)
		}
	}
}

func TestMarginalsMSE(t *testing.T) {
	cases := []struct {
		e, a [][]float64
		d    float64
	}{{
		[][]float64{
			{.9, .1},
			{.1, .7, .2},
			{.8, .20001},
			{.5, .5},
		},
		[][]float64{
			{.9, .1},
			{.1, .7, .2},
			{.8, .20001},
			{.5, .5},
		},
		0,
	}, {
		[][]float64{
			{.9, .1},
			{.1, .7, .2},
		},
		[][]float64{
			{.8, .2},
			{.2, .7, .1},
		},
		0.008333333,
	}}
	for _, tt := range cases {
		got := marginalsMSE(tt.e, tt.a)
		if !floats.AlmostEqual(tt.d, got, 1e-5) {
			t.Errorf("want %v, got %v", tt.d, got)
		}
	}
}
