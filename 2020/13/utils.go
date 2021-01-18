package main

import (
	"bufio"
	"fmt"
	"log"
	"os"
	"strconv"
	"strings"
	"syscall"
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

func getFile(fptr *string) *os.File {

	f, err := os.Open(*fptr)
	if err != nil {
		log.Fatal(err)
	}
	return f
}

func getMaxBus(busArr []*int) (int, int){
	maxBus := 1
	maxIdx := -1
	for idx, bus := range busArr {
		if bus == nil {
			continue
		}
		if *bus > maxBus {
			maxBus = *bus
			maxIdx = idx
		}
	}
	return maxBus, maxIdx
}


func printAnswer(b BusTimestamp, busArr []*int, maxBusNum *int, firstBusMod *int) {
	if firstBusMod == nil {
		fmt.Printf("\n\ncheckBusTime():\tEARLIEST: %d FOR BUS %d\n", b.timestamp, b.bus)
		fmt.Printf("checkBusTime():\tHOWEVER, firstBusMod is %v, so exiting program.\n", firstBusMod)
		syscall.Exit(1)
	}
	busModInfoArr := getBusModInfoArr(b.timestamp, busArr, maxBusNum)
	fmt.Printf("\n\ncheckBusTime():\tEARLIEST: %d FOR BUS %d\n", b.timestamp, b.bus)

	for i, busInfo := range busModInfoArr {
		if i == 0 {
			fmt.Printf("printAnswer():\tEARLIEST + mod %d for bus %d of index %d: %d\n", *firstBusMod, busInfo.bus, busInfo.index, b.timestamp + *firstBusMod)
		}
		fmt.Printf("printAnswer():\t%d\t%d\t%v*\n", busInfo.index, busInfo.bus, busInfo.modArr)

	}
}