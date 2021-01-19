package main

import (
	"bufio"
	"flag"
	"fmt"
	"log"
	"net/http"
	_ "net/http/pprof"
	"time"
)

var VERBOSE = false
var SHOWTIMES = false

type BusTimestamp struct {
	bus int
	timestamp int
	index int
	mod *int
}


func main() {
	startTime := time.Now()
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

	busArr := getBusArr(bufio.NewScanner(f))

	maxBus := getMaxBus(busArr)

	checkBusMultiples(maxBus, *startAtTimestampPtr, busArr)

	endTime := time.Now()
	elapsed := endTime.Sub(startTime)
	fmt.Printf("elapsed time: %v\n", elapsed)
}

func checkBusMultiples(maxBus BusTimestamp, startAtTimestamp int, busArr []BusTimestamp) {
	i := 1
	modCorrection := (maxBus.bus * i + startAtTimestamp) % maxBus.bus
	maxBus.timestamp = startAtTimestamp + (maxBus.bus * i) - modCorrection

	if VERBOSE {
		fmt.Printf("checkBusMultiples():\tfor bus %d starting time %d\n", maxBus.bus, maxBus.timestamp)
	}

	finished := checkBusTimes(maxBus, busArr)

	if finished == true {
		printAnswer(maxBus, busArr)
	}

	elapsedTotal := time.Duration(0)
	startTime := time.Now()

	for !finished {
		if SHOWTIMES {
			startTime = time.Now()
		}

		maxBus.timestamp += maxBus.bus
		finished = checkBusTimes(maxBus, busArr)

		if finished == true {
			printAnswer(maxBus, busArr)
		}

		i++
		if SHOWTIMES {
			endTime := time.Now()
			elapsed := endTime.Sub(startTime)
			elapsedTotal += elapsed
		}

		if VERBOSE && i % 100_000_000 == 0 {
			// every 100,000,000 iterations, print an update

			if SHOWTIMES {

				avgElapsed := elapsedTotal.Nanoseconds()/int64(i)
				fmt.Printf("checkBusMultiples():\tfor bus %d checking time %d avg elapsed: %d ns\n", maxBus.bus, maxBus.timestamp, avgElapsed)
			} else {
				fmt.Printf("checkBusMultiples():\tfor bus %d checking time %d\n", maxBus.bus, maxBus.timestamp)
			}
			fmt.Printf("checkBusMultiples():\tfor bus %d checking time %d\n", maxBus.bus, maxBus.timestamp)

		}
	}
}


func shouldContinue(busMod int, currBus BusTimestamp, maxBus BusTimestamp) bool {
	maxBusModDiff := maxBus.bus - maxBus.index
	return busMod - currBus.index < maxBusModDiff

}

func isMatch(busMod int, currBus BusTimestamp, maxBus BusTimestamp) bool {
	maxBusModDiff := maxBus.bus - maxBus.index
	return busMod - currBus.index == maxBusModDiff
}

// busModInfoArr is an array of BusTimestamp instances.
// The busModInfoArr doesn't have nil values because
// each BusTimestamp instance has its own index property.
func checkBusTimes(maxBus BusTimestamp, busArr []BusTimestamp) bool {

	for currIdx, currBus := range busArr {
		currBus.mod = nil
		busArr[currIdx].mod = nil
		j := 1

		busMod := currBus.bus * j - maxBus.timestamp % currBus.bus

		if isMatch(busMod, currBus, maxBus) {
			busArr[currIdx].mod = &busMod
			continue
		}

		for shouldContinue(busMod, currBus, maxBus) {
			busMod = currBus.bus * j - maxBus.timestamp % currBus.bus
			if isMatch(busMod, currBus, maxBus) {
				busArr[currIdx].mod = &busMod
				break // this bus is consecutive; move on to the next bus
			}

			if !shouldContinue(busMod, currBus, maxBus) {
				return false
			}
			j++
		}

		if busArr[currIdx].mod == nil {
			return false
		}
	}
	return true
}