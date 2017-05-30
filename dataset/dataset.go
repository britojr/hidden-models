package dataset

import (
	"bufio"
	"fmt"
	"io"
	"log"
	"strings"
	"time"

	"github.com/britojr/kbn/assignment"
	"github.com/britojr/kbn/list"
	"github.com/britojr/kbn/utl"
	"github.com/britojr/kbn/utl/conv"
	"github.com/willf/bitset"
)

// HdrFlags type is used to store the flags indicating the kind header lines the file has
// each kind of header line must appear in the fixed order of the constants declared bellow
type HdrFlags uint

const (
	// HdrName indicates that there is a line with names of variables
	HdrName HdrFlags = 1 << iota
	//HdrCardin indicates that there is a line with cardinality of each variable
	HdrCardin
	//HdrNameCard indicates that there is one line in the format <Name>_<Card>
	HdrNameCard
)

// Dataset is used  to read a file and handle the data
// including counting occurrences of sets of variables
type Dataset struct {
	delimiter rune
	headerlns HdrFlags

	varNames []string
	cardin   []int // cardinality of each variable
	data     [][]int

	values []*valToLine     // all assignable values for each variable
	cache  map[string][]int // cached occurrence counting slices for different varlists
}

type valToLine map[int]*bitset.BitSet

// New creates new dataset from a reader
func New(r io.Reader, delimiter rune, headerlns HdrFlags) *Dataset {
	d := new(Dataset)
	d.delimiter = delimiter
	d.headerlns = headerlns
	d.read(r)
	if len(d.cardin) == 0 {
		d.calcCardinality()
	}
	d.initCount()
	return d
}

// NewFromFile creates new from a given file
func NewFromFile(fileName string, delimiter rune, headerlns HdrFlags) *Dataset {
	log.Printf("Loading dataset: %v\n", fileName)
	start := time.Now()
	f := utl.OpenFile(fileName)
	defer f.Close()
	d := New(f, delimiter, headerlns)
	elapsed := time.Since(start)
	log.Printf("Variables: %v, Instances: %v, Time: %v\n",
		d.NCols(), d.NLines(), elapsed,
	)
	return d
}

// read the complete file and stores the data in memory
func (d *Dataset) read(r io.Reader) {
	splitFunc := func(c rune) bool {
		return c == d.delimiter
	}
	scanner := bufio.NewScanner(r)
	if d.headerlns&HdrNameCard > 0 {
		scanner.Scan()
		nameCard := strings.FieldsFunc(scanner.Text(), splitFunc)
		for _, v := range nameCard {
			x := strings.FieldsFunc(v, func(c rune) bool {
				return c == '_'
			})
			d.varNames = append(d.varNames, x[0])
			d.cardin = append(d.cardin, conv.Atoi(x[1]))
		}
	} else {
		if d.headerlns&HdrName > 0 {
			scanner.Scan()
			d.varNames = strings.FieldsFunc(scanner.Text(), splitFunc)
		}
		if d.headerlns&HdrCardin > 0 {
			scanner.Scan()
			d.cardin = conv.Satoi(strings.FieldsFunc(scanner.Text(), splitFunc))
		}
	}
	for i := 0; scanner.Scan(); i++ {
		cells := strings.FieldsFunc(scanner.Text(), splitFunc)
		d.data = append(d.data, conv.Satoi(cells))
	}
}

// initCount initializes value counter
func (d *Dataset) initCount() {
	lin, col := d.NLines(), d.NCols()
	d.values = make([]*valToLine, col)
	for i, c := range d.cardin {
		d.values[i] = new(valToLine)
		*d.values[i] = make(map[int]*bitset.BitSet)
		for j := 0; j < c; j++ {
			(*d.values[i])[j] = bitset.New(uint(lin))
		}
	}
	for i := 0; i < lin; i++ {
		for j := 0; j < col; j++ {
			(*d.values[j])[d.data[i][j]].Set(uint(i))
		}
	}
	// initialize empty cache
	d.cache = make(map[string][]int)
}

// Cardin returns cardinality slice
func (d *Dataset) Cardin() []int {
	return d.cardin
}

// Data returns the whole dataset as an int matrix
func (d *Dataset) Data() [][]int {
	return d.data
}

// NLines returns the number of lines of the stored dataset
func (d *Dataset) NLines() int {
	return len(d.data)
}

// NCols returns the number of columns of the dataset stored
func (d *Dataset) NCols() int {
	return len(d.data[0])
}

// calcCardinality calculates the cardinality of each variable by scaning the dataset
func (d *Dataset) calcCardinality() {
	d.cardin = make([]int, len(d.data[0]))
	for j := range d.cardin {
		// m := make(map[int]bool)
		// for i := 0; i < len(d.data); i++ {
		// 	m[d.data[i][j]] = true
		// }
		// d.cardin[j] = len(m)
		m := 0
		for i := 0; i < len(d.data); i++ {
			if d.data[i][j] > m {
				m = d.data[i][j]
			}
		}
		d.cardin[j] = m + 1
	}
}

// Count returns number of occurrences for a particular assignment
func (d *Dataset) Count(assig *assignment.Assignment) (n int, ok bool) {
	setlist := make([]*bitset.BitSet, 0, len(assig.Variables()))
	for i := range assig.Variables() {
		if assig.Var(i) < len(d.cardin) {
			setlist = append(setlist, (*d.values[assig.Var(i)])[assig.Value(i)])
		}
	}
	if len(setlist) > 0 {
		return int(list.IntersectionBits(setlist).Count()), true
	}
	return -1, false
}

// CountAssignments returns a slice with the counting
// of each possible assignment of the given set of variables
func (d *Dataset) CountAssignments(varlist []int) (v []int) {
	if len(varlist) <= 0 {
		return
	}
	strvarlist := fmt.Sprint(varlist)
	v, ok := d.cache[strvarlist]
	if !ok {
		assig := assignment.New(varlist, d.cardin)
		for assig.Next() {
			if count, ok := d.Count(assig); ok {
				v = append(v, count)
			} else {
				break
			}
		}
		d.cache[strvarlist] = v
	}
	return
}
