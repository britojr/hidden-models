// Package filehandler implements file handling
package filehandler

import (
	"bufio"
	"fmt"
	"os"
	"strings"

	"github.com/britojr/playgo/utils"
)

// DataSet ...
type DataSet struct {
	fileName       string
	delimiter      rune
	hasCardinality bool
	cardinality    []int
	data           [][]int
}

// NewDataSet creates new dataset
func NewDataSet(fileName string, delimiter rune, hasCardinality bool) *DataSet {
	d := &DataSet{}
	d.fileName = fileName
	d.delimiter = delimiter
	d.hasCardinality = hasCardinality
	return d
}

// Read reads the complete file
func (d *DataSet) Read() {
	file := openFile(d.fileName)
	defer file.Close()
	scanner := bufio.NewScanner(file)
	f := func(c rune) bool {
		return c == d.delimiter
	}
	if d.hasCardinality {
		scanner.Scan()
		cells := strings.FieldsFunc(scanner.Text(), f)
		d.cardinality = utils.SliceAtoi(cells)
	}
	for i := 0; scanner.Scan(); i++ {
		cells := strings.FieldsFunc(scanner.Text(), f)
		d.data = append(d.data, utils.SliceAtoi(cells))
	}
}

// Cardinality returns cardinality slice
func (d *DataSet) Cardinality() []int {
	return d.cardinality
}

// Data returns the whole dataset
func (d *DataSet) Data() [][]int {
	return d.data
}

func openFile(name string) *os.File {
	fp, err := os.Open(name)
	utils.ErrCheck(err, fmt.Sprintf("Can't open file %v", name))
	return fp
}
