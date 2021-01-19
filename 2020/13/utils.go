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
func getBusArr(s *bufio.Scanner) []BusTimestamp {
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

	var busModInfoArr []BusTimestamp

	for idx, busPtr := range busArr {
		if busPtr == nil {
			continue
		}
		bus := *busPtr

		busInfo := BusTimestamp{
			bus: bus,
			timestamp: 0,
			index: idx,
		}

		busModInfoArr = append(busModInfoArr, busInfo)
	}

	return busModInfoArr
}

func getFile(fptr *string) *os.File {

	f, err := os.Open(*fptr)
	if err != nil {
		log.Fatal(err)
	}
	return f
}

func getMaxBus(busArr []BusTimestamp) BusTimestamp {
	maxBus := BusTimestamp{}
	for _, currBus := range busArr {
		if currBus.bus > maxBus.bus {
			maxBus = currBus
		}
	}
	return maxBus
}


func printAnswer(b BusTimestamp, busArr []BusTimestamp) {

	fmt.Printf("\n\ncheckBusTime():\tEARLIEST: %d FOR BUS %d\n", b.timestamp, b.bus)

	for i, busTimestamp := range busArr {
		if i == 0 {
			fmt.Printf("printAnswer():\tEARLIEST + mod %d for bus %d of index %d: %d\n", *(busArr[0].mod), busTimestamp.bus, busTimestamp.index, b.timestamp + *(busArr[0].mod))
		}
		fmt.Printf("printAnswer():\t%d\t%d\t%d\n", busTimestamp.index, busTimestamp.bus, *(busTimestamp.mod))
	}
}