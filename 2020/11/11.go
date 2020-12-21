package main

import (
	"bufio"
	"flag"
	"fmt"
)

var VERBOSE = false


type WaitingArea struct{
	strRowArr [][]string
}

func (wa *WaitingArea) getCopy() [][]string {
	oldStrRowArr := make([][]string, 0)

	for _, strRow := range wa.strRowArr {
		oldStrRow := make([]string, len(strRow))
		_ = copy(oldStrRow, strRow)
		oldStrRowArr = append(oldStrRowArr, oldStrRow)
	}
	return oldStrRowArr
}

func (wa *WaitingArea) getSpaceAt (x int, y int)  *string {
	if y < 0 || x < 0 || y >= len(wa.strRowArr) || x >= len(wa.strRowArr[y]){
		return nil
	}
	return &wa.strRowArr[y][x]
}

func (wa *WaitingArea) getFirstEmptyOrOccupiedSeatAt (x int, y int) {

}


func (wa *WaitingArea) getVisibleNeighbors (x int, y int) []*string {
	// top
	nextY := y-1
	topStr := wa.getSpaceAt(x, nextY)
	for topStr !=nil && *topStr == "." {
		nextY--
		topStr = wa.getSpaceAt(x, nextY)
	}

	// top right
	nextY = y-1
	nextX := x+1
	topRightStr := wa.getSpaceAt(nextX, nextY)
	for topRightStr !=nil && *topRightStr == "." {
		nextY--
		nextX++
		topRightStr = wa.getSpaceAt(nextX, nextY)
	}

	// right
	nextX = x+1
	rightStr := wa.getSpaceAt(nextX, y)
	for rightStr !=nil && *rightStr == "." {
		nextX++
		rightStr = wa.getSpaceAt(nextX, y)
	}

	// bottom right
	nextX = x+1
	nextY = y+1
	bottomRightStr := wa.getSpaceAt(nextX, nextY)
	for bottomRightStr !=nil && *bottomRightStr == "." {
		nextY++
		nextX++
		bottomRightStr = wa.getSpaceAt(nextX, nextY)
	}

	// bottom
	nextY = y+1
	bottomStr := wa.getSpaceAt(x, nextY)
	for bottomStr !=nil && *bottomStr == "." {
		nextY++
		bottomStr = wa.getSpaceAt(x, nextY)
	}

	// bottom left
	nextY = y+1
	nextX = x-1
	bottomLeftStr := wa.getSpaceAt(nextX, nextY)
	for bottomLeftStr !=nil && *bottomLeftStr == "." {
		nextY++
		nextX--
		bottomLeftStr = wa.getSpaceAt(nextX, nextY)
	}

	// left
	nextX = x-1
	leftStr := wa.getSpaceAt(nextX, y)
	for leftStr !=nil && *leftStr == "." {
		nextX--
		leftStr = wa.getSpaceAt(nextX, y)
	}

	nextY = y-1
	nextX = x-1
	topLeftStr := wa.getSpaceAt(nextX, nextY)
	for topLeftStr !=nil && *topLeftStr == "." {
		nextX--
		nextY--
		topLeftStr = wa.getSpaceAt(nextX, nextY)
	}

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

func (wa *WaitingArea) emptySpaceAround (x int, y int) bool {
	for _, str := range wa.getVisibleNeighbors(x, y) {
		if str != nil && *str == "#" {
			return false
		}
	}
	return true
}

func (wa *WaitingArea) hasAtLeastNVisibleAdjacentOccupied (x int, y int, n int) bool {
	occupiedCount := 0
	for _, str := range wa.getVisibleNeighbors(x, y) {
		if str != nil && *str == "#" {
			occupiedCount++
		}
	}
	return occupiedCount >= n
}

func (wa *WaitingArea) fillSeats() {
	newStrRowArr := wa.getCopy()
	for yIdx, strRow := range wa.strRowArr {
		for xIdx, seatOrSpace := range strRow {
			switch seatOrSpace {
			case "L":
				if wa.emptySpaceAround(xIdx, yIdx) {
					newStrRowArr[yIdx][xIdx] = "#"
				}
			case ".":
				continue
			case "#":
				if wa.hasAtLeastNVisibleAdjacentOccupied(xIdx, yIdx, 5) {
					newStrRowArr[yIdx][xIdx] = "L"
				}
			}
		}
	}
	wa.strRowArr = newStrRowArr
}

func (wa *WaitingArea) isEqTo (oldStrRowArr [][]string) bool {
	for yIdx, row := range wa.strRowArr {
		for xIdx, space := range row {
			if space != oldStrRowArr[yIdx][xIdx] {
				return false
			}
		}
	}
	return true
}

func (wa *WaitingArea) iterateUntilDone (oldStrRowArr [][]string, numIterations *int) (*WaitingArea, int){
	if !wa.isEqTo(oldStrRowArr) {
		*numIterations++
		nextOldStrRowArr := wa.getCopy()
		fmt.Printf("\nnumIterations %d and nextOldStrRowArr is...\n", numIterations)
		for _, elem := range nextOldStrRowArr {
			fmt.Println(elem)
		}
		wa.fillSeats()
		wa.iterateUntilDone(nextOldStrRowArr, numIterations)
	}
	return wa, *numIterations
}

func main() {
	fptr := flag.String("file", "input.txt", "file path to read from")
	vptr := flag.Bool("v", false, "verbose")
	flag.Parse()
	VERBOSE = *vptr

	f := getFile(fptr)
	defer f.Close()
	waitingArea := &WaitingArea{
		strRowArr: makeDataArr(bufio.NewScanner(f)),
	}

	for _, elem := range waitingArea.strRowArr {
		fmt.Println(elem)
	}
	oldStrRowArr := waitingArea.getCopy()
	waitingArea.fillSeats()
	one := 1
	finalWaitingArea, numIterations := waitingArea.iterateUntilDone(oldStrRowArr, &one)
	fmt.Printf("\nnumIterations %d and final configuration is:\n", numIterations)
	for _, elem := range finalWaitingArea.strRowArr {
		fmt.Println(elem)
	}
	fmt.Printf("\nnumber of occupied seats is %d\n", finalWaitingArea.getNumOccupiedSeats())
}