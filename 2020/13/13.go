package main

import (
	"bufio"
	"flag"
	"fmt"
	"log"
	"net/http"
	_ "net/http/pprof"
)

var VERBOSE = false
var SHOWTIMES = false

type BusTimestamp struct {
	bus int
	timestamp int
	index int
	multipleOfMaxBus int
	minInflation int
	modArr []int
}


func main() {
	//startTime := time.Now()
	fptr := flag.String("file", "input.txt", "file path to read from")
	vptr := flag.Bool("v", false, "verbose")
	pprofPtr := flag.Bool("pprof", false, "profiling")
	showTimePtr := flag.Bool("time", false, "show runtimes")
	startAtTimestampPtr := flag.Int("min", 1, "start checking at this timestamp")
	flag.Parse()
	VERBOSE = *vptr
	SHOWTIMES = *showTimePtr

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

	busArr, maxBusNum := getBusArr(bufio.NewScanner(f))

	checkBusMultiples(maxBusNum, *startAtTimestampPtr, busArr)

	//endTime := time.Now()
	//elapsed := endTime.Sub(startTime)
	//fmt.Printf("elapsed time: %v\n", elapsed)
}

func checkBusMultiples(maxBusNum int, startAtTimestamp int, busArr []BusTimestamp) {
	i := 1
	modCorrection := (maxBusNum * i + startAtTimestamp) % maxBusNum
	earliest := startAtTimestamp + (maxBusNum * i) - modCorrection

	if VERBOSE {
		fmt.Printf("checkBusMultiples():\tfor bus %d starting time %d\n", maxBusNum, earliest)
	}
	busPtr := &BusTimestamp{bus: maxBusNum, timestamp: earliest}
	finished := checkBusTime(*busPtr, busArr)

	//elapsedTotal := time.Duration(0)
	//startTime := time.Now()

	for !finished {
		//if SHOWTIMES {
		//	startTime = time.Now()
		//}
		earliest = earliest + maxBusNum
		busPtr := &BusTimestamp{bus: maxBusNum, timestamp: earliest}
		finished = checkBusTime(*busPtr, busArr)
		i++
		//if SHOWTIMES {
		//	endTime := time.Now()
		//	elapsed := endTime.Sub(startTime)
		//	elapsedTotal += elapsed
		//}
		if VERBOSE && i % 1000000 == 0 {
			// every 1,000,000 iterations, print an update

			//if SHOWTIMES {
			//	avgElapsed := elapsedTotal.Nanoseconds()/int64(i)
			//	fmt.Printf("checkBusMultiples():\tfor bus %d checking time %d avg elapsed: %d ns\n", maxBusNum, earliest, avgElapsed)
			//} else {
			//	fmt.Printf("checkBusMultiples():\tfor bus %d checking time %d\n", maxBusNum, earliest)
			//}
			fmt.Printf("checkBusMultiples():\tfor bus %d checking time %d\n", maxBusNum, earliest)

		}
	}
}

func checkBusTime(b BusTimestamp, busArr []BusTimestamp) bool {
	fillBusModArrays(b.timestamp, busArr)

	finished, firstBusMod := busesAreConsecutive(busArr)

	if finished == true {
		printAnswer(b, busArr, firstBusMod)
	}
	return finished
}

// busModInfoArr is an array of BusTimestamp instances.
// The busModInfoArr doesn't have nil values because
// each BusTimestamp instance has its own index property.
func fillBusModArrays(earliest int, busArr []BusTimestamp) { // []BusTimestamp {

	for _, busTimestamp := range busArr {

		// get multipleOfMaxBus possible values to compare against others
		for j := 1; j<busTimestamp.multipleOfMaxBus+1; j++ {
			minuend := busTimestamp.minInflation * j * busTimestamp.bus
			subtrahend := earliest % busTimestamp.bus
			busTimestamp.modArr[j-1] = minuend - subtrahend
		}
	}
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
