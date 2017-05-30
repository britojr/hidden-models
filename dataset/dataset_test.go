package dataset

import (
	"reflect"
	"strings"
	"testing"
)

func TestNew(t *testing.T) {
	cases := []struct {
		name     string
		delim    rune
		hdr      HdrFlags
		lin, col int
		cardin   []int
		data     [][]int
		content  string
	}{{
		"file1.txt", ' ', HdrName | HdrCardin,
		3, 11, []int{2, 2, 3, 2, 2, 4, 2, 2, 3, 2, 2},
		[][]int{
			{0, 1, 1, 1, 1, 1, 0, 0, 1, 0, 0},
			{0, 0, 1, 0, 1, 2, 1, 1, 0, 1, 0},
			{0, 0, 0, 0, 0, 1, 0, 0, 0, 1, 0},
		},
		`A B C D E F G H I J K
2 2 3 2 2 4 2 2 3 2 2
0 1 1 1 1 1 0 0 1 0 0
0 0 1 0 1 2 1 1 0 1 0
0 0 0 0 0 1 0 0 0 1 0`,
	}, {
		"file2.txt", ',', HdrCardin,
		4, 11, []int{2, 2, 3, 2, 2, 4, 2, 2, 3, 2, 2},
		[][]int{
			{0, 1, 1, 1, 1, 1, 0, 0, 1, 0, 0},
			{1, 0, 0, 0, 1, 2, 0, 1, 2, 0, 1},
			{0, 0, 0, 0, 1, 3, 0, 1, 0, 0, 1},
			{1, 1, 2, 1, 0, 3, 1, 0, 2, 1, 0},
		},
		`2,2,3,2,2,4,2,2,3,2,2
0,1,1,1,1,1,0,0,1,0,0
1,0,0,0,1,2,0,1,2,0,1
0,0,0,0,1,3,0,1,0,0,1
1,1,2,1,0,3,1,0,2,1,0`,
	}, {
		"file3.txt", ',', HdrName,
		3, 11, []int{2, 2, 2, 2, 2, 2, 2, 2, 2, 2, 1},
		[][]int{
			{0, 1, 1, 1, 1, 1, 0, 0, 1, 0, 0},
			{1, 0, 1, 1, 1, 0, 1, 1, 1, 0, 0},
			{0, 0, 0, 0, 0, 1, 0, 0, 0, 1, 0},
		},
		`A,B,C,D,E,F,G,H,I,J,K
0,1,1,1,1,1,0,0,1,0,0
1,0,1,1,1,0,1,1,1,0,0
0,0,0,0,0,1,0,0,0,1,0`,
	}, {
		"file4.txt", ',', HdrNameCard,
		3, 11, []int{2, 2, 3, 2, 2, 4, 2, 2, 3, 2, 2},
		[][]int{
			{0, 1, 1, 1, 1, 1, 0, 0, 1, 0, 0},
			{1, 0, 1, 1, 1, 0, 1, 1, 1, 0, 0},
			{0, 0, 0, 0, 0, 1, 0, 0, 0, 1, 0},
		},
		`A_2,B_2,C_3,D_2,E_2,F_4,G_2,H_2,I_3,J_2,K_2
0,1,1,1,1,1,0,0,1,0,0
1,0,1,1,1,0,1,1,1,0,0
0,0,0,0,0,1,0,0,0,1,0`,
	}}

	for _, tt := range cases {
		r := strings.NewReader(tt.content)
		d := New(r, tt.delim, tt.hdr)
		if tt.lin != d.NLines() {
			t.Errorf("wrong num lines %v != %v", tt.lin, d.NLines())
		}
		if tt.col != d.NCols() {
			t.Errorf("wrong num cols %v != %v", tt.col, d.NCols())
		}
		if !reflect.DeepEqual(tt.cardin, d.Cardin()) {
			t.Errorf("wrong cardinality %v != %v", tt.cardin, d.Cardin())
		}
		if !reflect.DeepEqual(tt.data, d.Data()) {
			t.Errorf("wrong data %v != %v", tt.data, d.Data())
		}

	}
}

func TestCountAssignments(t *testing.T) {
	cases := []struct {
		data    [][]int
		cardin  []int
		varlist []int
		result  []int
	}{{
		[][]int{
			{0, 0},
			{1, 0},
			{0, 1},
			{1, 0},
		},
		[]int{2, 2}, []int{0},
		[]int{2, 2},
	}, {
		[][]int{
			{0, 0},
			{1, 0},
			{0, 1},
			{1, 0},
		},
		[]int{2, 2}, []int{1},
		[]int{3, 1},
	}, {
		[][]int{
			{0, 0},
			{1, 0},
		},
		[]int{2, 2}, []int{3},
		[]int(nil),
	}, {
		[][]int{
			{0, 0},
			{1, 0},
		},
		[]int{2, 2}, []int{},
		[]int(nil),
	}}

	for _, tt := range cases {
		d := new(Dataset)
		d.cardin = tt.cardin
		d.data = tt.data
		d.initCount()
		got := d.CountAssignments(tt.varlist)
		if !reflect.DeepEqual(tt.result, got) {
			t.Errorf("want(%v); got(%v)", tt.result, got)
		}
	}
}
