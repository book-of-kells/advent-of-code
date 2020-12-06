package main

import (
	"bufio"
	"flag"
	"log"
	"os"
)

func mungeData(dataStr string) interface{} {
	// do stuff here like
	//dataInt, err := strconv.Atoi(dataStr)
	return dataStr
}

func makeDataArr(s *bufio.Scanner) []interface{} {
	dataArr := make([]interface{}, 0)
	for s.Scan() {
		dataLineStr := s.Text()
		mungedDataLine := mungeData(dataLineStr)
		dataArr = append(dataArr, mungedDataLine)
	}
	if err := s.Err(); err != nil {
		log.Fatal(err)
	}
	return dataArr
}

func getFile() *os.File {
	fptr := flag.String(".", "input.txt", "file path to read from")
	flag.Parse()

	f, err := os.Open(*fptr)
	if err != nil {
		log.Fatal(err)
	}
	return f
}