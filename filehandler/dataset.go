// Package filehandler implements file handling
package filehandler

import (
	"bufio"
	"fmt"
	"os"
	"strings"

	"github.com/britojr/kbn/utils"
)

// DataHandler interface, returns data as an int matrix
type DataHandler interface {
	Data() [][]int
}

// HeaderFlags type is used to store the flags indicating the kind header lines the file has
// each kind of header line must appear in the fixed order of the constants declared bellow
type HeaderFlags byte

const (
	// NameHeader indicates that there is a line with names of variables
	NameHeader HeaderFlags = 1 << iota
	//CardinHeader indicates that there is a line with cardinality of each variable
	CardinHeader
)

// DataSet is used  to read a file and store the data
type DataSet struct {
	fileName    string
	delimiter   rune
	headerlns   HeaderFlags
	splitFunc   func(c rune) bool
	varNames    []string
	cardinality []int
	data        [][]int
}

// NewDataSet creates new dataset with the given parametes
func NewDataSet(fileName string, delimiter rune, headerlns HeaderFlags) (d *DataSet) {
	d = new(DataSet)
	d.fileName = fileName
	d.delimiter = delimiter
	d.headerlns = headerlns
	d.splitFunc = func(c rune) bool {
		return c == d.delimiter
	}
	return
}

// Read reads the complete file and stores the data in memory
func (d *DataSet) Read() {
	file := openFile(d.fileName)
	defer file.Close()
	scanner := bufio.NewScanner(file)
	if d.headerlns&NameHeader != 0 {
		scanner.Scan()
		d.varNames = strings.FieldsFunc(scanner.Text(), d.splitFunc)
	}
	if d.headerlns&CardinHeader != 0 {
		scanner.Scan()
		cells := strings.FieldsFunc(scanner.Text(), d.splitFunc)
		d.cardinality = utils.SliceAtoi(cells)
	}
	for i := 0; scanner.Scan(); i++ {
		cells := strings.FieldsFunc(scanner.Text(), d.splitFunc)
		d.data = append(d.data, utils.SliceAtoi(cells))
	}
}

// Cardinality returns cardinality slice
func (d *DataSet) Cardinality() []int {
	if len(d.cardinality) == 0 {
		d.calcCardinality()
	}
	return d.cardinality
}

// Data returns the whole dataset
func (d *DataSet) Data() [][]int {
	return d.data
}

// Size returns the number of lines on the dataset stored
func (d *DataSet) Size() int {
	return len(d.data)
}

func (d *DataSet) calcCardinality() {
	d.cardinality = make([]int, len(d.data[0]))
	for j := 0; j < len(d.data[0]); j++ {
		// m := make(map[int]bool)
		// for i := 0; i < len(d.data); i++ {
		// 	m[d.data[i][j]] = true
		// }
		// d.cardinality[j] = len(m)
		m := 0
		for i := 0; i < len(d.data); i++ {
			if d.data[i][j] > m {
				m = d.data[i][j]
			}
		}
		d.cardinality[j] = m + 1
	}
}

func openFile(name string) *os.File {
	fp, err := os.Open(name)
	utils.ErrCheck(err, fmt.Sprintf("Can't open file %v", name))
	return fp
}
