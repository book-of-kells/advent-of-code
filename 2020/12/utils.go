package main

import (
	"bufio"
	"fmt"
	"log"
	"os"
	"strconv"
)

type Command struct {
	cmd string
	num int
}

func mungeData(dataStr string) Command {
	dataInt, err := strconv.Atoi(dataStr[1:])
	if err != nil {
		log.Fatalf("Error converting %s to int: %v\n", dataStr[1:], err)
	}
	return Command{
		cmd: string(dataStr[0]),
		num: dataInt,
	}
}

func makeDataArr(s *bufio.Scanner) []Command {
	dataArr := make([]Command, 0)
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