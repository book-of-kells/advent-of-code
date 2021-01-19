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


func getBusArr(s *bufio.Scanner) []BusData {
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

	var busModInfoArr []BusData

	for idx, busPtr := range busArr {
		if busPtr == nil {
			continue
		}
		bus := *busPtr

		busInfo := BusData{
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


func getMaxBus(busArr []BusData) BusData {
	maxBus := BusData{}
	for _, currBus := range busArr {
		if currBus.bus > maxBus.bus {
			maxBus = currBus
		}
	}
	return maxBus
}


func printAnswer(b BusData, busArr []BusData) {

	fmt.Printf("\n\ncheckBusTime():\tEARLIEST: %d FOR BUS %d\n", b.timestamp, b.bus)
	fmt.Printf("printAnswer():\t%d + mod %d for bus %d of index %d = %d\n", b.timestamp, *(busArr[0].mod), busArr[0].bus, busArr[0].index, b.timestamp + *(busArr[0].mod))

	fmt.Printf("printAnswer():\tidx\tbusNum\tmod\ttimestamp\n")
	for _, busData := range busArr {
		fmt.Printf("printAnswer():\t%d\t%d\t%d\t%d\n", busData.index, busData.bus, *(busData.mod), busData.timestamp)
	}
}