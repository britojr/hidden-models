/*
run experiments in batch

*/
package main

import (
	"fmt"
	"path/filepath"
)

func main() {
	fmt.Println("running tests...")
	files, _ := filepath.Glob("*")
	fmt.Println(files)
}
