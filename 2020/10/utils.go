package main

import (
	"bufio"
	"fmt"
	"log"
	"os"
	"strconv"
)



func mungeData(dataStr string) int {
	// do stuff here like
	dataInt, err := strconv.Atoi(dataStr)
	if err != nil {
		log.Fatalf("error converting %s to integer: %v\n", dataStr, err)
	}
	return dataInt
}

func makeDataArr(s *bufio.Scanner) []int {
	dataArr := make([]int, 0)
	for s.Scan() {
		dataLineStr := s.Text()
		mungedDataLine := mungeData(dataLineStr)
		dataArr = append(dataArr, mungedDataLine)
	}
	if err := s.Err(); err != nil {
		log.Fatal(err)
	}

	if VERBOSE {
		fmt.Printf("length of data array: %d\n", len(dataArr))
	}
	return dataArr
}

func getFile(fptr *string) *os.File {

	f, err := os.Open(*fptr)
	if err != nil {
		log.Fatal(err)
	}
	return f
}