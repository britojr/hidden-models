package learn

import (
	"bytes"
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
