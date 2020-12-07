package main

import (
	"bufio"
	"flag"
	"fmt"
	"log"
	"os"
	"strconv"
)

// helper func for debugging
func getFormatters (n int) (string, string) {
	tabs := ""
	for j:=1; j<n; j++ {
		tabs += "\t"
	}

	ordinal := "th"
	switch n % 10 {
	case 1:
		ordinal = "st"
	case 2:
		ordinal = "nd"
	case 3:
		ordinal = "rd"
	default:
		ordinal = "th"
	}

	return tabs, ordinal
}

// helper func for debugging
func getSumString(n int, mthInt int, nthInt int, sumOfInts int) string {
	mthAndNthStr := fmt.Sprintf("%d", nthInt)
	if n > 1 {
		mthAndNthStr = fmt.Sprintf("%d + %d", mthInt, nthInt)
	}

	for i := n; i<NUMBER_OF_INTEGERS_TO_FIND; i++ {
		mthAndNthStr += " + _"
	}

	mthAndNthStr += fmt.Sprintf(" = %d", sumOfInts)
	return mthAndNthStr
}

func mungeData(dataStr string) int {
	currInt, err := strconv.Atoi(dataStr)
	if err != nil {
		log.Fatalf("err trying to parse data %s: %v\n", dataStr, err)
	}
	return currInt
}

// gets a line of text from a scanner, converts it to an int, and returns it (if err == nil)
func getCurrIntFromFromScannerText(s *bufio.Scanner) int {
	currNumStr := s.Text()
	currInt := mungeData(currNumStr)
	return currInt
}

func makeDataArr(s *bufio.Scanner) []int {
	dataArr := make([]int, 0)
	for s.Scan() {
		// these two are only here in the outermost func
		potentialFirstInt := getCurrIntFromFromScannerText(s)
		dataArr = append(dataArr, potentialFirstInt)
	}
	if err := s.Err(); err != nil {
		log.Fatal(err)
	}
	return dataArr
}

func getFile() *os.File {
	fptr := flag.String("file", "input.txt", "file path to read from")
	flag.Parse()

	f, err := os.Open(*fptr)
	if err != nil {
		log.Fatal(err)
	}
	return f
}
