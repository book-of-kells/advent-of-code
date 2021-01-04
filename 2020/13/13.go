package main

import (
	"bufio"
	"flag"
	"fmt"
	"math"
	"syscall"
	"time"
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
	minptr := flag.Int("min", 1, "start checking at this timestamp")
	flag.Parse()
	VERBOSE = *vptr

	f := getFile(fptr)
	defer f.Close()
	busArr := getBusArr(bufio.NewScanner(f))
	maxBus, maxBusIdx := getMaxBus(busArr)
	chanArr := getChanArr(busArr, &maxBusIdx)

	quit := make(chan bool)
	start := time.Now()
	//DEPRECATED: go sendTimes(busArr, chanArr, *minptr, quit)

	// only send timestamps to the channel for the maximum bus number
	go sendToChan(maxBus, chanArr[maxBusIdx], maxBusIdx, *minptr)
	go monitorChan(chanArr[maxBusIdx], maxBusIdx, busArr, quit, &maxBus)

	finished := false
	if VERBOSE {
		fmt.Println("main():\t\tmonitoring quit channel ")
	}
	for !finished {
		select {
		case <-quit:
			finished = true
			quit <-finished
		}
	}
	elapsed := time.Since(start)
	if VERBOSE {
		fmt.Println("main():\t\toutside of finished loop with elapsed time", elapsed)
	}
}

// DEPRECATED
func sendTimes(busArr []*int, chanArr []chan *BusTimestamp, min int, q chan bool) {
	if VERBOSE {
		fmt.Println("sendTimes():\tsending timestamps to channels with min", min)
	}
	for i, bus := range busArr {
		if bus == nil {
			continue
		}
		go sendToChan(*bus, chanArr[i], i, min)
		go monitorChan(chanArr[i], i, busArr, q, nil)
	}
}


func sendToChan(busNum int, bchan chan *BusTimestamp, chanIdx int, minTime int) {
	i := 1
	modCorrection := (busNum * i + minTime) % busNum
	earliest := minTime + (busNum * i) - modCorrection

	if VERBOSE {
		fmt.Printf("sendToChan():\tfor bus %d at index %d starting time %d\n", busNum, chanIdx, earliest)
	}
	bchan <- &BusTimestamp{bus: busNum, timestamp: earliest, elapsed: 0}
	for {
		if VERBOSE && i % 100000 == 0 {
			// every 10,000 iterations, print an update
			fmt.Printf("sendToChan():\tfor bus %d at index %d sending time %d\n", busNum, chanIdx, earliest)
		}
		earliest = earliest + busNum
		bchan <- &BusTimestamp{bus: busNum, timestamp: earliest, elapsed: 0}
		i++
	}
}


func monitorChan(bchan chan *BusTimestamp, chanIdx int, busArr []*int, quit chan bool, maxBusNum *int) {
	if VERBOSE {
		fmt.Println("monitorChan():\tcheck chan at index", chanIdx)
	}
	finished := false
	for !finished {
		select {
		case bPtr := <-bchan:
			go checkBusTime(*bPtr, busArr, quit, maxBusNum)

		case <-quit:
			finished = true
			quit <- finished
		}
	}
	if VERBOSE {
		fmt.Println("monitorChan():\toutside of finished loop for chanIdx", chanIdx)
	}
}


func checkBusTime(b BusTimestamp, busArr []*int, q chan bool, maxBusNum *int) {
	start := time.Now()
	finished, firstBusMod := isFinished(b.timestamp, busArr, maxBusNum)
	b.elapsed = time.Since(start).Microseconds()

	if finished == true {
		printAnswer(b, busArr, q, maxBusNum, firstBusMod)
	}
}


func printAnswer(b BusTimestamp, busArr []*int, q chan bool, maxBusNum *int, firstBusMod *int) {
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

	if VERBOSE {
		fmt.Println("printAnswer():\tsending 'true' to channel q now")
	}
	q <- true
}


func isFinished (earliest int, busArr []*int, maxBusNum *int) (bool, *int) {
	busModInfoArr := getBusModInfoArr(earliest, busArr, maxBusNum)
	return busesAreConsecutive(busModInfoArr)
}

func getBusModInfoArr(earliest int, busArr[]*int, maxBusNum *int) []BusTimestamp { //map[int][]int {
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
