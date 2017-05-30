package utl

import (
	"fmt"
	"os"

	"github.com/britojr/kbn/utl/errchk"
)

// OpenFile returns a pointer to an open file
func OpenFile(name string) *os.File {
	f, err := os.Open(name)
	errchk.Check(err, fmt.Sprintf("Can't open file %v", name))
	return f
}
