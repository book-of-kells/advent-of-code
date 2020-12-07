package main

import (
	"bufio"
	"flag"
	"fmt"
	"log"
	"os"
)

func mungeData(dataStr string) string {
	// do stuff here like
	//dataInt, err := strconv.Atoi(dataStr)
	return dataStr
}

func makeDataArr(s *bufio.Scanner) []string {
	dataArr := make([]string, 0)
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

func getFile() *os.File {
	fptr := flag.String("file", "input.txt", "file path to read from")
	flag.Parse()

	f, err := os.Open(*fptr)
	if err != nil {
		log.Fatal(err)
	}
	return f
}