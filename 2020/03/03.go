package main

import (
	"bufio"
	"flag"
	"fmt"
	"log"
	"os"
	"strings"
)

type Slope struct {
	right int
	down int
}

var slopes = []Slope{
	Slope{1, 1},
	Slope{3, 1},
	Slope{5, 1},
	Slope{7, 1},
	Slope{1, 2},
}

func buildArrOfRows() [][]string{
	fptr := flag.String("file", "input.txt", "file path to read from")
	flag.Parse()

	f, err := os.Open(*fptr)
	if err != nil {
		log.Fatal(err)
	}
	defer func() {
		if err = f.Close(); err != nil {
			log.Fatal(err)
		}
	}()
	s := bufio.NewScanner(f)

	arrOfRows := make([][]string, 0)

	for s.Scan() {
		rowStr := s.Text()
		rowStr = strings.Repeat(rowStr, 323)

		rowArr := make([]string, 0)
		for _, c := range rowStr {
			rowArr = append(rowArr, string(c))
		}

		arrOfRows = append(arrOfRows, rowArr)
	}
	return arrOfRows
}

type Location struct {
	X int
	Y int
}

func getNumTreesHitForSlope(arrOfRows [][]string, slope Slope) int {
	numTreesHit := 0
	currLoc := Location{0,0}

	for currLoc.Y < len(arrOfRows) - slope.down && currLoc.X < len(arrOfRows[0]) - slope.right {
		newX := currLoc.X + slope.right
		newY := currLoc.Y + slope.down
		newLoc := Location{newX, newY}

		if arrOfRows[newY][newX] == "#" {
			arrOfRows[newY][newX] = "X"
			numTreesHit++
		} else {
			arrOfRows[newY][newX] = "O"
		}
		currLoc = newLoc
	}
	return numTreesHit
}

func main() {
	arrOfRows := buildArrOfRows()
	product := 1

	for _, slope := range slopes {
		treesHit := getNumTreesHitForSlope(arrOfRows, slope)
		product = product*treesHit
		fmt.Printf("right %d | down %d | treesHit %d | product %d\n", slope.right, slope.down, treesHit, product)
	}
	fmt.Println(slopes)
}