package main

import (
	"fmt"

	"github.com/britojr/playgo/jtree"
	"github.com/britojr/tcc/characteristic"
)

func main() {
	fmt.Println("Generate JTree")
	jt := jtree.Generate(
		&characteristic.Tree{P: []int{-1, 0, 0, 2, 3}, L: []int{-1, -1, -1, 2, 2}},
		[]int{4, 1, 6, 3, 0, 2, 5},
	)
	fmt.Printf("JT: %v\n", jt)
}
