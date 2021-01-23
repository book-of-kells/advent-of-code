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

type BusData struct {
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

	solve(maxBus, *startAtTimestampPtr, busArr)

	endTime := time.Now()
	elapsed := endTime.Sub(startTime)
	fmt.Printf("elapsed time: %v\n", elapsed)
}

func solve(maxBus BusData, startAtTimestamp int, busArr []BusData) {

	modCorrection := (maxBus.bus + startAtTimestamp) % maxBus.bus
	maxBus.timestamp = startAtTimestamp + maxBus.bus - modCorrection

	if VERBOSE {
		fmt.Printf("solve():\tfor bus %d starting time %d\n", maxBus.bus, maxBus.timestamp)
	}

	elapsedTotal := time.Duration(0)
	startTime := time.Now()
	finished := false
	i := 1

	for !finished {
		if SHOWTIMES {
			startTime = time.Now()
		}

		finished = checkBusTimes(maxBus, busArr)
		if finished == true {
			printAnswer(maxBus, busArr)
		}

		maxBus.timestamp += maxBus.bus

		if SHOWTIMES {
			endTime := time.Now()
			elapsed := endTime.Sub(startTime)
			elapsedTotal += elapsed
		}

		// every 100,000,000 iterations, print an update
		if VERBOSE && i % 100_000_000 == 0 {
			if SHOWTIMES {
				avgElapsed := elapsedTotal.Nanoseconds()/int64(i)
				fmt.Printf("solve():\tfor bus %d checking time %d avg elapsed: %d ns\n", maxBus.bus, maxBus.timestamp, avgElapsed)
			} else {
				fmt.Printf("solve():\tfor bus %d checking time %d\n", maxBus.bus, maxBus.timestamp)
			}
		}
		i++
	}
}


func checkBusTimes(maxBus BusData, busArr []BusData) bool {
	maxBusModDiff := maxBus.bus - maxBus.index

	for currIdx, currBus := range busArr {
		originalBusMod := currBus.bus - maxBus.timestamp % currBus.bus
		currBus.mod = nil
		busArr[currIdx].mod = nil

		if (maxBusModDiff - (originalBusMod - currBus.index)) % currBus.bus != 0 {
			return false
		}

		busArr[currIdx].mod = &originalBusMod
		busArr[currIdx].timestamp = maxBus.timestamp + originalBusMod

		if currIdx != 0 {
			continue
		}

		j := 1
		busMod := currBus.bus * j - maxBus.timestamp % currBus.bus
		for busMod - currBus.index <= maxBusModDiff {
			if busMod - currBus.index == maxBusModDiff {
				busArr[currIdx].mod = &busMod
				busArr[currIdx].timestamp = maxBus.timestamp + busMod
				break // this bus is consecutive; move on to the next bus
			}
			busMod = currBus.bus * j - maxBus.timestamp % currBus.bus
			j++
		}
	}
	return true
}