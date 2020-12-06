package main

import "bufio"

func main() {
	f := getFile()
	defer f.Close()
	_ = makeDataArr(bufio.NewScanner(f))

	// do stuff with array of input data
}