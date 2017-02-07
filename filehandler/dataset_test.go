package filehandler

import (
	"reflect"
	"testing"
)

type filepack struct {
	fileName  string
	separator rune
	headerlns HeaderFlags
	cardin    []int
	data      [][]int
}

var testFiles = []filepack{
	filepack{
		"dataset_test1.txt",
		' ',
		NameHeader | CardinHeader,
		[]int{2, 2, 3, 2, 2, 4, 2, 2, 3, 2, 2},
		[][]int{
			[]int{0, 1, 1, 1, 1, 1, 0, 0, 1, 0, 0},
			[]int{1, 0, 1, 1, 1, 0, 1, 1, 1, 0, 0},
			[]int{1, 1, 1, 0, 0, 0, 1, 1, 0, 0, 0},
			[]int{1, 0, 0, 0, 1, 2, 0, 1, 2, 0, 1},
			[]int{0, 0, 0, 0, 1, 3, 0, 1, 0, 0, 1},
			[]int{0, 1, 2, 1, 0, 0, 0, 0, 2, 0, 1},
			[]int{1, 1, 2, 1, 0, 3, 1, 0, 2, 1, 0},
			[]int{0, 0, 1, 0, 1, 2, 1, 1, 0, 1, 0},
			[]int{0, 0, 0, 0, 0, 1, 0, 0, 0, 1, 0},
		},
	},
	filepack{
		"dataset_test2.txt",
		',',
		CardinHeader,
		[]int{2, 2, 3, 2, 2, 4, 2, 2, 3, 2, 2},
		[][]int{
			[]int{0, 1, 1, 1, 1, 1, 0, 0, 1, 0, 0},
			[]int{1, 0, 1, 1, 1, 0, 1, 1, 1, 0, 0},
			[]int{1, 1, 1, 0, 0, 0, 1, 1, 0, 0, 0},
			[]int{1, 0, 0, 0, 1, 2, 0, 1, 2, 0, 1},
			[]int{0, 0, 0, 0, 1, 3, 0, 1, 0, 0, 1},
			[]int{0, 1, 2, 1, 0, 0, 0, 0, 2, 0, 1},
			[]int{1, 1, 2, 1, 0, 3, 1, 0, 2, 1, 0},
			[]int{0, 0, 1, 0, 1, 2, 1, 1, 0, 1, 0},
			[]int{0, 0, 0, 0, 0, 1, 0, 0, 0, 1, 0},
		},
	},
	filepack{
		"dataset_test3.txt",
		',',
		NameHeader,
		[]int{2, 2, 3, 2, 2, 4, 2, 2, 3, 2, 2},
		[][]int{
			[]int{0, 1, 1, 1, 1, 1, 0, 0, 1, 0, 0},
			[]int{1, 0, 1, 1, 1, 0, 1, 1, 1, 0, 0},
			[]int{1, 1, 1, 0, 0, 0, 1, 1, 0, 0, 0},
			[]int{1, 0, 0, 0, 1, 2, 0, 1, 2, 0, 1},
			[]int{0, 0, 0, 0, 1, 3, 0, 1, 0, 0, 1},
			[]int{0, 1, 2, 1, 0, 0, 0, 0, 2, 0, 1},
			[]int{1, 1, 2, 1, 0, 3, 1, 0, 2, 1, 0},
			[]int{0, 0, 1, 0, 1, 2, 1, 1, 0, 1, 0},
			[]int{0, 0, 0, 0, 0, 1, 0, 0, 0, 1, 0},
		},
	},
}

func TestNewDataSet(t *testing.T) {
	for _, f := range testFiles {
		d := NewDataSet(f.fileName, f.separator, f.headerlns)
		d.Read()
	}
}

func TestSize(t *testing.T) {
	for _, f := range testFiles {
		d := NewDataSet(f.fileName, f.separator, f.headerlns)
		d.Read()
		got := d.Size()
		if !reflect.DeepEqual(len(f.data), got) {
			t.Errorf("want(%v); got(%v)", len(f.data), got)
		}
	}
}

func TestCardinality(t *testing.T) {
	for _, f := range testFiles {
		d := NewDataSet(f.fileName, f.separator, f.headerlns)
		d.Read()
		got := d.Cardinality()
		if !reflect.DeepEqual(f.cardin, got) {
			t.Errorf("want(%v); got(%v)", f.cardin, got)
		}
	}
}

func TestData(t *testing.T) {
	for _, f := range testFiles {
		d := NewDataSet(f.fileName, f.separator, f.headerlns)
		d.Read()
		got := d.Data()
		if !reflect.DeepEqual(f.data, got) {
			t.Errorf("want(%v); got(%v)", f.data, got)
		}
	}
}
