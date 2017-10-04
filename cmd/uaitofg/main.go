package main

import (
	"fmt"
	"os"

	"github.com/britojr/kbn/mrf"
	"github.com/britojr/utl/errchk"
)

var (
	mkfile string
	fgfile string
)

func parseArgs() {
	if len(os.Args) < 3 {
		fmt.Println("Please enter both file names.")
		return
	}
	mkfile = os.Args[1]
	fgfile = os.Args[2]
}

func main() {
	parseArgs()
	// read MRF
	f, err := os.Open(mkfile)
	errchk.Check(err, fmt.Sprintf("Can't create file %v", mkfile))
	mk := mrf.LoadFromUAI(f)
	f.Close()
	if mk == nil {
		fmt.Printf("an error occurred while loading file %v\n", mkfile)
		return
	}

	// save FG
	f, err = os.Create(fgfile)
	errchk.Check(err, fmt.Sprintf("Can't create file %v", fgfile))
	mk.SaveOnLibdaiFormat(f)
	f.Close()
}
