package main

import (
	"bufio"
	"flag"
	"fmt"
	"math"
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

var elapsedChan = make(chan BusTimestamp)

func main() {
	fptr := flag.String("file", "input.txt", "file path to read from")
	vptr := flag.Bool("v", false, "verbose")
	minptr := flag.Int("min", 1, "start checking at this timestamp")
	flag.Parse()
	VERBOSE = *vptr

	f := getFile(fptr)
	defer f.Close()
	busArr := getBusArr(bufio.NewScanner(f))
	//minBus, minBusIdx := getMinBus(busArr)
	chanArr := getChanArr(busArr, nil)

	quit := make(chan bool)
	start := time.Now()
	go sendTimes(busArr, chanArr, *minptr, quit)

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


func sendTimes(busArr []*int, chanArr []chan *BusTimestamp, min int, q chan bool) {
	if VERBOSE {
		fmt.Println("sendTimes():\tsending timestamps to channels with min", min)
	}
	for i, bus := range busArr {
		if bus == nil {
			continue
		}
		go sendToChan(*bus, chanArr[i], i, min)
		go monitorChan(chanArr[i], i, busArr, q)
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
		earliest = earliest + busNum
		bchan <- &BusTimestamp{bus: busNum, timestamp: earliest, elapsed: 0}
		i++
	}
}


func monitorChan(bchan chan *BusTimestamp, chanIdx int, busArr []*int, quit chan bool) {
	if VERBOSE {
		fmt.Println("monitorChan():\tcheck chan at index", chanIdx)
	}
	finished := false
	for !finished {
		select {
		case bPtr := <-bchan:
			go checkBusTime(*bPtr, busArr, quit)

		case <-quit:
			finished = true
			quit <- finished
		}
	}
	if VERBOSE {
		fmt.Println("monitorChan():\toutside of finished loop for chanIdx", chanIdx)
	}
}


func checkBusTime(b BusTimestamp, busArr []*int, q chan bool) {
	start := time.Now()
	finished := isFinished(b.timestamp, busArr)
	b.elapsed = time.Since(start).Microseconds()
	elapsedChan<-b

	if finished {
		printAnswer (b, busArr, q)
	}
}


func printAnswer (b BusTimestamp, busArr []*int, q chan bool) {
	busModInfoArr := getBusModInfoArr(b.timestamp, busArr)
	fmt.Printf("\n\ncheckBusTime():\tEARLIEST: %d FOR BUS %d\n", b.timestamp, b.bus)

	for i, busInfo := range busModInfoArr {
		if i == 0 {
			fmt.Printf("printAnswer():\tEARLIEST + mod %v for bus %d of index %d: ??\n", busInfo.modArr, busInfo.bus, i)
		}
		fmt.Printf("printAnswer():\t%d\t%d\t%v*\n", i, busInfo.bus, busInfo.modArr)

	}

	if VERBOSE {
		fmt.Println("printAnswer():\tsending 'true' to channel q now")
	}
	q <- true
}


func isFinished (earliest int, busArr []*int) bool {
	/*
	// DEPRECATED
	busModArr := getBusModArr(earliest, busArr)
	return isConsecutive(busModArr)
	*/

	busModInfoArr := getBusModInfoArr(earliest, busArr)
	return busesAreConsecutive(busModInfoArr)
}

func getBusModInfoArr(earliest int, busArr[]*int) []BusTimestamp { //map[int][]int {
	var busModInfoArr []BusTimestamp

	for idx, busPtr := range busArr {
		if busPtr == nil {
			continue
		}
		bus := *busPtr
		busInfo := BusTimestamp{
			bus: bus,
			timestamp: earliest,
			elapsed: 0,
			index: idx,
			modArr: make([]int, 5),
		}

		i := 1
		minInflation := int(math.Ceil(float64(idx/bus))) + i

		for minInflation*bus < idx {
			i++
			minInflation = int(math.Ceil(float64(idx/bus))) + i
		}

		// get five possible values to compare against others
		for j := 1; j<6; j++ {
			minuend := minInflation * j * bus
			subtrahend := earliest % bus
			busMod := minuend - subtrahend
			busInfo.modArr[j-1] = busMod
		}
		busModInfoArr = append(busModInfoArr, busInfo)
	}
	return busModInfoArr
}


func busesAreConsecutive(busInfoArr []BusTimestamp) bool {

	firstBus := busInfoArr[0] // assumes first element in busArr != nil

	for _, firstBusMod := range firstBus.modArr {
		firstBusModDiff := firstBusMod - firstBus.index

		for _, nextBus := range busInfoArr[1:] {
			nextIsConsecutive := false
			for _, nextBusMod := range nextBus.modArr {
				if nextBusMod - nextBus.index == firstBusModDiff {
					nextIsConsecutive = true
					break
				}
				if nextBusMod - nextBus.index > firstBusModDiff {
					break
				}
			}
			if nextIsConsecutive == false {
				return false
			}
		}
		return true
	}
	return false
}

/*
// DEPRECATED

func printAnswerOld(b BusTimestamp, busArr []*[]int, q chan bool) {
	busModArr := getBusModArr(b.timestamp, busArr)
	fmt.Printf("\n\ncheckBusTime():\tEARLIEST: %d FOR BUS %d\n", b.timestamp, b.bus)
	for i, busMod := range busModArr {
		if busMod == nil {
			continue
		}
		if i == 0 {
			modFactor := (*busMod)[1]
			fmt.Printf("checkBusTime():\tEARLIEST + mod %d for bus %d of index %d: %d\n", modFactor, (*busModArr[0])[0], i, b.timestamp+modFactor)
		}
		busNum := (*busMod)[0]
		if b.bus == busNum {
			fmt.Printf("checkBusTime():\t%d\t%v*\n", i, *busMod)
		} else {
			fmt.Printf("checkBusTime():\t%d\t%v\n", i, *busMod)
		}
	}
	if VERBOSE {
		fmt.Println("checkBusTime():\tsending 'true' to channel q now")
	}
	q <- true
}

func getBusModArr(earliest int, busArr []*int) []*[]int {
	var busModArr []*[]int
	for idx, busPtr := range busArr {
		if busPtr == nil {
			busModArr = append(busModArr, nil)
			continue
		}
		i := 1
		inflation := int(math.Ceil(float64(idx/(*busPtr)))) + i

		for inflation*(*busPtr) < idx {
			i++
			inflation = int(math.Ceil(float64(idx/(*busPtr)))) + i
		}

		busMod := inflation*(*busPtr) - (earliest % *busPtr)
		busModArr = append(busModArr, &[]int{*busPtr, busMod})
	}
	printArr(busModArr)
	return busModArr
}


// todo: isConsecutive can add to a chan int for that bus, representing num times checked
func isConsecutive(busModArr []*[]int) bool {
	busModDiff := (*busModArr[0])[1]
	for idx, busMod := range busModArr {
		if busMod == nil {
			continue
		}
		if (*busMod)[1] - idx != busModDiff {
			return false
		}
		busModDiff = (*busMod)[1]
	}
	return true
}

 */

