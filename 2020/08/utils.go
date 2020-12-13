package main

import (
	"bufio"
	"flag"
	"fmt"
	"log"
	"os"
	"strconv"
	"strings"
)

type LineCommand struct {
	index int
	command string
	direction string // todo later
	num int
	executed bool
}

func resetAll(arr []*LineCommand) {
	for _, cmd := range arr {
		cmd.executed = false
	}
}

func mungeData(dataStr string, index int) *LineCommand {
	// do stuff here like
	//dataInt, err := strconv.Atoi(dataStr)
	arr := strings.Split(dataStr, " ")
	numStr := arr[1][1:]
	num, err := strconv.Atoi(numStr)
	if err != nil {
		log.Fatalf("couldn't parse number from %s: %v\n", numStr, err)
	}
	dirStr := arr[1][0:1] // either '+' or '-'
	if dirStr != "+" && dirStr != "-" {
		log.Fatalf("direction should be '+' or '-' but instead was %s\n", dirStr)
	}
	lineCmd := &LineCommand{
		index: index,
		command: arr[0],
		direction: dirStr,
		num: num,
		executed: false,
	}


	return lineCmd
}

func makeDataArr(s *bufio.Scanner) []*LineCommand {
	dataArr := make([]*LineCommand, 0)
	i := 0
	for s.Scan() {
		dataLineStr := s.Text()
		mungedDataLine := mungeData(dataLineStr, i)
		dataArr = append(dataArr, mungedDataLine)
		i++
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