package main

import (
	"bufio"
	"flag"
	"fmt"
	"log"
	"math"
	"net/http"
	_ "net/http/pprof"
	"time"
)

var VERBOSE = false


type BusTimestamp struct {
	bus int
	timestamp int
	index int
	modArr []int
}


func main() {
	startTime := time.Now()
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

	endTime := time.Now()
	elapsed := endTime.Sub(startTime)
	fmt.Printf("elapsed time: %v\n", elapsed)
}

func checkBusMultiples(busNum int, startAtTimestamp int, busArr []*int) {
	i := 1
	modCorrection := (busNum * i + startAtTimestamp) % busNum
	earliest := startAtTimestamp + (busNum * i) - modCorrection

	if VERBOSE {
		fmt.Printf("checkBusMultiples():\tfor bus %d starting time %d\n", busNum, earliest)
	}
	busPtr := &BusTimestamp{bus: busNum, timestamp: earliest}
	finished := checkBusTime(*busPtr, busArr, &busNum)

	elapsedTotal := time.Duration(0)

	for !finished {
		startTime := time.Now()
		earliest = earliest + busNum
		busPtr := &BusTimestamp{bus: busNum, timestamp: earliest}
		finished = checkBusTime(*busPtr, busArr, &busNum)
		i++
		endTime := time.Now()
		elapsed := endTime.Sub(startTime)
		elapsedTotal += elapsed
		if VERBOSE && i % 1000000 == 0 {
			// every 1,000,000 iterations, print an update
			avgElapsed := elapsedTotal.Nanoseconds()/int64(i)
			fmt.Printf("checkBusMultiples():\tfor bus %d checking time %d avg elapsed: %d ns\n", busNum, earliest, avgElapsed)
		}
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

// busModInfoArr is an array of BusTimestamp instances.
// The busModInfoArr doesn't have nil values because
// each BusTimestamp instance has its own index property.
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
	isTrue := true
	isFalse := false

	for _, firstBusMod := range firstBus.modArr {
		firstBusModDiff := firstBusMod - firstBus.index
		isConsecutiveArr := make([]*bool, len(busInfoArr[1:]))

		// check the nextBus in the busInfoArr
		for nextBusIdx, nextBus := range busInfoArr[1:] {
			if nextBusIdx > 0 && isConsecutiveArr[nextBusIdx-1] != nil && *(isConsecutiveArr[nextBusIdx-1]) == false {
				// if even one value in isConsecutiveArr is false, stop checking the rest of the `nextBus`es
				break
			}
			// check all the mods for the nextBus
			for _, nextBusMod := range nextBus.modArr {
				if nextBusMod - nextBus.index == firstBusModDiff {
					isConsecutiveArr[nextBusIdx] = &isTrue
					break // skip to checking the next nextBus
				} else {
					isConsecutiveArr[nextBusIdx] = &isFalse // continue checking this nextBus's mods
				}
				if nextBusMod - nextBus.index > firstBusModDiff {
					// each BusTimestamp's mod in its modArr increases with the modArr index
					// so if the current nextBusMod is so large that
					// nextBusMod - nextBus.index > firstBusModDiff
					// then there is no need to check the rest of the mods in the modArr
					if *(isConsecutiveArr[nextBusIdx]) == false {
						// this will always be the case
					}
					break
				}

			} // check the next mod of this nextBus

		} // check the next nextBus

		allConsecutive := true
		for _, isConsecutive := range isConsecutiveArr {
			if *isConsecutive == false {
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
