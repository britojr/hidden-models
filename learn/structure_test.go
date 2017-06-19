package learn

import (
	"bytes"
	"fmt"
	"reflect"
	"strings"
	"testing"
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
			{0.11656955248666273, 0.883430674265923},
		},
		// `MAR
		// 4 2 0.90000 0.10000 3 0.10000 0.70000 0.20000 2 0.80000 0.20001 2 0.50000 0.50000`,
		"MAR\n4" +
			fmt.Sprintf(" 2 %g %g", .9, .1) +
			fmt.Sprintf(" 3 %g %g %g", .1, .7, .2) +
			fmt.Sprintf(" 2 %g %g", .8, .200013) +
			fmt.Sprintf(" 2 %g %g ", 0.11656955248666273, 0.883430674265923),
		// 	}, {
		// 		[][]float64{
		// 			{0.11656955248666273, 0.883430674265923},
		// 			{8.972273e-01, 1.027727e-01},
		// 			{8.800014e-01, 1.199986e-01},
		// 			{3.091610e-01, 6.908392e-01},
		// 			{5.309539e-02, 9.469051e-01},
		// 			{1.191980e-01, 8.808022e-01},
		// 			{1.193487e-01, 8.806514e-01},
		// 		},
		// 		`MAR
		// 7 2 0.11657 0.88343 2 0.897227 0.102773 2 0.880001 0.119999 2 0.309161 0.690839 2 0.0530953 0.946905 2 0.119198 0.880802 2 0.119349 0.880651 `,
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
