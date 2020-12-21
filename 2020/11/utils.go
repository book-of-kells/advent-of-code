package main

import (
	"bufio"
	"fmt"
	"log"
	"os"
	"strings"
)

func mungeData(dataStr string) []string {
	// split into array of characters
	return strings.SplitN(dataStr, "", len(dataStr))
}

func makeDataArr(s *bufio.Scanner) [][]string {
	dataArr := make([][]string, 0)
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

func getFile(fptr *string) *os.File {

	f, err := os.Open(*fptr)
	if err != nil {
		log.Fatal(err)
	}
	return f
}

func (wa *WaitingArea) getNeighborsAt (x int, y int) []*string {
	topStr := wa.getSpaceAt(x, y-1)
	topRightStr := wa.getSpaceAt(x+1, y-1)
	rightStr := wa.getSpaceAt(x+1, y)
	bottomRightStr := wa.getSpaceAt(x+1, y+1)
	bottomStr := wa.getSpaceAt(x, y+1)
	bottomLeftStr := wa.getSpaceAt(x-1, y+1)
	leftStr := wa.getSpaceAt(x-1, y)
	topLeftStr := wa.getSpaceAt(x-1, y-1)

	return []*string{
		topStr,
		topRightStr,
		rightStr,
		bottomRightStr,
		bottomStr,
		bottomLeftStr,
		leftStr,
		topLeftStr,
	}
}

func (wa *WaitingArea) emptySpaceAroundPartOne (x int, y int) bool {
	for _, str := range wa.getNeighborsAt(x, y) {
		if str != nil && *str == "#" {
			return false
		}
	}
	return true
}

// for final answer
func (wa *WaitingArea) getNumOccupiedSeats () int {
	numOccupiedSeats := 0
	for _, strRow := range wa.strRowArr {
		for _, seatOrSpace := range strRow {
			if &seatOrSpace != nil && seatOrSpace == "#" {
				numOccupiedSeats++
			}
		}
	}
	return numOccupiedSeats
}

func (wa *WaitingArea) hasAtLeastNAdjacentOccupied (x int, y int, n int) bool {
	occupiedCount := 0
	for _, str := range wa.getNeighborsAt(x, y) {
		if str != nil && *str == "#" {
			occupiedCount++
		}
	}
	return occupiedCount >= n
}