package main

import (
	"bufio"
	"flag"
	"fmt"
	"log"
	"math"
	"net/http"
	_ "net/http/pprof"
)

var VERBOSE = false


type BusTimestamp struct {
	bus int
	timestamp int
	elapsed int64
	index int
	modArr []int
}


func main() {

	fptr := flag.String("file", "input.txt", "file path to read from")
	vptr := flag.Bool("v", false, "verbose")
	pprofPtr := flag.Bool("pprof", false, "profiling")
	startAtTimestampPtr := flag.Int("min", 1, "start checking at this timestamp")
	flag.Parse()
	VERBOSE = *vptr

	if *pprofPtr == true {
		if VERBOSE {
			fmt.Println("Running pprof on http://localhost:6060")
		}
		go func() {
			log.Println(http.ListenAndServe("localhost:6060", nil))
		}()
	}

	f := getFile(fptr)
	defer f.Close()

	busArr := getBusArr(bufio.NewScanner(f))
	maxBus, _ := getMaxBus(busArr)

	checkBusMultiples(maxBus, *startAtTimestampPtr, busArr)
}

func checkBusMultiples(busNum int, startAtTimestamp int, busArr []*int) {
	i := 1
	modCorrection := (busNum * i + startAtTimestamp) % busNum
	earliest := startAtTimestamp + (busNum * i) - modCorrection

	if VERBOSE {
		fmt.Printf("checkBusMultiples():\tfor bus %d starting time %d\n", busNum, earliest)
	}
	busPtr := &BusTimestamp{bus: busNum, timestamp: earliest, elapsed: 0}
	finished := checkBusTime(*busPtr, busArr, &busNum)

	for !finished {
		if VERBOSE && i % 100000 == 0 {
			// every 10,000 iterations, print an update
			fmt.Printf("checkBusMultiples():\tfor bus %d checking time %d\n", busNum, earliest)
		}
		earliest = earliest + busNum
		busPtr := &BusTimestamp{bus: busNum, timestamp: earliest, elapsed: 0}
		finished = checkBusTime(*busPtr, busArr, &busNum)
		i++
	}
}

func checkBusTime(b BusTimestamp, busArr []*int, maxBusNum *int) bool {
	busModInfoArr := getBusModInfoArr(b.timestamp, busArr, maxBusNum)

	finished, firstBusMod := busesAreConsecutive(busModInfoArr)

	if finished == true {
		printAnswer(b, busArr, maxBusNum, firstBusMod)
	}
	return finished
}

func getBusModInfoArr(earliest int, busArr []*int, maxBusNum *int) []BusTimestamp {
	var busModInfoArr []BusTimestamp

	for idx, busPtr := range busArr {
		if busPtr == nil {
			continue
		}
		bus := *busPtr
		m := int(math.Ceil(float64(*maxBusNum)/float64(bus)))
		busInfo := BusTimestamp{
			bus: bus,
			timestamp: earliest,
			elapsed: 0,
			index: idx,
			modArr: make([]int, m),
		}

		i := 1
		minInflation := int(math.Ceil(float64(idx/bus))) + i

		for minInflation*bus < idx {
			i++
			minInflation = int(math.Ceil(float64(idx/bus))) + i
		}

		// get m possible values to compare against others
		for j := 1; j<m+1; j++ {
			minuend := minInflation * j * bus
			subtrahend := earliest % bus
			busMod := minuend - subtrahend
			busInfo.modArr[j-1] = busMod
		}
		busModInfoArr = append(busModInfoArr, busInfo)
	}
	return busModInfoArr
}

func busesAreConsecutive(busInfoArr []BusTimestamp) (bool, *int) {

	firstBus := busInfoArr[0] // assumes first element in busArr != nil

	for _, firstBusMod := range firstBus.modArr {
		firstBusModDiff := firstBusMod - firstBus.index
		boolArr := make([]bool, len(busInfoArr[1:]))

		// check the nextBus in the busInfoArr
		for nextBusIdx, nextBus := range busInfoArr[1:] {

			// check all the mods for the nextBus
			for _, nextBusMod := range nextBus.modArr {
				if nextBusMod - nextBus.index == firstBusModDiff {
					boolArr[nextBusIdx] = true
					break // skip to checking the next nextBus
				}
				if nextBusMod - nextBus.index > firstBusModDiff {
					boolArr[nextBusIdx] = false // continue checking this nextBus's mods
				}
			} // check the next mod of this nextBus

		} // check the next nextBus

		allConsecutive := true
		for _, isConsecutive := range boolArr {
			if isConsecutive == false {
				allConsecutive = false
				break
			}
		}
		if allConsecutive == true {
			return allConsecutive, &firstBusMod // return true if all of the buses are consecutive
		}
	}
	return false, nil
}