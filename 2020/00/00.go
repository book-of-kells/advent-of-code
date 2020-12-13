package main

import (
	"bufio"
	"flag"
	"fmt"
)

var VERBOSE = false

func main() {
	fptr := flag.String("file", "input.txt", "file path to read from")
	vptr := flag.Bool("v", false, "verbose")
	flag.Parse()
	VERBOSE = *vptr

	f := getFile(fptr)
	defer f.Close()
	dataArr := makeDataArr(bufio.NewScanner(f))

	for i, elem := range dataArr {
		_ = fmt.Sprintf("\nelement %d\n%s\n", i, elem)
	}

}