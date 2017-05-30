package main

import (
	"bufio"
	"fmt"
	"os"
	"strings"

	"github.com/britojr/kbn/utl/conv"
	"github.com/britojr/kbn/utl/errchk"
)

var (
	file1 string
	file2 string
)

func parseArgs() {
	if len(os.Args) < 3 {
		fmt.Println("Please enter both file names.")
		return
	}
	file1 = os.Args[1]
	file2 = os.Args[2]
}

func main() {
	parseArgs()
	a := loadMarg(file1)
	b := loadMarg(file2)
	fmt.Printf("%v\n", margVariance(a, b))
}

func loadMarg(name string) [][]float64 {
	f, err := os.Open(name)
	errchk.Check(err, fmt.Sprintf("Can't create file %v", name))

	scanner := bufio.NewScanner(f)
	scanner.Scan()
	scanner.Scan()
	line := strings.Fields(scanner.Text())
	values := make([][]float64, conv.Atoi(line[0]))
	j := 1
	for i := range values {
		n := conv.Atoi(line[j])
		j++
		for k := 0; k < n; k++ {
			values[i] = append(values[i], conv.Atof(line[j+k]))
		}
		j += n
	}
	return values
}

func margVariance(a, b [][]float64) float64 {
	c, d := 0, float64(0)
	for i := range a {
		for j, v := range a[i] {
			d += (v - b[i][j]) * (v - b[i][j])
			c++
		}
	}
	return d / float64(c)
}
