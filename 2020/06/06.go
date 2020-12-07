package main

import (
	"bufio"
	"fmt"
)

func main() {
	f := getFile()
	defer f.Close()
	dataArr := makeDataArr(bufio.NewScanner(f))

	totalYes := 0
	for groupIdx, groupSet := range dataArr {
		totalYes += groupSet.Len()
		fmt.Printf("groupIdx %d answered yes to %d answers\n", groupIdx, groupSet.Len())
	}
	fmt.Printf("totalYes is %d\n", totalYes)
}