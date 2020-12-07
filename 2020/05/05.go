package main

import (
	"bufio"
	"fmt"
)

type Seat struct {
	row int
	col int
	id *int
}

func (s *Seat) print () {
	fmt.Printf("row %d, column %d, seat ID %d\n", s.row, s.col, s.getId())
}

func (s Seat) getId() int {
	sid := s.id
	if s.id == nil || *sid != 8*s.row + s.col {
		id := 8*s.row + s.col
		s.id = &id
	}
	return *s.id
}



func getMaxId(sarr []*Seat) int {
	max := 0
	for _, s := range sarr {
		if s != nil && s.getId() > max {
			max = s.getId()
		}
	}
	return max
}

//// this seat array is sorted
//func getEmptySeats(sarr []Seat) int {
//
//}

func main() {
	f := getFile()
	defer f.Close()
	seatArr := makeDataArr(bufio.NewScanner(f))
	fmt.Printf("%d seats\n", len(seatArr)) // 756 from part 1

	maxId := getMaxId(seatArr)
	fmt.Printf("maxId: %d\n", maxId)
	//826, with max possible = 8*127 + 7 = 1023

	for i := maxId - 756; i < maxId; i++ {
		if seatArr[i] == nil {
			fmt.Printf("%d is empty\n", i)  // 678
		}
	}
}