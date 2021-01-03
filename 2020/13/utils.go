package main

import (
	"bufio"
	"fmt"
	"log"
	"os"
	"strconv"
	"strings"
)


func mungeData(dataStr string) string {
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

/*
 * If input is an array of strings such as [ "7","13","x","x","59","x","31","19" ]
 * then output is array of []*ints such as [ 7, 13, nil, nil, 59, nil, 31, 19 ]
 */
func getBusArr(s *bufio.Scanner) []*int {
	dataArr := makeDataArr(s)
	var busArr []*int
	for _, busNumStr := range strings.Split(dataArr[1], ",") {
		if busNumStr == "x" {
			busArr = append(busArr, nil)
			continue
		}
		busInt, err := strconv.Atoi(busNumStr)
		if err != nil {
			log.Fatalf("Error converting %s to integer: %v\n", busNumStr, err)
		}
		busArr = append(busArr, &busInt)
	}
	return busArr
}


func getChanArr(busArr []*int, minBusIdx *int) []chan *BusTimestamp {
	var chanArr []chan *BusTimestamp
	for currIdx, bus := range busArr {
		if bus == nil {
			chanArr = append(chanArr, nil)
			continue
		}
		if minBusIdx != nil && *minBusIdx != currIdx {
			chanArr = append(chanArr, nil)
			continue
		}
		b := make(chan *BusTimestamp)
		chanArr = append(chanArr, b)
	}
	return chanArr
}

func getFile(fptr *string) *os.File {

	f, err := os.Open(*fptr)
	if err != nil {
		log.Fatal(err)
	}
	return f
}

func getMinBus(busArr  []*int ) (int, int){
	minBus := 1000
	minIdx := -1
	for idx, bus := range busArr {
		if bus == nil {
			continue
		}
		if *bus < minBus {
			minBus = *bus
			minIdx = idx
		}
	}
	return minBus, minIdx
}


func printArr(busModArr []*[]int) {
	for idx, busMod := range busModArr {
		if busMod == nil {
			continue
		}
		fmt.Printf("%d\t%v\n", idx, *busMod)
	}
}
