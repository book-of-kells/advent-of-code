package main

import (
	"bufio"
	"flag"
	"fmt"
)

var VERBOSE = false

func main() {
	vptr := flag.Bool("v", false, "verbose")
	flag.Parse()
	VERBOSE = *vptr

	f := getFile()
	defer f.Close()
	dataArr := makeDataArr(bufio.NewScanner(f))

	numToSumTo := 0

	for i, elem := range dataArr {
		if i < 25 {
			continue
		}
		//fmt.Printf("\nelement %d\n%d\n", i, elem)
		if isEnd := isSumOfTwoOfPrev25Numbers(i, dataArr); isEnd == true {
			fmt.Printf("dataArr[%d] (%d) is not the sum of any two of the previous 25 numbers\n", i, elem)
			numToSumTo = elem
			break
		}
	}


	intArrChan := make(chan []int)
	go getSumOfArr(dataArr, numToSumTo, intArrChan)

	select {
	case intArr := <-intArrChan:
		fmt.Printf("%d = sum(%v)\n", numToSumTo, intArr)
		min := getMinElemInArr(intArr)
		max := getMaxElemInArr(intArr)
		fmt.Printf("min %d + max %d = %d\n", min, max, min+max)
	}
}


func getSumOfArr(dataArr []int, sumToThisNum int, sumChan chan []int) {
	for i, _ := range dataArr {
		acc := 0
		for acc < sumToThisNum {
			for j := i;j<len(dataArr);j++{
				acc += dataArr[j]
				if acc == sumToThisNum {
					sumChan <- dataArr[i:j+1]
				}
			}
		}
	}
}

func getMinElemInArr(arr []int) int {
	memo := arr[0]
	for i:=1;i<len(arr);i++ {
		if i+1 == len(arr) {
			return memo
		}
		if arr[i] < memo {
			memo = arr[i]
		}
	}
	return memo
}

func getMaxElemInArr(arr []int) int {
	memo := arr[0]
	for i:=1;i<len(arr);i++ {
		if i+1 == len(arr) {
			return memo
		}
		if arr[i] > arr[i-1] {
			memo = arr[i]
		}
	}
	return memo
}

func isSumOfTwoOfPrev25Numbers(index int, dataArr []int) bool {

	prev25Arr := dataArr[index-25:index]
	for i, outerElem := range prev25Arr {
		for _, innerElem := range prev25Arr[i+1:] {
			if outerElem + innerElem == dataArr[index] {
				return false
			}
		}
	}
	return true
}